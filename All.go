package promiselike

import (
	"reflect"
)

/*All works just like Promise.all() in Javascript and waits for all channels to return a value.
After returning a list of values output channel is closed.
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
