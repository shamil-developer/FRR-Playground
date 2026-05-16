package utils

import (
	"fmt"

	"github.com/TylerBrock/colorjson"
)

func PrintAny(data any) error {

	formatter := colorjson.NewFormatter()

	formatter.Indent = 2

	result, err := formatter.Marshal(data)

	if err != nil {
		return err
	}

	fmt.Println(string(result))

	return nil
}
