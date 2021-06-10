package iplinkext

import (
	"strconv"
	"strings"
)

// GetIfaceVNUM : Generate VNUM of the interface by IP address
func GetIfaceVNUM(ip string) (vnum int) {
	var ifaceVNUM = 0

	ipSplit := strings.Split(ip, ".")
	for _, ipSplited := range ipSplit {
		ipSplitedInt, _ := strconv.Atoi(ipSplited)
		ifaceVNUM += ipSplitedInt
	}

	return ifaceVNUM
}
