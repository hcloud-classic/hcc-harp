package dhcpd

import (
	"sync"
	"time"
)

var writeDHCPConfigLocked = make(map[string]bool)
var lockStateListMutex = sync.Mutex{}

func isWriteDHCPConfigLocked(serverUUID string) bool {
	lockStateListMutex.Lock()
	isLocked, exist := writeDHCPConfigLocked[serverUUID]
	if !exist {
		lockStateListMutex.Unlock()
		return false
	}

	lockStateListMutex.Unlock()
	return isLocked
}

func writeDHCPConfigLock(serverUUID string) {
	lockStateListMutex.Lock()
	writeDHCPConfigLocked[serverUUID] = true
	lockStateListMutex.Unlock()
}

func writeDHCPConfigUnlock(serverUUID string) {
	lockStateListMutex.Lock()
	delete(writeDHCPConfigLocked, serverUUID)
	lockStateListMutex.Unlock()
}

func waitWriteDHCPConfigUnlock(serverUUID string) {
	for true {
		if !isWriteDHCPConfigLocked(serverUUID) {
			break
		}

		time.Sleep(1 * time.Second)
	}
}
