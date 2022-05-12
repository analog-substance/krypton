all:: amd64 386

amd64::
	GOARCH=amd64 go build -o build/krypton_amd64
386::
	GOARCH=386 go build -o build/krypton_386