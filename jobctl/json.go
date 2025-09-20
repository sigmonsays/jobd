package main

import (
	"encoding/json"
	"fmt"
)

func PrintJson(res any) error {
	// fmt.Printf("%T\n", res)
	b, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		return fmt.Errorf("PrintJson: Error: %s", err)
	}
	fmt.Printf("%s\n", b)
	return nil
}
