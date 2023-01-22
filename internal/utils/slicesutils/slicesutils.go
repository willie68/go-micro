// Package slicesutils some convenient functions on slices
package slicesutils

/*
Contains checking if the e string is present in the slice s
*/
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

/*
Remove removes the entry with the index i from the slice
*/
func Remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

// RemoveString removes the e entry from the s slice, if e is not present in the slice, nothing will happen
func RemoveString(s []string, e string) []string {
	index := Find(s, e)
	if index >= 0 {
		return Remove(s, index)
	}
	return s
}

/*
Find finding the index of the e string in the s slice
*/
func Find(s []string, e string) int {
	for i, n := range s {
		if e == n {
			return i
		}
	}
	return -1
}
