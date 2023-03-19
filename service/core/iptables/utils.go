package iptables

import (
	"net"
	"strconv"

	"github.com/v2rayA/v2rayA/common"
	"github.com/v2rayA/v2rayA/common/cmds"
	"github.com/v2rayA/v2rayA/conf"
	"golang.org/x/net/nettest"
)

func IPNet2CIDR(ipnet *net.IPNet) string {
	ones, _ := ipnet.Mask.Size()
	return ipnet.IP.String() + "/" + strconv.Itoa(ones)
}

func GetLocalCIDR() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	var cidrs []string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			cidrs = append(cidrs, IPNet2CIDR(ipnet))
		}
	}
	return cidrs, nil
}

func IsIPv6Supported() bool {
	switch conf.GetEnvironmentConfig().IPV6Support {
	case "on":
		return true
	case "off":
		return false
	default:
	}
	if common.IsDocker() {
		return false
	}
	if !nettest.SupportsIPv6() {
		return false
	}
	if _, isNft := Redirect.(*nftRedirect); isNft {
		return true
	}
	return cmds.IsCommandValid("ip6tables") || cmds.IsCommandValid("ip6tables-nft")
}

func IsNFTablesSupported() bool {

	switch conf.GetEnvironmentConfig().NFTablesSupport {
	// Warning:
	// This is an experimental feature for nftables support.
	// The default value is "off" for now but may be changed to "auto" in the future
	case "on":
		return true
	case "off":
		return false
	default:
	}
	if common.IsDocker() {
		return false
	}
	return cmds.IsCommandValid("nft")
}
