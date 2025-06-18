package util

import (
	"errors"
	"fmt"
)

func DebugError(err error) {
	fmt.Printf("DebugError %T %s\n", err, err)
	for i := 0; i < 10; i++ {
		fmt.Printf("unwrap%d: %T %s\n", i, err, err)
		if err == nil {
			break
		}
		err = errors.Unwrap(err)
	}
}
