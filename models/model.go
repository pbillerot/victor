package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/yaml.v2"
)

// Config de config.yaml
var Config AppConfig

// La liste des répertoires et fichiers du contenu du site Hugo
var hugo []HugoFile

// HugoFile propriétés d'un fichier dans le dossier hugoDir
type HugoFile struct {
	Path        string
	Prefix      string
	PathAbsolu  string
	PathReal    string // path réel du fichier Ex. /data/config.yaml
	Base        string
	Dir         string
	Ext         string
	IsDir       int
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

// HugoFileMeta meta données
type HugoFileMeta struct {
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
	HugoURL     string
	HugoDeploy  string
}

// Breadcrumb as
type Breadcrumb struct {
	Base   string
	Path   string
	IsLast bool
}

func init() {
	// Initialisation de models.Config
	if val, ok := config.String("hugo_racine"); ok == nil {
		Config.HugoRacine = val
	}
	if val, ok := config.String("hugo_url"); ok == nil {
		Config.HugoURL = val
	}
	if val, ok := config.String("hugo_deploy"); ok == nil {
		Config.HugoDeploy = val
	}
	if val, ok := config.String("title"); ok == nil {
		Config.Title = val
	}
	if val, ok := config.String("description"); ok == nil {
		Config.Description = val
	}
	if val, ok := config.String("favicon"); ok == nil {
		Config.Favicon = val
	}
	if val, ok := config.String("icon"); ok == nil {
		Config.Icon = val
	}
	logs.Info("Config", Config)
	// Enregistrement de content an tant que répertoire statis
	beego.SetStaticPath("/content", Config.HugoRacine+"/content")
	loadHugo()
}

// GetFilesFolder retourne la liste des fichiers du <folder>
func GetFilesFolder(folder string) (hugoFiles []HugoFile, err error) {
	// suppression du / à la fin
	hugoFolder := strings.TrimSuffix(Config.HugoRacine+"/content"+folder, "/")
	var pis []HugoPathInfo
	err = readDir(hugoFolder, &pis)
	if err != nil {
		return
	}
	for _, pi := range pis {
		record := fileRecord(hugoFolder, pi.Path, pi.Info)
		// ajout dans hugo
		hugoFiles = append(hugoFiles, record)
	}

	return
}

func fileRecord(hugoContent string, pathAbsolu string, info os.FileInfo) (record HugoFile) {

	// On elève le chemin absolu du path
	lenPrefixe := len(Config.HugoRacine + "/content")
	path := pathAbsolu[lenPrefixe:]
	if path == "" {
		return
	}

	record.PathAbsolu = pathAbsolu
	record.Path = path // on enlève la partie hugoDirectory du chemin
	record.Dir = filepath.Dir(path)
	record.Base = filepath.Base(path)
	if info.IsDir() {
		record.IsDir = 1
		if record.Dir == "/" {
			record.Dir += record.Base
		} else {
			record.Dir += "/" + record.Base
		}
	} else {
		record.IsDir = 0
	}
	record.Ext = filepath.Ext(path)
	record.SRC = fmt.Sprintf("%s/content%s", Config.HugoURL, record.Path)
	// record.URL = fmt.Sprintf("%s/%d", Config.HugoURL) TODO

	ext := filepath.Ext(path)
	if record.Base == "config.yaml" {
		// le fichier a son clone dans /data
		// lecture du fichier yaml
		content, err := ioutil.ReadFile(pathAbsolu)
		if err != nil {
			logs.Error(err)
		}
		record.Content = string(content[:])
		record.PathReal = strings.Replace(record.PathAbsolu, "/content/site/", "/data/", 1)
	} else if ext == ".md" || ext == ".yaml" {
		// lecture des metadata du fichier markdown
		content, err := ioutil.ReadFile(pathAbsolu)
		if err != nil {
			logs.Error(err)
		}
		// Extraction des meta entre les --- meta ---
		var meta HugoFileMeta
		err = yaml.Unmarshal(content, &meta)
		if err != nil {
			logs.Error(err)
		}
		record.Title = meta.Title
		record.Date = meta.Date
		record.Action = meta.Action
		record.DatePublish = meta.DatePublish
		record.DateExpiry = meta.DateExpiry
		record.Inline = true
		if record.DatePublish != "" && record.DatePublish > time.Now().Format("2006-01-02") {
			record.Inline = false
			record.Planified = true
			// metaPlanified = append(metaPlanified, id)
		}
		if record.DateExpiry != "" && record.DateExpiry <= time.Now().Format("2006-01-02") {
			record.Inline = false
			record.Expired = true
			// metaExpired = append(metaExpired, id)
		}
		if meta.Draft {
			record.Draft = "1"
			record.Inline = false
			// metaDraft = append(metaDraft, id)
		} else {
			record.Draft = "0"
		}
		record.Tags = strings.Join(meta.Tags, ",")
		record.Categories = strings.Join(meta.Categories, ",")
		record.Content = string(content[:])
		record.HugoPath = strings.Replace(record.Path, ".md", "", 1)
		// maj meta
		// for _, v := range meta.Categories {
		// 	metaCat[v] = append(metaCat[v], id)
		// }
		// for _, v := range meta.Tags {
		// 	metaTag[v] = append(metaTag[v], id)
		// }

	}
	return
}

// readDir retourne la liste des fichiers dans HugoPathInfo
func readDir(dirname string, info *[]HugoPathInfo) (err error) {
	// ouverture du dossier
	f, err := os.Open(dirname)
	if err != nil {
		return
	}
	// lecture ds fichiers et dossiers du dossier courant
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return
	}
	// tri des dossiers sur le nom inversé si numérique
	sort.Slice(list, func(i, j int) bool {
		if _, err := strconv.Atoi(list[i].Name()); err == nil {
			if _, err := strconv.Atoi(list[j].Name()); err == nil {
				return list[i].Name() > list[j].Name()
			}
			return list[i].Name() < list[j].Name()
		}
		return list[i].Name() < list[j].Name()
	})
	// // tri des fichiers sur le nom
	// sort.Slice(list, func(i, j int) bool {
	// 	return list[i].Name() < list[j].Name()
	// })
	// Rangement des dossiers au début
	for _, file := range list {
		if file.IsDir() {
			var pi HugoPathInfo
			pi.Path = dirname + "/" + file.Name()
			pi.Info = file
			*info = append(*info, pi)
		}
	}
	// Rangement des fichiers à la fin
	for _, file := range list {
		if !file.IsDir() {
			var pi HugoPathInfo
			pi.Path = dirname + "/" + file.Name()
			pi.Info = file
			*info = append(*info, pi)
		}
	}
	// réentrance sur les sous-répertoires
	for _, file := range list {
		if file.IsDir() {
			var pi HugoPathInfo
			pi.Path = dirname + "/" + file.Name()
			pi.Info = file
			// *info = append(*info, pi)
			// appel récursif des répertoires
			readDir(dirname+"/"+file.Name(), info)
		}
	}
	return
}

// loadHugo retourne la liste des HugoFile
func loadHugo() {
	hugoFolder := Config.HugoRacine + "/content"
	var pis []HugoPathInfo
	err := readDir(hugoFolder, &pis)
	if err != nil {
		return
	}
	for _, pi := range pis {
		record := fileRecord(hugoFolder, pi.Path, pi.Info)
		// ajout dans hugo
		hugo = append(hugo, record)
	}
	logs.Info("Hugo", Config.HugoRacine, "rechargé")
	return
}

// HugoGetFolder return les HugoFile correspondant au folder
func HugoGetFolder(folder string) (hugoFiles []HugoFile) {
	if hugo == nil {
		loadHugo()
	}
	qSlashMax := strings.Count(folder, "/") + 1
	for _, record := range hugo {
		qSlash := strings.Count(record.Path, "/")
		if folder == "/" {
			if qSlash == 1 {
				hugoFiles = append(hugoFiles, record)
			}
		} else if strings.HasPrefix(record.Dir, folder) && record.Path != folder && qSlash <= qSlashMax {
			hugoFiles = append(hugoFiles, record)
		}
	}
	return
}

// HugoGetRecord return le HugoFile correspondant au path
func HugoGetRecord(path string) (hugoFile HugoFile) {
	if hugo == nil {
		loadHugo()
	}
	for _, record := range hugo {
		if record.Path == path {
			hugoFile = record
		}
	}
	return
}

// HugoGetFolders return seulement les répertoires de Hugo
func HugoGetFolders() (hugoFiles []HugoFile) {
	if hugo == nil {
		loadHugo()
	}
	for _, record := range hugo {
		if record.IsDir == 1 {
			hugoFiles = append(hugoFiles, record)
		}
	}
	return
}

// reloadHugo demande de rechargement de hugo
func reloadHugo() {
	hugo = nil
	loadHugo()
}
