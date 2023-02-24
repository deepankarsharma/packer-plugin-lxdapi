package lxdapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"errors"
	"log"
	"net"
	"os"
	"path"
	"regexp"
	"strings"
	"syscall"
	"time"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

var lxdLogFilenameRegex = regexp.MustCompile(`/1.0/instances/\S+/logs/(\S+)`)

func GetLXDInstanceServer(lxd_uri string) (lxd.InstanceServer, error) {
	c, err := lxd.ConnectLXDUnix(lxd_uri, nil)
	if err != nil {
		return c, err
	}
	return c, nil
}

type LXDUtils struct {
	lxd lxd.InstanceServer
}

func NewLXDUtilStruct(lxd_uri string) (*LXDUtils, error) {
	c, err := GetLXDInstanceServer(lxd_uri)
	if err != nil {
		return nil, err
	}
	n := new(LXDUtils)
	n.lxd = c
	return n, nil 
}


func (u *LXDUtils) CreateInstance(requestStruct api.InstancesPost) error {
	
	op, err := u.lxd.CreateInstance(requestStruct)
	if err != nil {
		return err
	}
	
	err = op.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (u *LXDUtils) StartInstance(instanceName string) error {
	reqState := api.InstanceStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err := u.lxd.UpdateInstanceState(instanceName, reqState, "")
	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (u *LXDUtils) StopInstance(instanceName string) error {
	reqState := api.InstanceStatePut{
		Action:  "stop",
		Timeout: -1,
	}
	op, err := u.lxd.UpdateInstanceState(instanceName, reqState, "")
	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		log.Println("Trying to stop instance that is already stopped")
	}

	return nil
}

func (u *LXDUtils) DeleteInstance(instanceName string) error {
	if err := u.StopInstance(instanceName); err != nil {
		return err
	}

	op, err := u.lxd.DeleteInstance(instanceName)
	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (u *LXDUtils) WaitInstanceIP(instanceName string, blacklist []string) (net.IP, error) {
	var ip net.IP
	var err error
	ip, err = u.GetInstanceLXDIP(instanceName, blacklist)
	for c := 0; c < 50 && err != nil; c++ {
		log.Default().Printf("waiting for %s to get an IP address...", instanceName)
		time.Sleep(2 * time.Second)
		ip, err = u.GetInstanceLXDIP(instanceName, blacklist)
	}
	if err != nil {
		return nil, err
	}
	return ip, nil
}

func (u *LXDUtils) GetInstanceLXDIP(instanceName string, blacklist []string) (net.IP, error) {
	in, _, err := u.lxd.GetInstanceFull(instanceName)
	if err != nil {
		return nil, fmt.Errorf("error getting instance: %w", err)
	}

	var ips []string
	for netName, net := range in.State.Network {
		if net.Type == "loopback" {
			continue
		}

		for _, addr := range net.Addresses {
			if addr.Scope == "link" || addr.Scope == "local" {
				continue
			}

			if strings.Contains(addr.Family, "inet") && netName != "cni0" && netName != "flannel.1" {
				blacklisted := false
				for _, black := range blacklist {
					if strings.Contains(addr.Address, black) {
						blacklisted = true
						break
					}
				}
				if strings.Count(addr.Address, ":") < 2 && !blacklisted {
					ips = append(ips, addr.Address)
				}
			}
		}
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("instance %s has no IP address", instanceName)
	}

	ip := net.ParseIP(ips[0])
	if ip == nil {
		return nil, fmt.Errorf("not a valid ip: %s", ips[0])
	}

	return ip, nil
}

func (u *LXDUtils) UploadFile(instanceName, fromFile, toDir string) error {
	var mode os.FileMode
	var toPath string
	UID := int64(0)
	GID := int64(0)

	stat, err := os.Stat(fromFile)
	if err != nil {
		return fmt.Errorf("cannot stat %s: %w", fromFile, err)
	}

	if linuxstat, ok := stat.Sys().(*syscall.Stat_t); ok {
		UID = int64(linuxstat.Uid)
		GID = int64(linuxstat.Gid)
	}
	mode = os.FileMode(0755)

	data, err := ioutil.ReadFile(fromFile)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", fromFile, err)
	}
	_, filename := path.Split(fromFile)
	toPath = path.Join(toDir, filename)

	err = u.RecursiveMkdir(instanceName, toDir, mode, UID, GID)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(data)

	args := lxd.InstanceFileArgs{
		Type:    "file",
		UID:     UID,
		GID:     GID,
		Mode:    int(mode.Perm()),
		Content: reader,
	}

	err = u.lxd.CreateInstanceFile(instanceName, toPath, args)
	if err != nil {
		return fmt.Errorf("cannot push %s to %s: %w", fromFile, toPath, err)
	}

	return nil
}

func (u *LXDUtils) RecursiveMkdir(instanceName, dir string, mode os.FileMode, UID, GID int64) error {
	if dir == "/" {
		return nil
	}

	dir = strings.TrimSuffix(dir, "/")

	split := strings.Split(dir[1:], "/")
	if len(split) > 1 {
		err := u.RecursiveMkdir(instanceName, "/"+strings.Join(split[:len(split)-1], "/"), mode, UID, GID)
		if err != nil {
			return err
		}
	}

	_, resp, err := u.lxd.GetInstanceFile(instanceName, dir)
	if err == nil && resp.Type == "directory" {
		return nil
	}
	if err == nil && resp.Type != "directory" {
		return fmt.Errorf("%s is not a directory", dir)
	}

	args := lxd.InstanceFileArgs{
		Type: "directory",
		UID:  UID,
		GID:  UID,
		Mode: int(mode.Perm()),
	}
	return u.lxd.CreateInstanceFile(instanceName, dir, args)
}

func (u *LXDUtils) UploadFiles(instanceName string, from []string, to string) error {
	for _, from := range from {
		err := u.UploadFile(from, to, instanceName)
		if err != nil {
			return err
		}
	}
	return nil
}

type ExecResult struct {
	Stdout string
	Stderr string
	ReturnCode float64
}

func (u *LXDUtils) Exec(instanceName, command string, environment map[string]string) (*ExecResult, error) {
	split := strings.Fields(command)

	op, err := u.lxd.ExecInstance(instanceName,
		api.InstanceExecPost{
			Command: split,
			Environment: environment,
			Interactive: false,
			RecordOutput: true,
		},
		&lxd.InstanceExecArgs{},
	)

	if err != nil {
		return nil, err
	}
	
	err = op.Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to exec %s: %w", command, err)
	}

	op_ := op.Get()
	output := op_.Metadata["output"].(map[string]interface{})
	outReader, err := u.lxd.GetInstanceLogfile(instanceName, lxdLogFilenameRegex.FindStringSubmatch(output["1"].(string))[1])
	if err != nil {
		return nil, err
	}
	errReader, err := u.lxd.GetInstanceLogfile(instanceName, lxdLogFilenameRegex.FindStringSubmatch(output["2"].(string))[1])
	if err != nil {
		return nil, err
	}

	outData, err := ioutil.ReadAll(outReader)
	if err != nil {
		return nil, err
	}
	outReader.Close()

	errData, err := ioutil.ReadAll(errReader)
	if err != nil {
		return nil, err
	}
	errReader.Close()

	stdout := string(outData)
	stderr := string(errData)

	code, ok := op_.Metadata["return"].(float64)
	if !ok {
		return &ExecResult{stdout, stderr, -1}, errors.New("failed to get exit code")
	}
	return &ExecResult{stdout, stderr, code}, nil
}

func (u *LXDUtils) CreateInstanceHL(req api.InstancesPost) error {
	err := u.CreateInstance(req)
	if err != nil {
		return err
	}

    err = u.StartInstance(req.Name)
	if err != nil {
		return err
	}

	_, err = u.WaitInstanceIP(req.Name, []string{})
	if err != nil {
		return err
	}
	return nil
}

func (u *LXDUtils) PublishInstanceAsImage(instanceName string, req api.ImagesPost) error {
	for _, alias := range req.Aliases {
		aliasEntry, _, _ := u.lxd.GetImageAlias(alias.Name)
		if aliasEntry != nil {
			err := u.lxd.DeleteImageAlias(alias.Name)
			if err != nil {
				return err
			}
		}
	} 

	op, err := u.lxd.CreateImage(req, nil)
	if err != nil {
		return err
	}

	// Wait for operation to finish
	err = op.Wait()
	if err != nil {
		return err
	}
	return nil
}