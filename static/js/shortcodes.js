/**
 * Mots clés des shortcodes
 */

 var $shortcodes = [
  { text: '{{< bouton [icone="nom icone"] [image="image"] titre="titre" [lien="lien"] >}}', displayText: 'bouton' },
  { text: '{{< cartes taille="s m l xl" >}}\n{{< carte image="image" [diapo="diapo"] [titre="titre"] [lien="lien"] [pdf="pdf"] [taille="s m l xl"] >}}\ntexte...\n{{< /carte >}}\n{{< /cartes >}}', displayText: 'carte' },
  { text: '{{< colonnes [two three]>}}\n...\n<--->\n...\n{{< /colonnes >}}', displayText: 'colonne' },
  { text: '{{<centre>}}...{{</centre>}}', displayText: 'colonne' },
  { text: '{{<diaporama>}}', displayText: 'diaporama' },
  { text: '{{< galerie [sous-repertoire] />}', displayText: 'galerie' },
  { text: '{{< icone nomIcone >}}', displayText: 'icone' },
  { text: '{{< image image="/..." [lien="lien"] [position="gauche droite"] [taille="s m l xl]" [forme="ronde"] >}}', displayText: 'image' },
  { text: '{{< label label="libellé" [icone=""] [lien="lien"] >}}', displayText: 'label' },
  { text: '[Je suis un lien](https://www.google.com)', displayText: 'lien' },
  { text: '{{< message [info success warning error] >}}\nmessage\n{{< /message >}}', displayText: 'message' },
  { text: '{{< player "lien_vers_audio" ["titre"] ["boucle"] >}}', displayText: 'player' },
  { text: '{{< players répertoire >}}', displayText: 'players' },
  { text: '{{< repertoire "répertoire" >}}', displayText: 'repertoire' },
  { text: '| Tables         | Sont            | Frais  |\n  | -------------- |:---------------:| ------:|\n  | col 3 est      | aligné à droite | 1600 € |\n  | col 2 est      | centré          |   12 € |\n  | lignes zébrées | harmonieuses    |    1 € |', displayText: 'table' },
  { text: '{{<toc>}}', displayText: 'table des matières toc' },
  { text: '{{< texte [xs s m l xl] [rouge orange vert cyan bleu pourpre violet marron gris] >}}...{{< /texte >}}', displayText: 'texte' },
 ]
