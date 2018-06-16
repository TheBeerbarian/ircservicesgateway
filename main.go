package main

import (
        "flag"
        "fmt"
        "log"
        "os"
        "os/signal"
        "syscall"

        "github.com/thebeerbarian/ircservicesgateway/pkg/ircservicesgateway"
)


const VERSION = "0.1.0"

func init() {
	ircservicesgateway.Version = VERSION
}

func main() {
	printVersion := flag.Bool("version", false, "Print the version")
	configFile := flag.String("config", "config.conf", "Config file location")
	
	if *printVersion {
		fmt.Println(ircservicesgateway.Version)
		os.Exit(0)
	}

        runGateway(*configFile)
}

func runGateway(configFile string) {
        // Print any ircservicesgateway logout to STDOUT
	go printLogOutput()

	ircservicesgateway.SetConfigFile(configFile)
	log.Printf("Using config %s", ircservicesgateway.CurrentConfigFile())

	err := ircservicesgateway.LoadConfig()
	if err != nil {
		log.Printf("Config file error: %s", err.Error())
		os.Exit(1)
	}

	watchForSignals()
	ircservicesgateway.Prepare()
	ircservicesgateway.Listen()

	justWait := make(chan bool)
	<-justWait
}

func watchForSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for {
			<-c
			fmt.Println("Recieved SIGHUP, reloading config file")
			ircservicesgateway.LoadConfig()
		}
	}()
}

func printLogOutput() {
	for {
		line, _ := <-ircservicesgateway.LogOutput
		log.Println(line)
	}
}
