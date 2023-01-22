package slicesutils

import (
	"testing"
)

func TestContains(t *testing.T) {
	mySlice := []string{"Willie", "Arthur", "Till"}
	value := Contains(mySlice, "Willie")
	if !value {
		t.Errorf("Willie was not in the slice")
	}
}
func TestRemoveString(t *testing.T) {
	mySlice := []string{"Willie", "Arthur", "Till"}
	value := RemoveString(mySlice, "Willie")
	if Contains(value, "Willie") {
		t.Errorf("Willie was not removed from the slice")
	}
	value = RemoveString(mySlice, "Herman")
	if len(value) != 3 {
		t.Errorf("slice not unchanged")
	}
}

func TestRemove(t *testing.T) {
	mySlice := []string{"Willie", "Arthur", "Till"}
	value := Remove(mySlice, 0)
	if Contains(value, "Willie") {
		t.Errorf("Willie was not removed from the slice")
	}
}

func TestFind(t *testing.T) {
	mySlice := []string{"Willie", "Arthur", "Till"}
	value := Find(mySlice, "Willie")
	if value != 0 {
		t.Errorf("Willie was not found in the slice: index: %d", value)
	}
	value = Find(mySlice, "Arthur")
	if value != 1 {
		t.Errorf("Arthur was not found in the slice: index: %d", value)
	}
	value = Find(mySlice, "Till")
	if value != 2 {
		t.Errorf("Till was not found in the slice: index: %d", value)
	}
	value = Find(mySlice, "till")
	if value >= 0 {
		t.Errorf("till was wrongly found in the slice: index: %d", value)
	}
	value = Find(mySlice, "Herman")
	if value >= 0 {
		t.Errorf("Herman was wrongly found in the slice: index: %d", value)
	}
}
