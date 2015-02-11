package main

import (
	"github.com/yookoala/gofan"
	"log"
)

// generate a stream of integer for testing
func generate(a int) (out chan int) {
	out = make(chan int)

	go func() {
		defer close(out)
		for i := 1; i <= a; i++ {
			out <- i
		}
	}()

	return
}

// a simple fan out / fan in funciton for test
func pipe(b int, in chan int) (out chan int) {
	out = make(chan int)
	fg := gofan.NewGroup(b)

	for i := range in {

		// clone variable for fan out
		j := i

		// store the function and run until ready
		fg.Run(func() {
			// returns only the even numbers
			if j%2 == 0 {
				out <- j
			}
		})
	}

	// start a goroutine to wait
	fg.OnFinish(func() {
		close(out)
	})

	return
}

func main() {

	// distribute works to 5 workers
	s := 0
	for n := range pipe(5, generate(100)) {
		s += n
	}

	if s == 2550 {
		log.Printf("The sum is 2550, as expected")
	} else {
		log.Panic("The sum is not 2550 as expected.")
	}
}
