package set3

// this type is for benchmark comparison only
type nativeSetType[T comparable] map[T]struct{}

//go:noinline
func emptyNativeWithCapacity[T comparable](size uint32) *nativeSetType[T] {
	result := make(nativeSetType[T], size)
	return &result
}

//go:noinline
func (thisSet *nativeSetType[T]) add(val T) {
	(*thisSet)[val] = struct{}{}
}

//go:noinline
func (thisSet *nativeSetType[T]) contains(val T) bool {
	_, b := (*thisSet)[val]
	return b
}
