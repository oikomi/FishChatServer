// +build deadlock

package syncs

import (
	"testing"
	"time"
)

func init() {
	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()
}

func deadlockTest(t *testing.T, callback func()) {
	testDone := make(chan interface{})

	go func() {
		defer func() {
			testDone <- recover()
		}()
		callback()
	}()

	select {
	case err := <-testDone:
		if err != nil {
			switch err.(type) {
			case DeadlockError:
				if testing.Verbose() {
					println(err.(DeadlockError).Error())
				}
			default:
				panic(err)
			}
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}

}

func Test_DeadLock1(t *testing.T) {
	deadlockTest(t, func() {
		var mutex1 Mutex

		mutex1.Lock()
		mutex1.Lock()
	})
}

func Test_DeadLock2(t *testing.T) {
	deadlockTest(t, func() {
		var (
			mutex1 Mutex
			mutex2 Mutex
		)

		mutex1.Lock()

		var wait1 WaitGroup
		wait1.Add(1)
		go func() {
			mutex2.Lock()
			wait1.Done()
			mutex1.Lock()
		}()
		wait1.Wait()

		mutex2.Lock()
	})
}

func Test_DeadLock3(t *testing.T) {
	deadlockTest(t, func() {
		var (
			mutex1 Mutex
			mutex2 Mutex
			mutex3 Mutex
		)

		mutex1.Lock()

		var wait1 WaitGroup
		wait1.Add(1)
		go func() {
			mutex2.Lock()

			var wait2 WaitGroup
			wait2.Add(1)
			go func() {
				mutex3.Lock()
				wait2.Done()
				mutex2.Lock()
			}()
			wait2.Wait()

			wait1.Done()
			mutex1.Lock()
		}()
		wait1.Wait()

		mutex3.Lock()
	})
}

func Test_DeadLock4(t *testing.T) {
	deadlockTest(t, func() {
		var (
			mutex1 RWMutex
			mutex2 RWMutex
			mutex3 RWMutex
		)

		mutex1.Lock()

		var wait1 WaitGroup
		wait1.Add(1)
		go func() {
			mutex2.RLock()

			var wait2 WaitGroup
			wait2.Add(1)
			go func() {
				mutex3.Lock()
				wait2.Done()
				mutex2.Lock()
			}()
			wait2.Wait()

			wait1.Done()
			mutex1.RLock()
		}()
		wait1.Wait()

		mutex3.Lock()
	})
}
