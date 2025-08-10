package integers

import (
	"fmt"
	"testing"
)

func TestAdder(t *testing.T) {
	actual := Add(2, 2)
	expected := 4

	if expected != actual {
		t.Errorf("expected %d, actual: %d", expected, actual)
	}
}

// Functions that start with Example are useful for
// examples outside code in documentation if you really
// want to go the extra mile.
func ExampleAdd() {
	sum := Add(5, 3)
	sum2 := Add(5, 5)
	fmt.Println(sum)
	fmt.Println(sum2)
	// Output: 8
	// 10
}
