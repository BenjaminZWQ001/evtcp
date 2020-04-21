package incoming

import (
	"context"
	"evtcp/message"
	"log"
	"net"
	"strconv"
)

var ConnList = make(map[string]*net.TCPConn)

//var connChannel chan *net.TCPConn = make(chan *net.TCPConn)
var connCount = 5000
var connPool = make(chan bool, connCount)

type AcceptEngine struct {
	Host     string
	Port     int
	listener *net.TCPListener
	cancel   context.CancelFunc
}

func (aEngine *AcceptEngine) Run(incomingChannel chan *message.MessagePack) {
	listener, _ := listenerInit(aEngine.Host, aEngine.Port)
	aEngine.listener = listener
	ctx, cancel := context.WithCancel(context.Background())
	aEngine.cancel = cancel
	for i := 0; i < connCount; i++ {
		connPool <- true
	}
	go ConnectionAccept(ctx, listener, incomingChannel)
	//go HandleConn(ctx, connChannel, incomingChannel)
}

func listenerInit(host string, port int) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host+":"+strconv.Itoa(port))
	if err != nil {
		panic("解析ip地址失败: " + err.Error())
	}
	log.Printf("Listening %s:%d ....\n", host, port)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic("监听TCP失败: " + err.Error())
	}
	log.Printf("Listen success on %s:%d with tcp4\n", host, port)
	return listener, nil
}

func ConnectionAccept(ctx context.Context, listener *net.TCPListener, incomingChannel chan *message.MessagePack) {
	log.Println("Wating connection ....")
	for {
		select {
		case <-ctx.Done():
			log.Println("Engine stop to accept new connection")
			return
		case <-connPool:
			connection, err := listener.AcceptTCP()
			if err != nil {
				log.Println("Accept 失败: " + err.Error())
			} else {
				remoteAddr := connection.RemoteAddr()
				remoteAddrStr := remoteAddr.String()
				ConnList[remoteAddrStr] = connection
				log.Println("Client " + remoteAddrStr + " connected")
				go readConn(ctx, connection, incomingChannel, connPool)
			}
		}
	}
}

//func HandleConn(ctx context.Context, connChannel chan *net.TCPConn, incomingChannel chan *message.MessagePack) {
//	fmt.Println("Wating connection ....")
//	for {
//		select {
//		case <-ctx.Done():
//			fmt.Println("Handle Conn stop process")
//			return
//		case conn := <-connChannel:
//			remoteAddr := conn.RemoteAddr()
//			remoteAddrStr := remoteAddr.String()
//			ConnList[remoteAddrStr] = conn
//			fmt.Println("Client " + remoteAddrStr + " connected")
//			go readConn(ctx, conn, incomingChannel)
//		}
//	}
//}

func (aEngine *AcceptEngine) Stop() {
	aEngine.cancel()
	listener := aEngine.listener
	(*listener).Close()
	log.Println("Listener closed")
	for remoteAddr, conn := range ConnList {
		log.Println("accept engine close read func with Client " + remoteAddr)
		(*conn).CloseRead()
	}
}
