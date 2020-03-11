package schedule

import "evtcp/message"

type ReadyNotifier interface {
	WorkerReady(chan *message.MessagePack)
	PusherReady(chan *message.MessagePack)
}
