#!/bin/bash

apt-get install golang-any clamav clamav-freshclam clamav-daemon clamav-libunrar

sed -i 's/^#Foreground .*$/Foreground true/g' /etc/clamav/clamd.conf
sed -i 's/^#TCPSocket .*$/TCPSocket 3310/g' /etc/clamav/clamd.conf 
sed -i 's/^#Foreground .*$/Foreground true/g' /etc/clamav/freshclam.conf
freshclam --quiet

iptables -A INPUT -p tcp --dport 9000 -j ACCEPT

mv ./clamrest /usr/bin/
mv ./service/entrypointClamAV.sh /usr/bin/
chmod +x /usr/bin/entrypointClamAV.sh
mv ./service/entrypointClamAV.service /usr/lib/systemd/system/

systemctl daemon-restart
systemctl enable entrypointClamAV.service
systemctl start entrypointClamAV.service

