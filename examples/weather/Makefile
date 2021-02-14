.PHONY: assets
assets:
	cd assets && yarn run parcel build -d ../dist --public-url /assets index.html components/*.html views/*.html

.PHONY: weather
weather:
	go build -o $@ ./*.go

.PHONY: server
server: assets weather
	./weather
