package promiselike

import (
	"fmt"
	"time"
)

type IntValue struct {
	Value int
}

func main() {
	chain := Select(
		func() <-chan *IntValue {
			one := make(chan *IntValue)

			go func() {
				for i := 0; i < 30; i++ {
					time.Sleep(100 * time.Millisecond)
					val := IntValue{i << 2}
					one <- &val
				}

				close(one)
			}()
			return one
		}(),
		func() <-chan *IntValue {
			two := make(chan *IntValue)

			go func() {
				for i := 0; i < 12; i++ {
					time.Sleep(250 * time.Millisecond)
					val := IntValue{i * 10000}
					two <- &val
				}

				close(two)
			}()

			return two
		}(),
	)

	for {
		value, ok := <-chain

		if !ok {
			break
		}

		fmt.Printf("Value: %v\n", *value.(*IntValue))
	}

	fmt.Println("Done.")
}
