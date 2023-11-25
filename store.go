package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CASpathTransformFunc(key string) PathKey{
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string,sliceLen)

	for i := 0; i < sliceLen; i++ {
		from,to := i*blockSize , (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}
	
	return PathKey{
		PathName: strings.Join(paths,"/"),
		FileName: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct{
	PathName	string
	FileName	string
}

func (p PathKey)FirstPathName() string{
	paths :=strings.Split(p.PathName,"/")
	if len(paths)==0{
		return ""
	}
	return paths[0]
}

func (p PathKey) FullPath() string{
	return fmt.Sprintf("%s/%s",p.PathName,p.FileName)
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) delete(key string) error{
	pathKey := s.PathTransformFunc(key)

	defer func(){
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	return os.RemoveAll(pathKey.FirstPathName())
}

func (s *Store) read(key string) (io.Reader, error){
	f,err := s.readStream(key)
	if err!=nil{
		return nil,err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_,err =io.Copy(buf,f)
	return buf,err
}

func (s *Store) readStream(key string)(io.ReadCloser,error){
	pathKey := s.PathTransformFunc(key)
	return os.Open(pathKey.FullPath())
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathKey.PathName,os.ModePerm);err!=nil{
		return err
	}

	// buf := new(bytes.Buffer)
	// io.Copy(buf,r)

	// filenameBytes := md5.Sum(buf.Bytes())
	// filename := hex.EncodeToString(filenameBytes[:])
	fullPath := pathKey.FullPath()

	f,err := os.Create(fullPath)
	if err!=nil{
		return err
	}

  n,err:= io.Copy(f,r)
	if err!=nil{
		return err
	} 
	
	fmt.Printf("Written (%d) bytes to disk: %s\n",n,fullPath)


	return nil
}