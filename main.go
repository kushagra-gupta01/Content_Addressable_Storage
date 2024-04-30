package main

import (
	"bytes"
	"log"
	"time"
	"github.com/kushagra-gupta01/Content_Addressable_Storage/p2p"
)

func makeServer(listenAddr string,nodes ...string) *FileServer{
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr: listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.Defaultdecoder{},
	}
	tcpTransport:=p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot: 				listenAddr+"_network",
		PathTransformFunc: 	CASpathTransformFunc,
		Transport: 					tcpTransport,
		BootstrapNodes: 		nodes,
	}

	s:=NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":3000","")
	s2 := makeServer(":4000",":3000")
	go func ()  {
		log.Fatal(s1.Start())
	}()
	time.Sleep(4*time.Second)
	go s2.Start()
	time.Sleep(4*time.Second)
	data:= bytes.NewReader([]byte("my data file is heree"))
	s2.StoreData("myprivdata",data)
	select{}
}