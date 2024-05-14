package Logger

import (
	"log"
	"os"
)

var infoLogFile *os.File
var errorLogFile *os.File
var infoLogger *log.Logger
var errorLogger *log.Logger

func Init() {
	logsPath := "./logs"
	_, err := os.Stat(logsPath)
	if os.IsNotExist(err) {
		err := os.Mkdir(logsPath, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	} else if err != nil {
		log.Fatalln(err)
	}

	infoLogFile, err := os.OpenFile(logsPath+"/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("Cannot open info log file:", err)
	}

	errorLogFile, err := os.OpenFile(logsPath+"/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("Cannot open error log file:", err)
	}

	infoLogger = log.New(infoLogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(errorLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Close() {
	infoLogFile.Close()
	errorLogFile.Close()
}

func Info(message string) {
	infoLogger.Println(message)
}

func Error(err error) {
	errorLogger.Println(err.Error())
}
