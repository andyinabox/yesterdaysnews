package util

// StringSliceRemove will return a new slice with the given element removed
func StringSliceSplice(slice []string, i int) (string, []string) {
	str := slice[i]
	newSlice := make([]string, i)
	copy(newSlice, slice[:i])
	return str, append(newSlice, slice[i+1:]...)
}
