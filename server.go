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
	Size int64
}
func (s *FileServer) StoreData(key string,r io.Reader) error{
	//1. Store this file to disk
	//2. broadcast this file to all known peers in the network
	fileBuffer := new(bytes.Buffer)
	tee:= io.TeeReader(r,fileBuffer)
	size,err:= s.store.Write(key,tee)
	if err!=nil{
		return err
	}
	msg:= Message{
		Payload: MessageStoreFile{
			Key: key,
			Size: size,
		},
	}
	
	msgBuf:= new(bytes.Buffer)
	if err:= gob.NewEncoder(msgBuf).Encode(msg);err!=nil{
		return err
	}

	for _,peer :=range s.peers{
		if err:= peer.Send(msgBuf.Bytes());err!=nil{
			return err
		}
	}
	time.Sleep(3*time.Second)
	for _,peer :=range s.peers{
		n,err:= io.Copy(peer,fileBuffer)
		if err!=nil{
			return err
		}
		fmt.Println("received and written bytes to disk: ",n)
	}

	return nil
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
				return
			}

			if err:= s.handleMessage(rpc.From,&msg);err!=nil{
				log.Println(err)
				return
			}
		case <-s.quitCh: 
			return
		}
	}
}

func(s *FileServer) handleMessage(from string,msg *Message)error{
	switch v := msg.Payload.(type){
	case MessageStoreFile:
		return s.handleMessageStoreFile(from,v)
	}
	return nil
}

func (s *FileServer) handleMessageStoreFile(from string,msg MessageStoreFile) error{
	peer,ok:= s.peers[from]
	if !ok{
		return fmt.Errorf("peer (%s) could not be found in peerlist",from)
	}
	if _,err:= s.store.Write(msg.Key,io.LimitReader(peer,msg.Size));err!=nil{
		return err
	}
	peer.(*p2p.TCPpeer).Wg.Done()
	return nil
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

func init(){
	gob.Register(MessageStoreFile{})
}