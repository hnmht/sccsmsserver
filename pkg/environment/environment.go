package environment

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"runtime"
	"sort"
	"strings"

	"github.com/jaypipes/ghw"
	"go.uber.org/zap"
)

// Get computer network card MAC address text
func GetMacArray() (mac string, err error) {
	net, err := ghw.Network()
	if err != nil {
		return "macarrayseacloud2023", nil
	}
	var macSlice = make([]string, 0)
	for _, nic := range net.NICs {
		if !nic.IsVirtual {
			macSlice = append(macSlice, nic.MacAddress)
		}
	}
	sort.StringSlice(macSlice).Sort()

	return strings.Join(macSlice, ","), nil
}

// Get computer motherboard serial number
func GetSerialNumber() (serialNumber string, err error) {
	baseboard, err := ghw.Baseboard()
	if err != nil {
		return "serialseacloud2023", nil
	}
	return baseboard.SerialNumber, nil
}

// Get computer hash value
func GetMachineHash(macArray string, serialNumber string) (machineahash string, err error) {
	var build strings.Builder
	build.WriteString("macArray:")
	build.WriteString(macArray)
	build.WriteString("serialNumber:")
	build.WriteString(serialNumber)

	h := sha256.New()

	_, err = h.Write([]byte(build.String()))
	if err != nil {
		zap.L().Error("GetMachineHash h.Write failed:", zap.Error(err))
		return "machinehashseacloud2023", nil
	}
	b := h.Sum(nil)
	machineahash = hex.EncodeToString(b)
	return
}

// Get computer IPV4 address list
func GetLocalIPs() (ipList []net.IP) {
	ifs, err := net.Interfaces()
	if err != nil {
		zap.L().Error("GetLocalIPsnet.Interfaces failed:", zap.Error(err))
		return
	}
	for _, interf := range ifs {
		addrs, err := interf.Addrs()
		if err != nil {
			continue
		}
		if runtime.GOOS == "windows" && interf.Flags&net.FlagUp == 0 {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ipList = append(ipList, ip)
		}
	}
	return ipList
}
