package ircservicesgateway

import (
	"errors"
	"math/rand"
	"net/http"

	"github.com/thebeerbarian/ircservicesgateway/pkg/ircnetservices"
)

var (
	NetservicesConfig ConfigNetServices
	DEBUG = 1
	INFO  = 2
	WARN  = 3
)

func ircservicesHTTPHandler(router *http.ServeMux) {
	var err error

	//Get ConfigNetServices
	NetservicesConfig, err = loadNetServices()
	if err != nil {
		logOut(WARN, "ircservicesgateway:No IRC Network Services available")
		return
	}
	
        //Pass log level to ircnetservices
        ircnetservices.ConfigLog(Config.LogLevel)
	
	//Pass config options to ircnetservices
	ircnetservices.ConfigOptions(NetservicesConfig.XmlrpcURL,
	        NetservicesConfig.NsCookieName,
		NetservicesConfig.NsCookieHashKey,
		NetservicesConfig.NsCookieDomain,
		NetservicesConfig.IRCservicesURI,
		NetservicesConfig.IRCservicesTest)
		
	//Setup Handler and pass http requests to ircnetservices
	router.HandleFunc(NetservicesConfig.IRCservicesURI, func(w http.ResponseWriter, r *http.Request) {
		ircnetservices.Netirc(w, r)
	})
}

// Load config options for ircnetservices.
func loadNetServices() (ConfigNetServices, error) {
	var ret ConfigNetServices

	if len(Config.NetServices) == 0 {
		return ret, errors.New("No IRC Network Services available")
	}

	randIdx := rand.Intn(len(Config.NetServices))
	ret = Config.NetServices[randIdx]
	return ret, nil
}
