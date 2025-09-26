package set

import (
	"reflect"
	"sort"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		items    []int
		expected []int
	}{
		{"empty set", []int{}, []int{}},
		{"single item", []int{1}, []int{1}},
		{"multiple items", []int{1, 2, 3}, []int{1, 2, 3}},
		{"duplicate items", []int{1, 2, 2, 3, 1}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		s := New(tt.items...)

		if s.Len() != len(tt.expected) {
			t.Errorf("%s: expected length %d, got %d", tt.name, len(tt.expected), s.Len())
		}

		for _, item := range tt.expected {
			if !s.Contains(item) {
				t.Errorf("%s: expected set to contain %v", tt.name, item)
			}
		}
	}
}

func TestAdd(t *testing.T) {
	s := New[int]()

	// Test adding to empty set
	s.Add(1)
	if !s.Contains(1) {
		t.Error("expected set to contain 1 after adding")
	}
	if s.Len() != 1 {
		t.Errorf("expected length 1, got %d", s.Len())
	}

	// Test adding duplicate
	s.Add(1)
	if s.Len() != 1 {
		t.Errorf("expected length to remain 1 after adding duplicate, got %d", s.Len())
	}

	// Test adding different item
	s.Add(2)
	if !s.Contains(2) {
		t.Error("expected set to contain 2 after adding")
	}
	if s.Len() != 2 {
		t.Errorf("expected length 2, got %d", s.Len())
	}
}

func TestRemove(t *testing.T) {
	s := New(1, 2, 3)

	// Test removing existing item
	s.Remove(2)
	if s.Contains(2) {
		t.Error("expected set to not contain 2 after removal")
	}
	if s.Len() != 2 {
		t.Errorf("expected length 2, got %d", s.Len())
	}

	// Test removing non-existing item
	s.Remove(4)
	if s.Len() != 2 {
		t.Errorf("expected length to remain 2 after removing non-existing item, got %d", s.Len())
	}

	// Test removing all items
	s.Remove(1)
	s.Remove(3)
	if s.Len() != 0 {
		t.Errorf("expected empty set, got length %d", s.Len())
	}
}

func TestContains(t *testing.T) {
	s := New(1, 2, 3)

	tests := []struct {
		item     int
		expected bool
	}{
		{1, true},
		{2, true},
		{3, true},
		{4, false},
		{0, false},
	}

	for _, tt := range tests {
		result := s.Contains(tt.item)
		if result != tt.expected {
			t.Errorf("Contains(%d): expected %v, got %v", tt.item, tt.expected, result)
		}
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name     string
		items    []int
		expected int
	}{
		{"empty set", []int{}, 0},
		{"single item", []int{1}, 1},
		{"multiple items", []int{1, 2, 3}, 3},
		{"duplicates", []int{1, 1, 2, 2, 3}, 3},
	}

	for _, tt := range tests {
		s := New(tt.items...)
		if s.Len() != tt.expected {
			t.Errorf("%s: expected length %d, got %d", tt.name, tt.expected, s.Len())
		}
	}
}

func TestItems(t *testing.T) {
	tests := []struct {
		name     string
		items    []int
		expected []int
	}{
		{"empty set", []int{}, []int{}},
		{"single item", []int{1}, []int{1}},
		{"multiple items", []int{1, 2, 3}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		s := New(tt.items...)
		result := s.Items()

		// Sort both slices for comparison since map iteration order is not guaranteed
		sort.Ints(result)
		sort.Ints(tt.expected)

		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, result)
		}
	}
}

func TestUnion(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected []int
	}{
		{"empty sets", []int{}, []int{}, []int{}},
		{"one empty set", []int{1, 2}, []int{}, []int{1, 2}},
		{"no overlap", []int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4}},
		{"partial overlap", []int{1, 2, 3}, []int{3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"complete overlap", []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		s1 := New(tt.set1...)
		s2 := New(tt.set2...)
		result := s1.Union(s2)

		if result.Len() != len(tt.expected) {
			t.Errorf("%s: expected length %d, got %d", tt.name, len(tt.expected), result.Len())
		}

		for _, item := range tt.expected {
			if !result.Contains(item) {
				t.Errorf("%s: expected union to contain %v", tt.name, item)
			}
		}
	}
}

func TestIntersection(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected []int
	}{
		{"empty sets", []int{}, []int{}, []int{}},
		{"one empty set", []int{1, 2}, []int{}, []int{}},
		{"no overlap", []int{1, 2}, []int{3, 4}, []int{}},
		{"partial overlap", []int{1, 2, 3}, []int{3, 4, 5}, []int{3}},
		{"complete overlap", []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		s1 := New(tt.set1...)
		s2 := New(tt.set2...)
		result := s1.Intersection(s2)

		if result.Len() != len(tt.expected) {
			t.Errorf("%s: expected length %d, got %d", tt.name, len(tt.expected), result.Len())
		}

		for _, item := range tt.expected {
			if !result.Contains(item) {
				t.Errorf("%s: expected intersection to contain %v", tt.name, item)
			}
		}
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected []int
	}{
		{"empty sets", []int{}, []int{}, []int{}},
		{"first set empty", []int{}, []int{1, 2}, []int{}},
		{"second set empty", []int{1, 2}, []int{}, []int{1, 2}},
		{"no overlap", []int{1, 2}, []int{3, 4}, []int{1, 2}},
		{"partial overlap", []int{1, 2, 3}, []int{3, 4, 5}, []int{1, 2}},
		{"complete overlap", []int{1, 2, 3}, []int{1, 2, 3}, []int{}},
	}

	for _, tt := range tests {
		s1 := New(tt.set1...)
		s2 := New(tt.set2...)
		result := s1.Difference(s2)

		if result.Len() != len(tt.expected) {
			t.Errorf("%s: expected length %d, got %d", tt.name, len(tt.expected), result.Len())
		}

		for _, item := range tt.expected {
			if !result.Contains(item) {
				t.Errorf("%s: expected difference to contain %v", tt.name, item)
			}
		}
	}
}

func TestSymmetricDifference(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected []int
	}{
		{"empty sets", []int{}, []int{}, []int{}},
		{"first set empty", []int{}, []int{1, 2}, []int{1, 2}},
		{"second set empty", []int{1, 2}, []int{}, []int{1, 2}},
		{"no overlap", []int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4}},
		{"partial overlap", []int{1, 2, 3}, []int{3, 4, 5}, []int{1, 2, 4, 5}},
		{"complete overlap", []int{1, 2, 3}, []int{1, 2, 3}, []int{}},
	}

	for _, tt := range tests {
		s1 := New(tt.set1...)
		s2 := New(tt.set2...)
		result := s1.SymmetricDifference(s2)

		if result.Len() != len(tt.expected) {
			t.Errorf("%s: expected length %d, got %d", tt.name, len(tt.expected), result.Len())
		}

		for _, item := range tt.expected {
			if !result.Contains(item) {
				t.Errorf("%s: expected symmetric difference to contain %v", tt.name, item)
			}
		}
	}
}

// Test with string type to verify generic functionality
func TestSetWithStrings(t *testing.T) {
	s := New("apple", "banana", "cherry")

	if s.Len() != 3 {
		t.Errorf("expected length 3, got %d", s.Len())
	}

	if !s.Contains("apple") {
		t.Error("expected set to contain 'apple'")
	}

	s.Add("date")
	if !s.Contains("date") {
		t.Error("expected set to contain 'date' after adding")
	}

	s.Remove("banana")
	if s.Contains("banana") {
		t.Error("expected set to not contain 'banana' after removal")
	}

	items := s.Items()
	sort.Strings(items)
	expected := []string{"apple", "cherry", "date"}
	sort.Strings(expected)

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("expected %v, got %v", expected, items)
	}
}

// Test set operations with strings
func TestSetOperationsWithStrings(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		set1      []string
		set2      []string
		expected  []string
	}{
		{"union no overlap", "union", []string{"a", "b"}, []string{"c", "d"}, []string{"a", "b", "c", "d"}},
		{"union with overlap", "union", []string{"a", "b", "c"}, []string{"c", "d", "e"}, []string{"a", "b", "c", "d", "e"}},
		{"intersection with overlap", "intersection", []string{"a", "b", "c"}, []string{"c", "d", "e"}, []string{"c"}},
		{"intersection no overlap", "intersection", []string{"a", "b"}, []string{"c", "d"}, []string{}},
		{"difference with overlap", "difference", []string{"a", "b", "c"}, []string{"c", "d", "e"}, []string{"a", "b"}},
		{"difference no overlap", "difference", []string{"a", "b"}, []string{"c", "d"}, []string{"a", "b"}},
		{"symmetric difference with overlap", "symmetric_difference", []string{"a", "b", "c"}, []string{"c", "d", "e"}, []string{"a", "b", "d", "e"}},
		{"symmetric difference no overlap", "symmetric_difference", []string{"a", "b"}, []string{"c", "d"}, []string{"a", "b", "c", "d"}},
	}

	for _, tt := range tests {
		s1 := New(tt.set1...)
		s2 := New(tt.set2...)

		var result Set[string]
		switch tt.operation {
		case "union":
			result = s1.Union(s2)
		case "intersection":
			result = s1.Intersection(s2)
		case "difference":
			result = s1.Difference(s2)
		case "symmetric_difference":
			result = s1.SymmetricDifference(s2)
		}

		if result.Len() != len(tt.expected) {
			t.Errorf("%s: expected length %d, got %d", tt.name, len(tt.expected), result.Len())
		}

		for _, item := range tt.expected {
			if !result.Contains(item) {
				t.Errorf("%s: expected result to contain %s", tt.name, item)
			}
		}
	}
}

// Test edge cases
func TestEmptySetOperations(t *testing.T) {
	empty1 := New[int]()
	empty2 := New[int]()

	tests := []struct {
		name      string
		operation func() Set[int]
	}{
		{"union", func() Set[int] { return empty1.Union(empty2) }},
		{"intersection", func() Set[int] { return empty1.Intersection(empty2) }},
		{"difference", func() Set[int] { return empty1.Difference(empty2) }},
		{"symmetric difference", func() Set[int] { return empty1.SymmetricDifference(empty2) }},
	}

	for _, tt := range tests {
		result := tt.operation()
		if result.Len() != 0 {
			t.Errorf("%s of empty sets should be empty, got length %d", tt.name, result.Len())
		}
	}
}

func TestOperationsPreserveOriginalSets(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := New(3, 4, 5)

	originalS1Len := s1.Len()
	originalS2Len := s2.Len()

	// Perform operations
	s1.Union(s2)
	s1.Intersection(s2)
	s1.Difference(s2)
	s1.SymmetricDifference(s2)

	// Original sets should be unchanged
	if s1.Len() != originalS1Len {
		t.Error("s1 was modified by set operations")
	}
	if s2.Len() != originalS2Len {
		t.Error("s2 was modified by set operations")
	}
}

// Benchmark tests
func BenchmarkSetAdd(b *testing.B) {
	s := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkSetContains(b *testing.B) {
	s := New[int]()
	for i := range 1000 {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Contains(i % 1000)
	}
}

func BenchmarkSetUnion(b *testing.B) {
	s1 := New[int]()
	s2 := New[int]()
	for i := range 500 {
		s1.Add(i)
		s2.Add(i + 250) // 50% overlap
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s1.Union(s2)
	}
}
