package utils

import (
	"fmt"
	"net"
)

func GetFirstNoneLoopIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		return "", fmt.Errorf("failed to fetch interfaces => %w", err)
	}

	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no first-none-loop ip found")
}
