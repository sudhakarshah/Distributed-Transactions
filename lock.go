package main
// package semaphore
import(
	"sync"
	_"errors"
	_"net"
	_"encoding/json"
	"fmt"
	"time"
	_"bufio"
	_"strings"
	"container/list"
	"context"
	"golang.org/x/sync/semaphore"
)

type Weighted struct {
	size    int64
	cur     int64
	mu      sync.Mutex
	waiters list.List
}

type Lock struct{
	readLock *semaphore.Weighted
	writeLock *semaphore.Weighted
	readHolders map[string]int
	writeHolder string
}


func (l *Lock)init()int {
	l.readLock = semaphore.NewWeighted(int64(10))
	l.writeLock = semaphore.NewWeighted(int64(1))
	l.readHolders = make(map[string]int)
	l.writeHolder = ""
	return 1;
}

func (l *Lock)isReader(cliNum string) bool {
		for val, _ := range l.readHolders {
				if val == cliNum {
					return true
				}
		}
    return false
}

func (l *Lock)isWriter(cliNum string) bool {
	return l.writeHolder == cliNum
}

func (l *Lock)readerExists() bool {
	return len(l.readHolders) != 0
}

func (l *Lock)writerExists() bool {
	return l.writeHolder != ""
}

func (l *Lock)lockReader(cliNum string) {
	// spin lock while writer exists
	for l.writerExists() {
		time.Sleep(1)
	}
	ctx := context.Background()
	l.readLock.Acquire(ctx, 1)
	l.readHolders[cliNum] = 1
}

func (l *Lock)lockWriter(cliNum string) {

	// spin lock while reader exists
	for l.readerExists() {
		if (len(l.readHolders) == 1 && l.isReader(cliNum)) {
			break
		}
		time.Sleep(1)
		fmt.Println("pooling")
	}
	ctx := context.Background()
	l.writeLock.Acquire(ctx, 1)
	l.writeHolder = cliNum
}

func (l *Lock)UnlockReader(cliNum string) {
	l.readLock.Release(1)
	delete(l.readHolders, cliNum)
}

func (l *Lock)UnlockWriter(cliNum string) {
	l.writeLock.Release(1)
	l.writeHolder = ""
}

func (l *Lock)removeReader(cliNum string) {
	if _, ok := l.readHolders[cliNum]; ok {
		l.UnlockReader(cliNum)
	}
}

func (l *Lock)removeWriter(cliNum string) {
	if l.writeHolder == cliNum {
		l.UnlockWriter(cliNum)
	}
}
