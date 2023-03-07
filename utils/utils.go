package utils

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func Chunk[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

func ExportFile(path string, data []byte) error {
	f, err := os.Create("./files/example_retrive.png")
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func VerifyFile(hash string, data []byte) bool {
	data_hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", data_hash) == hash
}

func DisplayHashTable(hash_table map[int]string) {
	for i := 0; i < len(hash_table); i++ {
		fmt.Printf("%v : %v\n", i, hash_table[i])
	}
}
