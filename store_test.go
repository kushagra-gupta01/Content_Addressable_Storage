package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T){
	
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefautPathTransformFunc,
	}

	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("myspecialpicture",data);err!=nil{
		t.Error(err)
	}
}