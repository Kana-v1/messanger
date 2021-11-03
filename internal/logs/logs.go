package logs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

func ErrorLog(fileName string, message string, comeInError error) error {
	_, projectPath, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(projectPath)

	if comeInError == nil {
		comeInError = errors.New("")
	}
	if fileName == "" {
		fileName = "error.log"
	}
	f, err := os.OpenFile(filepath.Join(basePath, "../", "logs", "logs", fileName), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	errLog := log.New(f, "", log.Ldate|log.Ltime)
	
	_, err = errLog.Writer().Write([]byte(fmt.Sprintf("%v\t%s\t%v\n",time.Now().Format("2006-01-02 15:04:05"), message, comeInError.Error())))
	fmt.Println(message, comeInError.Error())
	
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
