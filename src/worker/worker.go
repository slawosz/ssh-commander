package worker

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type SshWorker struct{}

func NewWorker() *SshWorker {
	return &SshWorker{}
}

func (w *SshWorker) Work(payload *WorkerPayload) *HostResult {
	// return fmt.Sprintf("executing %v on %v\n", cmd, hostname)
	config := &ssh.ClientConfig{
		User: payload.User,
		Auth: []ssh.AuthMethod{ssh.Password(payload.Password)},
		//Config: ssh.Config{Ciphers: []string{"aes192-ctr"}},
	}

	fmt.Printf("executing %v on %v:%v\n", payload.Command, payload.Host, payload.Port)
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", payload.Host.Host, payload.Port), config)
	if err != nil {
		msg := fmt.Sprintf("Connection to %v:%v failed", payload.Host, payload.Port)
		return w.handleError(payload, msg, err)
	}
	session, err := conn.NewSession()
	if err != err {
		return w.handleError(payload, "ssh session creation failed", err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Start(payload.Command)
	if err != nil {
		return w.handleError(payload, "remote command execution failed", err)
	}
	err = session.Wait()
	if err != nil {
		// fmt.Println(stdoutBuf.String())
		// return w.handleError(payload, "wait for remote command execution failed", err)
		// TODO: distinguish between error
		return &HostResult{stdoutBuf.String(), payload.Host.Host, payload.Port, "-1"}
	}

	return &HostResult{stdoutBuf.String(), payload.Host.Host, payload.Port, "0"}
}

func (w *SshWorker) handleError(payload *WorkerPayload, msg string, err error) *HostResult {
	return &HostResult{fmt.Sprintf("Execution on host FAILED \"%v\" with error \"%v\"", msg, err), payload.Host.Host, payload.Port, "-1"}
}
