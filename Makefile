OUTDIR ?= dist

default: test build

test:
	go test ./bc ./bo

APIKEY=
run:
	go run main.go --assets ./static --apikey $(APIKEY)

build: clean
	mkdir $(OUTDIR)
	cp -a static $(OUTDIR)/
	go build -o $(OUTDIR)/bcbo main.go
	cd $(OUTDIR) && zip -qr bcbo.zip bcbo static

clean:
	rm -rf $(OUTDIR)
