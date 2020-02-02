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
