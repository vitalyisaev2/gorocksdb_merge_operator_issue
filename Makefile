build-go:
	go build -o mve-go ./go/...

build-cpp:
	g++ cpp/main.cpp cpp/merge_operator.cpp -o mve-cpp -I/usr/local/include/rocksdb/ -I./cpp -lrocksdb -O2 -std=c++11

clean:
	rm ./mve-go ./mve-cpp || true
	rm massif* || true
