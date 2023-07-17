package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"golang.org/x/sync/semaphore"
)

type QueueElement struct {
	identifier any
	turnSem    *semaphore.Weighted
}

type MutexLog struct {
	avaialble    bool
	requestQueue []QueueElement
	history      []any

	mtx sync.Mutex

	ctx       context.Context
	cancelCtx context.CancelFunc
}

func CreateMutexLog() *MutexLog {
	requestQueue := make([]QueueElement, 0)

	ctx, cancelCtx := context.WithCancel(context.Background())

	mtxLog := &MutexLog{
		avaialble:    true,
		requestQueue: requestQueue,
		history:      make([]any, 0),
		ctx:          ctx,
		cancelCtx:    cancelCtx,
	}

	return mtxLog
}

func (mutex *MutexLog) Lock(identifier any) error {
	mutex.mtx.Lock()
	var err error
	writeOnLog("1", identifier)

	err = mutex.ctx.Err()
	if err != nil {
		writeOnLog("4", identifier)
		return mutex.ctx.Err()
	}

	if mutex.avaialble {
		mutex.avaialble = false
		mutex.history = append(mutex.history, identifier)
		mutex.mtx.Unlock()
		writeOnLog("2", identifier)
		return nil
	}

	sem := semaphore.NewWeighted(1)
	sem.Acquire(mutex.ctx, 1) // make semafore unavailable

	mutex.requestQueue = append(mutex.requestQueue, QueueElement{identifier, sem})
	mutex.mtx.Unlock()

	err = sem.Acquire(mutex.ctx, 1) // wait for turn
	if err != nil {
		writeOnLog("4", identifier)
	} else {
		writeOnLog("2", identifier)
	}

	return err
}

func (mutex *MutexLog) Unlock() {
	mutex.mtx.Lock()

	lastIndex := len(mutex.history) - 1
	identifier := mutex.history[lastIndex]
	writeOnLog("3", identifier)

	if len(mutex.requestQueue) > 0 { // someone is waiting
		element := mutex.requestQueue[0]
		fmt.Println("release:", element.identifier)
		element.turnSem.Release(1)
		mutex.history = append(mutex.history, element.identifier)

		mutex.requestQueue = mutex.requestQueue[1:]
	} else {
		mutex.avaialble = true
	}

	mutex.mtx.Unlock()
}

func (mutex *MutexLog) Cancel() {
	mutex.cancelCtx()

}

func (mutex *MutexLog) GetQueue() []any {
	mutex.mtx.Lock()

	copiedQueue := make([]any, len(mutex.requestQueue))
	for i, v := range mutex.requestQueue {
		copiedQueue[i] = v.identifier
	}

	mutex.mtx.Unlock()

	return copiedQueue
}

func (mutex *MutexLog) GetHistory() []any {
	mutex.mtx.Lock()

	copiedHistory := make([]any, len(mutex.history))
	copy(copiedHistory, mutex.history)

	mutex.mtx.Unlock()

	return copiedHistory
}

// WRITE ON log.txt
var FILEPATH string = "log.txt"
var fileMtx sync.Mutex

func writeOnLog(code string, identifier any) {
	var text string

	switch code {
	case "1":
		text = "[R] Request-" + identifier.(string) + "\n"
	case "2":
		text = "[S] Grant-" + identifier.(string) + "\n"
	case "3":
		text = "[R] Release-" + identifier.(string) + "\n"
	case "4":
		text = "[S] Close-" + identifier.(string) + "\n"
	case "9":
		text = "[R] Exit-" + identifier.(string) + "\n"
	}

	fileMtx.Lock()

	_, err := os.Stat(FILEPATH)
	if os.IsNotExist(err) {
		createFile(FILEPATH)
	}

	appendToFile(FILEPATH, text)

	fileMtx.Unlock()
}
