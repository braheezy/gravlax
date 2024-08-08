package main

import (
	"errors"
	"slices"
)

type LinkedList struct {
	nodes []Node
}

type Node struct {
	Data string
	Next *Node
	Prev *Node
}

func (list *LinkedList) Insert(newValue string) {
	newNode := Node{Data: newValue}
	list.nodes = append(list.nodes, newNode)

	if list.Size() > 1 {
		list.nodes[list.Size()-1].Prev = &list.nodes[list.Size()-2]
		list.nodes[list.Size()-2].Next = &newNode
	}
}

func (list *LinkedList) Size() int {
	return len(list.nodes)
}
func (list *LinkedList) Find(search string) *Node {
	for _, node := range list.nodes {
		if node.Data == search {
			return &node
		}
	}
	return nil
}
func (list *LinkedList) Delete(target string) error {
	found := -1
	for i, node := range list.nodes {
		if node.Data == target {
			list.nodes = slices.Delete(list.nodes, i, i+1)
			found = i
			break
		}
	}
	if found == -1 {
		return errors.New("target does not exist")
	} else {

		if found == 0 {
			list.nodes[0].Prev = nil
		} else {
			list.nodes[found-1].Next = &list.nodes[found+1]
			list.nodes[found+1].Prev = &list.nodes[found-1]
		}
		return nil
	}
}
func (list *LinkedList) Head() *Node {
	return &list.nodes[0]
}
func (list *LinkedList) Tail() *Node {
	return &list.nodes[list.Size()-1]
}
