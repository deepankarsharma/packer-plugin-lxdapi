package lxdapi

import (
	"io/fs"
)

func GetUidGid(stat fs.FileInfo) (int64, int64, error) {
	return 0, 0, nil
}