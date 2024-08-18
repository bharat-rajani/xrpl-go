package utils

func PointerOf[T any](t T) *T {
	return &t
}
