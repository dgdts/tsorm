package log

import (
	"io"
	"log"
	"os"
	"sync"
)

// errorLog and infoLog are the loggers for error and info messages respectively.
var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog} // loggers holds references to errorLog and infoLog
	mu       sync.Mutex                         // mu is a mutex for synchronization
)

// Error and Errorf are functions for logging error messages.
var (
	Error  = errorLog.Println // Error logs a message with a newline.
	Errorf = errorLog.Printf  // Errorf logs a formatted message.
)

// Info and Infof are functions for logging info messages.
var (
	Info  = infoLog.Println // Info logs a message with a newline.
	Infof = infoLog.Printf  // Infof logs a formatted message.
)

// Constants representing different log levels.
const (
	InfoLevel  = iota // InfoLevel represents the info log level.
	ErrorLevel        // ErrorLevel represents the error log level.
	Disabled          // Disabled represents the disabled log level.
)

// SetLevel sets the log level for both errorLog and infoLog.
func SetLevel(level int) {
	mu.Lock()         // Acquire the lock to ensure exclusive access.
	defer mu.Unlock() // Release the lock when done.

	// Set the output of both loggers to os.Stdout initially.
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	// If the log level is higher than ErrorLevel, set the output of errorLog to discard.
	if ErrorLevel < level {
		errorLog.SetOutput(io.Discard)
	}

	// If the log level is higher than InfoLevel, set the output of infoLog to discard.
	if InfoLevel < level {
		infoLog.SetOutput(io.Discard)
	}
}
