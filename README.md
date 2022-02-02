# Clamav Rest

API REST pour ClamAV, en GoLang

Cette API est un fork de celle de niilo : https://github.com/niilo/clamav-rest qui contient une image docker, toute la partie docker a été retirée pour n'avoir que l'API en go seule.

## **Installation** :
Pour installer cette API, executer le script install.sh

Ce script crée et démarre un service nommé entrypointClamAV.service qui écoutera sur le port 9000 les input de fichier pour les scanner.


## **Utilisation** :

2 Fonctions sont présentes dans cette API : Le scan de fichier ainsi que la mise en quarantaine.

### Scan 

Pour scanner un fichier à l'aide de cette API, voici un exemple d'appel cURL:

`$ curl -i -X POST -F FILES=@./eicar3.com http://clamav.atgpedi.net:9000/scan`

L'API renvoi : 
- Un code http : 406 si le fichier est vérolé, 200 sinon*.
- En JSON La valeur de la clé "Status" FOUND si vérolé, OK sinon.
- En JSON La valeur de la clé "Description", la description du virus si détecté.


**Vérolé** :

```
HTTP/1.1 406 Not Acceptable
Content-Type: application/json; charset=utf-8
Date: Thu, 27 Jan 2022 09:17:50 GMT
Content-Length: 56

{"Status":"FOUND","Description":"Win.Test.EICAR_HDB-1"}
```
**Non vérolé** :

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 27 Jan 2022 12:37:26 GMT
Content-Length: 33

{"Status":"OK","Description":""}
```

### Quarantaine

Cette fonction permet d'upload un fichier en quarantaine, voici un exemple :

`$ curl -i -X POST -F FILES=@./eicar3.com http://clamav.atgpedi.net:9000/update`

Le fichier en question sera déplacé dans le dossier '/home/web.app/data.clamav/quarantine/date-du-jour/' avec comme nom l'heure au format '15:04:05'-nom_du_fichier.

L'API renvoi alors le header (code http:200 si pas d'erreur) ainsi que "uploaded file:eicar3.com;length:68".


*Les différents codes HTTP que peut retourner l'API lors d'un scan :

**HTTP Status codes:**

| code HTTP | Description |
| ------ | ------ |
| 200 | clean file = no KNOWN infections |
| 406 | INFECTED |
| 400 | ClamAV returned general error for file |
| 412 | unable to parse file |
| 501 | unknown request |
