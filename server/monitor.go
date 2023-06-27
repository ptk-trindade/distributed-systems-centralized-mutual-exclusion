package main

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type QueueElement struct {
	processId    string
	handlerGrant *semaphore.Weighted
}

type Monitor struct {
	mtxQueue        sync.Mutex
	mtxAttCnt       sync.Mutex
	requestQueue    []QueueElement
	attendanceCount map[string]int
	queueSem        *semaphore.Weighted
	releaseSem      *semaphore.Weighted
	wg              sync.WaitGroup
	ctx             context.Context
	cancelCtx       context.CancelFunc
}

func createMonitor() *Monitor {
	requestQueue := make([]QueueElement, 0)
	attendanceCount := make(map[string]int)
	queueSem := semaphore.NewWeighted(100) // 100 is the number of processes that can be in the queue at the same time
	queueSem.Acquire(context.Background(), 100)

	releaseSem := semaphore.NewWeighted(1)
	releaseSem.Acquire(context.Background(), 1)

	ctx, cancelCtx := context.WithCancel(context.Background())

	monitor := &Monitor{
		requestQueue:    requestQueue,
		attendanceCount: attendanceCount,
		queueSem:        queueSem,
		releaseSem:      releaseSem,
		ctx:             ctx,
		cancelCtx:       cancelCtx,
	}

	return monitor
}

// --- QUEUE ---
// called by handler when a request is received from a client
func (monitor *Monitor) sendRequest(processId string, handlerGrant *semaphore.Weighted) {
	monitor.mtxQueue.Lock()
	monitor.requestQueue = append(monitor.requestQueue, QueueElement{processId, handlerGrant})
	monitor.mtxQueue.Unlock()
	monitor.queueSem.Release(1)
}

// called by controller when "critical section" is free
func (monitor *Monitor) receiveRequest() (QueueElement, error) {
	err := monitor.queueSem.Acquire(context.Background(), 1)
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
