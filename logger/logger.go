package logger

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Purple  = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
	Pretext = "[Buckette 🦌]: "
	NewLine = "\n"
)

func LogPrint(s string) string {
	var b strings.Builder

	b.WriteString(White)
	b.WriteString(Pretext)
	b.WriteString(s)
	b.WriteString(Reset)
	b.WriteString("\n")

	log.Print(b.String())

	return b.String()
}

func LogConnection(req *http.Request) string {
	var b strings.Builder

	b.WriteString(White)
	b.WriteString(Pretext)
	b.WriteString(Reset)
	b.WriteString("Connection from ")
	b.WriteString(Cyan)
	b.WriteString(req.RemoteAddr)
	b.WriteString(Reset)
	b.WriteString(" using ")
	b.WriteString(Yellow)
	b.WriteString(req.RequestURI)
	b.WriteString(Reset)
	b.WriteString(" - ")
	b.WriteString(Purple)
	b.WriteString(req.Method)
	b.WriteString(Reset)
	b.WriteString(NewLine)

	log.Print(b.String())

	return b.String()
}

func LogSuccess(s string) string {
	var b strings.Builder

	b.WriteString(Green)
	b.WriteString(Pretext)
	b.WriteString(s)
	b.WriteString(Reset)
	b.WriteString("\n")

	log.Print(b.String())

	return b.String()
}

func LogWarning(what string, where string, err error) string {
	var b strings.Builder

	b.WriteString(Yellow)
	b.WriteString(Pretext)
	b.WriteString("An error has occured:")
	b.WriteString(Reset)
	b.WriteString(NewLine)

	b.WriteString(White)
	b.WriteString("\t|What: ")
	b.WriteString(what)
	b.WriteString(Reset)
	b.WriteString(NewLine)

	b.WriteString(White)
	b.WriteString("\t|Where: ")
	b.WriteString(where)
	b.WriteString(Reset)
	b.WriteString(NewLine)

	b.WriteString(Yellow)
	b.WriteString("\t|Error: ")
	b.WriteString(err.Error())
	b.WriteString(Reset)
	b.WriteString(NewLine)

	b.WriteString("\n")

	log.Print(b.String())

	return b.String()
}

func LogFatal(what string, where string, err error) string {
	var b strings.Builder

	b.WriteString(Yellow)
	b.WriteString(Pretext)
	b.WriteString("An error has occured:")
	b.WriteString(Reset)
	b.WriteString(NewLine)

	b.WriteString(White)
	b.WriteString("\t|What: ")
	b.WriteString(what)
	b.WriteString(Reset)
	b.WriteString(NewLine)

	b.WriteString(White)
	b.WriteString("\t|Where: ")
	b.WriteString(where)
	b.WriteString(Reset)
	b.WriteString(NewLine)

	b.WriteString(Yellow)
	b.WriteString("\t|Error: ")
	b.WriteString(err.Error())
	b.WriteString(Reset)
	b.WriteString(NewLine)

	b.WriteString("\n")

	log.Print(b.String())

	if flag.Lookup("test.v") != nil {
		log.Print(b.String())
	} else {
		log.Fatal(b.String())
	}

	return b.String()
}
