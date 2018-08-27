package main

import (
	"encoding/hex"
	"log"
	"strconv"

	"github.com/tecbot/gorocksdb"
)

// performIteration runs iterator several times
func performIteration(db *gorocksdb.DB) error {
	const n = 5
	for i := 0; i < n; i++ {
		if err := iterate(db); err != nil {
			return err
		}
	}
	return nil
}


// iterate creates new iterator and walks through the whole database
func iterate(db *gorocksdb.DB) error {
	log.Println("Iteration started")
	defer log.Println("Iteration finished")

	// estimate number of keys
	keysTotalString := db.GetProperty("rocksdb.estimate-num-keys")
	keysTotal, err := strconv.Atoi(keysTotalString)
	if err != nil {
		return err
	}

	// create iterator over whole database
	ro := gorocksdb.NewDefaultReadOptions()
	ro.SetFillCache(false)
	ro.SetTailing(true)
	it := db.NewIterator(ro)
	defer it.Close()

	// perform iteration
	var count int
	for it.SeekToFirst(); it.Valid(); it.Next() {
		step(it, keysTotal, &count)
	}

	if err := it.Err(); err != nil {
		return err
	}

	return nil
}

// step makes some fake work with every database key and value
func step(it *gorocksdb.Iterator, keysTotal int, keyCount *int) {
	key := it.Key()
	defer key.Free()
	value := it.Value()
	defer value.Free()

	k := key.Data()
	v := value.Data()

	if *keyCount%(keysTotal/10) == 0 {
		log.Printf("Progress: %v\n", int(100*float64(*keyCount)/float64(keysTotal)))
		log.Printf("Data example: key=%v value_len=%v\n", hex.EncodeToString(k), len(v))
	}
	*keyCount++
}
