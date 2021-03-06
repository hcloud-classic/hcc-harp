package iplink

import (
	"fmt"
	"hcc/harp/lib/iplinkext"
	"strconv"
)

func generateMACAddress(ip string) string {
	// iplinkext.getIfaceVNUM() will return between from 0 to 1020
	var vnumStr = fmt.Sprintf("%04d", iplinkext.GetIfaceVNUM(ip))
	bytes := []byte(vnumStr)

	newMAC := "68:61:" + // h:a: (2 letters from harp)
		strconv.Itoa(int(bytes[0])) + ":" +
		strconv.Itoa(int(bytes[1])) + ":" +
		strconv.Itoa(int(bytes[2])) + ":" +
		strconv.Itoa(int(bytes[3]))

	return newMAC
}
