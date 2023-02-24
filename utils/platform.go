// +build linux darwin freebsd openbsd netbsd solaris

package lxdapi


import (
	"io/fs"
	"syscall"
)

func GetUidGid(stat fs.FileInfo) (int64, int64, error) {
	if linuxstat, ok := stat.Sys().(*syscall.Stat_t); ok {
		UID := int64(linuxstat.Uid)
		GID := int64(linuxstat.Gid)
		return UID, GID, nil

	}
	return 0, 0, nil
}