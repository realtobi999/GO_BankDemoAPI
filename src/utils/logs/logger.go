package logs

import (
	"fmt"
	"log"
	"os"
	"time"
)

const PathToLogs string = "src/utils/logs/logs.txt"
const PathToTestLogs string = "src/utils/logs/logs_test.txt"

type Logger struct {
	logger *log.Logger
	file   *os.File
}

func NewLogger(path string) *Logger {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return &Logger{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		file:   file,
	}
}

func (l *Logger) LogEvent(message any) {
	l.logToConsole("EVENT", message)
	l.logToFile("EVENT", message)
}

func (l *Logger) LogDebug(message any) {
	l.logToConsole("DEBUG", message)
	l.logToFile("DEBUG", message)
}

func (l *Logger) LogError(message any) {
	l.logToConsole("ERROR", message)
	l.logToFile("ERROR", message)
}

func (l *Logger) LogWarning(message any) {
	l.logToConsole("WARNING", message)
	l.logToFile("WARNING", message)
}

func (l *Logger) logToConsole(bracketTxt, message any) {
	l.logger.Printf("[%s]:\t %s\n", bracketTxt, message)
}

func (l *Logger) logToFile(bracketTxt, message any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(l.file, "%v [%s]:\t %s\n", timestamp, bracketTxt, message)
}

func (l *Logger) Fatal(data ...interface{}) {
	for _, d := range data {
		l.logToConsole("ERROR", d)
		l.logToFile("ERROR", d)
	}

	os.Exit(1) // Exit with a non-zero status code to indicate an error.
}