package ircservicesgateway

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	// Version - The current version of ircservicesgateway
	Version     = "-"
	HttpRouter  *http.ServeMux
	LogOutput   chan string
)

func init() {
	HttpRouter = http.NewServeMux()
	LogOutput = make(chan string, 5)
}

func Prepare() {
	initHttpRoutes()
}

func initHttpRoutes() error {
	// Add the transport route
	ircservicesHTTPHandler(HttpRouter)
			
	// Add some general server info about this ircservicesgateway instance
	HttpRouter.HandleFunc("/webirc/", func(w http.ResponseWriter, r *http.Request) {
		out, _ := json.Marshal(map[string]interface{}{
			"name":    "ircservicesgateway",
			"version": Version,
		})

		w.Write(out)
	})

	return nil
}

func Listen() {
	for _, server := range Config.Servers {
		go startServer(server)
	}
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

func startServer(conf ConfigServer) {
	addr := fmt.Sprintf("%s:%d", conf.LocalAddr, conf.Port)

	logOut(2, "Listening on %s", addr)
	err := http.ListenAndServe(addr, HttpRouter)
	logOut(3, err.Error())
}
