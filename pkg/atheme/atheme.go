// Package atheme implements an Atheme XMLRPC client and does all the
// horrifyingly ugly scraping of the raw output to machine-usable structures.
package atheme

import (
	//"encoding/xml"
	"fmt"
	//"io"
	"net/http"
	"strings"
	"time"

	"github.com/thebeerbarian/ircservicesgateway/pkg/atheme/xmlrpc"
	//"golang.org/x/net/html/charset"
	//"github.com/kolo/xmlrpc"
)

// Atheme is an Atheme context. This contains everything a client needs to access Atheme
// data remotely.
type Atheme struct {
	Privset     []string // Privilege set of the user
	Account     string   // Account Atheme is logged in as
	serverProxy *xmlrpc.Client
	url         string
	Authcookie  string
	Ipaddr      string
	NickServ    *NickServ
	ChanServ    *ChanServ
	OperServ    *OperServ
	HostServ    *HostServ
	MemoServ    *MemoServ
	LastUsed    time.Time // When the last RPC call was made
}

// NewAtheme returns a new Atheme instance or raises an error.
func NewAtheme(url string) (atheme *Atheme, err error) {
	var serverproxy *xmlrpc.Client
	serverproxy, err = xmlrpc.NewClient(url, nil)

	if err != nil {
		return nil, err
	}

	atheme = &Atheme{
		Account:     "*",
		serverProxy: serverproxy,
		url:         url,
		Authcookie:  "*",
		Ipaddr:      "0",
		LastUsed:    time.Now(),
	}

	atheme.NickServ = &NickServ{a: atheme}
	atheme.ChanServ = &ChanServ{a: atheme}
	atheme.OperServ = &OperServ{a: atheme}
	atheme.HostServ = &HostServ{a: atheme}
	atheme.MemoServ = &MemoServ{a: atheme}

	return atheme, nil
}

// Command runs an Atheme command and gives the output or an error.
func (a *Atheme) Command(args ...string) (string, error) {
	var result string

	fullcommand := []string{a.Authcookie, a.Account, a.Ipaddr}

	for _, arg := range args {
		fullcommand = append(fullcommand, arg)
	}

	err := a.serverProxy.Call("atheme.command", &fullcommand, &result)

	a.LastUsed = time.Now()

	return result, err
}

// Cmd ... Info gets raw NickServ info on an arbitrary user or returns an error.
func (a *Atheme) Cmd(target []string, w http.ResponseWriter, r *http.Request) (res string, err error) {
	var result string
	fullcommand := []string{a.Authcookie, a.Account, a.Ipaddr}

	for _, arg := range target {
		fullcommand = append(fullcommand, arg)
	}

	err = a.serverProxy.Call("atheme.command", &fullcommand, &result)

	a.LastUsed = time.Now()

	fmt.Fprintln(w, "=== Begin DEBUG ===")
	fmt.Fprint(w, "err: ", err, "\ntarget: ", target, "\nfullcommand: ", fullcommand, "\n")
	fmt.Fprint(w, "==== End DEBUG ====\n\n\n")

	return result, err
}

// Login attempts to log a user into Atheme. It returns true or false
func (a *Atheme) Login(username string, password string, w http.ResponseWriter, r *http.Request) (err error) {
	var authcookie string
	//var reader io.Reader
	//decoder := xml.NewDecoder(reader)
	//decoder.CharsetReader = charset.NewReaderLabel
	//err = decoder.Decode(&parsed)

	//err = a.serverProxy.Call("atheme.login", []string{username, password, "::1"}, &authcookie) //Atheme
	//err = a.serverProxy.Call("checkAuthentication", []string{username, password}, decoder.Decode(&authcookie)) //Anope
	err = a.serverProxy.Call("checkAuthentication", []string{username, password}, &authcookie)

	if err != nil {
		fmt.Fprintln(w, "=== Begin DEBUG ===")
		fmt.Fprint(w, "err: ", err, "\nusername: ", username, "\npassword: *", "\n")
		fmt.Fprint(w, "==== End DEBUG ====\n\n")
		return err
	}

	a.Authcookie = authcookie
	a.Account = username

	return
}

// Logout logs a user out of Atheme. There is no return.
func (a *Atheme) Logout() {
	var res string

	a.serverProxy.Call("atheme.logout", []string{a.Authcookie, a.Account}, &res)

	a.Account = "*"
	a.Authcookie = "*"

	return
}

// GetPrivset returns the privset of a user.
func (a *Atheme) GetPrivset() (privs []string) {
	if a.Privset == nil {
		var res string

		a.serverProxy.Call("atheme.privset", []string{a.Authcookie, a.Account}, &res)

		a.Privset = strings.Split(res, " ")
	}

	return a.Privset
}
