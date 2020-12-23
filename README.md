# Golang History

## Usage

```
go install "github.com/thiagonache/golang-history"
```

## Summary

Go History is a concurrent Go program that records the commands and its output
executed by the user.

It does not implement terminal, so commands that requires terminal like vi
cannot be execute through this recorder.

## TODO

- Handle run-time errors such as no disk space when writing history file
