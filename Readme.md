# ligo - scheme like language interpreter in golang

+ This is just a hobby project.
+ No documentation
+ Just a single example to demonstrate the working.

## Try it?

This is just a package.

```shell
$ go get -u github.com/aki237/ligo
```

To try the interpreter :
```shell
$ cd $GOPATH/src/github.com/aki237/ligo/example/interpreter/
$ go build
$ ./interpreter sample/main.lg
Hello,
world!
This
is
hello
from
ligo!!
```

## FAQ

+ **Scheme?**

  This is a scheme like but little different in syntax.
+ **So it is not scheme?**

  Don't worry it has all the paranthesis goodness of scheme.
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
