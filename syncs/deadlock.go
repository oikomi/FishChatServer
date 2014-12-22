// +build deadlock

package syncs

import (
	"bytes"
	"container/list"
	"github.com/funny/debug"
	"github.com/funny/goid"
	"strconv"
	"sync"
)

type Mutex struct {
	monitor
	sync.Mutex
}

func (m *Mutex) Lock() {
	waitInfo := m.monitor.wait('w')
	m.Mutex.Lock()
	m.monitor.using(waitInfo)
}

func (m *Mutex) Unlock() {
	m.monitor.release('w')
	m.Mutex.Unlock()
}

type RWMutex struct {
	monitor
	sync.RWMutex
}

func (m *RWMutex) Lock() {
	waitInfo := m.monitor.wait('w')
	m.RWMutex.Lock()
	m.monitor.using(waitInfo)
}

func (m *RWMutex) Unlock() {
	m.monitor.release('w')
	m.RWMutex.Unlock()
}

func (m *RWMutex) RLock() {
	waitInfo := m.monitor.wait('r')
	m.RWMutex.RLock()
	m.monitor.using(waitInfo)
}

func (m *RWMutex) RUnlock() {
	m.monitor.release('r')
	m.RWMutex.RUnlock()
}

var (
	globalMutex = new(sync.Mutex)
	waitingList = make(map[int32]*lockUsage)
	titleStr    = []byte("[DEAD LOCK]\n")
	goStr       = []byte("goroutine ")
	waitStr     = []byte(" wait")
	holdStr     = []byte(" hold")
	readStr     = []byte(" read")
	writeStr    = []byte(" write")
	lineStr     = []byte{'\n'}
)

type lockUsage struct {
	monitor *monitor
	mode    byte
	goid    int32
	stack   debug.StackInfo
}

type monitor struct {
	holders *list.List
}

func (m *monitor) wait(mode byte) *lockUsage {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	waitInfo := &lockUsage{m, mode, goid.Get(), debug.StackTrace(3, 0)}
	waitingList[waitInfo.goid] = waitInfo

	if m.holders == nil {
		m.holders = list.New()
	}

	m.diagnose(mode, []*lockUsage{waitInfo})

	return waitInfo
}

func (m *monitor) diagnose(mode byte, waitLink []*lockUsage) {
	for i := m.holders.Front(); i != nil; i = i.Next() {
		holder := i.Value.(*lockUsage)
		if mode != 'r' || holder.mode != 'r' {
			// deadlock detected
			if holder.goid == waitLink[0].goid {
				deadlockPanic(waitLink)
			}
			// the lock holder is waiting for another lock
			if waitInfo, exists := waitingList[holder.goid]; exists {
				waitInfo.monitor.diagnose(waitInfo.mode, append(waitLink, waitInfo))
			}
		}
	}
}

func (m *monitor) using(waitInfo *lockUsage) {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	delete(waitingList, waitInfo.goid)
	m.holders.PushBack(waitInfo)
}

func (m *monitor) release(mode byte) {
	id := goid.Get()
	for i := m.holders.Back(); i != nil; i = i.Prev() {
		if info := i.Value.(*lockUsage); info.goid == id && info.mode == mode {
			m.holders.Remove(i)
			break
		}
	}
}

func deadlockPanic(waitLink []*lockUsage) {
	buf := new(bytes.Buffer)
	buf.Write(titleStr)
	for i := 0; i < len(waitLink); i++ {
		buf.Write(goStr)
		buf.WriteString(strconv.Itoa(int(waitLink[i].goid)))
		buf.Write(waitStr)
		if waitLink[i].mode == 'w' {
			buf.Write(writeStr)
		} else {
			buf.Write(readStr)
		}
		buf.Write(lineStr)
		buf.Write(waitLink[i].stack.Bytes("  "))

		// lookup waiting for who
		n := i + 1
		if n == len(waitLink) {
			n = 0
		}
		waitWho := waitLink[n]

		for j := waitLink[i].monitor.holders.Front(); j != nil; j = j.Next() {
			waitHolder := j.Value.(*lockUsage)
			if waitHolder.goid == waitWho.goid {
				buf.Write(goStr)
				buf.WriteString(strconv.Itoa(int(waitHolder.goid)))
				buf.Write(holdStr)
				if waitHolder.mode == 'w' {
					buf.Write(writeStr)
				} else {
					buf.Write(readStr)
				}
				buf.Write(lineStr)
				buf.Write(waitHolder.stack.Bytes("  "))
				break
			}
		}
	}
	panic(DeadlockError(buf.String()))
}
