
# CHANGELOG

### À venir :
- activation de la recherche

3.0.4 du 28 juin 2021
----------------------
- `changed` aide en ligne actualisée

3.0.3 du 24 juin 2021
----------------------
- `changed` maj extension + go mod tidy 

3.0.2 du 17 juin 2021
----------------------
- `fixed` correction du titre de la modal déployer
- `added` lien changelog dans news
- `added` doc goaccess

3.0.1 du 22 mai 2021
----------------------
- `changed` eddy : completion shortcode - mémo curseur - aide des commandes 

3.0.0 du 10 mai 2021
----------------------
- `changed` Gestion d'un contexte hugo par session pour permettre le travail à plusieurs sur la même instance Victor sur des webapp Hugo différentes.
- `removed` suppression gestion ctx viper - n'est plus utile car ctx géré dans session

2.9.5 du 9 mai 2021
----------------------
- `changed` Fin version mono session

2.9.4 du 7 mai 2021
----------------------
- `changed` Documentation terminée

2.9.3 du 6 mai 2021
----------------------
- `changed` Documentation bien avancée

2.9.2 du 4 mai 2021
----------------------
- `fixed` téléchargement sur adresse incomplète (/content)

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
