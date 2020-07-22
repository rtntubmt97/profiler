package utils

import (
	"fmt"
	"time"
)

func WrapError(pkgName, msg string, a interface{}) error {
	if a == nil {
		return nil
	}
	return fmt.Errorf("%s: %s\n%w", pkgName, msg, a)
}

func PrintError(err error) {
	timeStr := time.Now().Format("2006/01/02 - 15:04:05")
	fmt.Printf("[%s] %s\n", timeStr, err.Error())
}
