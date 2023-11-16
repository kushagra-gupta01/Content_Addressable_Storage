package main

import (
	"fmt"
	"log"
	"github.com/kushagra-gupta01/Content_Addressable_Storage/p2p"
)

func OnPeer(peer p2p.Peer) error{
	fmt.Printf("doing logic outside TCPTransport\n\n")
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:			":3000",
		HandshakeFunc: 	p2p.NOPHandshakeFunc,
		Decoder: 				p2p.Defaultdecoder{},
		OnPeer: 				OnPeer,		
	}
	
	tr := p2p.NewTCPTransport(tcpOpts)	

	go func ()  {
		for{
			msg := <-tr.Consume()
			fmt.Printf("%+v\n",msg)
		}
	}()

	if err:= tr.ListenAndAccept();err!=nil{
		log.Fatal(err)
	}

	select{}
}