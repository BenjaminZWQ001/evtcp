package main

import (
	"evtcp/engine"
	"evtcp/incoming"
	"evtcp/schedule"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var signChannel = make(chan os.Signal, 1)

func main() {
	installSign()
	accEngin := new(incoming.AcceptEngine)
	accEngin.Host = "127.0.0.1"
	accEngin.Port = 9529
	conEngin := new(engine.ConcurrentEngine)
	conEngin.WorkerNumber = 8
	conEngin.PusherNumber = 8
	conEngin.AcceptEngine = *accEngin
	conEngin.Scheduler = new(schedule.QueueScheduler)

	go conEngin.Run()

	for {
		select {
		case <-signChannel:
			fmt.Println("Get shutdown sign")
			goto EXIT
		}
	}

EXIT:
	conEngin.Stop()
}

func installSign() {
	signal.Notify(signChannel, os.Interrupt, os.Kill, syscall.SIGTSTP)
}
