package dhcpdext

import (
	"hcc/harp/lib/config"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetSubnetConfFiles : Get paths of harp's DHCPD configuration files
func GetSubnetConfFiles() ([]string, error) {
	var files []string

	if WritingSubnetConfigCounter != 0 {
		for true {
			if WritingSubnetConfigCounter == 0 {
				break
			}

			time.Sleep(time.Millisecond * 500)
		}
	}

	folder := config.DHCPD.ConfigFileLocation
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if !strings.Contains(path, "harp_dhcpd.conf") &&
			path != config.DHCPD.ConfigFileLocation {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
