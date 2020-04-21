package incoming

import (
	"context"
	"evtcp/message"
	"fmt"
	"net"
)

func readConn(ctx context.Context, conn *net.TCPConn, incomingChannel chan *message.MessagePack, connPool chan bool) {
	var protocolPack *message.Protocol
	var err error
	for {
		select {
		case <-ctx.Done():
			remoteAddr := conn.RemoteAddr()
			remoteAddrStr := remoteAddr.String()
			connPool <- true
			fmt.Println("Stop reading from client " + remoteAddrStr)
			return
		default:
			protocolPack, err = message.UnPacket(conn)
			if err != nil {
				remoteAddr := (*conn).RemoteAddr()
				connPool <- true
				fmt.Println("message unpack error from connection " + remoteAddr.String() + " with message: " + err.Error())
				break
			}
			messagePack := &message.MessagePack{}
			messagePack.ConnPoint = conn
			messagePack.Length = protocolPack.Length
			messagePack.Content = protocolPack.Content
			incomingChannel <- messagePack
		}
	}
}
