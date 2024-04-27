package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T){
	key := "heyKushagrathisSide"
	pathKey:= CASpathTransformFunc(key)
	expectedFileName := "5040d03b8f3185a5e84e397d86a468dc448cb3a1"
	expectedPathName := "5040d/03b8f/3185a/5e84e/397d8/6a468/dc448/cb3a1"

	if pathKey.PathName != expectedPathName{
		t.Errorf("Have %s , want %s",pathKey.PathName,expectedPathName)
	}

	if pathKey.FileName != expectedFileName{
		t.Errorf("Have %s , want %s",pathKey.FileName,expectedFileName)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t,s)

	for i:=0;i<50;i++{
		key := fmt.Sprintf("foo_%d",i)
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

		if err:= s.Delete(key);err!=nil{
			t.Error(err)
		}

		if ok:= s.Has(key);ok{
			t.Errorf("expected to not have key %s",key)
		}
	}
}

func newStore() *Store{
	opts:= StoreOpts{
		PathTransformFunc: CASpathTransformFunc,
	}	
	return NewStore(opts)
}

func teardown(t *testing.T,s *Store){
	if err:=s.Clear();err!=nil{
		t.Error(err)
	}
}