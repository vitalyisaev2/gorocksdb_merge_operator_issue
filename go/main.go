package main

import (
	"log"
	"os"
	"runtime"

	"github.com/tecbot/gorocksdb"
)

const (
	mergeDummy = "dummy"
	mergeReal = "real"
	actionGenerate = "generate"
	actionIterate = "iterate"
)

func main() {

	if len(os.Args) != 3 {
		log.Fatalf("Wrong amount of arguments: %v\n", os.Args)
	}

	// initialize merge operator
	var mergeOperator gorocksdb.MergeOperator
	switch os.Args[1] {
	case mergeDummy:
		mergeOperator = &dummyMergeOperator{}
	case mergeReal:
		mergeOperator = &realMergeOperator{}
	default:
		log.Fatalf("Wrong merge operator mode: %v\n", os.Args[1])
	}

	var (
		dbFactory func(gorocksdb.MergeOperator) (*gorocksdb.DB, func(), error)
		action func(db *gorocksdb.DB) error
	)
	switch os.Args[2] {
	case actionGenerate:
		log.Fatal("Don't use data generation, use prepared dump")
		// dbFactory = openDBForWriting
		// action = performGeneration
	case actionIterate:
		dbFactory = openDBForReading
		action = performIteration
	default:
		log.Fatalf("Wrong action: %v\n", os.Args[2])
	}

	// create database instance
	db, free, err := dbFactory(mergeOperator)
	if err != nil {
		log.Fatal(err)
	}
	defer free()

	// perform action
	if err := action(db); err != nil {
		log.Fatal(err)
	}

}

func openDBForWriting(mo gorocksdb.MergeOperator) (*gorocksdb.DB, func(), error) {
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	filter := gorocksdb.NewBloomFilter(10)
	bbto.SetFilterPolicy(filter)

	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	opts.SetBlockBasedTableFactory(bbto)
	opts.PrepareForBulkLoad()
	opts.SetMergeOperator(mo)

	db, err := gorocksdb.OpenDb(opts, "segments")
	if err != nil {
		return nil, nil, err
	}

	free := func() {
		db.Close()
		opts.Destroy()
		bbto.Destroy()
	}
	return db, free, nil
}

func openDBForReading(mo gorocksdb.MergeOperator) (*gorocksdb.DB, func(), error) {

	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	filter := gorocksdb.NewBloomFilter(10)
	bbto.SetFilterPolicy(filter)

	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetMergeOperator(mo)
	opts.IncreaseParallelism(runtime.NumCPU())
	opts.OptimizeLevelStyleCompaction(512 * 1024 * 1024)

	db, err := gorocksdb.OpenDb(opts, "segments")
	if err != nil {
		return nil, nil, err
	}

	free := func() {
		db.Close()
		opts.Destroy()
		bbto.Destroy()
	}
	return db, free, nil
}
