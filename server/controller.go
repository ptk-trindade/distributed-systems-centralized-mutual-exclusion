package main

import "context"

func Controller(monitor *Monitor) {

	req, err := monitor.receiveRequest()
	for err == nil {
		monitor.addAttendance(req.processId) // add to attendance count

		req.handlerGrant.Release(1) // signal the process that it can continue

		monitor.releaseSem.Acquire(context.Background(), 1) // wait for the process to finish

		req, err = monitor.receiveRequest()
	}

}
