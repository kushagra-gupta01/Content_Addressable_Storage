package main

import (
	"log"
	"time"

	"github.com/kushagra-gupta01/Content_Addressable_Storage/p2p"
)

func main() {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.Defaultdecoder{},
		//todo: onpeer func
	}
	tcpTransport:=p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot: "3000_network",
		PathTransformFunc: CASpathTransformFunc,
		Transport: tcpTransport,
	}

	s:= NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(time.Second*5)
		s.Stop()
	}()

	if err:= s.Start();err!=nil{
		log.Fatal(err)
	}
}