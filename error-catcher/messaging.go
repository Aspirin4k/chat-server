package error_catcher

import (
	"os"
	"fmt"
	"time"
)

var (
	LocalReader *os.File
	LocalWriter *os.File
)

func Init() {
	var err error
	LocalReader, LocalWriter, err = os.Pipe()
	CheckError(err)
}

func PushMessage(msg string) {
	LocalWriter.WriteString(fmt.Sprintf("\u001b[36m[%s]\u001b[0m: %s\n", time.Now().Format("15:04:05"), msg))
}