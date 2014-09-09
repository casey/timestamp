serve:
	goapp serve

deploy:
	appcfg.py --oauth2 update .

vet:
	go vet

format:
	gofmt -w *.go

lint:
	@golint | sed "/don't use underscores/d"
