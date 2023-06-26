package main

func Controller(monitor *Monitor) {

	go Server(monitor)

	req, err := monitor.receiveRequest()
	for err == nil {
		monitor.addAttendance(req.processId) // add to attendance count

		req.handlerGrant <- true // signal the process that it can continue

		monitor.receiveRelease() // wait for the process to finish

		req, err = monitor.receiveRequest()
	}

}
