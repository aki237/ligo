# ligo - scheme like language interpreter in golang

[![Join the chat at https://gitter.im/hackingligo/Lobby](https://badges.gitter.im/hackingligo/Lobby.svg)](https://gitter.im/hackingligo/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

+ This is just a hobby project.
+ No documentation
+ Just a single example to demonstrate the working.

## Building

Build the interpreter first. This builds the main package automatically.

```shell
go install github.com/aki237/ligo/cmd/ligo/
```

Now without the packages the interpreter is not functional at all.
This builds some basic dl packages and copies them to the ligo package search directory. (`$HOME/ligo/`)
```go
cd $GOPATH/github.com/aki237/ligo/packages/
./build.sh
```

## FAQ

+ **Scheme?**

  This is a scheme like but little different in syntax.
+ **So it is not scheme?**

  Don't worry it has all the parenthesis goodness of scheme.
+ **How different?**

  This is scheme :
  ```scheme
  (define sum(lambda (x y)
      (+ x y)))
  ```
  This is ligo :
  ```lisp
  (fn sum |x y|
      (+ x y))
  ```
