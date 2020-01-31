package promiselike

import (
  "reflect"
)

func checkForChan(inputs []interface{}) {
  for index, input := range inputs {
    kind := reflect.TypeOf(input).Kind().String()
    if kind != "chan" {
      panic("All parameters must be channels, but received " + kind + " for #" + string(index))
    }
  }
}

/*All works just like Promise.all() in Javascript and waits for all channels to return a value.
After returning a list of values output channel is closed.

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
*/
func All(inputs ...interface{}) <-chan []interface{} {
  checkForChan(inputs)

  out := make(chan []interface{})

  go func() {
    length := len(inputs)
    var values []interface{} = make([]interface{}, length)
    var filled []bool = make([]bool, length)

    for {
      hasEmpty := false

      for index := 0; index < length; index++ {
        if filled[index] {
          continue
        }

        selectCases := []reflect.SelectCase{
          {
            Dir:  reflect.SelectRecv,
            Chan: reflect.ValueOf(inputs[index]),
          },
          {
            Dir: reflect.SelectDefault,
          },
        }

        switch caseIndex, value, ok := reflect.Select(selectCases); caseIndex {

        case 0:
          if ok {
            values[index] = value.Interface()
          }

          filled[index] = true
        case 1:
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

/*Select checks for values from all channels at once until all of them are closed.

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
*/
func Select(inputs ...interface{}) <-chan interface{} {
  checkForChan(inputs)

  out := make(chan interface{})

  go func() {
    length := len(inputs)
    closed := make([]bool, length)

    for {
      hasOpen := false

      for index := 0; index < length; index++ {
        if closed[index] {
          continue
        }

        selectCases := []reflect.SelectCase{
          {
            Dir:  reflect.SelectRecv,
            Chan: reflect.ValueOf(inputs[index]),
          },
          {
            Dir: reflect.SelectDefault,
          },
        }

        switch caseIndex, value, ok := reflect.Select(selectCases); caseIndex {

        case 0:
          if ok {
            out <- value.Interface()
            hasOpen = true
          } else {
            closed[index] = true
          }
        case 1:
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

/*SelectUntil checks for values from all channels at once until all of them are closed or value passed into <-done channel that signals to stop recieving.

type IntValue struct {
  Value int
}

func main() {
  done := make(chan int)

  chain := SelectUntil(
    done,
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

  go func() {
    time.Sleep(time.Second)
    fmt.Println("Force stream stop.")
    done <- 1
  }()

  for {
    value, ok := <-chain

    if !ok {
      break
    }

    fmt.Printf("Value: %v\n", *value.(*IntValue))
  }

  fmt.Println("Done.")
}
*/
func SelectUntil(done <-chan int, inputs ...interface{}) <-chan interface{} {
  checkForChan(inputs)

  doneVal := reflect.ValueOf(done)
  out := make(chan interface{})

  go func() {
    length := len(inputs)
    closed := make([]bool, length)

    for {
      hasOpen := false
      for index := 0; index < length; index++ {
        if closed[index] {
          continue
        }

        selectCases := []reflect.SelectCase{
          {
            Dir:  reflect.SelectRecv,
            Chan: reflect.ValueOf(inputs[index]),
          },
          {
            Dir:  reflect.SelectRecv,
            Chan: doneVal,
          },
          {
            Dir: reflect.SelectDefault,
          },
        }

        switch caseIndex, value, ok := reflect.Select(selectCases); caseIndex {

        case 0:
          if ok {
            out <- value.Interface()
            hasOpen = true
          } else {
            closed[index] = true
          }
        case 1:
          close(out)
          return
        case 2:
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

/*Race works just like Promise.race from JavaScript and returns a value first received from input channels, then closes output channel.

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
*/
func Race(inputs ...interface{}) <-chan interface{} {
  checkForChan(inputs)

  out := make(chan interface{})

  go func() {
    length := len(inputs)
    closed := make([]bool, length)

    for {
      hasOpen := false

      for index := 0; index < length; index++ {
        selectCases := []reflect.SelectCase{
          {
            Dir:  reflect.SelectRecv,
            Chan: reflect.ValueOf(inputs[index]),
          },
          {
            Dir: reflect.SelectDefault,
          },
        }

        switch caseIndex, value, ok := reflect.Select(selectCases); caseIndex {
        case 0:
          if ok {
            out <- value.Interface()
            close(out)
            return
          } else {
            closed[index] = true
          }
        case 1:
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
