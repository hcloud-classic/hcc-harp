package dhcpdext

import (
	"sync"
	"sync/atomic"
)

// WritingSubnetConfigCounter : Counter will increase while creating subnet config file.
// Must be 0 when there are no creating works.
var WritingSubnetConfigCounter int64

// IncWritingSubnetConfigCounter : Increase the value of WritingSubnetConfigCounter
func IncWritingSubnetConfigCounter() {
	atomic.AddInt64(&WritingSubnetConfigCounter, 1)
}

// DecWritingSubnetConfigCounter : Decrease the value of WritingSubnetConfigCounter
func DecWritingSubnetConfigCounter() {
	atomic.AddInt64(&WritingSubnetConfigCounter, -1)
}

// HarpDHCPDConfigWriteLock : Lock for writing harp_dhcpd.conf file
var HarpDHCPDConfigWriteLock sync.Mutex
