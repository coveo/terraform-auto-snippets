install:
	go install

deploy:
	GOOS=linux go build -o .pkg/tfautosnip_linux
	GOOS=darwin go build -o .pkg/tfautosnip_darwin
	GOOS=windows go build -o .pkg/tfautosnip.exe
