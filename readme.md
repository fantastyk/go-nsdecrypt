# go-nsdecrypt

This is a Go port of the popular netscaler cookie decryptor [nsscookiedecrypt.py](https://github.com/catalyst256/Netscaler-Cookie-Decryptor/blob/master/nsccookiedecrypt.py)

# Building
### Windows 
```bash
GOOS=windows GOARCH=amd64 go build -o nsdecrypt main.go
```

### Linux
```bash
GOOS=linux GOARCH=amd64 go build -o nsdecrypt main.go
```

### Mac
```bash
GOOS=darwin GOARCH=amd64 go build -o nsdecrypt main.go
```

