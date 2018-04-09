#!/bin/sh

APP_PATH="/opt/ssws/"
BIN="ssws"
SERVICE="sswsservice"
FULL_PATH="${APP_PATH}${BIN}"

# prepare dir and service user
if [ -f $FULLPATH ]; then
    sudo systemctl stop $BIN 
fi

sudo mkdir -p ${APP_PATH}static
sudo mkdir -p ${APP_PATH}certs

sudo useradd $SERVICE -s /sbin/nologin -M

# add as a proper service with logging
sudo cp ${BIN}.service /lib/systemd/system/${BIN}.service
sudo cp 30-${BIN}.conf /etc/rsyslog.d/30-${BIN}.conf
sudo chmod 755 /lib/systemd/system/${BIN}.service

# copy bin & config
cd .. && go build && sudo cp ${BIN} $FULL_PATH
sudo chown $SERVICE -R $APP_PATH
sudo setcap 'cap_net_bind_service=+ep' $FULL_PATH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# start her up.
sudo systemctl enable ${BIN}.service
sudo systemctl restart rsyslog 
sudo systemctl start ${BIN}
