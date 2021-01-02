# Golang History

![Go](https://github.com/thiagonache/go-history/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thiagonache/go-history)](https://goreportcard.com/report/github.com/thiagonache/go-history)

## Usage

```shell
git clone https://github.com/thiagonache/go-history.git
cd go-history
go run cmd/main.go
```

## Summary

This project is my playground project to practice Go. The main idea is to implement stuff that already exist for learning purpose.

Go History is concurrent Go program that records the commands and its output executed by the user. It does not implement terminal, so commands that uses terminal like vi will get control over the command output and it won't be recorded.
