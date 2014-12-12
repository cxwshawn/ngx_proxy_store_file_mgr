package srvlog

import (
	"config"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const (
	LogError, LogAccess = 0, 1
)

func InitLog(serverName string) error {
	logError := config.IsLogError()
	logAccess := config.IsLogAccess()
	if logError == true {
		err := initLog(os.Args[0], serverName, LogError)
		if err != nil {
			return err
		}
	}
	if logAccess == true {
		err := initLog(os.Args[0], serverName, LogAccess)
		if err != nil {
			return err
		}
	}
	return nil
}

var DefaultLogger []*log.Logger
var DefaultLogFile []*os.File

func initLog(exe string, logName string, logType int) error {
	exePath, err := exec.LookPath(exe)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	exeFullPath, err := filepath.Abs(exePath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	var logFullPath string
	if logType == LogError {
		logFullPath = fmt.Sprintf("%s/%s_%d_error.log", filepath.Dir(exeFullPath), logName, syscall.Getpid())
	} else if logType == LogAccess {
		logFullPath = fmt.Sprintf("%s/%s_%d_access.log", filepath.Dir(exeFullPath), logName, syscall.Getpid())
	}

	file, err := os.Create(logFullPath)
	//LogErrorFile = file
	DefaultLogFile = append(DefaultLogFile, file)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	srvLogger := log.New(file, logName, log.Ldate|
		log.Ltime|log.Lmicroseconds|log.Llongfile)
	DefaultLogger = append(DefaultLogger, srvLogger)

	return nil
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	if config.IsLogAccess() == true {
		DefaultLogger[LogAccess].Output(2, fmt.Sprintf(format, v...))
	}
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	if config.IsLogAccess() == true {
		DefaultLogger[LogAccess].Output(2, fmt.Sprint(v...))
	}
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	if config.IsLogAccess() == true {
		DefaultLogger[LogAccess].Output(2, fmt.Sprintln(v...))
	}
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	if config.IsLogAccess() == true {
		DefaultLogger[LogAccess].Output(2, fmt.Sprint(v...))
		os.Exit(1)
	}
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func EPrintf(format string, v ...interface{}) {
	DefaultLogger[LogError].Output(2, fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func EPrint(v ...interface{}) { DefaultLogger[LogError].Output(2, fmt.Sprint(v...)) }

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func EPrintln(v ...interface{}) { DefaultLogger[LogError].Output(2, fmt.Sprintln(v...)) }

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func EFatal(v ...interface{}) {
	DefaultLogger[LogError].Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	DefaultLogger[LogError].Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	DefaultLogger[LogError].Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}
