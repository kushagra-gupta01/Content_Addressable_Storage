package main

import (
	// "bytes"
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
	r,err:= s2.Get("coolPicture.jpg")
	if err!=nil{
		log.Fatal(err)
	}
	b,err:= io.ReadAll(r)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(string(b))
	
	// data:= bytes.NewReader([]byte("my big data file here!"))
	// s2.Store("coolPicture.jpg",data)
	// time.Sleep(time.Millisecond*5)	

}