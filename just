serve:
	goapp serve

deploy:
	appcfg.py --oauth2 update .

format:
	gofmt -w -tabs=false -tabwidth=2 *.go
