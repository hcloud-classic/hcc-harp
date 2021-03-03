package dhcpdext

import (
	"sync"
	"sync/atomic"
)

// CreatingSubnetConfigCounter : Counter will increase while creating subnet config file.
// Must be 0 when there are no creating works.
var CreatingSubnetConfigCounter int64 = 0

// IncCreatingSubnetConfigCounter : Increase the value of CreatingSubnetConfigCounter
func IncCreatingSubnetConfigCounter() {
	atomic.AddInt64(&CreatingSubnetConfigCounter, 1)
}

// DecCreatingSubnetConfigCounter : Decrease the value of CreatingSubnetConfigCounter
func DecCreatingSubnetConfigCounter() {
	atomic.AddInt64(&CreatingSubnetConfigCounter, -1)
}

// HarpDHCPDConfigWriteLock : Lock for writing harp_dhcpd.conf file
var HarpDHCPDConfigWriteLock sync.Mutex
