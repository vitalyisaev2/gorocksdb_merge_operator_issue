### gorocksdb_merge_operator_issue

We are experiencing uncontrolled process memory growth during iteration over the whole RocksDB database. After some tests we've find out that the memory allocated within `MergeOperator` is actually never freed. This is the minimal working example reproducing this issue. We prepared two implementations of `MergeOperator`:
1. `dummy` that does nothing;
2. `real` which allocates some memory, emulating the behaviour of real-life `MergeOperator` implementation;
The use of `real` implementaion results in memory leak on the `C++` side of the application. 

Please follow these steps to reproduce:

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

#### Results
