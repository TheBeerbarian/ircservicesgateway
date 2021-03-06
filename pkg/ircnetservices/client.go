package ircnetservices

import (
	"errors"
	"net/http"

	"github.com/amfranz/go-xmlrpc-client"  	//forked from "github.com/nilshell/xmlrpc"
)

// Process post requests for the xmlrpm calls
func ircservicesCommand(r *http.Request) (output xmlrpc.Struct, err error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	//DEBUG logOut(1, "\u001b[31mI am here. \u001b[0m")

	//Setup client writer/responder to send/receive data from XMLRPC server source.
	client, _ := xmlrpc.NewClient(NetservicesConfig.XmlrpcURL, nil)
	result := xmlrpc.Struct{}
	
   	//Method being defined for testing, method will be from HTTP POST.
	method := r.PostFormValue("method")

	//Have method. Attempt to process request.
	if method != "" {

		logOut(DEBUG, "ircservicesgateway:XMLRPC method: %s", method)
		
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
		default:
			defaultMap := xmlrpc.Struct{"result": "error", "error": "Invalid Method"}
			logOut(DEBUG, "ircservicesgateway:Not a valid xmlrpc method: %s", method)
			return defaultMap, nil
		}
	} else {
	        if NetservicesConfig.IRCservicesTest {
		logOut(DEBUG, "ircservicesgateway:No method. Resending XMLRPC POST request page.")
		        return xmlrpc.Struct{"result": "error", "error": "No Method"},
		                errors.New("No method")
		}
		logOut(DEBUG, "ircservicesgateway:No method.")
		return xmlrpc.Struct{"result": "error", "error": "No Method"},
		        errors.New("No method")
	}
	return
}

func mergeMaps(map1 xmlrpc.Struct, map2 xmlrpc.Struct) xmlrpc.Struct {
	for k, v := range map2 {
	        map1[k] = v
	}
	return map1
}

func BytesToString(data []byte) string {
     return string(data[:])
}
