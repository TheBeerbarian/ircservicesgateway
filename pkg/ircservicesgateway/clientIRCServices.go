package ircservicesgateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	//"strings"

	"github.com/gorilla/securecookie"
	//"github.com/thebeerbarian/ircservicesgateway/pkg/atheme"
)

// Struct ircservices contains all fields for
var (
	//Atheme            *atheme.Atheme
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
	var (
		err        error
		authcookie = "*"
		account    = ""
		ipaddr     = "0"
		hashkey    = []byte(netservicesConfig.nsCookieHashKey)
		s          = securecookie.New(hashkey, nil)
	)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error reading POST form data", http.StatusInternalServerError)
		return
	}

	//DEBUG logOut(1, "\u001b[31mI am here. \u001b[0m")

	// Check for Secure Cookie.
	if cookie, err := r.Cookie(netservicesConfig.nsCookieName); err == nil {
		value := make(map[string]string)
		// try to decode it
		if err = s.Decode(netservicesConfig.nsCookieName, cookie.Value, &value); err == nil {
			authcookie = value["authcookie"]
			account = value["account"]
			ipaddr = value["ipaddr"]
			logOut(DEBUG, "Retreived nsCookieName. authcookie: '%s' account: '%s' ipaddr: '%s'", value["authcookie"], value["account"], value["ipaddr"])
		}
	}

	if ipaddr == "0" {
		//ipaddr = r.Header.Get("X-Forwarded-For") //TODO: may not be proxied. Need fallback.
	}

	//Setup client writer/responder to send/receive data from XMLRPC server source.
	//Atheme, err := atheme.NewAtheme(netservicesConfig.XmlrpcURL)

	// No valid authcookie, login required from form data.
	if authcookie == "*" {

		//nick := r.PostFormValue("nick")
		//password := r.PostFormValue("password")

		if err != nil {
			logOut(WARN, "%s", err)
			return
		}

		//No data back from source.
		//if Source == nil {
		//	logOut(WARN, "Data source is nil")
		//	return
		//}

		//process Login 
		//err = Login(nick, password, w, r)

		if err != nil {
			logOut(WARN, "Login error: %s", err.Error())
			return
		}

		// Valid auth.  Generate and store encoded cookie.
		if authcookie != "*" {
		//	authcookie = Authcookie
		//	account = Account
		//	value := map[string]string{
		//		"authcookie": Authcookie,
		//		"account":    Account,
		//		"ipaddr":     ipaddr,
		//	}

		//	if encoded, err := s.Encode(netservicesConfig.nsCookieName, value); err == nil {
		//		cookie := &http.Cookie{
		//			Name:   netservicesConfig.nsCookieName,
		//			Value:  encoded,
		//			Domain: netservicesConfig.nsCookieDomain,
		//		}
		//		logOut(DEBUG, "cookie ", cookie)
		//		http.SetCookie(w, cookie)
		//		logOut(DEBUG, "Stored nsCookieName. authcookie: '%s', account: '%s' ipaddr: '%s'", Authcookie, Account, ipaddr)
		//	}
		}
		
	} else {
	
		//Setup vars from decoded cookie
		//Account = account
		//Ipaddr = ipaddr

		//If requesting data similar to `/privmsg ServiceName Command args`
		//commands := strings.Split(r.PostFormValue("command"), " ")
		//logOut(DEBUG, "Commands: %s", commands)
		//result, err := Cmd(commands, w, r)

		//If was sending to service name functions. Following would just need to pass the Nick.
		//var result map[string]string
		//command := r.PostFormValue("command")
		//logOut(DEBUG, "Command: %s", command)
		//result, err := NickServ.Info(command)

		if err != nil {
			fmt.Fprint(w, "\n", err.Error(), "\n")
			logOut(WARN, "error: %s", err.Error())
			return
		}

		//fmt.Fprint(w, result, "\n")

	}

	out, _ := json.Marshal(map[string]interface{}{
		"authcookie": authcookie,
		"account":    account,
		"ipaddr":     ipaddr,
	})

	w.Write(out)
}

// Generate a temporary developer pages to post data to verify services functions are working.
// Checks will be added to check if user has valid authcookie and give status information
//   for their registration otherwise prompt for login.  Maybe have a config to enable/disable.

func ircservicesRespond(w http.ResponseWriter, r *http.Request) {
	var (
		//err error
		authcookie = "*"
		account    = ""
		ipaddr     = "0" //r.Header.Get("X-Forwarded-For") //TODO: may not be proxied. Need fallback.
		hashkey    = []byte(netservicesConfig.nsCookieHashKey)
		s          = securecookie.New(hashkey, nil)
	)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error reading POST form data", http.StatusInternalServerError)
		return
	}

	// Check for Secure Cookie.
	if cookie, err := r.Cookie(netservicesConfig.nsCookieName); err == nil {
		value := make(map[string]string)
		// try to decode it
		if err = s.Decode(netservicesConfig.nsCookieName, cookie.Value, &value); err == nil {
			authcookie = value["authcookie"]
			account = value["account"]
			ipaddr = value["ipaddr"]
			logOut(DEBUG, "Retreived nsCookieName. authcookie: '%s' account: '%s' ipaddr: '%s'", value["authcookie"], value["account"], value["ipaddr"])
		}
	}

	logOut(DEBUG, "authcookie: '%s'", authcookie)

	//No cookie. Need to Auth
	switch authcookie {
	case "*":
		loginpage(w, r)

	//Have cookie
	default:
		//if expired or invalid cookie {
		//       loginpage(w, r)
		//} else {
		//present Irc Network Services command entry form
		postpage(w, r)

		out, _ := json.Marshal(map[string]interface{}{
			"authcookie": authcookie,
			"account":    account,
			"ipaddr":     ipaddr,
		})

		fmt.Fprintln(w, "Authcookie found. No login required.\n")
		w.Write(out)
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

func loginpage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "  <head>")
	fmt.Fprintln(w, "    <title>IRC Services Test</title>")
	fmt.Fprintln(w, "  </head>")
	fmt.Fprintln(w, "  DEBUG: Missing or expired authcookie. Login required.<br><br>")
	fmt.Fprintln(w, "  <div style=\"margin: 0 auto;padding: 0;width: 800px;\"><body>")
	fmt.Fprintln(w, "    <h1 style=\"color: black; font-family: verdana; text-align: center;\">")
	fmt.Fprintln(w, "      IRC Services Test")
	fmt.Fprintln(w, "    </h1>")
	fmt.Fprintln(w, "    <form action=\""+netservicesConfig.IRCservicesURI+"\" method=\"post\">")
	fmt.Fprintln(w, "    <div style=\"text-align: center;margin: 0 auto;\">")
	fmt.Fprintln(w, "      <label for=\"nick\">Nickname</label>")
	fmt.Fprintln(w, "      <input type=\"text\" name=\"nick\">&nbsp;&nbsp;")
	fmt.Fprintln(w, "      <label for=\"password\">Password</label>")
	fmt.Fprintln(w, "      <input type=\"password\" name=\"password\"><br><br>")
	fmt.Fprintln(w, "      <button type=\"Submit\" value=\"Submit\">Submit</button>")
	fmt.Fprintln(w, "    </div></form>")
	fmt.Fprintln(w, "  </body></div>")
	fmt.Fprintln(w, "</html>")
}

func postpage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "  <head>")
	fmt.Fprintln(w, "    <title>IRC Services Test</title>")
	fmt.Fprintln(w, "  </head>")
	fmt.Fprintln(w, "  <div style=\"margin: 0 auto;padding: 0;width: 800px;\"><body>")
	fmt.Fprintln(w, "    <h1 style=\"color: black; font-family: verdana; text-align: center;\">")
	fmt.Fprintln(w, "      IRC Services Test</h1>")
	fmt.Fprintln(w, "    <form action=\""+netservicesConfig.IRCservicesURI+"\" method=\"post\">")
	fmt.Fprintln(w, "    <div style=\"text-align: center;margin: 0 auto;\">")
	fmt.Fprintln(w, "      <label for=\"command\">Enter an IRC Network Services Command<br>")
	fmt.Fprintln(w, "        Example: `NickServ Help`")
	fmt.Fprintln(w, "      </label><br>")
	fmt.Fprintln(w, "      <input type=\"text\" name=\"command\"><br><br>")
	fmt.Fprintln(w, "      <button type=\"Submit\" value=\"Submit\">Submit</button>")
	fmt.Fprintln(w, "    </div></form>")
	fmt.Fprintln(w, "  </body></div>")
	fmt.Fprintln(w, "</html>")
}
