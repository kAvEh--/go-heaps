// Package pairing implements a Pairing heap Data structure
//
// Structure is not thread safe.
//
// Reference: https://en.wikipedia.org/wiki/Pairing_heap
package pairing

import (
	"github.com/theodesp/go-heaps"
)

// PairHeap represents a Pairing Heap.
// The zero value for PairHeap Root is an empty Heap.
type PairHeap struct {
	Root       *PairHeapNode
}

// PairHeapNode contains the current Value and the list if the sub-heaps
type PairHeapNode struct {
	// for use by client; untouched by this library
	Value go_heaps.Item
	// List of children PairHeapNodes all containing values less than the Top of the heap
	children []*PairHeapNode
	// A reference to the parent Heap Node
	parent *PairHeapNode
	heap   *PairHeap // The heap to which this node belongs.
}

func (n *PairHeapNode) detach() []*PairHeapNode {
	if n.parent == nil {
		return nil // avoid detaching root
	}
	for _, node := range n.children {
		node.parent = n.parent
	}
	var idx int
	for i, node := range n.parent.children {
		if node == n {
			idx = i
			break
		}
	}
	n.parent.children = append(n.parent.children[:idx], n.parent.children[idx+1:]...)
	n.parent = nil
	return n.children
}

// Init initializes or clears the PairHeap
func (p *PairHeap) Init() *PairHeap {
	p.Root = &PairHeapNode{}
	return p
}

// New returns an initialized PairHeap with the provided Comparator.
func New() *PairHeap { return new(PairHeap).Init() }


// IsEmpty returns true if PairHeap p is empty.
// The complexity is O(1).
func (p *PairHeap) IsEmpty() bool {
	return p.Root.Value == nil
}

// Resets the current PairHeap
func (p *PairHeap) Clear() {
	p.Root = &PairHeapNode{}
}

// Find the smallest item in the priority queue.
// The complexity is O(1).
func (p *PairHeap) FindMin() go_heaps.Item {
	if p.IsEmpty() {
		return nil
	}
	return p.Root.Value
}

// Inserts the value to the PairHeap and returns the PairHeapNode
// The complexity is O(1).
func (p *PairHeap) Insert(v go_heaps.Item) *PairHeapNode {
	n := PairHeapNode{Value: v, heap: p}
	merge(&p.Root, &n)
	return &n
}

// DeleteMin removes the top most value from the PairHeap and returns it
// The complexity is O(log n) amortized.
func (p *PairHeap) DeleteMin() interface{} {
	if p.IsEmpty() {
		return nil
	}
	result := mergePairs(&p.Root, p.Root.children)
	return result.Value
}

// Adjusts the value to the PairHeapNode Value and returns it
// The complexity is O(n) amortized.
func (p *PairHeap) Adjust(node *PairHeapNode, v go_heaps.Item) *PairHeapNode {
	if node == nil || node.heap != p {
		return nil
	}

	if node == p.Root {
		p.DeleteMin()
		return p.Insert(v)
	} else {
		children := node.detach()
		node.Value = v
		mergePairs(&p.Root, append(p.Root.children, append([]*PairHeapNode{node}, children...)...))
		return node
	}
}

// Deletes a PairHeapNode from the heap and returns the Value
// The complexity is O(log n) amortized.
func (p *PairHeap) Delete(node *PairHeapNode) interface{} {
	if node == nil || node.heap != p {
		return nil
	}
	if node == p.Root {
		p.DeleteMin()
	} else {
		children := node.detach()
		mergePairs(&p.Root, append(p.Root.children, children...))
	}
	return node.Value
}

// Do calls function cb on each element of the PairingHeap, in order of appearance.
// The behavior of Do is undefined if cb changes *p.
func (p *PairHeap) Do(cb func(v interface{})) {
	if p.IsEmpty() {
		return
	}
	// Call root first
	cb(p.Root.Value)
	// Then continue to children
	visitChildren(p.Root.children, cb)
}

// Exhausting search of the element that matches v. Returns it as a PairHeapNode
// The complexity is O(n) amortized.
func (p *PairHeap) Find(v go_heaps.Item) *PairHeapNode {
	if p.IsEmpty() {
		return nil
	}

	if  p.Root.Value.Compare(v) == 0 {
		return p.Root
	} else {
		return p.findInChildren(p.Root.children, v)
	}
}

func (p *PairHeap) findInChildren(children []*PairHeapNode, v go_heaps.Item) *PairHeapNode {
	if len(children) == 0 {
		return nil
	}
	var node *PairHeapNode
loop:
	for _, heapNode := range children {
		cmp := heapNode.Value.Compare(v)
		switch {
		case cmp == 0: // found
			node = heapNode
			break loop
		default:
			node = p.findInChildren(heapNode.children, v)
			if node != nil {
				break loop
			}
		}
	}
	return node
}

func visitChildren(children []*PairHeapNode, cb func(v interface{})) {
	if len(children) == 0 {
		return
	}
	for _, heapNode := range children {
		cb(heapNode.Value)
		visitChildren(heapNode.children, cb)
	}
}

func merge(first **PairHeapNode, second *PairHeapNode) *PairHeapNode {
	q := *first
	if q.Value == nil { // Case when root is empty
		*first = second
		return *first
	}

	cmp := q.Value.Compare(second.Value)
	if cmp < 0 {
		// put 'second' as the first child of 'first' and update the parent
		q.children = append([]*PairHeapNode{second}, q.children...)
		second.parent = *first
		return *first
	} else {
		// put 'first' as the first child of 'second' and update the parent
		second.children = append([]*PairHeapNode{q}, second.children...)
		q.parent = second
		*first = second
		return second
	}
}

// Merges heaps together
func mergePairs(root **PairHeapNode, heaps []*PairHeapNode) *PairHeapNode {
	q := *root
	if len(heaps) == 0 {
		*root = &PairHeapNode{
			parent: nil,
		}
		return q
	}
	if len(heaps) == 1 {
		*root = heaps[0]
		heaps[0].parent = nil
		return q
	}
	var merged *PairHeapNode
	for { // iteratively merge heaps
		if len(heaps) == 0 {
			break
		}
		if len(heaps) == 1 {
			// merge odd one out
			merged = merge(&merged, heaps[0])
			break
		}
		if merged == nil {
			merged = merge(&heaps[0], heaps[1])
			heaps = heaps[2:]
		} else {
			merged = merge(&merged, heaps[0])
			heaps = heaps[1:]
		}
	}
	*root = merged
	merged.parent = nil

	return q
}
