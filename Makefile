build:
	go build

clean:
	rm ./gorocksdb_merge_operator_issue || true
	rm massif* || true
