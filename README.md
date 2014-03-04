## Beancounter Backoffice

This is a prototype of a backoffice for Beancounter written in
[Go](http://golang.org) + [AngularJS](http://angularjs.org).

![Backoffice screenshot](https://lh3.googleusercontent.com/-EsjrXEIgetY/UxXiRz-nauI/AAAAAAAAVBY/6ua80crinK0/w599-h515-no/bo-screenshot.png)


### Setup

You will need Go language installed on your platform.

Also, make sure you have `make` command. Alternatively, you could also just
copy&paste commands from Makefile to cmd line terminal.


### Running locally

Execute

  ```make run APIKEY=<customer-key>```

or

  ```go run main.go --assets ./static --apikey key```

Currenly, there's no authentication so you'll need to supply your
customer API key from the command line.

Then navigate to http://localhost:9090.

The app assumes there's a Beancounter Platform instance running on
http://localhost:8080. To use a different instance specify `--bcapi` flag, e.g.

  ```go run main.go --assets ./static --apikey key --bcapi https://beancounter.io/rest```

For all cmd line args run `go run main.go -h`.


### Testing

To run backend tests execute `make test` or `go test ./bc ./bo`.

TODO: front-end Angular app testing with Karma.


### Build

`make build` will create `dist` dir with static assets folder and compiled
binary app.
