package jspromises

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
