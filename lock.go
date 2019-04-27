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
	"strconv"
	"container/list"
	_"context"
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
	kill KillSwitch
}

type KillSwitch struct{
	kill map[string]bool
	mux sync.Mutex
}

func (ks * KillSwitch) On(id string){
	ks.mux.Lock()
	ks.kill[id] = true
	ks.mux.Unlock()
}

func (ks * KillSwitch) Off(id string){
	ks.mux.Lock()
	ks.kill[id] = false
	ks.mux.Unlock()
}

func (ks * KillSwitch) isOn(id string)bool{
	ks.mux.Lock()
	defer ks.mux.Unlock()
	if ks.kill[id]{
		ks.kill[id] = false
		return true
	}
	return false
}

func (l *Lock)init()int {
	l.readLock = semaphore.NewWeighted(int64(10))
	l.writeLock = semaphore.NewWeighted(int64(1))
	l.readHolders = make(map[string]int)
	l.writeHolder = ""
	l.kill.kill = make(map[string]bool)
	for i:=0;i<10;i++{
		l.kill.Off(strconv.Itoa(i))
	}
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

func (l *Lock)lockReader(cliNum string) bool{
	// spin lock while writer exists
	for l.writerExists() {
		if l.kill.isOn(cliNum){
			return false
		}
		time.Sleep(1)
	}
	//ctx := context.Background()
	for !l.readLock.TryAcquire(1){
		if l.kill.isOn(cliNum){
			return false
		}
	}
	l.readHolders[cliNum] = 1
	return true
}

func (l *Lock)lockWriter(cliNum string)bool{

	// spin lock while reader exists
	for l.readerExists() {
		if (len(l.readHolders) == 1 && l.isReader(cliNum)) {
			break
		}
		if l.kill.isOn(cliNum){
			return false
		}
		time.Sleep(1)
		fmt.Println("pooling")
	}
	// ctx := context.Background()
	for !l.writeLock.TryAcquire(1){
		if l.kill.isOn(cliNum){
			return false
		}
	}
	l.writeHolder = cliNum
	return true
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
