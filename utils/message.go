package utils

import (
	"fmt"
)

func WrapWithSeparators(msg string) string {
	return fmt.Sprintf("=================\n%s\n=================", msg)
}
