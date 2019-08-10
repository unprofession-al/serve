# serve

`serve` is a simple http server which allows so serve a local directory. 
It also comes with a few handy features for my personal needs: 

* Watch the served directory and live-reload the website in the browser when files are changed
* Talk to no-ip in order do serve a web site from a dynamic IP address
* If a `README.md` file is present in the directory you are serving, navigate to `/render/README` and see the rendered markdown

See `serve -h` for infos.

## Installation

Installation is as easy as: 

```
go get -u github.com/unprofession-al/serve
```

Go version >= 1.11.x is required.
