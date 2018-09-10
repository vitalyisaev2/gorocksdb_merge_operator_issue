## gorocksdb_merge_operator_issue

We have `Go` application that uses `Rocksdb` via `CGO` wrappers (https://github.com/tecbot/gorocksdb). Currently we are experiencing uncontrolled process memory growth at the time of the iteration over the whole RocksDB database. After some tests we've find out that memory allocated within `MergeOperator` is actually never freed. Two minimal working examples in `Go` and `C++` are provided to explore this issue. Both examples are just iterating several times over the database. Please follow these steps to reproduce:

Install build dependencies the way you prefer:
1. `go` 1.11
2. `g++` 7.3.1
3. `librocksdb.so`: for example, headers should go to `/usr/local/include/rocksdb`, binaries should go to `/usr/loca/lib`; you can [build shared library](https://github.com/facebook/rocksdb/blob/master/INSTALL.md) on your own, or, assuming that you're on Linux x86_64, you can download [this archive](https://drive.google.com/open?id=1LDpjtE0rS5j_MlAVezB8sG5hsLYR2uUO)
4. `valgrind`
5. `massif-visualizer`

Get project sources:
```
go get -v https://github.com/vitalyisaev2/gorocksdb_merge_operator_issue || true
cd $GOPATH/src/github.com/vitalyisaev2/gorocksdb_merge_operator_issue
```

Download database dump `segments.tar.gz` from Google Drive:
https://drive.google.com/file/d/13pn0ZW2qt4Tb9c5hPYer0HjgGt_rJtNR/view?usp=sharing

Unpack database dump:
```
tar xzvf segments.tar.gz`
```

Compile binaries:
```
make build
```
____

### Go

We provide two implementations of `MergeOperator`:
1. `dummy` that actually does nothing;
2. `real` which allocates some memory during `[]byte` concatenation, emulating the behaviour of real-life `MergeOperator` implementation;

Run database iterators in two different modes:
```
valgrind --tool=massif ./mve-go dummy iterate
valgrind --tool=massif ./mve-go real iterate
```

Launch GUI tool to visualize heap profile:
```
massif-visualizer massif.out.$PID
```

Profiler HTTP server is also available on http://localhost:6060/debug/pprof.

Everything is fine for `dummy` operator:
![dummy](https://github.com/vitalyisaev2/gorocksdb_merge_operator_issue/blob/master/go/profile.dummy.jpeg)

With `real` operator the heap is leaking. It turns out that the huge amount of memory is allocated within `CGO` parts of the application code (because it is hidden behind `runtime.asmcgocall`), and this memory is never freed.
![real](https://github.com/vitalyisaev2/gorocksdb_merge_operator_issue/blob/master/go/profile.real.jpeg)

____

### C++

For C++ we provide only `real` implementation.
```
valgrind --tool=massif ./mve-cpp
```
Everything is fine here:
![real](https://github.com/vitalyisaev2/gorocksdb_merge_operator_issue/blob/master/cpp/profile.real.jpeg)
