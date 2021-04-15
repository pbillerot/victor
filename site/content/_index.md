---
title: "Guide d'utilisation de Victor"
date: 2021-04-07
draft: false
categories:
tags:
cover:
---
<!--more-->
{{< diaporama >}}
{{< image image="/site/captures/page-site.png" position="droite" taille="m" >}}

**VICTOR** est une application pour gérer le contenu d'un site [HUGO](https://gohugo.io/).  
Un site **Hugo** n'est constituée que de fichiers statiques et par conséquent sans utilisation de base de données.

**Victor** est une application WEB, il va permettre de créer, modifier, supprimer les fichiers du site directement en ligne. 

En gros **Victor** est un gestionnaire de fichiers sur le web avec la particularités de pouvoir générer publier directement le site hugo en ligne. Il permettra de gérer 2 environnements : un environnement de test et l'environnement de production.

L'objectif de ce guide est de vous présenter :

- comment installer l'application sur un serveur dans un container [DOCKER](https://fr.wikipedia.org/wiki/Docker_(logiciel))
- l'organisation des répertoires du site
- le gestionnaire de fichiers
- l'éditeur de contenu et de configuration du site
- l'éditeur d'image
- l'outil de dessin [DRAWIO](https://github.com/jgraph/drawio-integration) qui a été intégré dans l'application
- les fonctions de publication du site

