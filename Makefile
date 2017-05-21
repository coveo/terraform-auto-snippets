install:
	go get -u github.com/jteeuwen/go-bindata/...
	(cd convert-tfdoc; go-bindata -o tempfile_resources.go resources)
	go install ./convert-tfdoc ./extract-tfdoc

deploy:
	GOOS=linux go build -o .pkg/tfautosnip_linux
	GOOS=darwin go build -o .pkg/tfautosnip_darwin
	GOOS=windows go build -o .pkg/tfautosnip.exe
