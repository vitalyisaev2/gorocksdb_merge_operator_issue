build-go:
	go build ./go -o mve-go

clean:
	rm ./mve-go ./mve-cpp || true
	rm massif* || true
