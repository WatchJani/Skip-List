package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/WatchJani/pool"
	st "github.com/WatchJani/stack"
)

func main() {
	sl := New(32, 2500, 0.5)
	// s := []int{136, 137, 119, 42, 134, 41, 145, 110, 44, 4}

	for index := range 12315 {
		sl.Insert(index, index)
	}

	fmt.Println(sl.Search(12315))
	sl.Clear()

	// sl.Read()
}

type SkipList struct {
	roots     []*Node
	rootIndex int
	st.Stack[st.Stack[*Node]]
	pool.Pool[Node]
	sync.RWMutex
	percentage float64
	height     int
}

func New(height, capacity int, percentage float64) *SkipList {
	//fix this part to be dynamic
	stack := st.New[st.Stack[*Node]](250) // max number of parallel readings 250

	for range 250 {
		stack.Push(st.New[*Node](height))
	}

	roots := make([]*Node, height)
	prevues := &Node{leaf: true}
	roots[0] = prevues

	for index := 1; index < len(roots); index++ {
		roots[index] = &Node{down: prevues}
		prevues = roots[index]
	}

	return &SkipList{
		roots:      roots,
		Stack:      stack,
		Pool:       pool.New[Node](capacity),
		percentage: percentage,
		height:     height,
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
	s.Lock()
	defer s.Unlock()

	current := s.roots[s.rootIndex]

	stack, err := s.Stack.Pop()
	if err != nil {
		stack = st.New[*Node](s.height)
	}

	for {
		for current.next != nil && current.next.key < key {
			current = current.next
		}

		if current.leaf {
			break
		}

		stack.Push(current)
		current = current.down
	}

	nextNode := current.next

	node := s.Pool.Insert()

	current.next = node
	*node = NewNode(nextNode, nil, value, key, true) // create new leaf node

	for flipCoin(s.percentage) {
		downNode := node
		leftNode, err := stack.Pop()

		if err != nil {
			s.rootIndex++
			leftNode = s.roots[s.rootIndex]
		}

		nextNode = leftNode.next

		node = s.Pool.Insert()
		leftNode.next = node
		*node = NewNode(nextNode, downNode, value, key, false) // create new internal node
	}

	stack.Clear()       //Clear stack
	s.Stack.Push(stack) // return to stack stack
}

func flipCoin(percentage float64) bool {
	return rand.Float64() < percentage
}

func (s *SkipList) Search(key int) (bool, int) {
	s.Lock()
	defer s.Unlock()

	current := s.roots[s.rootIndex]

	for {
		for current.next != nil && current.next.key <= key {
			current = current.next
		}

		if current.leaf {
			break
		}

		current = current.down
	}

	return current.key == key, current.value
}

func (s *SkipList) Read() {
	for startNode := s.roots[0].next; startNode != nil; startNode = startNode.next {
		fmt.Println(startNode)
	}
}

func (s *SkipList) Clear() {
	for index := range s.roots {
		s.roots[index].next = nil
	}

	s.Pool.Clear()
}
