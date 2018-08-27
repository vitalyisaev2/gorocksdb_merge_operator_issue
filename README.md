## gorocksdb_merge_operator_issue

We are experiencing uncontrolled process memory growth during iteration over the whole RocksDB database in one of our Go applications. After some tests we've find out that the memory allocated within `MergeOperator` is actually never freed. This is the minimal working example reproducing this issue. The application just iterates several times over the database.

There are two implementations of `MergeOperator`:
1. `dummy` that actually does nothing;
2. `real` which allocates some memory, emulating the behaviour of real-life `MergeOperator` implementation;

The use of `real` implementaion results in memory leak on the `C++` side of the application. 

Please follow these steps to reproduce:


### Go

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

##### `Dummy`
Everything is fine here.
![dummy](https://github.com/vitalyisaev2/gorocksdb_merge_operator_issue/blob/master/go/profile.dummy.jpeg)

##### `Real`
It turns out that the huge amount of memory is allocated within `CGO` parts of the application code (it is hidden behind `runtime.asmcgocall`), and this memory is never freed.
![real](https://github.com/vitalyisaev2/gorocksdb_merge_operator_issue/blob/master/go/profile.real.jpeg)

### C++

