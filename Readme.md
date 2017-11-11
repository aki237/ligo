# ligo - scheme like language interpreter in golang

[![Join the chat at https://gitter.im/hackingligo/Lobby](https://badges.gitter.im/hackingligo/Lobby.svg)](https://gitter.im/hackingligo/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/aki237/ligo)](https://goreportcard.com/report/github.com/aki237/ligo)
[![GoDoc](https://godoc.org/github.com/aki237/ligo?status.svg)](https://godoc.org/github.com/aki237/ligo)
[![Build Status](https://travis-ci.org/aki237/ligo.svg?branch=master)](https://travis-ci.org/aki237/ligo)

## Introduction

ligo is implemented as a subset of lisp language with some changes in
syntax. The syntax and usage is described in the documentation included
in the [`doc`](doc) folder.

## Building
  + First of all fetch all the packages without installing.
    ```shell
    go get -d github.com/aki237/ligo
    ```
  + `cd` into the project directory
    ```shell
    cd $GOPATH/src/github.com/aki237/ligo
    ```
  + Build the interpreter
    ```shell
    go install ./cmd/ligo
    ```
  + Build the ligo plugin packages
    ```
    cd $GOPATH/src/github.com/aki237/ligo/packages/
    ./build.sh
    ```
    This installs the ligo plugins in `$HOME/ligo`

The ligo interpreter is installed in your `$GOPATH/bin`.

## Usage
A commandline call without any arguments starts a interactive interpreter session.
In that process it also initializes a interpreter by running a start script from the file
`$HOME/.ligorc` (like `.bashrc`, in case of `bash`). Any argument passed is treated
as a file and executes the contents in the file.

## Simple Example

Simple example to get an input from the shell and
```clojure
;; include all the required libraries
(require "base")

(printf "Hello %s!!" (input "Enter your name :"))
```

The bare interpreter has no functionalities (not even basic functionalities).
The `"base"` package includes all the bare minimum functionality.

## Extending the interpreter

Writing packages for the interpreter in Go is very simple and is discussed
in [`this file`](doc/writing_packages/0_Inroduction.md).

### Contributing
Please refer [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.
