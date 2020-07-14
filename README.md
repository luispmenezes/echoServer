# ECHO SERVER
Echo server is a simple http server that outputs incoming request headers and formatted body to console.

# Usage
Run server listening on port **8080** and responding with status code **200**, printing formatted **json**,**xml** or **html** request bodies.
```
./echoServer
```

Run server listening on port **300** and responding with status code **201** and printing raw request bodies. 
```
./echoServer -port 3000 -respCode 201 -raw
```

