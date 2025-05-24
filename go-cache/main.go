package main

import "fmt"

const SIZE = 5

// Node struct
type Node struct {
	Val   string // Node Value
	Left  *Node  // Pointer to the Left Node
	Right *Node  // Pointer to the Right Node
}

// Queue struct
type Queue struct {
	Head   *Node // Queue's Head Node
	Tail   *Node // Queue's Tail Node
	Length int   // Queue's Length
}

// Hash Type
type Hash map[string]*Node

// LRU Cache struct
type Cache struct {
	Queue Queue
	Hash  Hash
}

// Initialize New Queue
func NewQueue() Queue {

	head := &Node{}
	tail := &Node{}

	head.Right = tail
	tail.Left = head

	return Queue{Head: head, Tail: tail}
}

// Initialize new LRU Cache
func NewCache() Cache {
	return Cache{
		Queue: NewQueue(),
		Hash:  Hash{},
	}
}

// Cache Add Function
func (c *Cache) Add(n *Node) {
	fmt.Printf("add: %s\n", n.Val)
	tmp := c.Queue.Head.Right

	c.Queue.Head.Right = n
	n.Left = c.Queue.Head
	n.Right = tmp
	tmp.Left = n

	c.Queue.Length++

	if c.Queue.Length > SIZE {
		c.Remove(c.Queue.Tail.Left)
	}
}

// Cache Remove Function
func (c *Cache) Remove(n *Node) *Node {
	fmt.Printf("remove: %s\n", n.Val)
	left := n.Left
	right := n.Right

	left.Right = right
	right.Left = left

	c.Queue.Length -= 1

	delete(c.Hash, n.Val)

	return n
}

// Cache Check Function
func (c *Cache) Check(str string) {
	node := &Node{}

	if val, ok := c.Hash[str]; ok {
		node = c.Remove(val)
	} else {
		node = &Node{Val: str}
	}

	c.Add(node)
	c.Hash[str] = node
}

// Cache Display Function
func (c *Cache) Display() {
	c.Queue.Display()
}

// Queue Display
func (q *Queue) Display() {
	node := q.Head.Right
	fmt.Printf("%d - [", q.Length)

	for i := 0; i < q.Length; i++ {
		fmt.Printf("{%s}", node.Val)
		if i < q.Length-1 {
			fmt.Printf("<-->")
		}
		node = node.Right
	}
	fmt.Println("]")
}

// Simple Go LRU Cache
func main() {
	fmt.Println("::: Go Cache :::")

	cache := NewCache()

	for _, word := range []string{"parrot", "avocado", "dragonfruit", "tree", "potato", "tomato", "tree"} {
		fmt.Println(word)
		cache.Check(word)
		cache.Display()
	}
}
