package main

import (
	"log"
	"math/rand"

	"github.com/tecbot/gorocksdb"
)

func performGeneration(db *gorocksdb.DB) error {
	const (
		keysTotal     = 1 * 1024 * 1024
		recordsPerKey = 10 // every key will be merged 10 times
		keyLength     = 32
		valueLength   = 64
	)

	wo := gorocksdb.NewDefaultWriteOptions()
	wo.SetSync(false)
	defer wo.Destroy()

	// generate keys
	log.Println("Generating keys")
	keys := make([][]byte, keysTotal)
	for i := 0; i < keysTotal; i++ {
		keys[i] = randBytes(keyLength)
	}

	// merge values
	totalOperations := 0
	for j := 0; j < recordsPerKey; j++ {
		for i := 0; i < keysTotal; i++ {

			key := keys[i]
			value := randBytes(valueLength)
			if err := db.Merge(wo, key, value); err != nil {
				return err
			}
			totalOperations++

			if totalOperations%(keysTotal*recordsPerKey/10) == 0 {
				log.Printf("Progress: %v\n", int(100*float64(totalOperations)/float64(keysTotal*recordsPerKey)))
			}
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
