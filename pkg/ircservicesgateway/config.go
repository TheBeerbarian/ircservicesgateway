package ircservicesgateway

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

// ConfigServer - A web server config
type ConfigServer struct {
	LocalAddr string
	Port      int
}

// Config IRCservices - A IRC Network Services config
type ConfigNetServices struct {
	XmlrpcURL       string
	nsCookieName    string
	nsCookieHashKey string
	//	nsCookieHashKey	     []byte
	nsCookieDomain string
	IRCservicesURI string
}

// Config Structures
var Config struct {
	ConfigFile  string
	LogLevel    int
	Servers     []ConfigServer
	NetServices []ConfigNetServices
}

func SetConfigFile(ConfigFile string) {
	// Config paths starting with $ is executed rather than treated as a path
	if strings.HasPrefix(ConfigFile, "$ ") {
		Config.ConfigFile = ConfigFile
	} else {
		Config.ConfigFile, _ = filepath.Abs(ConfigFile)
	}
}

// CurrentConfigFile - Return the full path or command for the config file in use
func CurrentConfigFile() string {
	return Config.ConfigFile
}

func LoadConfig() error {
	var configSrc interface{}

	if strings.HasPrefix(Config.ConfigFile, "$ ") {
		cmdRawOut, err := exec.Command("sh", "-c", Config.ConfigFile[2:]).Output()
		if err != nil {
			return err
		}

		configSrc = cmdRawOut
	} else {
		configSrc = Config.ConfigFile
	}

	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, configSrc)
	if err != nil {
		return err
	}

	// Clear the existing config
	Config.Servers = []ConfigServer{}
	Config.NetServices = []ConfigNetServices{}

	for _, section := range cfg.Sections() {
		if strings.Index(section.Name(), "DEFAULT") == 0 {
			Config.LogLevel = section.Key("logLevel").MustInt(3)
			if Config.LogLevel < 1 || Config.LogLevel > 3 {
				logOut(3, "Config option logLevel must be between 1-3. Setting default value of 3.")
				Config.LogLevel = 3
			}
		}

		if strings.Index(section.Name(), "server.") == 0 {
			server := ConfigServer{}
			server.LocalAddr = confKeyAsString(section.Key("bind"), "127.0.0.1")
			server.Port = confKeyAsInt(section.Key("port"), 80)

			Config.Servers = append(Config.Servers, server)
		}

		if strings.Index(section.Name(), "ircservices") == 0 {
			ircservices := ConfigNetServices{}
			ircservices.XmlrpcURL = section.Key("xmlrpc_url").MustString("http://127.0.0.1:8080/xmlrpc")
			ircservices.nsCookieName = section.Key("nscookiename").MustString("IRCSERVICEAUTH")
			ircservices.nsCookieHashKey = section.Key("nscookiehashkey").MustString("MY_IRCSERVICEAUTH_HASH_KEY")
			ircservices.nsCookieDomain = section.Key("nscookiedomain").MustString("")
			ircservices.IRCservicesURI = section.Key("ircservices_uri").MustString("/webirc/ircservices/")

			Config.NetServices = append(Config.NetServices, ircservices)
		}
	}
	return nil
}

func confKeyAsString(key *ini.Key, def string) string {
	val := def

	str := key.String()
	if len(str) > 1 && str[:1] == "$" {
		val = os.Getenv(str[1:])
	} else {
		val = key.MustString(def)
	}

	return val
}

func confKeyAsInt(key *ini.Key, def int) int {
	val := def

	str := key.String()
	if len(str) > 1 && str[:1] == "$" {
		envVal := os.Getenv(str[1:])
		envValInt, err := strconv.Atoi(envVal)
		if err == nil {
			val = envValInt
		}
	} else {
		val = key.MustInt(def)
	}

	return val
}
