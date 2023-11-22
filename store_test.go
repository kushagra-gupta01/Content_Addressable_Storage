package main

import (
	"bytes"
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

	if pathKey.Original != expectedOriginalKey{
		t.Errorf("Have %s , want %s",pathKey.Original,expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASpathTransformFunc,
	}

	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("myspecialpicture",data);err!=nil{
		t.Error(err)
	}
}