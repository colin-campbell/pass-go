## Pass-Go

A complete rewrite of [SnapPass](https://github.com/pinterest/snappass) in Golang.  
SnapChat for passwords! - share passwords and other confidential data quickly and securely.

### Features
* 13.5MB Docker image! (Original is 268MB)
* Redis-compatible database built-in (No need to run a separate container for Redis) 
* One single static binary with no external dependencies whatsoever. 
* Automatic built-in LetsEncrypt support if needed. (No need for SSL terminator)

### How to use
To build container (Docker and Local Golang installation required)  
`make docker`  
`docker run --rm -it -p 5000:5000 digitalist-se/pass-go`  
Visit http://localhost:5000 in your browser.

To run on a different port set the environment variable PASSGO_HTTP_PORT

To enable LetsEncrypt support, set environment variable PASSGO_HTTP_HOSTS to a comma-separated list of hostnames you wish to obtain certificates for.  
Enabling LetsEncrypt sets the listening ports automatically to 80 & 443