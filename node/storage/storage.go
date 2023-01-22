package storage

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type Storage struct {
	DB *leveldb.DB
}

func (s *Storage) Init(storage_path string) {
	db, err := leveldb.OpenFile(storage_path, nil)

	if err != nil {
		log.Fatal(err)
	}
	s.DB = db
}

func (s *Storage) Close() {
	s.DB.Close()
}

func (s *Storage) Get(key []byte) ([]byte, error) {
	return s.DB.Get(key, nil)
}

func (s *Storage) Put(key []byte, value []byte) error {
	return s.DB.Put(key, value, nil)
}

func (s *Storage) Delete(key []byte) error {
	return s.DB.Delete(key, nil)
}
