package promiselike

import (
	"fmt"
	"time"
)

type IntValue struct {
	Value int
}

func main() {
	values := <-All(
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

		func() <-chan bool {
			two := make(chan bool)

			go func() {
				time.Sleep(2 * time.Second)
				fmt.Println("Write second")
				two <- true
			}()

			return two
		}(),
		func() <-chan *string {
			three := make(chan *string)

			go func() {
				time.Sleep(1 * time.Second)
				fmt.Println("Write third")
				val := "my string"
				three <- &val
			}()

			return three
		}(),
	)

	var val1 IntValue = *values[0].(*IntValue)
	var val2 bool = values[1].(bool)
	var val3 string = *values[2].(*string)

	fmt.Println("Done", val1, val2, val3)
}
