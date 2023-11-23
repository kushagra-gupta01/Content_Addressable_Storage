package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
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
		Original: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct{
	PathName	string
	Original	string
}

func (p PathKey) Filename() string{
	return fmt.Sprintf("%s/%s",p.PathName,p.Original)
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

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathKey.PathName,os.ModePerm);err!=nil{
		return err
	}

	// buf := new(bytes.Buffer)
	// io.Copy(buf,r)

	// filenameBytes := md5.Sum(buf.Bytes())
	// filename := hex.EncodeToString(filenameBytes[:])
	pathAndFileName := pathKey.Filename()

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