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
			ircservicesRespond(w)
		case "POST":
			logOut(DEBUG, "Request method: %s", r.Method)
			output, err := 	ircservicesCommand(r)
			if err != nil {
			        logOut(DEBUG, "Error: %s", err)
				if netservicesConfig.IRCservicesTest {
				        loadPage(w)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(output)
		default:
			logOut(DEBUG, "Invalid request method: %s", r.Method)
		}
	})
}

func ircservicesCommand(r *http.Request) (output xmlrpc.Struct, err error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	//DEBUG logOut(1, "\u001b[31mI am here. \u001b[0m")

	//Setup client writer/responder to send/receive data from XMLRPC server source.
	client, _ := xmlrpc.NewClient(netservicesConfig.XmlrpcURL, nil)
	result := xmlrpc.Struct{}
	
   	//Method being defined for testing, method will be from HTTP POST.
	method := r.PostFormValue("method")

	//Have method. Attempt to process request.
	if method != "" {

		logOut(DEBUG, "XMLRPC method: %s", method)
		
		// Default minimum map setup.
		methodMap := xmlrpc.Struct{ "method": method }
		
		switch method {
		case "checkAuthentication": //Anope
			username := r.PostFormValue("nick")
			password := r.PostFormValue("password")
			if err := client.Call(method, []string{username, password}, &result); err != nil {
				return xmlrpc.Struct{"result": "error", "error": "Internal Client Call Error"}, err
			}
			result = mergeMaps(result, methodMap)
			return result, nil
			
		case "command": //Anope command
		        service := r.PostFormValue("service")
		        command := r.PostFormValue("command")
		        if err := client.Call(method, []string{service, "ircservicesgateway", command}, &result); err != nil {
				return xmlrpc.Struct{"result": "error", "error": "Internal Client Call Error"}, err
			}
			methodMap := xmlrpc.Struct{ "method": method, "service": service, "command": command,
				  "user": "ircservicesgateway" } //replace with requesting user when possible?
			result = mergeMaps(result, methodMap)
			return result, nil
			
		case "stats":  //Anope stats
		        if err := client.Call(method, nil, &result); err != nil {
				return xmlrpc.Struct{"result": "error", "error": "Internal Client Call Error"}, err
			}
			result = mergeMaps(result, methodMap)
			return result, nil

		case "channel":  //Anope channel
		        channel := r.PostFormValue("channel")
		        if err := client.Call("channel", []string{channel}, &result); err != nil {
				return xmlrpc.Struct{"result": "error", "error": "Internal Client Call Error"}, err
			}
			methodMap := xmlrpc.Struct{ "method": method, "channel": channel }
			result = mergeMaps(result, methodMap)
			return result, nil
			
		case "user":  //Anope user
		        user := r.PostFormValue("user")
		        if err := client.Call("user", []string{user}, &result); err != nil {
				return xmlrpc.Struct{"result": "error", "error": "Internal Client Call Error"}, err
			}
			methodMap := xmlrpc.Struct{ "method": method, "user": user }
			result = mergeMaps(result, methodMap)
			return result, nil
			
		case "opers":  //Anope opers
		        if err := client.Call("opers", nil, &result); err != nil {
				return xmlrpc.Struct{"result": "error", "error": "Internal Client Call Error"}, err
			}
			result = mergeMaps(result, methodMap)
			return result, nil
			
		case "notice":  //Anope notice
		        //err := client.Call("notice", []string{source, target, message}, &result)
		        if err := client.Call("notice", []string{"CtB", "CtB", "Test message."}, &result); err != nil {
				return xmlrpc.Struct{"result": "error", "error": "Internal Client Call Error"}, err
			}
			result = mergeMaps(result, methodMap)
			return result, nil
			
		default:
			defaultMap := xmlrpc.Struct{"result": "error", "error": "Invalid Method"}
			logOut(DEBUG, "Not a valid xmlrpc method: %s", method)
			return defaultMap, nil
		}
		
	} else {
		logOut(DEBUG, "No method. Resending XMLRPC POST request page.")
	        if netservicesConfig.IRCservicesTest {
		        return xmlrpc.Struct{"result": "error", "error": "No Method"},
		                errors.New("No method")
		}
		return xmlrpc.Struct{"result": "error", "error": "No Method"},
		        errors.New("No method")
	}
	return
}

// Generate a temporary developer pages to post data to verify services functions are working.
// Config.conf option to enable/disable.

func ircservicesRespond(w http.ResponseWriter) {
        if netservicesConfig.IRCservicesTest { 
	        loadPage(w)
		logOut(DEBUG, "Present XMLRPC POST request page.")
		return
	} else {
        	w.Header().Set("Content-Type", "application/json")
		temp := map[string]string{"status": "ready", "info": "Post a valid method"}
	        json.NewEncoder(w).Encode(temp)
		logOut(DEBUG, "Ready for XMLRPC request.")
		return
	}
}

// loadPage - I didn't want to setup a template.
func loadPage(w http.ResponseWriter) {
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
	fmt.Fprintln(w, "        <hr>&nbsp;&nbsp;Test to check for invalid POST Method.")
	fmt.Fprintln(w, "      </th>")
	fmt.Fprintln(w, "      <tr>")
	fmt.Fprintln(w, "        <td>")
	fmt.Fprintln(w, "          <input type=\"radio\" name=\"method\" value=\"notvalid\">")
	fmt.Fprintln(w, "        </td><td>")
	fmt.Fprintln(w, "        </td><td>")
	fmt.Fprintln(w, "        </td>")
	fmt.Fprintln(w, "      </tr>")

	fmt.Fprintln(w, "      <th colspan=\"3\" style=\"margin: 0 auto;padding: 0;text-align: left;\">")
	fmt.Fprintln(w, "          <br><button type=\"Submit\" value=\"Submit\">Submit</button>")
	fmt.Fprintln(w, "      </th>")
	fmt.Fprintln(w, "    </table>")
	fmt.Fprintln(w, "  </form>")
	fmt.Fprintln(w, "  </body>")
	fmt.Fprintln(w, "</html>")
}

func BytesToString(data []byte) string {
     return string(data[:])
}

func mergeMaps(map1 xmlrpc.Struct, map2 xmlrpc.Struct) xmlrpc.Struct {
	for k, v := range map2 {
	        map1[k] = v
	}
	return map1
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
