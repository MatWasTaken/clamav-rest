# Clamav Rest

API REST for ClamAV, in GoLang

This API is forked from Niilo's : https://github.com/niilo/clamav-rest that contains a Docker image. The whole docker part was deleted so that we just keep the Go API updated with an installer.

## **Installation** :
Execute install.sh to install

This script creates and launches a systemctl service "entrypointClamAV.service" that listens to the port 9000 to scan files in input.

## **Usage** :

This API contains 2 functions: A scan that...(drumrolls) scans the file you send and a quarantine that uploads the file in a quarantine folder.

### Scan 

To scan a file with this API, here is an example of a cURL POST request:

`$ curl -i -X POST -F FILES=@./eicar3.com 172.16.1.100:9000/scan`

The API returns : 
- An http code : 406 if the file is infected, or else 200*.
- In JSON, the value "FOUND" is affected to the key "Status" if infected, or "OK" if not.
- In JSON, in the key "Description", the virus description if the file is infected.


**Infected** :

```
HTTP/1.1 406 Not Acceptable
Content-Type: application/json; charset=utf-8
Date: Thu, 27 Jan 2022 09:17:50 GMT
Content-Length: 56

{"Status":"FOUND","Description":"Win.Test.EICAR_HDB-1"}
```
**Safe** :

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 27 Jan 2022 12:37:26 GMT
Content-Length: 33

{"Status":"OK","Description":""}
```

### Quarantine

This function uploads a file in a quarantine folder :

`$ curl -i -X POST -F FILES=@./eicar3.com 172.16.1.100:9000/update`

The file will be moved to /home/web.app/data.clamav/quarantine/date-of-upload/ entitled with "hour-of-upload-filename". (hour of upload : 24h format i.e. 15h05)

The API returns the header (http code 200 if no error) and "uploaded file:eicar3.com;length:68"

*Different HTTP codes that the scan returns:

**Status codes:**

- 200 - clean file = no KNOWN infections
- 406 - INFECTED
- 400 - ClamAV returned general error for file
- 412 - unable to parse file
- 501 - unknown request
