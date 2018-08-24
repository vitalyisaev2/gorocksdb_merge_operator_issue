package main

import (
	"github.com/tecbot/gorocksdb"
)

// dummyMergeOperator actually doesn't do any useful work
type dummyMergeOperator struct{}

var _ gorocksdb.MergeOperator = (*dummyMergeOperator)(nil)

func (mo *dummyMergeOperator) FullMerge(_, _ []byte, _ [][]byte) ([]byte, bool) {
	return []byte{}, true
}

func (mo *dummyMergeOperator) PartialMerge(_, _, _ []byte) ([]byte, bool) {
	return []byte{}, true
}

func (mo *dummyMergeOperator) Name() string { return mergeDummy }

// realMergeOperator emulates the behavior of real merger;
// consider a contrived example when merge operation is just a concatenation of bytes
type realMergeOperator struct{}

var _ gorocksdb.MergeOperator = (*realMergeOperator)(nil)

func (mo *realMergeOperator) FullMerge(key, existingValue []byte, operands [][]byte) ([]byte, bool) {
	merged := append([]byte{}, existingValue...)
	for _, op := range operands {
		merged = append(merged, op...)
	}
	return merged, true
}

func (mo *realMergeOperator) PartialMerge(key, leftOperand, rightOperand []byte) ([]byte, bool) {
	var merged []byte
	merged = append(merged, leftOperand...)
	merged = append(merged, rightOperand...)
	return merged, true
}

func (mo *realMergeOperator) Name() string { return mergeReal }
