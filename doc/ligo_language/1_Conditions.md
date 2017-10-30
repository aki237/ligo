# Conditions

`if else` is the only condition used in the language for now.

The `if...else` construct is very easy to use.

**Syntax**

```clojure
(if CONDITION
    (SUCCESS_CLAUSE) ;; single expression
  (FAILURE_CLAUSE)   ;; single clause, optional.
  )
```

**Example**

```clojure
(require "base")
(var age 22)
(if (> age 18)   ;; ie., (age > 18)
    (println "You are an adult")   ;; in this case this clause will be executed.
  (println "You are not an adult"))
```

If you want to do more than a single expression in any of the clauses,
throw a `progn` in there... like this :

```clojure
(if (> age 18)
    (progn
      (println "You are an adult")
      (set allowed true))
  (println "You are not an adult"))
```

If the `else` clause is not required just leave the `FAILURE_CLAUSE` empty or `()`
**Like this**

```clojure
(if (> age 18)
    (progn
      (println "You are an adult")
      (set allowed true)))
```

### Returns

In ligo `if..else` can also return a value based on a condition.
So consider the previous example.
The variable `allowed` is set to true. Here is another way to do it.

```clojure
(var allowed (if (> age 18) true false))
```

This can be translated as :

```clojure
(if (> age 18)
    (var allowed true)
  (var allowed false))
```

Next Section : ~~[Loops](2_Loops.md)~~
