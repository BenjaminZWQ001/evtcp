package engine

import (
	"evtcp/incoming"
	"evtcp/message"
	"evtcp/pusher"
	"evtcp/schedule"
	"evtcp/worker"
	"fmt"
	"time"
)

type ConcurrentEngine struct {
	WorkerNumber int
	PusherNumber int
	AcceptEngine incoming.AcceptEngine
	Scheduler    schedule.Scheduler
}

var incomingChannel = make(chan *message.MessagePack)
var outcomingChannel = make(chan *message.MessagePack)

var workerExitChanel = make(chan int)
var pusherExitChanel = make(chan int)

func (e *ConcurrentEngine) Run() {
	e.AcceptEngine.Run(incomingChannel)
	e.Scheduler.Run()

	for i := 0; i < e.WorkerNumber; i++ {
		newWorker := worker.NewWorker()
		newWorker.StandBy(workerExitChanel, outcomingChannel, e.Scheduler)
	}

	for i := 0; i < e.PusherNumber; i++ {
		newPusher := pusher.NewPusher()
		newPusher.StandBy(pusherExitChanel, e.Scheduler)
	}

	for {
		select {
		case incomingMessage := <-incomingChannel:
			e.Scheduler.Submit(incomingMessage)
		case outcomingMessage := <-outcomingChannel:
			e.Scheduler.Dispatch(outcomingMessage)
		}
	}
}

func (e *ConcurrentEngine) Stop() {
	e.AcceptEngine.Stop()
	fmt.Println("notice worker gorouting exit ....")
	for i := 0; i < e.WorkerNumber; i++ {
		workerExitChanel <- 1
	}
	close(incomingChannel)
	time.Sleep(time.Millisecond * 20)
	fmt.Println("notice pusher gorouting exit ....")
	for i := 0; i < e.PusherNumber; i++ {
		pusherExitChanel <- 1
	}
	close(outcomingChannel)
	time.Sleep(time.Millisecond * 20)
	for remoteAddr, conn := range incoming.ConnList {
		fmt.Println("ConcurrentEngine closed Client " + remoteAddr)
		(*conn).Close()
	}
}
