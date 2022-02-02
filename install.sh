#!/bin/bash

apt-get install golang-any clamav clamav-freshclam clamav-daemon

cp ./clamrest /usr/bin/clamav-rest

sed -i 's/^Foreground .*$/Foreground true/g' /etc/clamav/clamd.conf
sed -i '/^#Please/a TCPSocket 3310' /etc/clamav/clamd.conf
sed -i 's/^Foreground .*$/Foreground true/g' /etc/clamav/freshclam.conf

service clamav-daemon stop
service clamav-freshclam stop
freshclam --quiet
service clamav-daemon start
service clamav-freshclam start

mv ./service/entrypointClamAV.sh /usr/bin/
chmod +x /usr/bin/entrypointClamAV.sh
chmod +x /usr/bin/clamav-rest
mv ./service/entrypointClamAV.service /usr/lib/systemd/system/

systemctl daemon-reload
systemctl enable entrypointClamAV.service
systemctl start entrypointClamAV.service
