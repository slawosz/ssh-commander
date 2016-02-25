```
make install
bin/http-with-console
curl -H "Content-Type: application/json" -X POST --data @example.json http://localhost:7632/new_job
```

Troubleshooting
* http://superuser.com/questions/302235/no-more-ptys-when-trying-to-start-screen

# Command line mode

You can run expect-worker as command line:
```
bin/expect-worker --mode=command --user=user --password=password --host=localhost --commands='ls,ls -l / | wc -l' --prompt='$'
```

And you will get plain text outpu:

```
```

# Http server mode

```
bin/expect-worker --listen=6666
```

and then send request like this:

```
POST /exe

[
  {
    "Command": "ls -l / | wc -l",
    "Host": "localhost",
    "Port": "2222",
    "User": "vagrant",
    "Password": "vagrant",
    "Prompt": "$",
    "JID": "JOB-bla-bla"
  }
]
```

with such response:

```
```
