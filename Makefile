export GOPATH=${CURDIR}
make:
		go install cmds/http-with-console
		go install cmds/ssh-execute
