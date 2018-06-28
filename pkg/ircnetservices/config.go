package ircnetservices

import (
	"fmt"
	"log"
)

// Config IRCservices - A IRC Network Services config
type ConfigNetServices struct {
	XmlrpcURL       string
	NsCookieName    string
	NsCookieHashKey string
	NsCookieDomain  string
	IRCservicesURI  string
	IRCservicesTest bool
}

// Config Structures
var Config struct {
	LogLevel    int
	NetServices []ConfigNetServices
}

var (
	DEBUG = 1
	INFO  = 2
	WARN  = 3
	LogOutput chan string
	NetservicesConfig ConfigNetServices
)

func init() {
	LogOutput = make(chan string, 5)

	// Print any ircservicesgateway logout to STDOUT
	go printLogOutput()
}

func ConfigLog(loglevel int) {
        Config.LogLevel = loglevel
	//log.Printf("Config.LogLevel: %d", Config.LogLevel)
	logOut(DEBUG, "Config.LogLevel: %d", Config.LogLevel)
	return
}

func ConfigOptions(url string, cname string, chash string, cdomain string, uri string, test bool) {
        NetservicesConfig.XmlrpcURL = url
        NetservicesConfig.NsCookieName = cname
        NetservicesConfig.NsCookieHashKey = chash
        NetservicesConfig.NsCookieDomain = cdomain
        NetservicesConfig.IRCservicesURI = uri
        NetservicesConfig.IRCservicesTest = test
	logOut(DEBUG, "NetservicesConfig - XmlrpcURL: '%s' NSCookieName: '%s' NsCookieHashKey: '%s' NsCookieDomain: '%s' IRCservicesURI: '%s' IRCservicesTest: '%t'",
	        NetservicesConfig.XmlrpcURL,
	        NetservicesConfig.NsCookieName,
		NetservicesConfig.NsCookieHashKey,
		NetservicesConfig.NsCookieDomain,
		NetservicesConfig.IRCservicesURI,
		NetservicesConfig.IRCservicesTest)
}

func logOut(level int, format string, args ...interface{}) {
	if level < Config.LogLevel {
		return
	}

	levels := [...]string{"L_DEBUG", "L_INFO", "L_WARN"}
	line := fmt.Sprintf(levels[level-1]+" "+format, args...)

	select {
	case LogOutput <- line:
	}
}

func printLogOutput() {
	for {
		line, _ := <-LogOutput
		log.Println(line)
	}
}
