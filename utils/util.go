package utils

func Use(vals ...interface{}) {}

func Perform[A Any](cont func() A) chan A {
	ch := make(chan A)
	go func() {
		ch <- cont()
	}()
	return ch
}

func Try(cont func() interface{}) chan interface{} {
	ch := make(chan interface{})
	go func() {
		ch <- cont()
	}()
	return ch
}
