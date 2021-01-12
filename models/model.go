package models

import "os"

// HugoFile propriétés d'un fichier dans le répertoire hugoDir
type HugoFile struct {
	ID          int
	Key         string
	Root        string
	Path        string
	Prefix      string
	PathAbsolu  string
	PathReal    string // path réel du fichier Ex. /data/config.yaml
	Base        string
	Dir         string
	Ext         string
	IsDir       int
	Level       int
	Title       string
	Draft       string
	Date        string
	Action      string
	DatePublish string
	DateExpiry  string
	Inline      bool // page en ligne et visible
	Planified   bool // page qui sera en ligne bientôt
	Expired     bool // page dont la date a expirée
	Tags        string
	Categories  string
	Content     string
	HugoPath    string // path de la page sur le site /exposant/expostants.md -> /exposant/exposant
	URL         string
	SRC         string
}

// HugoMeta meta données
type HugoMeta struct {
	Title       string   `yaml:"title"`
	Draft       bool     `yaml:"draft"`
	Date        string   `yaml:"date"`
	Action      string   `yaml:"action"`
	DatePublish string   `yaml:"publishDate"`
	DateExpiry  string   `yaml:"expiryDate"`
	Tags        []string `yaml:"tags"`
	Categories  []string `yaml:"categories"`
}

// HugoPathInfo as
type HugoPathInfo struct {
	Path string
	Info os.FileInfo
}

// AppConfig structure du fichier de configuration de l'application app.conf
type AppConfig struct {
	Title       string
	Description string
	Version     string
	Favicon     string
	Icon        string
	HugoRacine  string
	HugoDeploy  string
}
