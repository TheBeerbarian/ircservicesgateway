package ircnetservices

import (
	//"encoding/json"
	//"errors"
	"fmt"
	//"log"
	//"math/rand"
	"net/http"

	//"github.com/amfranz/go-xmlrpc-client"  	//forked from "github.com/nilshell/xmlrpc"
)

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
	fmt.Fprintln(w, "  <div style=\"margin: 0 auto;padding: 0;width: 600px;text-align: center;\">")
	fmt.Fprintln(w, "    View page source for the values and names of the associated POST methods.<br><br>")
	fmt.Fprintln(w, "  </div>")
	fmt.Fprintln(w, "    <form action=\""+NetservicesConfig.IRCservicesURI+"\" method=\"post\">")
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

