package jspromise

//Any type used as general type for values passed via channels.
type Any interface{}

/*All works just like Promise.all() in Javascript and waits for all channels to return a value.
After returning a list of values output channel is closed. All values passed through must be of type *Any.

func main() {
	values := <-All(
		func() <-chan *Any {
			one := make(chan *Any)

			go func() {
				time.Sleep(3 * time.Second)
				fmt.Println("Write first")
				val := Any(124)
				one <- &val
			}()
			return one
		}(),

		func() <-chan *Any {
			two := make(chan *Any)

			go func() {
				time.Sleep(2 * time.Second)
				fmt.Println("Write second")
				val := Any(true)
				two <- &val
			}()

			return two
		}(),
		func() <-chan *Any {
			three := make(chan *Any)

			go func() {
				time.Sleep(1 * time.Second)
				fmt.Println("Write third")
				val := Any("my string")
				three <- &val
			}()

			return three
		}(),
	)

	var val1 int = (*values[0]).(int)
	var val2 bool = (*values[1]).(bool)
	var val3 string = (*values[2]).(string)

	fmt.Println("Done", val1, val2, val3)
}
*/
func All(inputs ...<-chan *Any) <-chan []*Any {
	out := make(chan []*Any)

	go func() {
		length := len(inputs)
		var values []*Any = make([]*Any, length)
		var filled []bool = make([]bool, length)

		for {
			hasEmpty := false

			for index := 0; index < length; index++ {
				if filled[index] {
					continue
				}

				select {
				case value := <-inputs[index]:
					values[index] = value
					filled[index] = true
				default:
					hasEmpty = true
				}
			}

			if !hasEmpty {
				break
			}
		}

		out <- values
		close(out)
	}()

	return out
}

/*Select selects values from all channels at once until all of them close.

func main() {
	chain := Select(
		func() <-chan *Any {
			one := make(chan *Any)

			go func() {
				for i := 0; i < 30; i++ {
					time.Sleep(100 * time.Millisecond)
					val := Any(i << 2)
					one <- &val
				}

				close(one)
			}()
			return one
		}(),
		func() <-chan *Any {
			two := make(chan *Any)

			go func() {
				val := Any(true)

				for i := 0; i < 12; i++ {
					time.Sleep(250 * time.Millisecond)
					val = Any(!val.(bool))
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

		fmt.Printf("Value: %v\n", *value)
	}

	fmt.Println("Done.")
}
*/
func Select(inputs ...<-chan *Any) <-chan *Any {
	out := make(chan *Any)

	go func() {
		length := len(inputs)
		closed := make([]bool, length)

		for {
			hasOpen := false

			for index := 0; index < length; index++ {
				if closed[index] {
					continue
				}

				select {
				case value, ok := <-inputs[index]:
					if ok {
						out <- value
						hasOpen = true
					} else {
						closed[index] = true
					}
				default:
					// This default is important because it allows to skip channel if it is not ready
					hasOpen = true
				}
			}

			if !hasOpen {
				close(out)
				return
			}
		}
	}()

	return out
}

/*SelectUntil selects values from all channels at once until all of them close or value passed into <-done channel that signals to stop recieving.

func main() {
	done := make(chan int)

	chain := SelectUntil(
		done,
		func() <-chan *Any {
			one := make(chan *Any)

			go func() {
				for i := 0; i < 30; i++ {
					time.Sleep(100 * time.Millisecond)
					val := Any(i << 2)
					one <- &val
				}

				close(one)
			}()
			return one
		}(),
		func() <-chan *Any {
			two := make(chan *Any)

			go func() {
				val := Any(true)

				for i := 0; i < 12; i++ {
					time.Sleep(250 * time.Millisecond)
					val = Any(!val.(bool))
					two <- &val
				}

				close(two)
			}()

			return two
		}(),
	)

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Force stream stop.")
		done <- 1
	}()

	for {
		value, ok := <-chain

		if !ok {
			break
		}

		fmt.Printf("Value: %v\n", *value)
	}

	fmt.Println("Done.")
}
*/
func SelectUntil(done <-chan int, inputs ...<-chan *Any) <-chan *Any {
	out := make(chan *Any)

	go func() {
		length := len(inputs)
		closed := make([]bool, length)

		for {
			hasOpen := false
			for index := 0; index < length; index++ {
				if closed[index] {
					continue
				}

				select {
				case value, ok := <-inputs[index]:
					if ok {
						out <- value
						hasOpen = true
					} else {
						closed[index] = true
					}
				case <-done:
					close(out)
					return
				default:
					hasOpen = true
				}
			}

			if !hasOpen {
				close(out)
				break
			}
		}
	}()

	return out
}

/*Race works just like Promise.race from JavaScript and returns a value returned first, then closes output channel.

func main() {
	value := *<-Race(
		func() <-chan *Any {
			one := make(chan *Any)

			go func() {
				time.Sleep(3 * time.Second)
				fmt.Println("Write first")
				val := Any(124)
				one <- &val
			}()
			return one
		}(),
		func() <-chan *Any {
			two := make(chan *Any)

			go func() {
				time.Sleep(2 * time.Second)
				fmt.Println("Write second")
				val := Any(true)
				two <- &val
			}()

			return two
		}(),
		func() <-chan *Any {
			three := make(chan *Any)

			go func() {
				time.Sleep(1 * time.Second)
				fmt.Println("Write third")
				val := Any("my string")
				three <- &val
			}()

			return three
		}(),
	)

	fmt.Println("Done", value)
}
*/
func Race(inputs ...<-chan *Any) <-chan *Any {
	out := make(chan *Any)

	go func() {
		length := len(inputs)
		for {
			for index := 0; index < length; index++ {
				select {
				case value := <-inputs[index]:
					out <- value
					close(out)
					return
				default:
				}
			}
		}
	}()

	return out
}
