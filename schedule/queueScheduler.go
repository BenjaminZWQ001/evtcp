package schedule

import (
	"evtcp/message"
)

type QueueScheduler struct {
	requestChan   chan *message.MessagePack
	workerChan    chan chan *message.MessagePack
	requestQueue  []*message.MessagePack
	workerQueue   []chan *message.MessagePack
	dispatchChan  chan *message.MessagePack
	pusherChan    chan chan *message.MessagePack
	dispatchQueue []*message.MessagePack
	pusherQueue   []chan *message.MessagePack
}

func (s *QueueScheduler) Submit(request *message.MessagePack) {
	s.requestChan <- request
}

func (s *QueueScheduler) WorkerReady(worker chan *message.MessagePack) {
	s.workerChan <- worker
}

func (s *QueueScheduler) Dispatch(dispatch *message.MessagePack) {
	s.dispatchChan <- dispatch
}

func (s *QueueScheduler) PusherReady(pusher chan *message.MessagePack) {
	s.pusherChan <- pusher
}

func (s *QueueScheduler) Run() {
	s.requestChan = make(chan *message.MessagePack)
	s.workerChan = make(chan chan *message.MessagePack)

	s.dispatchChan = make(chan *message.MessagePack)
	s.pusherChan = make(chan chan *message.MessagePack)

	go func() {
		for {
			var activeRequest *message.MessagePack
			var activeWorker chan *message.MessagePack

			var activeDispatch *message.MessagePack
			var activePusher chan *message.MessagePack

			if len(s.requestQueue) > 0 && len(s.workerQueue) > 0 {
				activeRequest = s.requestQueue[0]
				activeWorker = s.workerQueue[0]
			}

			if len(s.dispatchQueue) > 0 && len(s.pusherQueue) > 0 {
				activeDispatch = s.dispatchQueue[0]
				activePusher = s.pusherQueue[0]
			}

			select {
			case newRequest := <-s.requestChan:
				s.requestQueue = append(s.requestQueue, newRequest)
			case readyWorker := <-s.workerChan:
				s.workerQueue = append(s.workerQueue, readyWorker)
			case activeWorker <- activeRequest:
				s.requestQueue = s.requestQueue[1:]
				s.workerQueue = s.workerQueue[1:]

			case newDispatch := <-s.dispatchChan:
				s.dispatchQueue = append(s.dispatchQueue, newDispatch)
			case readyPusher := <-s.pusherChan:
				s.pusherQueue = append(s.pusherQueue, readyPusher)
			case activePusher <- activeDispatch:
				s.dispatchQueue = s.dispatchQueue[1:]
				s.pusherQueue = s.pusherQueue[1:]
			}
		}
	}()
}
