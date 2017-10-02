package storage

import "github.com/ernado/selectel/storage"

func New(user,key string) (storage.API,error){
	return storage.New(user, key)
}