package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"
)

func TestPathTransformFunc(t *testing.T){
	key := "heyKushagrathisSide"
	pathKey:= CASpathTransformFunc(key)
	expectedOriginalKey := "d9e06924cbe4f7c5f59269e6267f971d02774564"
	expectedPathName := "d9e06/924cb/e4f7c/5f592/69e62/67f97/1d027/74564"

	if pathKey.PathName != expectedPathName{
		t.Errorf("Have %s , want %s",pathKey.PathName,expectedPathName)
	}

	if pathKey.FileName != expectedOriginalKey{
		t.Errorf("Have %s , want %s",pathKey.FileName,expectedOriginalKey)
	}
}

func TestStoreDeleteKey(t *testing.T){
	opts:= StoreOpts{
		PathTransformFunc: CASpathTransformFunc,
	}	

	s := NewStore(opts)
	key := "43r43frerf"
	data := []byte("hola hola")

	if err:= s.writeStream(key,bytes.NewReader(data));err!=nil{
		t.Error(err)
	}

	time.Sleep(100 * time.Millisecond)
	if err := s.Delete(key);err!=nil{
		t.Error(err)	
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASpathTransformFunc,
	}

	s := NewStore(opts)
	key := "omm>>>>>>"
	data := []byte("some jpg bytes")
	
	if err := s.writeStream(key,bytes.NewReader(data));err!=nil{
		t.Error(err)
	}

	if ok := s.Has(key); !ok{
		t.Errorf("expected to have key %s",key)
	}

	r,err := s.Read(key)
	if err!=nil{
		t.Error(err)
	}

	b,err := io.ReadAll(r)
	if err!=nil{
		t.Error(err)
	}

	fmt.Println(string(b))
	if string(b) != string(data){
		t.Errorf("want %s, have %s",data,b)
	}
	s.Delete(key)
}