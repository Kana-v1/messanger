package logs

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
)

func ErrorLog(fileName string, message string, comeInError error) error {
	if comeInError == nil {
		comeInError = errors.New("")
	}
	if fileName == "" {
		fileName = "error.log"
	}
	f, err := os.OpenFile("messanger/"+fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	errLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	_, err = errLog.Writer().Write([]byte(fmt.Sprintf("%s\n%v", message, comeInError.Error())))

	return errors.Wrap(err, fmt.Sprintf("Can not log message: '%s'", message))
}

func InfoLog(fileName string, message string, err error) error {
	if fileName == "" {
		fileName = "info.log"
	}

	return ErrorLog(fileName, message, err)
}

func FatalLog(fileName string, message string, err error) {
	if fileName == "" {
		fileName = "fatal.log"
	}
	ErrorLog(fileName, message, err)
	os.Exit(1)
}
