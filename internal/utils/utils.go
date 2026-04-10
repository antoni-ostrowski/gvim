package utils

import (
	"fmt"
	"os"
)

func Debuglog(format string, args ...any) {
	f, err := os.OpenFile("/tmp/gvim.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf(format, args...) + "\n")
}
