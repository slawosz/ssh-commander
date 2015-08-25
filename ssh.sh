#!/bin/bash

/usr/bin/expect<<EOF

set timeout 10
spawn ssh -StrictHostKeyChecking=no -p $PORT $USER@$HOST {$COMMAND}
expect "*?assword: "
send -- "$PASSWORD\r"
expect eof
EOF

exit 0
