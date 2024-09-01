package main

import (
	"math/rand/v2"
	"sync"
	"time"

	"github.com/WatchJani/pool"
	s "github.com/WatchJani/stack"
)

func main() {
	sl := New(8, 250, 0.25)

	sl.Insert(5, 23)
}

type SkipList struct {
	roots     []*Node
	rootIndex int
	s.Stack[s.Stack[*Node]]
	pool.Pool[Node]
	sync.RWMutex
	percentage float64
}

func New(height, capacity int, percentage float64) *SkipList {
	//fix this part to be dynamic
	stack := s.New[s.Stack[*Node]](250)

	for range 250 {
		stack.Push(s.New[*Node](height))
	}

	roots := make([]*Node, height)
	for index := range roots {
		roots[index] = &Node{}
	}

	return &SkipList{
		roots:      roots,
		Stack:      stack,
		Pool:       pool.New[Node](capacity),
		percentage: percentage,
	}
}

type Node struct {
	next  *Node
	down  *Node
	value int
	key   int
	time  time.Time
	leaf  bool
}

func NewNode(next, down *Node, value, key int, leaf bool) Node {
	return Node{
		time:  time.Now(),
		next:  next,
		down:  down,
		value: value,
		key:   key,
		leaf:  leaf,
	}
}

func (s *SkipList) Insert(key, value int) {
	current := s.roots[s.rootIndex]
	stack, err := s.Stack.Pop()
	if err != nil {
		//Create new stack
	}

	for {
		for current.next != nil && current.next.key < key {
			s.RLock()
			current = current.next
			s.RUnlock()
		}

		if current.leaf || s.rootIndex == 0 {
			break
		}

		stack.Push(current)

		s.RLock()
		current = current.down
		s.RUnlock()
	}

	rightSide := current.next

	node := s.Pool.Insert()
	current.next = node
	*node = NewNode(rightSide, nil, value, key, true) // create new node

	for FlipCoin(s.percentage) {
		temp := node
		leftNode, err := stack.Pop()
		if err != nil {
			//new big height
		}

		rightSide = leftNode.next
		node = s.Pool.Insert()
		leftNode.next = rightSide

		*node = NewNode(rightSide, temp, key, value, false)
	}
}

func FlipCoin(percentage float64) bool {
	return rand.Float64() < percentage
}
