package fileutils

import (
	"io/fs"
	"os"
	"strings"
)

type Callback func(fileinfo fs.DirEntry) bool

func GetFiles(rootpath, prefix string, callback Callback) error {
	infos, err := os.ReadDir(rootpath)
	if err != nil {
		return err
	}
	for _, i := range infos {
		if !strings.HasPrefix(i.Name(), prefix) {
			ok := callback(i)
			if !ok {
				return nil
			}
		}
	}
	return nil
}
