package main

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type Monitor struct {
	mtxQueue        sync.Mutex
	mtxAttCnt       sync.Mutex
	requestQueue    []QueueElement
	attendanceCount map[string]int
	semaphore       *semaphore.Weighted
	releaseChan     chan bool

	ctx       context.Context
	cancelCtx context.CancelFunc
}

func createMonitor() *Monitor {
	requestQueue := make([]QueueElement, 0)
	attendanceCount := make(map[string]int)
	queueSem := semaphore.NewWeighted(100) // 100 is the number of processes that can be in the queue at the same time
	queueSem.Acquire(context.Background(), 100)
	release := make(chan bool)
	ctx, cancelCtx := context.WithCancel(context.Background())

	monitor := &Monitor{
		requestQueue:    requestQueue,
		attendanceCount: attendanceCount,
		semaphore:       queueSem,
		releaseChan:     release,
		ctx:             ctx,
		cancelCtx:       cancelCtx,
	}

	return monitor
}

// --- QUEUE ---
// called by handler when a request is received from a client
func (monitor *Monitor) sendRequest(processId string, handlerGrant chan bool) {
	monitor.mtxQueue.Lock()
	monitor.requestQueue = append(monitor.requestQueue, QueueElement{processId, handlerGrant})
	monitor.mtxQueue.Unlock()
	monitor.semaphore.Release(1)
}

// called by controller when "critical section" is free
func (monitor *Monitor) receiveRequest() (QueueElement, error) {
	err := monitor.semaphore.Acquire(context.Background(), 1)
	if err != nil {
		return QueueElement{}, err
	}

	monitor.mtxQueue.Lock()

	req := monitor.requestQueue[0]
	monitor.requestQueue = monitor.requestQueue[1:]

	monitor.mtxQueue.Unlock()

	return req, nil
}

// called by main when the user wants to see the queue
func (monitor *Monitor) getRequestQueue() []QueueElement {
	monitor.mtxQueue.Lock()

	copiedQueue := make([]QueueElement, len(monitor.requestQueue))
	copy(copiedQueue, monitor.requestQueue)

	monitor.mtxQueue.Unlock()

	return copiedQueue
}

// --- ATTENDANCE COUNT ---
// called by controller when a request is attended
func (monitor *Monitor) addAttendance(processId string) {
	monitor.mtxAttCnt.Lock()
	monitor.attendanceCount[processId]++
	monitor.mtxAttCnt.Unlock()
}

// called by main when the user wants to see the attendance count
func (monitor *Monitor) getAttendanceCount() map[string]int {
	monitor.mtxAttCnt.Lock()

	copiedMap := make(map[string]int)
	for k, v := range monitor.attendanceCount {
		copiedMap[k] = v
	}

	monitor.mtxAttCnt.Unlock()

	return copiedMap
}

// --- RELEASE ---
// called by handler when client finishes using critical section
func (monitor *Monitor) sendRelease() {
	monitor.releaseChan <- true
}

// called by controller when waiting for a handler to finish
func (monitor *Monitor) receiveRelease() {
	<-monitor.releaseChan
}

// --- STOP MONITOR ---
// called by main when the program is closed
func (monitor *Monitor) exit() {
	monitor.cancelCtx()
}
