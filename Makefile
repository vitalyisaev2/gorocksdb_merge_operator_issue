build-go:
	go build -o mve-go ./go/...

build-cpp:
	g++ cpp/*.cpp -I/usr/local/include/rocksdb/ -I./cpp -lrocksdb -O2 -g -std=c++11 -fno-rtti -o mve-cpp

build: build-go build-cpp

clean:
	rm ./mve-go ./mve-cpp || true
	rm massif* || true
