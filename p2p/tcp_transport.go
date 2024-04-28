package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
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

//Close implements the Peer interface
func (p *TCPpeer)Close() error{
	return p.conn.Close()
}

type TCPTransportOpts struct{
	ListenAddr		string
	HandshakeFunc	HandshakeFunc
	Decoder				Decoder
	OnPeer				func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener
	rpcch					chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport{
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch: 						make(chan RPC),	
	}
}

//Consume implements trasport interface, which will return read-only channel
//for reading the incomming messages recieved from another peer in the network.
func (t *TCPTransport)Consume() <- chan RPC{
	return t.rpcch
} 

//Close implements Transport interface
func (t *TCPTransport) Close() error{
	return t.listener.Close()
}

func(t *TCPTransport)ListenAndAccept() error{
	var err error
	t.listener,err = net.Listen("tcp",t.ListenAddr)
	if err !=nil{
		return err
	}

	go t.startAcceptLoop()
	log.Printf("TCP transport listening on port %s\n",t.ListenAddr)
	return nil
}

func (t *TCPTransport)startAcceptLoop(){
	for{
		conn,err := t.listener.Accept()
		if errors.Is(err,net.ErrClosed){
			return 
		}
		if err!=nil{
			fmt.Printf("TCP accept error: %s\n", err)
		}
		fmt.Printf("new incomming connection: %+v\n",conn)
		go t.handleConn(conn)
	}
}

func (t *TCPTransport)handleConn(conn net.Conn){
	var err error

	defer func ()  {
		fmt.Printf("dropping peer connection: %s\n",err)
		conn.Close()
	}()

	peer:= NewTCPpeer(conn,true)

	if err := t.HandshakeFunc(peer);err!=nil{
		return	
	}

	if t.OnPeer !=nil{
		if err = t.OnPeer(peer);err!=nil{
			return
		}
	}

	//Read Loop
	rpc :=RPC{}
	for{
		if err := t.Decoder.Decode(conn,&rpc);err!=nil{
			fmt.Printf("TCP error: %s\n",err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
		
	}
}