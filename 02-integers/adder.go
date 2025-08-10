package integers

import (
	"fmt"
)

// Add takes two integers and returns the sum of them
func Add(a, b int) (sum int) {
	// Named return value!

	// In the future, we'll learn about property-based testing
	// to stop sham TDD (eg. `return 4`)
	return a + b
}

func main() {
	fmt.Println(Add(2, 2))
}
