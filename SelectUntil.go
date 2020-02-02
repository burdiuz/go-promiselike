package promiselike

import (
	"reflect"
)

/*SelectUntil checks for values from all channels at once until all of them are closed or value passed into <-done channel that signals to stop recieving.
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
