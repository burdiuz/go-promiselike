package promiselike

import (
	"reflect"
)

/*Race works just like Promise.race from JavaScript and returns a value first received from input channels, then closes output channel.
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
