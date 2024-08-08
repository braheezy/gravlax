package main

import (
	"testing"
)

// TestDoublyLinkedList tests the operations of the doubly-linked list.
func TestDoublyLinkedList(t *testing.T) {
	list := LinkedList{}

	// Test Insert
	list.Insert("Hello")
	list.Insert("World")
	if list.Size() != 2 {
		t.Errorf("Expected list size of 2, got %d", list.Size())
	}

	// Test Find
	if node := list.Find("Hello"); node == nil || node.Data != "Hello" {
		t.Errorf("Expected to find 'Hello'")
	}
	if node := list.Find("NotExist"); node != nil {
		t.Errorf("Expected not to find 'NotExist'")
	}

	// Test Delete
	if err := list.Delete("Hello"); err != nil {
		t.Errorf("Failed to delete 'Hello': %s", err)
	}
	if list.Size() != 1 {
		t.Errorf("Expected list size of 1 after deletion, got %d", list.Size())
	}
	if node := list.Find("Hello"); node != nil {
		t.Errorf("Expected 'Hello' to be deleted")
	}

	// Insert more to test order and linkage
	list.Insert("First")
	list.Insert("Second")
	if list.Size() != 3 {
		t.Errorf("Expected list size of 3, got %d", list.Size())
	}
	if list.Head().Data != "World" || list.Tail().Data != "Second" {
		t.Errorf("Incorrect head or tail after multiple inserts")
	}

	// Test deleting non-existent element
	if err := list.Delete("NotExist"); err == nil {
		t.Errorf("Expected error when deleting non-existent element")
	}
}
