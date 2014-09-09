serve:
	goapp serve

deploy:
	appcfg.py --oauth2 update .

check:
	go fmt *.go
	for file in *.go; (echo expand -t 2 $$file > tmp.go && mv tmp.go $$file)
	go vet
	@echo go lint
	@golint | sed "/don't use underscores/d"
