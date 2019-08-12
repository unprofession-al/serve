# serve

![serve](./serve.svg "serve")

`serve` is a simple http server which allows so serve a local directory. 
It also comes with a few handy features for my personal needs: 

* Watch the served directory and live-reload the website in the browser when files are changed
* Talk to no-ip in order do serve a web site from a dynamic IP address
* Markdown files will be rendered as html. Append `?raw` to your request to avoid.

See `serve -h` for infos.

## Disclaimer

Please to not use this for any other reason but local testing. The source is terrible and there is
most likely a ton of bugs allover the place. `serve` is intended for my personal needs and everytime
a new requirement comes up it is quickly applied and made things worse... 

## Installation

### From soucre: 

Installation is as easy as: 

```
go get -u github.com/unprofession-al/serve
```

Go version >= 1.11.x is required.

## Run it

Navigate the directory you want to serve and run `serve`. The address to open in your browser will
be printed to STDOUT:

```
# cd ~/Projects/serve 
# serve
2019/08/12 19:29:53 Listening at http://127.0.0.1:8989
```
