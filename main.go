package main

import (
	"bytes"
	"fmt"
	"io"
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
		EncKey:							newEncryptionKey(),
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
	s2 := makeServer(":4000","")
	s3 := makeServer(":5000",":3000",":4000")

	go func (){ log.Fatal(s1.Start())}()
	time.Sleep(500*time.Millisecond)
	
	go func (){ log.Fatal(s2.Start())}()
	time.Sleep(500*time.Millisecond)

	go s3.Start()
	time.Sleep(2*time.Second)
	
	for i := 0; i < 20; i++ {
		key:="coolPicture.jpg"	
		data:= bytes.NewReader([]byte("my big data file here!"))
		s3.Store(key,data)

		if err := s3.store.Delete(s3.ID,key); err!=nil{
			log.Fatal(err)
		}

		r,err:= s3.Get(key)
		if err!=nil{
			log.Fatal(err)
		}
		b,err:= io.ReadAll(r)
		if err!=nil{
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}

	
}