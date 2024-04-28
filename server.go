package main

import (
	"fmt"
	"log"

	"github.com/kushagra-gupta01/Content_Addressable_Storage/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
}

type FileServer struct {
	FileServerOpts
	store *Store
	quitCh chan struct{}
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
	}
}

func (s *FileServer) Stop(){
	close(s.quitCh)
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

func (s *FileServer) Start() error{
	if err:= s.Transport.ListenAndAccept();err!=nil{
		return err
	}
	s.loop()
	return  nil
}
