package doublemap

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

type DoubleMap struct {
	tree *tree
}

type Value interface{}
type short [2]byte
type tree [65536]node

type node struct {
	key  short
	leaf [2]Value
	tree *tree
}

func singlehash(a Value) [20]byte {
	ha := sha1.New()
	fmt.Fprintf(ha, "%#v", a)
	sa := ha.Sum(nil)

	var ret [20]byte
	for i := 0; i < 20; i++ {
		j := 19 - i
		ret[i] = sa[j]
	}
	
	return ret
}

func doublehash(a, b Value) [20]short {
	ha, hb := sha1.New(), sha1.New()

	fmt.Fprintf(ha, "%#v", a)
	fmt.Fprintf(hb, "%#v", b)

	sa, sb := ha.Sum(nil), hb.Sum(nil)

	var ret [20]short
	for i := 0; i < 20; i++ {
		j := 19 - i
		ret[i] = short{sa[j], sb[j]}
	}
	
	return ret
}

func New() DoubleMap {
	return DoubleMap{new(tree)}
}

func (dm DoubleMap) Add(a, b Value) {
	hash := doublehash(a, b)
	space := dm.tree
	nilshort := short([2]byte{0, 0})

Walker:
	for h, pair := range hash {
		l := 0
		for i, n := range space {
			if n.key != nilshort {
				l = i
			}

			if n.key == pair {
				if n.tree == nil {
					n.leaf = [2]Value{a, b}
					return
				}

				space = n.tree
				continue Walker
			}
		}

		l++
		if h == 19 {
			space[l] = node{
				key:  pair,
				leaf: [2]Value{a, b},
			}
			
		} else {
			space[l] = node{
				key:  pair,
				tree: new(tree),
			}

			space = space[l].tree
		}
	}
}

func (dm DoubleMap) GetFromA(a Value) Value {
	hash := singlehash(a)
	space := dm.tree

Walker:
	for _, single := range hash {
		for _, n := range space {
			if n.key[0] == single {
				if n.tree == nil {
					return n.leaf[1]
				}

				space = n.tree
				continue Walker
			}
		}
	}

	return nil
}

func (dm DoubleMap) GetFromB(b Value) Value {
	hash := singlehash(b)
	space := dm.tree

Walker:
	for _, single := range hash {
		for _, n := range space {
			if n.key[1] == single {
				if n.tree == nil {
					return n.leaf[0]
				}

				space = n.tree
				continue Walker
			}
		}
	}

	return nil
}

func (dm DoubleMap) GetFromEither(c Value) Value {
	a, b := dm.GetFromA(c), dm.GetFromB(c)
	
	if (a != nil) {
		return a
	}
	
	return b
}

func (dm DoubleMap) Verify(a, b Value) bool {
	vb, va := dm.GetFromA(a), dm.GetFromB(b)
	return va == a && vb == b
}


func (t *tree) dump(indent int) {
	nilshort := short([2]byte{0, 0})

	for _, n := range *t {
		if n.key != nilshort {
			if n.tree == nil {
				fmt.Printf("%s- %x: %#v\n", strings.Repeat("  ", indent), n.key, n.leaf)
				continue
			} else {
				fmt.Printf("%s+ %x:\n", strings.Repeat("  ", indent), n.key)
				n.tree.dump(indent+1)
			}
		}
	}
}

func (dm DoubleMap) Dump() {
	dm.tree.dump(0)
}