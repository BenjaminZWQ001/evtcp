package schedule

import "evtcp/message"

type Scheduler interface {
	ReadyNotifier
	Submit(*message.MessagePack)
	Dispatch(*message.MessagePack)
	Run()
}
