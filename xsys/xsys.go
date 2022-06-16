package xsys

import (
	"os"
	"sort"
)

func ReadDirAscByTime(name string) ([]os.DirEntry, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dirs, err := f.ReadDir(-1)
	sort.Slice(dirs, func(i, j int) bool {
		fileInfo1, err1 := dirs[i].Info()
		fileInfo2, err2 := dirs[j].Info()
		if (nil == err1) && (nil == err2) {
			return fileInfo1.ModTime().Before(fileInfo2.ModTime())
		}
		return dirs[i].Name() < dirs[j].Name()
	})
	return dirs, err
}

func ReadDirDescByTime(name string) ([]os.DirEntry, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dirs, err := f.ReadDir(-1)
	sort.Slice(dirs, func(i, j int) bool {
		fileInfo1, err1 := dirs[i].Info()
		fileInfo2, err2 := dirs[j].Info()
		if (nil == err1) && (nil == err2) {
			return fileInfo1.ModTime().After(fileInfo2.ModTime())
		}
		return dirs[i].Name() < dirs[j].Name()
	})
	return dirs, err
}
