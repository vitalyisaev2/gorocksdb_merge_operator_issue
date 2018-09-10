build-go:
	go get -u -v -d ./src/mve-go/...
	go build -o mve-go ./src/mve-go/...

build-cpp:
	g++ ./src/mve-cpp/*.cpp -I/usr/local/include/rocksdb/ -I./src/mve-cpp -lrocksdb -O2 -g -std=c++11 -fno-rtti -o mve-cpp

build: build-go build-cpp

clean:
	rm ./mve-go ./mve-cpp || true
	rm massif* || true
