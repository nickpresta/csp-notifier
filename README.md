# csp-notifier (Content Security Policy)

Provides an endpoint to `POST` [Reports](http://www.html5rocks.com/en/tutorials/security/content-security-policy/#reporting) and a simple interface to view all the reports.


## Requirements

* Go Version 1.0+ (not required if you're deploying -- build the binary as explained below)
* MongoDB server (2.4+)

## How to run

1. Export an environment variable called `MONGO_CONNECT`
	`MONGO_CONNECT="mongodb://user:password@your.mongo.server.com:PORT/Database"`
2. If you're developing locally: `go run main.go`
3. If you're deploying to production: `go build` will produce a `csp-notifier` binary. Copy it to your server.
