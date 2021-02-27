# Go-Keyre
Go-Keyre is an in-memory key value store that supports thread safe concurrent read and write operations.
It also has a persistence feature which stores the gob of the database in the disk.

## Running Locally
* Prerequisites
  * [Go](https://golang.org/dl/)

1. **Clone the repository**
```
    git clone https://github.com/salmanahmed404/go-keyre.git
```

2. **Build & Execute**
```
    go build go-keyre
    ./go-keyre
```

3. **Check if the server is up & running**
```
    telnet localhost 1200
    > PING
    PONG
```

