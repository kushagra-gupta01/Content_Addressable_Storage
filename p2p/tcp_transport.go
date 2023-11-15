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

type TCPTransportOpts struct{
	ListenAddr		string
	HandshakeFunc	HandshakeFunc
	Decoder				Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener

	mu 						sync.RWMutex
	peers 				map[string]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport{
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}


func(t *TCPTransport)ListenAndAccept() error{
	var err error
	t.listener,err = net.Listen("tcp",t.ListenAddr)
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
		fmt.Printf("new incomming connection: %+v\n",conn)
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport)handleConn(conn net.Conn){

	peer:= NewTCPpeer(conn,true)

	if err := t.HandshakeFunc(peer);err!=nil{
		conn.Close()
		fmt.Printf("TCP handshake error:%s\n",err)
		return	
	}

	//Read Loop
	msg :=&Message{}
	for{
		if err := t.Decoder.Decode(conn,msg);err!=nil{
			fmt.Printf("TCP error: %s\n",err)
			continue
		}

		msg.From = conn.RemoteAddr()
		fmt.Printf("message: %+v\n",msg)
	}
}