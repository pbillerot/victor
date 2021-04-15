---
title: "Installation"
date: 2021-04-08
draft: false
#categories:
tags:
- technique
cover: "/site/media/installation.jpg"
style: bee-doc
menu:
  page:
    parent: technique
    weight: 10
---
*Guide pour l'administrateur technique*
<!--more-->
{{< toc >}}

Je vous propose d'installer une plateforme complète pour héberger notre application **Victor**.

Nous utiliserons une VM (Machine Virtuelle) [DEBIAN Buster](https://fr.wikipedia.org/wiki/Debian) pour accueillir notre plateforme avec le gestionnaire de conteneur [Docker](https://fr.wikipedia.org/wiki/Docker_(logiciel)) installé.

Ce document ne décrit pas l'installation d'une VM Debian.  
Pour ma part je loue une VPS Debian chez l'hébergeur [OVH](https://www.ovhcloud.com/fr/vps/)

## Prérequis du système hôte

Système : `Debian Buster`
```shell
>uname -a
Linux vps-7d2d773f 4.19.0-16-cloud-amd64 #1 SMP Debian 4.19.181-1 (2021-03-19) x86_64 GNU/Linux
```
Gestionnaire de conteneur Docker :
```shell
>docker version
Client: Docker Engine - Community
 Version:           20.10.5
 API version:       1.41
 Go version:        go1.13.15
 ...
 
Server: Docker Engine - Community
 Engine:
  Version:          20.10.5
  API version:      1.41 (minimum version 1.12)
  Go version:       go1.13.15
  ...
```

```shell
>docker-compose version
docker-compose version 1.21.0, build unknown
docker-py version: 3.4.1
CPython version: 3.7.3
OpenSSL version: OpenSSL 1.1.1d  10 Sep 2019
```


## La plateforme Docker

{{< image image="/technique/docker.png" >}}

Notre plateforme sera composée de 4 containers :

- [Caddy Server](https://caddyserver.com/docs/) le frontal web, c'est l'élément le plus important. Il sera chargé :  
   - de contrôler le trafic http (:80) et https (:443)
   - de renouveller le certificat lié au nom de domaine
   - de gérer les authentifications pour certaines [URI](https://fr.wikipedia.org/wiki/Uniform_Resource_Identifier)
   - de rediriger les flux vers les autres containers en fonction des URI
   - de journaliser les accès et les erreurs
- [Bivouac](https://www.billerot.eu) le container de notre application qui utilisera une image de **Victor**

Pour plus de confort, j'utilise  
- [Portainer](https://korben.info/portainer-io-un-outil-graphique-pour-gerer-vos-environnements-docker-en-toute-securite.html) pour gérer graphiquement l'environnement Docker  
- [Filebrowser](https://filebrowser.org/features) pour manipuler les fichiers du répertoire partagé (volshare)

Les 4 containers ont accès à la même ressource de fichiers `volshare` et les échanges entre **Caddy Server** et les autres containers se feront à travers le réseau privé `web`. Ces containers ne seront pas  accessibles de l'extérieur.

La configuration de **Docker** se fera à travers le fichier `/volshare/dockker/docker-compose.yaml`, 
**Caddy Server**, via `/volshare/docker/caddy/caddyfile.conf`

Nous allons les détailler ci-aprés.

## Volume partagé /volshare

`/volshare` est le répertoire partagé entre tous les containers.

Il aura la structure suivante
```
/volshare
  /logs
    access.log access.0.log ... access.9.log
  /etc
    (les certificats du domaine)
  /data (le répertoire des données à sauvegarder)
    /store
      (le répertoire des fichiers statiques servi par Caddy)
    /filebrowser
      database.db
  /docker (les fichiers de configuration des containers)
    docker-compose.yaml
    /caddy
      caddyfile.conf
    /filebrowser
      filebrowser.conf
    /bivouac
      custom.conf
    /victor
      dockerfile
```

## Container Filebrowser

### /volshare/dockker/docker-compose.yaml

```yaml
  filebrowser:
    image: filebrowser/filebrowser:latest
    container_name: filebrowser
    restart: unless-stopped
    volumes:
    - /volshare:/srv
    - /volshare/data/filebrowser/database.db:/database.db
    - ./filebrowser/filebrowser.json:/.filebrowser.json    
    networks:
    - web
```

### /volshare/docker/caddy/caddyfile.conf

```shell
# filebrowser /fb
redir /fb /fb/
reverse_proxy /fb/* filebrowser:80
```

### /volshare/docker/filebrowser/filebrowser.json

```json
{
  "port": 80,
  "baseURL": "/fb",
  "address": "",
  "log": "stdout",
  "database": "/database.db",
  "root": "/srv",
}
```

## Container Portainer

### /volshare/dockker/docker-compose.yaml

```yaml
  portainer:
    image: portainer/portainer-ce
    container_name: portainer
    command: -H unix:///var/run/docker.sock
    restart: unless-stopped
    volumes:
    - /var/run/docker.sock:/var/run/docker.sock
    networks:
    - web
```

### /volshare/docker/caddy/caddyfile.conf

```shell
# portainer 
# on supprime le préfix /portainer après le routage
redir /portainer /portainer/
route /portainer/* {
    uri strip_prefix /portainer
    reverse_proxy portainer:9000
}
```

## Container Bivouac (Victor)

### /volshare/dockker/docker-compose.yaml

```yaml
  bivouac:
    build:
      context: victor
    image: victor:latest
    container_name: bivouac
    restart: unless-stopped
    user: 1000:1000
    volumes:
    - /volshare:/volshare
    - ./bivouac/custom.conf:/src/victor/conf/custom.conf
    networks:
    - web
```

### /volshare/docker/bivouac/custom.conf

```shell
# custom.conf -  Personnalisation du site déployé

# dev / production pour ne pas afficher les erreurs en détail
runmode = production

# Session
EnableXSRF = true # mettre true en HTTPS

# Titre de l'application
title = "BiVouac Admin"

# Répertoire de la webapp Hugo
hugo_dir = "/volshare/data/bivouac"
```

### /volshare/docker/caddy/caddyfile.conf

```shell
# bivouac
# est l'application d'accueil du site
reverse_proxy /* bivouac:8080

# PROTECTION par mot de passe pour les uri /victor /hugo
# user: admin password: admin (qui a été haché par la commande ci-dessous)
# docker container exec -it caddy /bin/sh
# puis
# caddy hash-password [--plaintext <password>]

basicauth /victor/* {
    admin JDJhJDE0JGRKNmt6a3g5L1BlSXRmbmVWV2RXeU9Lc2NlZzFGUnV6eHFyYlVYOUFGc3FmL3NsSS5zRXdt
}
basicauth /hugo/* {
    admin JDJhJDE0JGRKNmt6a3g5L1BlSXRmbmVWV2RXeU9Lc2NlZzFGUnV6eHFyYlVYOUFGc3FmL3NsSS5zRXdt
}

```

## Image Victor

### /volshare/dockker/victor/dockerfile

```dockerfile
# IMAGE VICTOR

# ETAPE COMPILATION
# Le GOPATH par défaut de cette image est /go.
FROM golang:alpine as goalpine
# Installation de GCC et GIT
RUN apk add build-base git
# Installation de victor
WORKDIR /src
RUN git clone https://github.com/pbillerot/victor.git
WORKDIR /src/victor
# Build avec CGO du binaire victor
RUN GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /src/victor/victor

# ETAPE GENERATION D'UN IMAGE 
# avec le projet
FROM alpine
# Installation de victor...
# l'environnement go ne sera pas installé car victor a été compilé dans l'étape compilation
# ce qui réduit consiérablement la taille de l'image finale
RUN mkdir -p /src/victor
copy --from=goalpine /src/victor /src/victor
# Uptade OS + hugo + git + nano
RUN apk add --update nano hugo git
# Ajout du user 1000
USER 1000:1000

# POINT D'ENTREE
WORKDIR /src/victor
ENTRYPOINT ./victor
# Le port sur lequel notre service écoute
EXPOSE 8080
```

## Container Caddy

### /volshare/dockker/docker-compose.yaml

```yaml
  caddy:
    # https://hub.docker.com/_/caddy?tab=description
    image: caddy:latest
    container_name: caddy
    restart: unless-stopped
    ports:
    - '80:80'
    - '443:443'
    volumes:
    - './caddy/caddyfile.conf:/etc/caddy/Caddyfile'
    - '/volshare/etc:/data'
    - '/volshare/data/store:/srv'
    - '/volshare:/volshare'
    networks:
    - web
```

### /volshare/docker/caddy/caddyfile.conf

```shell
# Configuration du serveur Caddy
# https://caddyserver.com/docs/

# GLOBAL option
# https://www.ssllabs.com/ssltest/analyze.html?d=mon.domaine.com
{
    email mon.email@gmail.com
}

# HOST
mon.domaine.com

# blacklist - 
@blaklist {
    remote_ip 94.130.212.180 134.119.20.10
}
handle @blaklist {
    respond "Refused!" 403
}

# Serveur de fichiers statics
redir /store /store/
handle_path /store/* {
    root * /volshare/data/store
    file_server browse
}

# Log du trafic (rotation automatique tous les 100 Mo (10 logs) 
log {
    output file /volshare/log/access.log
    format single_field common_log
}

```

