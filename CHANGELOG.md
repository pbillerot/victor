---
# lien symbolyque victor/CHANGELOG.md
title: "Quoi de neuf dans Victor ?"
date: 2021-04-08
draft: false
tags:
- technique
categories:
- news
style: bee-doc
_build:
   list: false
---
<!--more-->
### À venir :
- [ ] doc sauvegarde data - docker
- [ ] Structure du contenu - content/site/config.yaml -> ../config.yaml
- [ ] Doc Environnements Test et Production

2.9.1 du 3 mai 2021
----------------------
- `added` label menu déploiement

2.9.0 du 30 avril 2021
----------------------
- `added` menu déploiement sur site externe via un script /bin/sh
- `fixed` exit si erreur config

2.8.6 du 26 avril 2021
----------------------
- `fixed` pb ctx inconnu lors install container

2.8.5 du 26 avril 2021
----------------------
- `alert` le nom des répertoires racine du site doivent être différent de baseURL (bug HUGO)

2.8.3 du 25 avril 2021
----------------------
- `fixed` ajout cacheDir /tmp

2.8.2 du 25 avril 2021
----------------------
- `changed` nettoyage conf
- `fixed` audio non sélectable

2.8.1 du 23 avril 2021
----------------------
- `changed` hugo.yaml avec theme et themehelp
- `added` creation automatique lien /document/site/config.yaml config.yaml

2.8.0 du 22 avril 2021
----------------------
- `added` fichier /conf/hugo.yaml pour indiquer les répertoires à administrer par Victor
- `added` fichier /conf/ctx.yaml fichier contexte interne à victor (utilisation de viper)
- `removed` conf/custom.conf n'est plus utile

2.7.0 du 19 avril 2021
----------------------
- `removed` référentiel de l'aide en ligne - Répertoire help alimenté par ../victor-doc/public
- `changed` maj du thème via le git ou submodule

2.6.4 du 18 avril 2021
----------------------
- avant suppression /victor

2.6.3 du 15 avril 2021
----------------------
- `changed` Doc install
- `added` Editeur étendu yaml toml json dockerfile conf sh ini properties
- `changed` intégration beedream 1.2.0

2.6.2 du 11 avril 2021
---------------------------
- `changed` intégration beedream 1.1.8

2.6.2 du 11 avril 2021
---------------------------
- `changed` intégration beedream 1.1.7

2.6.1 du 11 avril 2021
---------------------------
- `changed` path /victorhelp au lieu de /help

2.6.0 du 10 avril 2021
---------------------------
- `added` Gestion du changelog
- `added` aide intégrée dans la webapp /help
- `fixed` le répertoire hugo peut maintenant être relatif à la webapp Victor

###### Types de changements:
`added` *pour les nouvelles fonctionnalités.*  
`changed` *pour les changements aux fonctionnalités préexistantes.*  
`deprecated` *pour les fonctionnalités qui seront bientôt supprimées*.  
`removed` *pour les fonctionnalités désormais supprimées.*  
`fixed` *pour les corrections de bugs.*  
`security` *en cas de vulnérabilités.*  
