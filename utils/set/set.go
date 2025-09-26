package set

// Set is a generic set implementation
type Set[T comparable] map[T]struct{}

// New creates a new set
func New[T comparable](items ...T) Set[T] {
	s := make(Set[T])
	for _, item := range items {
		s.Add(item)
	}
	return s
}

// Add adds an item to the set
func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

// Remove removes an item from the set
func (s Set[T]) Remove(item T) {
	delete(s, item)
}

// Contains checks if an item is in the set
func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

// Len returns the number of items in the set
func (s Set[T]) Len() int {
	return len(s)
}

// Items returns a slice of all items in the set
func (s Set[T]) Items() []T {
	items := make([]T, 0, len(s))
	for item := range s {
		items = append(items, item)
	}
	return items
}

// Union returns a new set containing all items from both sets
func (s Set[T]) Union(other Set[T]) Set[T] {
	union := New(s.Items()...)
	for item := range other {
		union.Add(item)
	}
	return union
}

// Intersection returns a new set containing all items common to both sets
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	intersection := New[T]()
	for item := range s {
		if other.Contains(item) {
			intersection.Add(item)
		}
	}
	return intersection
}

// Difference returns a new set containing all items in the first set that are not in the second set
func (s Set[T]) Difference(other Set[T]) Set[T] {
	difference := New[T]()
	for item := range s {
		if !other.Contains(item) {
			difference.Add(item)
		}
	}
	return difference
}

// SymmetricDifference returns a new set containing all items in either set, but not both
func (s Set[T]) SymmetricDifference(other Set[T]) Set[T] {
	return s.Difference(other).Union(other.Difference(s))
}
