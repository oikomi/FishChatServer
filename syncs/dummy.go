// +build !deadlock

package syncs

import "sync"

type Mutex struct {
	sync.Mutex
}

type RWMutex struct {
	sync.RWMutex
}
