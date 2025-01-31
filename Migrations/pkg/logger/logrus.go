package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()

	Logger.SetLevel(logrus.DebugLevel)

	/*Logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})*/

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.SetOutput(file)
	} else {
		Logger.SetOutput(os.Stdout)
		return
	}

	multiWriter := io.MultiWriter(file, os.Stdout)
	Logger.SetOutput(multiWriter)
}

func LogWithLocation(level logrus.Level, msg string) {
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		color := getColorForLevel(level, colorRed, colorGreen)
		fmt.Printf("%s%s%s\n", color, msg, colorReset)
		Logger.Log(level, msg)
		return
	}

	rootDir, err := os.Getwd()
	if err != nil {
		color := getColorForLevel(level, colorRed, colorGreen)
		fmt.Printf("%s%s%s\n", color, msg, colorReset)
		Logger.Log(level, msg)
		return
	}

	relPath, err := filepath.Rel(rootDir, file)
	if err != nil {
		color := getColorForLevel(level, colorRed, colorGreen)
		fmt.Printf("%s%s%s\n", color, msg, colorReset)
		Logger.Log(level, msg)
		return
	}

	relPath = filepath.ToSlash(relPath)

	color := getColorForLevel(level, colorRed, colorGreen)
	fmt.Printf("%s%s%s\n", color, fmt.Sprintf("%s:%d\t-\t%s\t:\t%s", relPath, line, strings.ToUpper(level.String()), msg), colorReset)

	if level == logrus.FatalLevel {
		os.Exit(1)
	}
}

func getColorForLevel(level logrus.Level, colorRed string, colorGreen string) string {
	switch level {
	case logrus.FatalLevel, logrus.ErrorLevel:
		return colorRed
	case logrus.InfoLevel, logrus.DebugLevel:
		return colorGreen
	default:
		return ""
	}
}
