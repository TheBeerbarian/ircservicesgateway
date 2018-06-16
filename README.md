
### Overview
ircservicesgateway is still under development.

The goal is for HTTP requests to be proxied by the ircservicesgateway to issue XMLRPC requests to a self-hosted irc network services (Atheme, Anope, etc.) then pass back data to the HTTP requestor.  An attempt will be made for the data to be passed back in raw format or a mapped format depending on the issued request format.

### Building and development
ircservicesgateway is built using Golang - make sure to have these programs installed and configured.

https://golang.org/dl/

To download the source:

`go get github.com/thebeerbarian/ircservicesgateway`

To update your existing source:

`go get -u github.com/thebeerbarian/ircservicesgateway`

Building from source:


cd ~/go/src/github.com/thebeerbarian/ircservicesgateway

mkdir -p ~/ircservicesgateway

go build -o ~/ircservicesgateway/ircservicesgateway main.go

cp config.conf.example ~/ircservicesgateway/config.conf

### Running
cd ~/ircservicesgateway

Run `./ircservicesgateway` to start it.

It is possible to start it up via system services like init.d or systemd but will not be included yet until further development has been completed.


Accessing the ircservicesgateway server:

During development the methods used to request the ircservicesgateway request may be in flux until fully determined with the help of other contributors.

Currently the current method is to issue a login POST method to `http://host/webirc/ircservices/` with form inputs of type text for `nick` and type password for `password`.  If successful, the network services will return an `authcookie string`, `account string`, and `ipaddr string` that the ircservicesgateway will issue an encrypted cookie to the end user.  After the end user has a valid cookie, the next POST method with a network services command will be processes as the authorized user for the login.  The post method is still under development, but may be as easy as just using a form input of type text of `nickserv info nick` for example.

Credit where credit is due:

### Derivative Work from "github.com/kiwiirc/webircgateway".
### License
~~~
   Copyright 2017 Kiwi IRC

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
~~~
