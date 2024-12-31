package conversions

type CoreConverter[T any] interface {
	ToCore() T
	FromCore(T)
}
