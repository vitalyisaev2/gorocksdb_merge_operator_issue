build-go:
	go build ./go -o mve-go

build-cpp:
	g++ cpp/main.cpp -o mve-cpp -I/usr/local/include/rocksdb/ -lrocksdb -O2 -std=c++11

clean:
	rm ./mve-go ./mve-cpp || true
	rm massif* || true
