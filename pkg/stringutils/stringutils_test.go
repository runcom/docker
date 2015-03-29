package stringutils

import "testing"

func TestRandomString(t *testing.T) {
	str := GenerateRandomString()
	if len(str) != 64 {
		t.Fatalf("Id returned is incorrect: %s", str)
	}
}

func TestRandomStringUniqueness(t *testing.T) {
	repeats := 25
	set := make(map[string]struct{}, repeats)
	for i := 0; i < repeats; i = i + 1 {
		str := GenerateRandomString()
		if len(str) != 64 {
			t.Fatalf("Id returned is incorrect: %s", str)
		}
		if _, ok := set[str]; ok {
			t.Fatalf("Random number is repeated")
		}
		set[str] = struct{}{}
	}
}

func TestTruncate(t *testing.T) {
	str := "teststring"
	newstr := Truncate(str, 4)
	if newstr != "test" {
		t.Fatalf("Expected test, got %s", newstr)
	}
	newstr = Truncate(str, 20)
	if newstr != "teststring" {
		t.Fatalf("Expected teststring, got %s", newstr)
	}
}

func TestInSlice(t *testing.T) {
	slice := []string{"test", "in", "slice"}

	test := InSlice(slice, "test")
	if !test {
		t.Fatalf("Expected string test to be in slice")
	}
	test = InSlice(slice, "SLICE")
	if !test {
		t.Fatalf("Expected string SLICE to be in slice")
	}
	test = InSlice(slice, "notinslice")
	if test {
		t.Fatalf("Expected string notinslice not to be in slice")
	}
}
