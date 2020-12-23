# Golang History

## Usage

```shell
go install "github.com/thiagonache/golang-history"
```

## Summary

Go History is a concurrent Go program that records the commands and its output
executed by the user.

It does not implement terminal, so commands that uses terminal like vi
will get control over the command output and it won't be recorded.
