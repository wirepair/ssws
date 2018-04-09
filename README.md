# Simple Static Web Service
This is a super simple static web service that includes lets encrypt support. It's meant to be installed in Ubuntu (see the install.sh script in /etc). 
Which will create a service user, a service, allow it to run on port 80 and 443 and log to syslog. You can of course just run it by itself. 

It will redirect all GET/HEAD requests from 80 to 443 using the autocert handler. If you want to add custom handler funcs to the HTTPS service, 
add them in the addRoutes function.

## Install
`
go get github.com/wirepair/ssws
cd $GOHOME/src/github.com/wirepair/ssws
vi main.go // replace hostname with your hostname
:q
cd etc
chmod +x install.sh && ./install.sh (will prompt for passwd)
`
