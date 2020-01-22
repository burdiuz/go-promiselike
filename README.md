### jspromise.Any
Any type used as general type for values passed via channels.

### jspromise.All(input ...<-chan *Any) <-chan []*Any
All works just like Promise.all() in Javascript and waits for all channels to return a value.
After returning a list of values output channel is closed. All values passed through must be of type *Any.
```go
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
```

### jspromise.Select(input ...<-chan *Any) <-chan *Any
Select checks for values from all channels at once until all of them are closed.
```go
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
```

### jspromise.SelectUntil(done <-chan int, input ...<-chan *Any) <-chan *Any
SelectUntil checks for values from all channels at once until all of them are closed or value passed into <-done channel that signals to stop recieving.
```go
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
```

### jspromise.Race(input ...<-chan *Any) <-chan *Any
Race works just like Promise.race from JavaScript and returns a value returned first, then closes output channel.
```go
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
```
