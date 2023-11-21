package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func CASpathTransformFunc(key string) string{
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	
}

type PathTransformFunc func(string) string

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

var DefautPathTransformFunc = func (key string) string  {
	return key
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathName,os.ModePerm);err!=nil{
		return err
	}

	filename := "somefilename"

	pathAndFileName := pathName + "/" + filename
	f,err := os.Create(pathAndFileName)
	if err!=nil{
		return err
	}

  n,err:= io.Copy(f,r)
	if err!=nil{
		return err
	} 
	
	fmt.Printf("Written (%d) bytes to disk: %s",n,pathAndFileName)


	return nil
}