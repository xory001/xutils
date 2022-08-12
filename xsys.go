package xutils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
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

func Md5sum(src string) string {
	h := md5.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

func UTCTime() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05.000 UTC")
}

func GetRoutineId() string {
	var buf [64]byte
	bytesRead := runtime.Stack(buf[:], false)
	if bytesRead > 0 {
		slice := strings.Split(string(buf[:]), " ")
		if len(slice) > 1 {
			return slice[1]
		}
	}
	return ""
}
