go build clamrest.go
mv clamrest /usr/bin
systemctl restart entrypointClamAV.service
echo "ClamAV-Rest API has been updated!"
