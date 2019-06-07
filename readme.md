SETUP
=====
```bash
echo "SPARKPOST_DOMAIN=mail.io" > .env
echo "ADDRESS=:1025"            >> .env
echo "SPARKPOST_API_KEY=***" 	>> .env
```

BUILD
=====
```bash
go build -o sparkgi.exe
```

TEST
====

in a shell run:
```bash
./sparkgi.exe
```


in an other shell run:
```bash
sendemail -f norply@foo.io -t bar@foo.io -u "Testing SparkGi" -m "message" -o message-charset=utf-8 -s localhost:1025
```