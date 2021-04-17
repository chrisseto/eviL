.PHONY: assets
assets:
	cd assets && yarn run parcel build -d ../dist --public-url /assets index.html components/*.html views/*.html

.PHONY: twitter
twitter:
	go build -o $@ ./*.go

.PHONY: server
server: assets twitter
	./twitter
