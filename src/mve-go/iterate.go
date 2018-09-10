package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/tecbot/gorocksdb"
)

// performIteration runs iterator several times
func performIteration(db *gorocksdb.DB) error {
	const n = 10
	for i := 0; i < n; i++ {
		log.Printf("Iteration start %d\n", i)
		if err := iterate(db); err != nil {
			return err
		}
		log.Printf("Iteration finished %d\n", i)
		printDatabaseStats(db)
	}

	log.Println("Waiting for a while in order to let runtime free memory")
	time.Sleep(time.Hour)

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
	var (
		keysCount   int
		valueLenSum int
	)
	for it.SeekToFirst(); it.Valid(); it.Next() {
		valueLenSum += step(it, keysTotal, keysCount)
		keysCount++
	}

	if err := it.Err(); err != nil {
		return err
	}

	log.Printf("Sum of value lengths: %d\n", valueLenSum)

	return nil
}

// step makes some fake work with every database key and value
func step(it *gorocksdb.Iterator, keysTotal int, keysCount int) int {
	key := it.Key()
	defer key.Free()
	value := it.Value()
	defer value.Free()

	_ = key.Data()
	_ = value.Data()

	if keysCount%(keysTotal/10) == 0 {
		log.Printf("Progress: %v\n", int(100*float64(keysCount)/float64(keysTotal)))
	}
	return value.Size()
}

func printDatabaseStats(db *gorocksdb.DB) {
	log.Println("Database stats:")
	properties := []string{
		"estimate-num-keys",
		"cur-size-all-mem-tables",
		"estimate-table-readers-mem",
	}
	for _, property := range properties {
		log.Printf("%s: %d\n", property, getIntProperty(db, property))
	}
}

func getIntProperty(db *gorocksdb.DB, property string) int {
	value := db.GetProperty(fmt.Sprintf("rocksdb.%s", property))
	result, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return result
}
