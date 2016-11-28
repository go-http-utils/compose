test:
	go test -v

cover:
	rm -rf *.coverprofile
	go test -coverprofile=compose.coverprofile
	gover
	go tool cover -html=compose.coverprofile
