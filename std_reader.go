// This program is designed to execute scripts
// that have long running processes in them and
// capture the STDOUT and STDERR and send it to 
// the syslog on the server.
// Author: Randy Coburn
// Date: 07 September 2016

package main

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"bufio"
)

// This function is used to call the logger command.
// It will take the line that you feed it and push it
// to the syslog.
func syslogger(txt string) {
	sysl := exec.Command("logger", txt)
	sysl.Env = []string{"PATH=/bin:/usr/bin"}
	_, err := sysl.Output()
	if err != nil {
		log.Fatal(os.Stderr, "Failed to log to logger.", err)
	}
}

func main() {
	// Check that the user has in deed supplied a script
	// run. Some further checking may be needed later.
	if len(os.Args) < 2 {
		log.Fatal("You have not supplied a script to run.")
		return
	}
	// find logger
	//loggger_location := exec.Output
	// Create the command that will run.
	longrunner := exec.Command("/bin/bash", os.Args[1])

	// Here we attach pipes to the STDOUT. We defer the close to
	// clean up later.
	longrunner_stdout,err := longrunner.StdoutPipe()
	if err != nil {
		log.Fatal(os.Stderr, "Error creating Stdout Pipe for cmd", err)
	}
	defer longrunner_stdout.Close()

	// The scanner looks for new items in the pipe.
	// The grot will read them.
	scanner := bufio.NewScanner(longrunner_stdout)
	go func(){
		for scanner.Scan(){
			syslogger(scanner.Text())
		}
	}()

	// Same as above just for STDERR
	longrunner_stderr, err := longrunner.StderrPipe()
	if err != nil{
		log.Fatal(os.Stderr, "Failed to create STDERR Pipe.", err)
	}
	defer longrunner_stderr.Close()
	stdErrScanner := bufio.NewScanner(longrunner_stderr)
	go func () {
		for stdErrScanner.Scan(){
			syslogger(stdErrScanner.Text())
		}
	}()

	// We start the process and then wait for it to end.
	// It is assumed that it is a very long running process.
	err = longrunner.Start()
	if err != nil{
		log.Fatal(os.Stderr, "Failed to start loop.", err)
	}
	err = longrunner.Wait()
	if err != nil {
		if err.Error() != "exit status 1" {
			fmt.Println(err)
			log.Fatal(err)
		}
	}
}