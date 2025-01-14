package util

// StringSliceRemove will return a slice with the given element removed
func StringSliceSplice(slice []string, i int) (string, []string) {
	return slice[i], append(slice[:i], slice[i+1:]...)
}
