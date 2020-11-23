install:
	go build -o /usr/local/bin/iparse pkg/cmd/cmd.go
	chmod u+x /usr/local/bin/iparse