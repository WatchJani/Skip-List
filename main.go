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
	sl := New(32, 25000, 0.5)
	s := []int{136, 135, 119, 42, 134, 41, 145, 110, 44, 4}

	for {
		for _, value := range s[:1] {
			fmt.Println("=>", value)
			id := sl.Insert(value, 23)

			fmt.Println(id)

			for index := id; index > -1; index-- {
				fmt.Println(sl.roots[id].next.key)
			}

			sl.Read()
			fmt.Println("new index", sl.rootIndex)
		}
	}
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
	stack := st.New[st.Stack[*Node]](250)

	for range 250 {
		go stack.Push(st.New[*Node](height))
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

func (s *SkipList) Insert(key, value int) int {
	current, startIndex := s.roots[s.rootIndex], s.rootIndex
	stack, err := s.Stack.Pop()

	if err != nil {
		stack = st.New[*Node](s.height)
	}

	for {
		for current.next != nil && current.next.key < key {
			current = current.next
		}

		if current.leaf || startIndex == 0 {
			break
		}

		stack.Push(current)

		current = current.down

		startIndex--
	}
	next := current.next

	node := s.Pool.Insert()
	current.next = node
	*node = NewNode(next, nil, value, key, true) // create new leaf node

	var counter int
	for flipCoin(s.percentage) && s.height > startIndex {
		downNode := node
		leftNode, err := stack.Pop()

		if err != nil {
			s.rootIndex++
			leftNode = s.roots[s.rootIndex]
		}

		counter++

		next = leftNode.next
		node = s.Pool.Insert()
		leftNode.next = node

		*node = NewNode(next, downNode, value, key, false) // create new internal node
	}

	stack.Clear()       //Clear stack
	s.Stack.Push(stack) // return to stack stack

	return counter
}

func flipCoin(percentage float64) bool {
	return rand.Float64() < percentage
}

func (s *SkipList) Read() {
	for startNode := s.roots[0].next; startNode != nil; startNode = startNode.next {
		fmt.Println(startNode)
	}
}
