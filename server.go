package main

import (
	"fmt"
	"log"
	"sync"
	"github.com/kushagra-gupta01/Content_Addressable_Storage/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes		[]string
}

type FileServer struct {
	FileServerOpts
	store 		*Store
	quitCh 		chan struct{}
	peers			map[string]p2p.Peer
	peerLock 	sync.Mutex
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitCh: make(chan struct{}),
		peers: make(map[string]p2p.Peer),
	}
}

func (s *FileServer) Stop(){
	close(s.quitCh)
}

func (s *FileServer) OnPeer(p p2p.Peer)error{
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with remote %s",p.RemoteAddr())
	return nil
}

func (s *FileServer) loop(){
	defer func(){
		log.Println("file server stopped due to quit action")
		s.Transport.Close()
	}()
	for{
		select{
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitCh: 
			return
		}
	}
}

func (s *FileServer) bootstrapNetwork() error{
	for _,addr := range s.BootstrapNodes{
		if len(addr)==0{continue}
		go func(addr string){
			fmt.Println("attempting to connect with remote: ",addr)
			if err:= s.Transport.Dial(addr);err!=nil{
				log.Println("dial error:",err)
			}
		}(addr)
	}
	return nil
}

func (s *FileServer) Start() error{
	if err:= s.Transport.ListenAndAccept();err!=nil{
		return err
	}

	s.bootstrapNetwork()
	s.loop()
	return  nil
}
