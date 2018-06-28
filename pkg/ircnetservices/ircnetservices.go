package ircnetservices

import (
	"encoding/json"
	"net/http"
)

func Netirc(w http.ResponseWriter, r *http.Request) {

	logOut(DEBUG, "ircservicesgateway:Connect from remote address: '%s'", r.RemoteAddr)
	if r.Header.Get("X-Forwarded-For") != "" {
	        logOut(DEBUG, "ircservicesgateway:Connection is proxy for address: '%s'",
	                r.Header.Get("X-Forwarded-For"))
	}
	switch r.Method {
	case "GET":
		logOut(DEBUG, "ircservicesgateway:Request method: %s", r.Method)
		ircservicesRespond(w)
	case "POST":
		logOut(DEBUG, "ircservicesgateway:Request method: %s", r.Method)
		output, err := 	ircservicesCommand(r)
		if err != nil {
		        logOut(DEBUG, "ircservicesgateway:Error: %s", err)
			if NetservicesConfig.IRCservicesTest {
			        loadPage(w)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	default:
		logOut(DEBUG, "ircservicesgateway:Invalid request method: %s", r.Method)
	}

}

// Generate a temporary developer pages to post data to verify services functions are working.
// Config.conf option to enable/disable.
func ircservicesRespond(w http.ResponseWriter) {
        if NetservicesConfig.IRCservicesTest { 
	        loadPage(w)
		logOut(DEBUG, "ircservicesgateway:Present XMLRPC POST request page.")
		return
	} else {
        	w.Header().Set("Content-Type", "application/json")
		temp := map[string]string{"status": "ready", "info": "Post a valid method"}
	        json.NewEncoder(w).Encode(temp)
		logOut(DEBUG, "ircservicesgateway:Ready for XMLRPC request.")
		return
	}
}
