#!/bin/bash

/usr/bin/expect<<EOF

set timeout -1
spawn ssh -p $PORT $USER@$HOST {$COMMAND}
expect "*?assword: "
send -- "$PASSWORD\r"
expect eof
EOF

exit 0
