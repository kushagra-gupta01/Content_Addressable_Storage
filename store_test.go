package main

import (
	"bytes"
	"io"
	"testing"
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

	r,err := s.read(key)
	if err!=nil{
		t.Error(err)
	}

	b,err := io.ReadAll(r)
	if err!=nil{
		t.Error(err)
	}

	if string(b) != string(data){
		t.Errorf("want %s, have %s",data,b)
	}
}