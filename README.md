### gorocksdb_merge_operator_issue

We are experiencing uncontrolled process memory growth during iteration over the whole RocksDB database. After some tests we've find out that the memory allocated within `MergeOperator` is actually held forewer.

#### Prerequisites
```
go >= 1.10
librocksdb.so >= 5.13
valgrind
massif-visualizer
```

#### Steps to reproduce

Get sources and compile binary:
```
go get -v https://github.com/vitalyisaev2/gorocksdb_merge_operator_issue
cd $GOPATH/src/github.com/vitalyisaev2/gorocksdb_merge_operator_issue
```

Download database dump `segments.tar.gz` from Google Drive:
https://drive.google.com/file/d/13pn0ZW2qt4Tb9c5hPYer0HjgGt_rJtNR/view?usp=sharing

Unpack database dump:
```
tar xzvf segments.tar.gz`
```

Run database iterators in two different modes:
```
valgrind --tool=massif ./gorocksdb_merge_operator_issue dummy iterate
valgrind --tool=massif ./gorocksdb_merge_operator_issue real iterate
```

Launch GUI tool to visualize heap profile:
```
massif-visualizer massif.out.$PID
```

#### Problem descrpition
