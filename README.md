# go-UDP-wejay

A UDP server for handling user-state

**Requirements**:

Add a .env (copy contents from .env-template) with your client-id and -secret from [spotify](https://developer.spotify.com/dashboard/applications/)

**To start**:

```sh
go run ./
```

**Tests**:

```sh
go test -gcflags=-l ./...
# or
make test
```
