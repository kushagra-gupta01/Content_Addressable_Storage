package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
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

func (s *FileServer) broadcast(msg *Message) error{
	peers:= []io.Writer{}
	for _,peer := range s.peers{
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}
type Message struct{
	Payload any
}

type MessageStoreFile struct{
	Key string
}
func (s *FileServer) StoreData(key string,r io.Reader) error{
	//1. Store this file to disk
	//2. broadcast this file to all known peers in the network

	buf := new(bytes.Buffer)
	msg:= Message{
		Payload: MessageStoreFile{
			Key: key,
		},
	}
	if err:= gob.NewEncoder(buf).Encode(msg);err!=nil{
		return err
	}

	for _,peer :=range s.peers{
		if err:= peer.Send(buf.Bytes());err!=nil{
			return err
		}
	}
	time.Sleep(3*time.Second)
	payload:= []byte("THIS IS A LARGE FILE")
	for _,peer :=range s.peers{
		if err:= peer.Send(payload);err!=nil{
			return err
		}
	}

	return nil
	// buf := new(bytes.Buffer)
	// tee := io.TeeReader(r,buf)

	// if err := s.store.Write(key,tee);err!=nil{
	// 	return err
	// }

	// p := &DataMessage{
	// 	Key: key,
	// 	Data: buf.Bytes(),
	// }

	// return s.broadcast(&Message{
	// 	From: "kush",
	// 	Payload: p,
	// })
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
		case rpc:= <-s.Transport.Consume():
			var msg Message
			if err:= gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg);err!=nil{
				log.Println(err)
			}

			fmt.Printf("payload: %+v\n",msg.Payload)
			peer,ok:= s.peers[rpc.From]
			if !ok{
				panic("peer not found")
			}
			
			b := make([]byte, 1000)
			if _,err := peer.Read(b);err!=nil{
				panic(err)
			}

			fmt.Printf("%s\n",string(b))
			peer.(*p2p.TCPpeer).Wg.Done()
			// if err:= s.handleMessage(&m);err!=nil{
			// 	log.Println(err)
			// }
		case <-s.quitCh: 
			return
		}
	}
}

// func(s *FileServer) handleMessage(msg *Message)error{
// 	switch v := msg.Payload.(type){
// 	case *DataMessage:
// 		fmt.Printf("recieved data %+v\n",v)
// 	}
// 	return nil
// }

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

func init(){
	gob.Register(MessageStoreFile{})
}