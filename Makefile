
all: client.go server.go cli.go 
	go build client.go cli.go
	go build server.go
clean:
	rm client server
