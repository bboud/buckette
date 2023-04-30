package main

import (
	"log"
	"net/http"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

type Color string

func Colorize(c Color, s string) string {
	return string(c) + s + Reset
}

func LogPrint(s string) {
	preString := Colorize(White, "[Buckette 🦌]:")

	log.Printf("%s %s\n", preString, Colorize(White, s))
}

func LogConnection(req *http.Request) {
	preString := Colorize(Blue, "[Buckette 🦌]:")

	log.Printf("%s %s %s - %s\n", preString, Colorize(White, "Connection from "+req.RemoteAddr+" using"), Colorize(Cyan, req.RequestURI), Colorize(Purple, req.Method))
}

func LogSucess(s string) {
	preString := Colorize(Green, "[Buckette 🦌]:")

	log.Printf("%s %s\n", preString, Colorize(Green, s))
}

func LogWarning(statement string, action string, err error) {
	preString := Colorize(Yellow, "[Buckette 🦌| *Warning*]: ")

	log.Printf("%s\n\t %s\n\t %s\n\t %s", preString, Colorize(White, "|What:"+statement), Colorize(Cyan, "|Where:"+action), Colorize(Yellow, "|Error Text:"+err.Error()))
}

func LogFatal(statement string, action string, err error) {
	preString := Colorize(Red, "[Buckette 🦌| *Error*]: ")

	log.Fatalf("%s\n\t %s\n\t %s\n\t %s", preString, Colorize(White, "|What:"+statement), Colorize(Cyan, "|Where:"+action), Colorize(Red, "|Error Text:"+err.Error()))
}
