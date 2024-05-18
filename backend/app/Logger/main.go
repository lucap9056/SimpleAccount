package Logger

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Manager struct {
	Info  *Logger
	Error *Logger
}

func New(logsPath string) (*Manager, error) {

	infoPath := filepath.Join(logsPath, "info.log")
	infoLogger, err := NewLogger(infoPath)
	if err != nil {
		return nil, err
	}

	errorPath := filepath.Join(logsPath, "error.log")
	errLogger, err := NewLogger(errorPath)
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		Info:  infoLogger,
		Error: errLogger,
	}

	return manager, nil
}

func (manager *Manager) Close() {
	manager.Error.Close()
	manager.Info.Close()
}

type Logger struct {
	file *os.File
	log  *log.Logger
	mux  sync.Mutex
}

func NewLogger(filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		file: file,
		log:  log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	return logger, nil
}

func (logger *Logger) Close() error {
	return logger.file.Close()
}

func (logger *Logger) Write(e interface{}) {
	logger.mux.Lock()
	defer logger.mux.Unlock()
  
	if err, ok := e.(error); ok {
		logger.log.Println(err.Error())
		return
	}
	logger.log.Println(e)
}
