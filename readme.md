folder dan file yang harus ada
- config/main.conf
- runtime/logs/

build
```
go build -ldflags "-s -w" -o goapp *.go
```

run
```
./goapp start | run
```