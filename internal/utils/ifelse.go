package utils

func IfElse[T any](cond bool, a, b T) T {
	if cond {
		return a
	}

	return b
}
