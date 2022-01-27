go build clamrest.go
cp clamrest /usr/bin
systemctl restart entrypointClamAV.service
