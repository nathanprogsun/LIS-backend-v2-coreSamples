package util

import (
	"testing"
)

func TestSwap(t *testing.T) {
	type StructA struct {
		IdA         int    `json:"id"`
		NameA       string `json:"name"`
		Description string `json:"descriptionA"`
		id          int
		name        string
	}

	type StructB struct {
		IdB         int    `json:"id"`
		NameB       string `json:"name"`
		Description string `json:"descriptionB"`
		id          int
	}

	a := StructA{
		IdA:         1,
		NameA:       "name",
		Description: "This is A",
		id:          1,
		name:        "name",
	}

	b := &StructB{}
	err := Swap(a, b)
	if err != nil {
		t.Fatalf("error when swap: %v", err)
	}
	if a.IdA != b.IdB {
		t.Fatalf("expect a.IdA and b.IdB is %d, getting %d and %d", a.IdA, a.IdA, b.IdB)
	}

	if a.NameA != b.NameB {
		t.Fatalf("expect a.NameA and b.NameB is %s, getting %s and %s", a.NameA, a.NameA, b.NameB)
	}

	if a.Description == b.Description {
		t.Fatalf("description with tags that don't match should not be the same")
	}

	if a.id == b.id {
		t.Fatalf("private field should not be swapped")
	}
}

func TestStringCompare(t *testing.T) {
	if !StringEqualIgnoreCase("_aBc!", "_Abc!") {
		t.Fatalf("string should be the same")
	}
}

func TestMin(t *testing.T) {
	var a1 int32 = 1
	var a2 int32 = 2
	if Min(a1, a2) != 1 {
		t.Fatalf("a1 should be smaller than a2")
	}

	var b1 float32 = 1.5
	var b2 float32 = 2.0
	if Min(b1, b2) != b1 {
		t.Fatalf("b1 should be smaller than b2")
	}
}

func TestElementsUniqueInt32(t *testing.T) {
	if ElementsUniqueInt32([]int32{2, 1, 1, 3}) {
		t.Fatalf("array should not be unique")
	}
	if !ElementsUniqueInt32([]int32{2, 1, 3}) {
		t.Fatalf("array should be unique")
	}
}
