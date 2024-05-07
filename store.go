package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "kknetwork"

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
	Root							string//Root is the folder name of the root,containing all the folders/files of the system.
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil{
		opts.PathTransformFunc=DefaultPathTransformFunc
	}
	
	if len(opts.Root)==0{
		opts.Root=defaultRootFolderName
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool{
	pathKey:=s.PathTransformFunc(key)
	fullPathWithRoot:= fmt.Sprintf("%s/%s",s.Root,pathKey.FullPath())
	_,err:= os.Stat(fullPathWithRoot)
	return !errors.Is(err,os.ErrNotExist)
}

func (s *Store)Clear() error{
	return os.RemoveAll(s.Root)
}

func (s *Store) Delete(key string) error{
	pathKey := s.PathTransformFunc(key)

	defer func(){
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()
	firstPathNameWithRoot:=fmt.Sprintf("%s/%s",s.Root,pathKey.FirstPathName())
	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *Store) Read(key string) (int64,io.Reader, error){
	return s.readStream(key)
}

func (s *Store) readStream(key string)(int64,io.ReadCloser,error){
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s",s.Root,pathKey.FullPath())
	file,err:= os.Open(fullPathWithRoot)
	if err!=nil{
		return 0,nil,err
	}
	
	fi,err:= file.Stat()
	if err!=nil{
		return 0,nil,err
	}
	return fi.Size(),file,nil
}

func (s *Store) Write(key string,r io.Reader) (int64,error){
	return s.writeStream(key,r)
}

func (s *Store) WriteDecrypt(encKey []byte,key string,r io.Reader)(int64,error){
	f,err := s.openfileForWriting(key)
	if err!=nil{
		return 0,err
	}

  n,err:= copyDecrypt(encKey,r,f)
	if err!=nil{
		return 0,err
	} 
	return int64(n),err
}

func (s *Store) openfileForWriting(key string)(*os.File,error){
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s",s.Root,pathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot,os.ModePerm);err!=nil{
		return nil,err
	}
	fullPathWithRoot := fmt.Sprintf("%s/%s",s.Root,pathKey.FullPath())

	return os.Create(fullPathWithRoot)
}

func (s *Store) writeStream(key string, r io.Reader) (int64,error) {

	f,err:= s.openfileForWriting(key)
	if err!=nil{
		return 0,err
	}
	// defer f.Close()
  return io.Copy(f,r)
}