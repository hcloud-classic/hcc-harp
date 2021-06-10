package vnstat

import (
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
	"time"
)

func getVnStatJSONDataByDay(harpIface string) (string, error) {
	cmd := exec.Command("sh", "-c", "vnstat -i "+harpIface+" --json d")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(output))
	}

	return string(output), nil
}

// GetTodayVnStatData : Get today vnStat data
func GetTodayVnStatData(harpIface string) (txKB int64, rxKB int64, err error) {
	jsonData, err := getVnStatJSONDataByDay(harpIface)
	if err != nil {
		return 0, 0, err
	}

	var vnStat vnStat
	err = json.Unmarshal([]byte(jsonData), &vnStat)
	if err != nil {
		return 0, 0, err
	}

	currentTime := time.Now()
	mm, _ := strconv.Atoi(currentTime.Format("01"))
	dd, _ := strconv.Atoi(currentTime.Format("02"))

	for _, dayData := range vnStat.Interfaces[0].Traffic.Days {
		if dayData.Date.Month == mm && dayData.Date.Day == dd {
			return dayData.Tx, dayData.Rx, nil
		}
	}

	return 0, 0, nil
}
