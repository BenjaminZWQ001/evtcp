package pusher

import (
	"evtcp/message"
	"evtcp/schedule"
	"fmt"
	"log"
)

type Pusher struct {
	requestChan chan *message.MessagePack
}

func NewPusher() Pusher {
	requestChan := make(chan *message.MessagePack)
	return Pusher{
		requestChan: requestChan,
	}
}

func (p *Pusher) StandBy(exitChanel chan int, notifier schedule.ReadyNotifier) {
	go func() {
		for {
			notifier.PusherReady(p.requestChan)
			select {
			case activeRequest := <-p.requestChan:
				err := p.processRequest(activeRequest)
				if err != nil {
					log.Printf("%v\n", err)
					continue
				}
			case <-exitChanel:
				fmt.Println("pusher get exit sign....")
				goto EXIT
			}
		}
	EXIT:
		fmt.Println("pusher gorouting exit ....")
	}()
}

func (p *Pusher) processRequest(request *message.MessagePack) error {
	connInfo := request.ConnPoint
	hostRevertStr := string(request.Content)
	hostRevertPackerByte := message.Packet(hostRevertStr)
	(*connInfo).Write(hostRevertPackerByte)
	remoteAddr := (*connInfo).RemoteAddr()
	log.Printf("revert to Client ID %s with message: %s\n", remoteAddr, hostRevertStr)
	return nil
}
