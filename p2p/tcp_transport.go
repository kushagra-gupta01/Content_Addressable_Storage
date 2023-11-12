package p2p

import (
	"fmt"
	"net"
	"sync"
)

//TCPpeer represents the remote node over a TCP established connection. 
type TCPpeer struct{
	//conn is the underlying connection of the peer
	conn net.Conn

	//if we dial and retreive a connection => outbound == true
	//if we accept and retreive a connection => outbound == false
	outbound bool
}

func NewTCPpeer(conn net.Conn, outbound bool) *TCPpeer{
	return &TCPpeer{
		conn: conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	handshakeFunc	HandshakeFunc
	decoder				Decoder

	mu 						sync.RWMutex
	peers 				map[string]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport{
	return &TCPTransport{
		handshakeFunc: NOPHandshakeFunc,
		listenAddress: listenAddr,
	}
}


func(t *TCPTransport)ListenAndAccept() error{
	var err error
	t.listener,err = net.Listen("tcp",t.listenAddress)
	if err !=nil{
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport)startAcceptLoop(){
	for{
		conn,err := t.listener.Accept()
		if err!=nil{
			fmt.Printf("TCP accept error: %s\n", err)
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport)handleConn(conn net.Conn){

	if err := t.handshakeFunc(conn);err!=nil{

	}
	fmt.Printf("New incomming connection: %+v\n",conn)
}