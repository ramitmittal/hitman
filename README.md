# hitman

A TUI HTTP client for debugging request/response headers inspired by [restclient.el](https://github.com/pashky/restclient.el) and [httpie](https://github.com/httpie/httpie).

## Features

* You choose the HTTP method and URL for the request *(obviously)*.
* Request + response are printed in a scrollable viewport.
* Your input is auto-saved on exit.

## How to use
* Type a request definition and press `TAB` to send.
```
GET jsonplaceholder.typicode.com/posts/2
Accept: application/json
```
* Use `#` to comment lines in input.
```
GET jsonplaceholder.typicode.com/posts/2
# an http request
```
* Use `"` for escaping spaces and colons.
```
GET "https://jsonplaceholder.typicode.com/posts/2"
Accept: "custom:value with spaces and colons"
```
* Use flags to modify the HTTP client. Flags should be placed after headers.
```
GET "https://jsonplaceholder.typicode.com/posts/2"
Accept: application/json
-insecure
```

![an image](docs/1.PNG)

# Supported Flags
* `-insecure`  
    Skip SSL cert checks.

## What's planned
* Releases.
* Configurable redirects.
* Something to do with request/response bodies.

## Meta
Issues + PRs are welcome!