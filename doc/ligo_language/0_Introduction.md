# ligo : language reference

ligo is a lisp (scheme) like language with implementation interpreter and packages in golang.
Similar to lisp, a simple command or function call is done by surrounding the function with
their parameters in a parenthesis.

```clojure
>>> (function param1 param2 ...)
```

A single function can return a single value which can be a integer, floating point, string,
a Go native struct, or even a nil value.

Unlike in C/C++, Go, Java, or Rust, ligo doesn't have to start from a main function.
Like python, the interpreter just reads the given file or any kind of input, and evaluates
them on the fly.

By default the interpreter doesn't have any functions built in except `require`,
`load-plugin` and `exit` (to be explained later). Not even basic arithmetic functions like
`+`, `-` are loaded.

Here's where the packages come into play. `base` package contains all the basic arithmetic
functions and console I/O functions (like `print`, `input`).

To load a package, just run `require` function with the package name (as a string).

Example :

```clojure
>>> (require "base")
```

Now you can call any other function from that library like `(+ 9 6)`.

### Where are these libraries located?

If this is built from source, the libraries will be located in `~/ligo/`, `$HOME/ligo`.
A simple library name denotes a system directory, which contains go native plugin libraries
as well as library files defined in ligo itself.

### Language Inbuilts

Defining variables, assignment, loops etc., are built into the ligo interpreter. Infact
those are keywords. There are a handful of keywords in ligo.

 + `var`
    - defining a new variable, if the passed name is already defined this throws an error.
    - **syntax** : `(var VARIABLE_NAME INIT_VALUE)`
    - **example** : `(var age 45)`
 + `set`
    - setting value to the variable. if the variable name passed is not defined, this throws an error.
    - **syntax** : `(set VARIABLE_NAME VALUE)`
    - **example** : `(set age 67)`
 + `fn` :
    - function declaration, to be discussed later.
 + `return` :
    - return a value
    - **syntax** : `(return VALUE|VARIABLE_NAME)`
    - **example** : `(return age)`, `(return 40)`, `(return "Lisp is awesome!!")`
 + `progn`
    - run a list of lisp expressions, to be discussed later.
 + `loop` :
    - `loop` is a type of loop construct.
    - This loop is like `while` loop in C. ie., `while (CONDITION) {}`
    - syntax to be discussed later.
 + `in` :
    - `in` is another kind of loop construct.
    - similar to the one in python. ie ., `for i in list: ...`
    - syntax to be discussed later
 + `if` :
    - condition construct of ligo.
    - syntax to be discussed later.
 + `match` :
    - `switch...case` condition construct of ligo
    - syntax to be discussed later.
 + `eval` :
    - `eval` keyword is used to evaluate a string as a ligo expression and return the evaluated value.
    - **syntax** : `(eval LIGO_EXPRESSION)`
    - **example** : `(eval "(+ 9 7)")` => 16
 + `fork` :
    - `fork` is used to start a parallel go routine.
 + `delete` :
    - `delete` is used to delete a variable from the interpreter's memory.
    - **syntax** : `(delete VARIABLE_NAME)`
    - **example** : `(delete age)`

### Basic Types
Like any other language, there are some inbuilt types like int, string, float etc.,
Defining them is very simple.
`1` is a simple integer. `3.14` is a simple floating point decimal. `"simple string"` is a
simple string. `true` or `false` can be used for denoting the Boolean.

**Example**

```clojure
>>> (var age 23)
```

ligo automatically recognizes the number and assigns the variable age with an integer value
of `23`

Test the type of age : (include `base` package first)

```clojure
>>> (type age)
Eval : int
```  

Similar to this all other variables can be set up.

Arrays can be set up like following :

```clojure
>>> (var fruits ["apple" "orange" "banana" "Papaya"])
```

The arrays defined can contain any types of variables.

```clojure
>>> (var person ["John Smith" 34 "john.smith@example.org" "212 2212486263" true])
```

Array handling functions are available in the base package.

Next Section : [Condition constructs](1_Conditions.md)
