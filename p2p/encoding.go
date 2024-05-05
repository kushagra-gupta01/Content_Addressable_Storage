package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader,*RPC) error
}

type GOBdecoder struct{}

func (dec GOBdecoder) Decode(r io.Reader, msg *RPC) error{
	return gob.NewDecoder(r).Decode(msg)
}

type Defaultdecoder struct{}

func (dec Defaultdecoder) Decode(r io.Reader,msg *RPC) error{
	peekBuf:= make([]byte,1)
	if _,err:= r.Read(peekBuf);err!=nil{
		return err
	}
	//Incase of a stream, we are not decoding what is being sent over the network.
	//We are just setting Stream true, so we can handle that in our logic.
	stream:= peekBuf[0] ==IncomingStream
	if stream{
		msg.Stream = true
		return nil
	}
	
	buf := make([]byte, 1024)
	n,err:= r.Read(buf)
	if err!=nil{
		return err
	}

	msg.Payload = buf[:n]

	return nil
}