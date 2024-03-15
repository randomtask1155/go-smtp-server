
# go-smtp-server

This basic smtp server will receive emails and emit them to stdout.  it will not actually relay the message anywhere and is meant for lab testing to verify email alert integartions


## usage

```
Usage of ./go-smtp-server:
  -cert string
    	path to tls cert: default is cert.pem (default "cert.pem")
  -k	allow insecure tls [ true | false ] : default true (default true)
  -key string
    	path to tls key: default is cert.key (default "cert.key")
  -l string
    	listen address: default is 0.0.0.0 for all interfaces (default "0.0.0.0")
  -p string
    	listen port: default is 1025 (default "1025")
  -tls string
    	listen port for TLS: default is 1026 (default "1026")
```


## examples 

start a simple insecure smtp server

```
./go-smtp-server 
```

start a tls smpt server.  By default go-smtp-server will look for cert.pem and cert.key in its local path and if found will start tls listener.  

```
./go-smtp-server -cert mycert.cert -key mycert.key
```
