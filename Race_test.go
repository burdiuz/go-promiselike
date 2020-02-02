package promiselike

import (
	"fmt"
	"time"
)

type IntValue struct {
	Value int
}

func main() {
	value := <-Race(
		func() <-chan *IntValue {
			one := make(chan *IntValue)

			go func() {
				time.Sleep(3 * time.Second)
				fmt.Println("Write first")
				val := IntValue{124}
				one <- &val
			}()
			return one
		}(),
		func() <-chan *IntValue {
			two := make(chan *IntValue)

			go func() {
				time.Sleep(2 * time.Second)
				fmt.Println("Write second")
				val := IntValue{10000}
				two <- &val
			}()

			return two
		}(),
		func() <-chan *IntValue {
			three := make(chan *IntValue)

			go func() {
				time.Sleep(time.Second)
				fmt.Println("Write third")
				val := IntValue{777}
				three <- &val
			}()

			return three
		}(),
	)

	fmt.Println("Done", *value.(*IntValue))
}
