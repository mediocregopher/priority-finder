/*
<!--
Copyright (c) 2016 Christoph Berger. Some rights reserved.
Use of this text is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->
*/
package main

import (
	"fmt"
	"strings"
)

// Node gets a new field, `bal`, to store the height difference between the node's subtrees.
type Node struct {
	Value string
	Left  *Node
	Right *Node
	bal   int // height(n.Right) - height(n.Left)
}

/* ### The modified `Insert` function
 */

// `Insert` takes a search value and inserts a new node (unless a node with the given
// search value already exists, in which case `Insert` only replaces the data).
//
// It returns:
//
// * `true` if the height of the tree has increased.
// * `false` otherwise.
func (n *Node) Insert(value string, less func(string, string) bool) bool {
	// The following actions depend on whether the new search value is equal, less, or greater than
	// the current node's search value.
	if less(value, n.Value) {
		// If there is no left child, create a new one.
		if n.Left == nil {
			// Create a new node.
			n.Left = &Node{Value: value}
			// If there is no right child, the new child node has increased the height of this subtree.
			if n.Right == nil {
				// The new left child is the only child.
				n.bal = -1
			} else {
				// There is a left and a right child. The right child cannot have children;
				// otherwise the tree would already have been out of balance at `n`.
				n.bal = 0
			}
		} else {
			// The left child is not nil. Continue in the left subtree.
			if n.Left.Insert(value, less) {
				// If the subtree's balance factor has become either -2 or 2, the subtree must be rebalanced.
				if n.Left.bal < -1 || n.Left.bal > 1 {
					n.rebalance(n.Left)
				} else {
					// If no rebalancing occurred, the left subtree has grown by one: Decrease the balance of the current node by one.
					n.bal--
				}
			}
		}
		// This case is analogous to `value < n.Value`, except that everything is mirrored.
	} else {
		if n.Right == nil {
			n.Right = &Node{Value: value}
			if n.Left == nil {
				n.bal = 1
			} else {
				n.bal = 0
			}
		} else {
			if n.Right.Insert(value, less) {
				if n.Right.bal < -1 || n.Right.bal > 1 {
					n.rebalance(n.Right)
				} else {
					n.bal++
				}
			}
		}
	}
	if n.bal != 0 {
		return true
	}
	// No more adjustments to the ancestor nodes required.
	return false
}

/* ### The new `rebalance()` method and its helpers `rotateLeft()`, `rotateRight()`, `rotateLeftRight()`, and `rotateRightLeft`.

 **Important note: Many of the assumptions about balances, left and right children, etc, as well as much of the logic usde in the functions below, apply to the `Insert` operation only. For `Delete` operations, different rules and operations apply.** As noted earlier, this article focuses on `Insert` only, to keep the code short and clear.
 */

// `rotateLeft` takes a child node and rotates the child node's subtree to the left.
func (n *Node) rotateLeft(c *Node) {
	// Save `c`'s right child.
	r := c.Right
	// `r`'s left subtree gets reassigned to `c`.
	c.Right = r.Left
	// `c` becomes the left child of `r`.
	r.Left = c
	// Make the parent node (that is, the current one) point to the new root node.
	if c == n.Left {
		n.Left = r
	} else {
		n.Right = r
	}
	// Finally, adjust the balances. After a single rotation, the subtrees are always of the same height.
	c.bal = 0
	r.bal = 0
}

// `rotateRight` is the mirrored version of `rotateLeft`.
func (n *Node) rotateRight(c *Node) {
	l := c.Left
	c.Left = l.Right
	l.Right = c
	if c == n.Left {
		n.Left = l
	} else {
		n.Right = l
	}
	c.bal = 0
	l.bal = 0
}

// `rotateRightLeft` first rotates the right child of `c` to the right, then `c` to the left.
func (n *Node) rotateRightLeft(c *Node) {
	// `rotateRight` assumes that the left child has a left child, but as part of the rotate-right-left process,
	// the left child of `c.Right` is a leaf. We therefore have to tweak the balance factors before and after
	// calling `rotateRight`.
	// If we did not do that, we would not be able to reuse `rotateRight` and `rotateLeft`.
	c.Right.Left.bal = 1
	c.rotateRight(c.Right)
	c.Right.bal = 1
	n.rotateLeft(c)
}

// `rotateLeftRight` first rotates the left child of `c` to the left, then `c` to the right.
func (n *Node) rotateLeftRight(c *Node) {
	c.Left.Right.bal = -1 // The considerations from rotateRightLeft also apply here.
	c.rotateLeft(c.Left)
	c.Left.bal = -1
	n.rotateRight(c)
}

// `rebalance` brings the (sub-)tree with root node `c` back into a balanced state.
func (n *Node) rebalance(c *Node) {
	switch {
	// Left subtree is too high, and left child has a left child.
	case c.bal == -2 && c.Left.bal == -1:
		n.rotateRight(c)
	// Right subtree is too high, and right child has a right child.
	case c.bal == 2 && c.Right.bal == 1:
		n.rotateLeft(c)
	// Left subtree is too high, and left child has a right child.
	case c.bal == -2 && c.Left.bal == 1:
		n.rotateLeftRight(c)
	// Right subtree is too high, and right child has a left child.
	case c.bal == 2 && c.Right.bal == -1:
		n.rotateRightLeft(c)
	}
}

// `Find` stays the same as in the previous article.
//func (n *Node) Find(s string) (string, bool) {
//
//	if n == nil {
//		return "", false
//	}
//
//	switch {
//	case s == n.Value:
//		return n.Data, true
//	case s < n.Value:
//		return n.Left.Find(s)
//	default:
//		return n.Right.Find(s)
//	}
//}

// `Dump` dumps the structure of the subtree starting at node `n`, including node search values and balance factors.
// Parameter `i` sets the line indent. `lr` is a prefix denoting the left or the right child, respectively.
func (n *Node) Dump(i int, lr string) {
	if n == nil {
		return
	}
	indent := ""
	if i > 0 {
		//indent = strings.Repeat(" ", (i-1)*4) + "+" + strings.Repeat("-", 3)
		indent = strings.Repeat(" ", (i-1)*4) + "+" + lr + "--"
	}
	fmt.Printf("%s%s[%d]\n", indent, n.Value, n.bal)
	n.Left.Dump(i+1, "L")
	n.Right.Dump(i+1, "R")
}

/*
## Tree

Changes to the Tree type:

* `Insert` now takes care of rebalancing the root node if necessary.
* A new method, `Dump`, exist for invoking `Node.Dump`.
* `Delete` is gone.

*/

//
type Tree struct {
	Root *Node
}

func (t *Tree) Insert(value string, less func(string, string) bool) {
	if t.Root == nil {
		t.Root = &Node{Value: value}
		return
	}
	t.Root.Insert(value, less)
	// If the root node gets out of balance,
	if t.Root.bal < -1 || t.Root.bal > 1 {
		t.rebalance()
	}
}

// `Node`'s `rebalance` method is invoked from the parent node of the node that needs rebalancing.
// However, the root node of a tree has no parent node.
// Therefore, `Tree`'s `rebalance` method creates a fake parent node for rebalancing the root node.
func (t *Tree) rebalance() {
	fakeParent := &Node{Left: t.Root, Value: "fakeParent"}
	fakeParent.rebalance(t.Root)
	// Fetch the new root node from the fake parent node
	t.Root = fakeParent.Left
}

//func (t *Tree) Find(s string) (string, bool) {
//	if t.Root == nil {
//		return "", false
//	}
//	return t.Root.Find(s)
//}

func (t *Tree) Traverse(n *Node, f func(*Node)) {
	if n == nil {
		return
	}
	t.Traverse(n.Left, f)
	f(n)
	t.Traverse(n.Right, f)
}

// `Dump` dumps the tree structure.
func (t *Tree) Dump() {
	t.Root.Dump(0, "")
}
