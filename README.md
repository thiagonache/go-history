# Golang History

## Usage

```
go install "github.com/thiagonache/golang-history"
```

## Summary

Golang History is a Go program that records the commands executed and its output
executed by the user.

It does not implement terminal, so commands that requires terminal like vi
cannot be execute through this recorder.

## TODO

* Recorder object
* ExecuteAndRecordCommand to take a single string for the command line
* Handle run-time errors such as no disk space when writing history file