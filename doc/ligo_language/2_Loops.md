# Loops

There are two loop constructs
 + `loop`
 + `in`

## `loop`

`loop` is a simple C/C++ while loop kind of construct.
It will run the loop body unless the condition is false.

**Syntax**

```clojure
(loop CONDITION
	(LOOP_BODY))
```

**Example**

```clojure
;; declaring variables
(var numbers [1 2 3 4 5 6 7])
(var sum 0)

(var i 0)
(var sum2 0)

(loop (< i (len numbers))
      (progn
        (set sum2 (+ sum2 (array-index numbers i)))
        (set i (+ i 1))
        ))
```

## `in`

`in` loop is used for traversing through an array easily. It is like python's `for i in range(5)` construct.

**Syntax**

```clojure
(in ARRAY ARRAY_VARIABLE LOOP_BODY)
```

**Example**

```clojure
(var numbers [1 2 3 4 5 6 7])
(var sum 0)

(in numbers number
    (progn
      (set sum (+ sum number))))
```


Next Section : ~~[Functions](3_Functions.md)~~