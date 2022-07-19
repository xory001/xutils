package xutils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math"
	"net"
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

func IPString2Long(ip string) (uint32, error) {
	b := net.ParseIP(ip).To4()
	if b == nil {
		return 0, errors.New("invalid ipv4 format")
	}

	return (uint32)(uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24), nil
}

func Long2IPString(i uint32) (string, error) {
	if i > math.MaxUint32 {
		return "", errors.New("beyond the scope of ipv4")
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip.String(), nil
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
