#!/bin/bash

echo $1
echo $2
/usr/bin/expect<<EOF

set timeout -1
spawn ssh -p 2222 $1@localhost {uname -a; df -h}
expect "*?assword: "
send -- "vagrant\r"
expect eof
EOF

echo "foo"
exit 0
