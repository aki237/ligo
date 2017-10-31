# Conditions

There are 2 condional constructs in ligo language.
 + `if...else`
 + `match`

## `if...else`

`if else` is one of the conditionals used in the language.

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

In this case the `if...else` evaluates to `true`. So `allowed` is assigned `true`.


This can be translated as :

```clojure
(if (> age 18)
    (var allowed true)
  (var allowed false))
```

## `match`

`match` conditional is similar to `switch...case` in `C/C++`. But there are no other keywords used unlike in `C/C++` (like `case`, `break`).
You can match any kind of variable like strings, floats, etc.,

The syntax is simple :

```clojure
(match MATCH_VARIABLE
    CASE_1_VARIABLE (CASE_1_BODY)
    CASE_2_VARIABLE (CASE_2_BODY)
    CASE_3_VARIABLE (CASE_3_BODY)
    CASE_4_VARIABLE (CASE_4_BODY)
    _               (DEFAULT_CASE_BODY)  ;; `default:` is here denoted as `_`
    )
```

**Example**

```clojure
(var number 1)

(match number
    0    (println "zero")
    1    (println "one")           ;; In this case this will succeed
    2    (println "two")
    3    (println "three")
    4    (println "four")
    5    (println "five")
    6    (println "six")
    7    (println "seven")
    8    (println "eight")
    9    (println "nine")
    _    (println "greater than 10")
)
```

### Returns

Similar to `if...else`, match case can return values.

**Example**

```clojure
(var number 1)
(var numberString (match number
                    0    "zero"
                    1    "one"
                    2    "two"
                    3    "three"
                    4    "four"
                    5    "five"
                    6    "six"
                    7    "seven"
                    8    "eight"
                    9    "nine"
                    _    "greater than 10"
                ))
```

In this case the `match` construct evaluates to `"one"`. That value is set to `numberString`.


Next Section : [Loops](2_Loops.md)
