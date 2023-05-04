package logger

import (
	"errors"
	"net/http"
	"testing"
)

func TestLogPrint(t *testing.T) {
	result := LogPrint("Wow this worked!")
	correct := "\033[97m[Buckette 🦌]: Wow this worked!\033[0m\n"
	if result == correct {
		t.Logf("Got %s expected %s", correct, result)
	} else {
		t.Errorf("Got %s expected %s", correct, result)
	}
}

type R struct{}

func (r R) Read(p []byte) (n int, err error) {
	return 0, nil
}

func TestConnection(t *testing.T) {
	var body R
	request, err := http.NewRequest("GET", "http://localhost:8080/test", body)
	request.RequestURI = "/hello"
	request.RemoteAddr = "localhost"
	if err != nil {
		t.Error(err)
	}
	result := LogConnection(request)
	correct := "\033[97m[Buckette 🦌]: \033[0mConnection from \033[36m" + request.RemoteAddr + "\033[0m using \033[33m" + request.RequestURI + "\033[0m - \033[35m" + request.Method + "\033[0m\n"

	if result == correct {
		t.Logf("Got %s expected %s", correct, result)
	} else {
		t.Errorf("Got %s expected %s", correct, result)
	}
}

func TestLogSuccess(t *testing.T) {
	result := LogSuccess("Wow this worked!")
	correct := Green + "[Buckette 🦌]: Wow this worked!" + Reset + "\n"
	if result == correct {
		t.Logf("Got %s expected %s", correct, result)
	} else {
		t.Errorf("Got %s expected %s", correct, result)
	}
}

func TestLogWarning(t *testing.T) {
	result := LogWarning("wha", "whe", errors.New("Error!"))
	correct := Yellow + "[Buckette 🦌]: An error has occured:" + Reset + NewLine +
		White + "\t|What: wha" + Reset + NewLine +
		White + "\t|Where: whe" + Reset + NewLine +
		Yellow + "\t|Error: Error!" + Reset + NewLine
	if result != correct {
		t.Logf("Got %s expected %s", correct, result)
	} else {
		t.Errorf("Got %s expected %s", correct, result)
	}
}

func TestLogFatal(t *testing.T) {
	result := LogFatal("wha", "whe", errors.New("Error!"))
	correct := Yellow + "[Buckette 🦌]: An error has occured:" + Reset + NewLine +
		White + "\t|What: wha" + Reset + NewLine +
		White + "\t|Where: whe" + Reset + NewLine +
		Yellow + "\t|Error: Error!" + Reset + NewLine
	if result != correct {
		t.Logf("Got %s expected %s", correct, result)
	} else {
		t.Errorf("Got %s expected %s", correct, result)
	}
}
