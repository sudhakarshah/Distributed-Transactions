
all: client.go server.go cli.go mp3_box.go lock.go
	go build client.go cli.go
	go build server.go mp3_box.go lock.go
clean:
	rm client server
