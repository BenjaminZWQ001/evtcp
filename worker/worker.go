package worker

import (
	"encoding/json"
	"evtcp/message"
	"evtcp/schedule"
	"fmt"
	"log"
	"strings"
)

type Worker struct {
	requestChan chan *message.MessagePack
}

func NewWorker() Worker {
	requestChan := make(chan *message.MessagePack)
	return Worker{
		requestChan: requestChan,
	}
}

func (w *Worker) StandBy(exitChanel chan int, resultChan chan *message.MessagePack, notifier schedule.ReadyNotifier) {
	go func() {
		for {
			notifier.WorkerReady(w.requestChan)
			select {
			case activeRequest := <-w.requestChan:
				result, err := w.processRequest(activeRequest)
				if err != nil {
					log.Printf("%v\n", err)
					continue
				}
				resultChan <- result
			case <-exitChanel:
				fmt.Println("worker get exit sign....")
				goto EXIT
			}
		}
	EXIT:
		fmt.Println("worker gorouting exit ....")
	}()
}

func (w *Worker) processRequest(request *message.MessagePack) (*message.MessagePack, error) {
	var contentInfo message.Content
	var messagePackBack message.MessagePack
	contentBytes := request.Content
	unmarshal := json.Unmarshal(contentBytes, &contentInfo)
	if unmarshal != nil {
		return &messagePackBack, unmarshal
	}
	connInfo := request.ConnPoint
	remoteAddr := (*connInfo).RemoteAddr()
	log.Printf("message received from Client ID %s serviceId: %s, content: %s\n", remoteAddr, contentInfo.ServiceId, contentInfo.Data)

	messagePackBack.ConnPoint = connInfo
	messagePackBack.Length = 0
	hostRevertJsonByte, _ := processContent(contentInfo)
	messagePackBack.Content = hostRevertJsonByte
	return &messagePackBack, nil
}

func processContent(contentInfo message.Content) ([]byte, error) {
	contentRevert := message.Content{}
	if contentInfo.ServiceId == "Heart.beat" {
		contentRevert.ServiceId = "Heart.beat"
		contentRevert.Data = "Pong"
	} else {
		contentRevert.ServiceId = "Host.revert"
		contentRevert.Data = strings.ToUpper(contentInfo.Data)
	}
	return json.Marshal(contentRevert)
}
