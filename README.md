# Live Upload POC

A lot of inspriation came from Phoenix LiveView.

* `s.Upload("file")` will build a `UploadConfig`
* It's the user's reponsiblity to add the `UploadConfig` to the state
* Files get uploaded over WebSockets in small chunks and reassmebled server side to a temporary file at `tmp/live-c85hnkjin56fohfd77b0.ext`
* `s.UploadConsume` is used to handle moving the temporary file to another destination returning back the public path of this new location

## Getting started

```
cd vendor/github.com/jfyne/live/web
npm run watch
```

```
fd --extension go | entr -r go run main.go
```

http://localhost:8080/upload
