package vnstat

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var updateVnStatLocked = make(map[string]bool)
var lockStateListMutex = sync.Mutex{}

func delaySecond(n time.Duration) {
	time.Sleep(n * time.Second)
}

func isUpdateVnStatQueued(harpIfaceName string) bool {
	lockStateListMutex.Lock()
	_, queued := updateVnStatLocked[harpIfaceName]
	lockStateListMutex.Unlock()
	return queued
}

func isUpdateVnStatLocked(harpIfaceName string) bool {
	lockStateListMutex.Lock()
	isLocked, exist := updateVnStatLocked[harpIfaceName]
	if !exist {
		lockStateListMutex.Unlock()
		return false
	}

	lockStateListMutex.Unlock()
	return isLocked
}

func updateVnStatLock(harpIfaceName string) {
	lockStateListMutex.Lock()
	updateVnStatLocked[harpIfaceName] = true
	lockStateListMutex.Unlock()
}

func updateVnStatUnlock(harpIfaceName string) {
	lockStateListMutex.Lock()
	updateVnStatLocked[harpIfaceName] = false
	lockStateListMutex.Unlock()
}

func updateVnStat(harpIfaceName string) {
	cmd := exec.Command("sh", "-c", "vnstat -u -i "+harpIfaceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println("Failed to update vnStat for " + harpIfaceName + " (" + string(output) + ")")
	}
}

func deleteVnStat(harpIfaceName string) {
	cmd := exec.Command("sh", "-c", "vnstat --delete --force -i "+harpIfaceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println("Failed to delete vnStat for " + harpIfaceName + " (" + string(output) + ")")
	}
}

func queueUpdateVnStat(harpIfaceName string) {
	go func() {
		if config.VnStat.Debug == "on" {
			logger.Logger.Println("queueUpdateVnStat(): Queued of running updateVnStat() after " + strconv.Itoa(int(config.VnStat.DatabaseUpdateIntervalSec)) + "sec" +
				" for harpIfaceName=" + harpIfaceName)
		}
		delaySecond(time.Duration(config.VnStat.DatabaseUpdateIntervalSec))
		ScheduleUpdateVnStat(harpIfaceName, false)
	}()
}

// ScheduleUpdateVnStat : Schedule database update for the Harp's internal interface
func ScheduleUpdateVnStat(harpIfaceName string, isNew bool) {
	if !isNew && !isUpdateVnStatQueued(harpIfaceName) {
		if config.VnStat.Debug == "on" {
			logger.Logger.Println("updateVnStat(): updateVnStat is canceled cause of queue is deleted for harpIfaceName=" + harpIfaceName)
		}
		logger.Logger.Println("updateVnStat(): Removing traffic stats from VnStat database for for harpIfaceName=" + harpIfaceName)
		deleteVnStat(harpIfaceName)

		return
	}

	if isUpdateVnStatLocked(harpIfaceName) {
		if config.VnStat.Debug == "on" {
			logger.Logger.Println("ScheduleUpdateVnStat(): Locked for harpIfaceName=" + harpIfaceName)
		}
		for true {
			if !isUpdateVnStatLocked(harpIfaceName) {
				break
			}
			if config.VnStat.Debug == "on" {
				logger.Logger.Println("ScheduleUpdateVnStat(): Rerun after " +
					strconv.Itoa(int(config.VnStat.DatabaseUpdateIntervalSec)) + "sec for harpIfaceName=" + harpIfaceName)
			}
			delaySecond(time.Duration(config.VnStat.DatabaseUpdateIntervalSec))
		}
	}

	go func() {
		updateVnStatLock(harpIfaceName)
		if config.VnStat.Debug == "on" {
			logger.Logger.Println("ScheduleUpdateVnStat(): Running UpdateVnStat() for harpIfaceName=" + harpIfaceName)
		}
		updateVnStat(harpIfaceName)
		updateVnStatUnlock(harpIfaceName)
	}()

	queueUpdateVnStat(harpIfaceName)
}

// RemoveUpdateVnStat : Remove Harp's internal interface from VnStat database
func RemoveUpdateVnStat(harpIfaceName string) {
	for true {
		if !isUpdateVnStatLocked(harpIfaceName) {
			lockStateListMutex.Lock()
			delete(updateVnStatLocked, harpIfaceName)
			lockStateListMutex.Unlock()

			break
		}

		delaySecond(1)
	}
}
