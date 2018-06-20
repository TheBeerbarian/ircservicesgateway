package ircservicesgateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/amfranz/go-xmlrpc-client"  	//forked from "github.com/nilshell/xmlrpc"
)

var (
	netservicesConfig ConfigNetServices

	DEBUG = 1
	INFO  = 2
	WARN  = 3
)

func ircservicesHTTPHandler(router *http.ServeMux) {
	var err error

	//Get ConfigNetServices
	netservicesConfig, err = loadNetServices()
	if err != nil {
		logOut(3, "No IRC Network Services available")
		return
	}

	router.HandleFunc(netservicesConfig.IRCservicesURI, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			logOut(DEBUG, "Request method: %s", r.Method)
			ircservicesRespond(w, r)
		case "POST":
			logOut(DEBUG, "Request method: %s", r.Method)
			ircservicesCommand(w, r)
		default:
			logOut(DEBUG, "Invalid request method: %s", r.Method)
			return
		}
	})
}

func ircservicesCommand(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error reading POST form data", http.StatusInternalServerError)
		return
	}

	//DEBUG logOut(1, "\u001b[31mI am here. \u001b[0m")

	//Setup client writer/responder to send/receive data from XMLRPC server source.
	client, _ := xmlrpc.NewClient(netservicesConfig.XmlrpcURL, nil)
	result    := xmlrpc.Struct{}
	
   	//Method being defined for testing, method will be from HTTP POST.
	method := r.PostFormValue("method")

	//Have method. Attempt to process request.
	if method != "" {
	

		logOut(DEBUG, "XMLRPC method: %s", method)

		methodMap := map[string]interface{}{ "method": method }
		
		switch method {
		case "checkAuthentication": //Anope
			username := r.PostFormValue("nick")
			password := r.PostFormValue("password")

			//Anope checkAuthentication
			err := client.Call("checkAuthentication", []string{username, password}, &result)
			check (err, w, r)
			for k, v := range methodMap {
			        result[k] = v
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
			
		case "command": //Anope command
		        service := r.PostFormValue("service")
		        command := r.PostFormValue("command")
		        err := client.Call("command", []string{service, "ircservicesgateway", command}, &result)
		        //err := client.Call("command", []string{"nickserv", "Test", "info CtB"}, &result)
			check (err, w, r)
			serviceMap := map[string]interface{}{ "service": service } //replace with service
			commandMap := map[string]interface{}{ "command": command } //replace with command
			for k, v := range methodMap {
			        result[k] = v
			}
			for k, v := range serviceMap {
			        result[k] = v
			}
			for k, v := range commandMap {
			        result[k] = v
			}
			userMap := map[string]interface{}{ "user": "Test" } //replace with user
			for k, v := range userMap {
			        result[k] = v
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
			
		case "stats":  //Anope stats
		        err := client.Call("stats", nil, &result)
			check (err, w, r)
			for k, v := range methodMap {
			        result[k] = v
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
			
		case "channel":  //Anope channel
		        channel := r.PostFormValue("channel")
		        err := client.Call("channel", []string{channel}, &result)
			check (err, w, r)
			for k, v := range methodMap {
			        result[k] = v
			}
			channelMap := map[string]interface{}{ "channel": channel }
			for k, v := range channelMap {
			        result[k] = v
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
			
		case "user":  //Anope user
		        user := r.PostFormValue("user")
		        err := client.Call("user", []string{user}, &result)
			check (err, w, r)
			for k, v := range methodMap {
			        result[k] = v
			}
			userMap := map[string]interface{}{ "user": user }
			for k, v := range userMap {
			        result[k] = v
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
			
		case "opers":  //Anope opers
		        err := client.Call("opers", nil, &result)
			check (err, w, r)
			for k, v := range methodMap {
			        result[k] = v
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
			
		case "notice":  //Anope notice
		        //err := client.Call("notice", []string{source, target, message}, &result)
		        err := client.Call("notice", []string{"CtB", "CtB", "Test message."}, &result)
			check (err, w, r)
			for k, v := range methodMap {
			        result[k] = v
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
			
		default:
	        	w.Header().Set("Content-Type", "application/json")
			temp := map[string]string{"result": "error", "error": "Invalid Method"}
		        json.NewEncoder(w).Encode(temp)
			for k, v := range methodMap {
			        result[k] = v
			}
			logOut(WARN, "Not a valid xmlrpc method: %s", method)
		        return
		}
		
	}
}

// Generate a temporary developer pages to post data to verify services functions are working.
// Config.conf option to enable/disable.

func ircservicesRespond(w http.ResponseWriter, r *http.Request) {
        if netservicesConfig.IRCservicesTest { 
	        loadPage(w, r)
		logOut(DEBUG, "Present XMLRPC POST request page.")
	} else {
        	w.Header().Set("Content-Type", "application/json")
		temp := map[string]string{"status": "ready", "info": "Post a valid method"}
	        json.NewEncoder(w).Encode(temp)
		logOut(DEBUG, "Ready for XMLRPC request.")
	}

}

func loadNetServices() (ConfigNetServices, error) {
	var ret ConfigNetServices

	if len(Config.NetServices) == 0 {
		return ret, errors.New("No IRC Network Services available")
	}

	randIdx := rand.Intn(len(Config.NetServices))
	ret = Config.NetServices[randIdx]

	return ret, nil
}

func loadPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "  <head>")
	fmt.Fprintln(w, "    <title>IRC Services Test</title>")
	fmt.Fprintln(w, "  </head>")
	fmt.Fprintln(w, "  <body>")
	fmt.Fprintln(w, "    <h1 style=\"color: black; font-family: verdana; text-align: center;\">")
	fmt.Fprintln(w, "      IRC Services Tests Enabled")
	fmt.Fprintln(w, "    </h1>")
	fmt.Fprintln(w, "    <form action=\""+netservicesConfig.IRCservicesURI+"\" method=\"post\">")
	fmt.Fprintln(w, "    <table style=\"margin: 0 auto;padding: 0;width: 450px;text-align: center;\">")
	fmt.Fprintln(w, "      <th colspan=\"3\" style=\"margin: 0 auto;padding: 0;text-align: left;\">")
	fmt.Fprintln(w, "        &nbsp;&nbsp;Test to check nickname and password authentication.")
	fmt.Fprintln(w, "      </th>")
	fmt.Fprintln(w, "      <tr>")
	fmt.Fprintln(w, "        <td>")
	fmt.Fprintln(w, "          <input type=\"radio\" name=\"method\" value=\"checkAuthentication\"> ")
	fmt.Fprintln(w, "        </td><td>")
	fmt.Fprintln(w, "          <input type=\"text\" placeholder=\"Nickname\" name=\"nick\">&nbsp;&nbsp;")
	fmt.Fprintln(w, "        </td><td>")
	fmt.Fprintln(w, "          <input type=\"password\" placeholder=\"Password\" name=\"password\">")
	fmt.Fprintln(w, "        </td>")
	fmt.Fprintln(w, "      </tr>")

	fmt.Fprintln(w, "      <th colspan=\"3\" style=\"margin: 0 auto;padding: 0;text-align: left;\">")
	fmt.Fprintln(w, "        <hr>&nbsp;&nbsp;Test to check command to a Network Service.")
	fmt.Fprintln(w, "      </th>")
	fmt.Fprintln(w, "      <tr>")
	fmt.Fprintln(w, "        <td>")
	fmt.Fprintln(w, "          <input type=\"radio\" name=\"method\" value=\"command\"> ")
	fmt.Fprintln(w, "        </td><td>")
	fmt.Fprintln(w, "          <input type=\"text\" placeholder=\"Service Name\" name=\"service\">&nbsp;&nbsp;")
	fmt.Fprintln(w, "        </td><td>")
	fmt.Fprintln(w, "          <input type=\"text\" placeholder=\"Command\" name=\"command\">")
	fmt.Fprintln(w, "        </td>")
	fmt.Fprintln(w, "      </tr>")

	fmt.Fprintln(w, "      <th colspan=\"3\" style=\"margin: 0 auto;padding: 0;text-align: left;\">")
	fmt.Fprintln(w, "         <hr>&nbsp;&nbsp;Test to request #channel information.")
	fmt.Fprintln(w, "      </th>")
	fmt.Fprintln(w, "      <tr>")
	fmt.Fprintln(w, "        <td>")
	fmt.Fprintln(w, "          <input type=\"radio\" name=\"method\" value=\"channel\"> ")
	fmt.Fprintln(w, "        </td><td>")
	fmt.Fprintln(w, "          <input type=\"text\" placeholder=\"#channel\" name=\"channel\">&nbsp;&nbsp;")
	fmt.Fprintln(w, "        </td><td>")
	fmt.Fprintln(w, "        </td>")
	fmt.Fprintln(w, "      </tr>")
	
	fmt.Fprintln(w, "      <th colspan=\"3\" style=\"margin: 0 auto;padding: 0;text-align: left;\">")
	fmt.Fprintln(w, "        <hr>&nbsp;&nbsp;Test to request user information.")
	fmt.Fprintln(w, "      </th>")
	fmt.Fprintln(w, "      <tr>")
	fmt.Fprintln(w, "        <td>")
	fmt.Fprintln(w, "          <input type=\"radio\" name=\"method\" value=\"user\"> ")
	fmt.Fprintln(w, "        </td><td>")
	fmt.Fprintln(w, "          <input type=\"text\" placeholder=\"Nickname\" name=\"user\">&nbsp;&nbsp;")
	fmt.Fprintln(w, "        </td><td>")
	fmt.Fprintln(w, "        </td>")
	fmt.Fprintln(w, "      </tr>")
	
	fmt.Fprintln(w, "      <th colspan=\"3\" style=\"margin: 0 auto;padding: 0;text-align: left;\">")
	fmt.Fprintln(w, "        <hr>&nbsp;&nbsp;Test to check status of Network Services.")
	fmt.Fprintln(w, "      </th>")
	fmt.Fprintln(w, "      <tr>")
	fmt.Fprintln(w, "        <td>")
	fmt.Fprintln(w, "          <input type=\"radio\" name=\"method\" value=\"stats\">")
	fmt.Fprintln(w, "        </td><td>")
	fmt.Fprintln(w, "        </td><td>")
	fmt.Fprintln(w, "        </td>")
	fmt.Fprintln(w, "      </tr>")

	fmt.Fprintln(w, "      <th colspan=\"3\" style=\"margin: 0 auto;padding: 0;text-align: left;\">")
	fmt.Fprintln(w, "        <br><button type=\"Submit\" value=\"Submit\">Submit</button>")
	fmt.Fprintln(w, "      </th>")
	fmt.Fprintln(w, "      </td>")
	fmt.Fprintln(w, "    </tr>")
	fmt.Fprintln(w, "    </table>")
	fmt.Fprintln(w, "  </form>")
	fmt.Fprintln(w, "  </body>")
	fmt.Fprintln(w, "</html>")
}

func check(err error, w http.ResponseWriter, r *http.Request) {
        if err != nil {
		w.Header().Set("Content-Type", "application/json")
		temp := map[string]string{"result": "Failure", "error": "error"}
		json.NewEncoder(w).Encode(temp)
	        logOut(2, "Error: %s", err)
		return
        }
}

func BytesToString(data []byte) string {
     return string(data[:])
}

