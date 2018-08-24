package main

import (
	"github.com/tecbot/gorocksdb"
	"math/rand"
	"log"
)


func performGeneration(db *gorocksdb.DB) error {
	const (
		keysTotal   = 10 * 1024 * 1024
		keyLength   = 32
		valueLength = 64
	)

	wo := gorocksdb.NewDefaultWriteOptions()
	wo.SetSync(false)
	defer wo.Destroy()

	for i := 0; i < keysTotal; i++ {
		key := randBytes(keyLength)
		value := randBytes(valueLength)
		if err := db.Put(wo, key, value); err != nil {
			return err
		}

		if i%(keysTotal/10) == 0 {
			log.Printf("Progress: %v\n", int(100*float64(i)/float64(keysTotal)))
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


