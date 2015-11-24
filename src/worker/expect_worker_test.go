package worker

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestWorker(t *testing.T) {
	b, err := ioutil.ReadFile("../../fixtures/vagrant_ssh_result")
	if err != nil {
		panic(err)
	}
	expected := string(b)

	commands := "whoami,ls -l"
	host := "localhost"
	port := "2222"
	user := "vagrant"
	password := "$vagrant"
	prompt := "$"

	p := &WorkerPayload{
		Host: &Host{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
		},
		Commands: commands,
		Prompt:   prompt,
	}

	w := NewExpectWorker()
	result := w.Work(p).Payload
	if result != expected {
		t.Error(fmt.Sprintf("Expected %v, got %v", expected, result))
	}
}
