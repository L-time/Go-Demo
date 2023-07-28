package foo

import "fmt"

func F() string {
	fmt.Println("I am a function called F.")
	return "F"
}

func QuickPow(a, b int64, mod int64) int64 {
	var result int64 = 1
	for b > 0 {
		if b&1 == int64(1) {
			result = result * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return result
}
