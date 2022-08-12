package xutils

import (
	"errors"
	"github.com/3th1nk/cidr"
	"math"
	"net"
)

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

func IpCidr2IpRange(ipCidr string) (ipStart, ipEnd string, err error) {
	c, err := cidr.ParseCIDR(ipCidr)
	if nil != err {
		return "", "", err
	}
	if c.IsIPv4() {
		nIpStart, _ := IPString2Long(c.Network())
		nIpStart++
		ipStart, _ := Long2IPString(nIpStart)

		nIpEnd, _ := IPString2Long(c.Broadcast())
		nIpEnd--
		ipEnd, _ := Long2IPString(nIpEnd)
		return ipStart, ipEnd, nil
	}
	return "", "", err
}
