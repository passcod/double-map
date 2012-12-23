DoubleMap
=========

__DoubleMap__ is a hashmap that can be accessed from both the key and the value.
As such, it is a `value<->value` store, not a `key->value` one. It is useful
in situations where one would usually keep two hashmaps in sync like so:

```go
aToB := make(map[A]B)
bToA := make(map[B]A)

func add(a A, b B) {
	aToB[a] = b
	bToA[b] = a
}
```

This has the disadvantage of storing the same content twice. Another approach
is to have a single hashmap and loop through to perform the reverse:

```go
type aToB map[A]B
dmap := make(aToB)

func (m aToB) BToA(b B) a A {
	for k, v := range m {
		if k == b {
			return v
		}
	}
}
```

This has the disadvantage of being *much* slower and intensive in one direction
compared to the other.


__DoubleMap__ solves this by storing a storing both values under their hashes,
in a tree like so:

```plain
+ abb0:
  + 2810:
    + 5490:
      + 3935:
        + e6a8:
          + 46c4:
            + 8dae:
              + c2b7:
                + 18ca:
                  + 4d60:
                    + 5707:
                      + 549c:
                        + 4cf1:
                          + b0cd:
                            + 13cc:
                              + 79ba:
                                + 2b37:
                                  + 1992:
                                    + 6a4b:
                                      - 35da: [2]Value{1, 2}
```

Here, `hash(1)` == `ab 28 54 39 e6 46 8d c2 18 4d 57 54 4c b0 13 79 2b 19 6a 35`
and `hash(2)` == `b0 10 90 35 a8 c4 ae b7 ca 60 07 9c f1 cd cc ba 37 92 4b da`.

Thus, access and setting always take twenty 65536-iteration loops, and presence
checking takes a best case of 1 loop and a worst case of 20. The hash function
used gives a rather good spread, and the theoretical space is (a bit less than)
2^320 value pairs (this is a bit abstract, so let's just say that it's larger
than IPv6's address space).


API
---

__DoubleMap__ is implemented in Go.

```go
import "github.com/passcod/double-map"
```


### Creation

```go
dm := doublemap.New()
```


### Add / Set

```go
dm.Add("abcde", -0.456)
```


### Access

```go
dm.GetFromA("abcde") //=> -0.456
dm.GetFromB(-0.456)  //=> "abcde"

dm.GetFromEither("abcde") //=> -0.456
dm.GetFromEither(-0.456)  //=> "abcde"
```

#### In case of pairs with common members:

```go
dm.Add("I", "me")
dm.Add("I", "myself")

dm.GetFromA("I")      //=> "me"
dm.GetFromB("me")     //=> "I"
dm.GetFromB("myself") //=> "I"
```


### Delete

	Not implemented (yet?)

### Dump

```go
dm.Dump() // Prints structure to STDOUT
```


Legal & Misc.
-------------

Released in the Public Domain or licensed under the
[Unlicense](http://unlicense.org), whichever gives *you* the
most rights in your legislation.

There's different ways to store the data structure. In this first Go
implementation, I chose to use SHA1 to hash, but reverse the hash to
avoid a few problems (e.g. first few bytes could be manipulated /
predicted). Go doesn't provide an easy way to manipulate nibbles
instead of bytes, so I was constrained into using a 20-level tree of
65536-element arrays (with short-sized keys), instead of the more
efficient 40-level tree of 256-element arrays (with byte-sized keys).
You are welcome to try different configurations.