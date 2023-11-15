package main

import (
	"log"
	"github.com/kushagra-gupta01/Content_Addressable_Storage/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:			":3000",
		HandshakeFunc: 	p2p.NOPHandshakeFunc,
		Decoder: 				p2p.Defaultdecoder{},
	}
	
	tr := p2p.NewTCPTransport(tcpOpts)	

	if err:= tr.ListenAndAccept();err!=nil{
		log.Fatal(err)
	}

	select{}
}