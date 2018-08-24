package main

import (
	"github.com/tecbot/gorocksdb"
	"math/rand"
	"log"
)


func performGeneration(db *gorocksdb.DB) error {
	const (
		keyTotal    = 1024 * 1024
		keyLength   = 32
		valueLength = 64
	)

	wo := gorocksdb.NewDefaultWriteOptions()
	wo.SetSync(false)
	defer wo.Destroy()

	for i := 0; i < keyTotal; i++ {
		key := randBytes(keyLength)
		value := randBytes(valueLength)
		if err := db.Put(wo, key, value); err != nil {
			return err
		}

		if i%(keyTotal/10) == 0 {
			log.Printf("Progress: %v\n", 100*float64(i)/float64(keyTotal))
		}
	}

	return nil
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}


