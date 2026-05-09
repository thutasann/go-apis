# Creating and accessing a package for advanced-go-course

In another Repo

```bash
go mod edit -replace github.com/thutasann/cryptit=../cryptit
```

## Installing the application

```bash
go env GOPATH # /Users/thutasann/go
```

```bash
export PATH=$PATH:/Users/thutasann/go/bin
```

```bash
cd cryptit
go install
```

## Go Doc

```bash
go doc decrypt
go doc decrypt.Nimbus
```
