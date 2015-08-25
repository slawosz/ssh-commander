#!/bin/bash

/usr/bin/expect<<EOF

set timeout -1
spawn ssh -p 2222 vagrant@localhost {uname -a; df -h}
expect "*?assword: "
send -- "vagrant\r"
expect eof
EOF

echo "foo"
exit 0
