# go_syslogger
A Simple Go program that will capture STDOUT and STDERR then send it to syslog.

## How it work?
This program will only run on a linux system. It has been tested on ubuntu 16.04 and CentOs 6.x.
It relies on bash being available in /bin/bash and logger being in either /usr/bin or /bin.

You will need to wrap your command that you wish to run in a shell script and give it execute rights.

vim exec.sh
```bash
#!/bin/bash
while true; do echo tick >&2; sleep 1; echo tock; sleep 1; done
```

Then run std_reader with the script as the argument.
```bash
std_reader exec.sh
```

It is primarily designed to wrap very long running processes that output as they run.
So you would not have to wait till the process stops to get the logs.

## What will it do?
std_reader will capture the STDOUT and STDERR from your script and send it to syslog by calling "logger $output".

