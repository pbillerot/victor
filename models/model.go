package models

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/yaml.v2"
)

// Config de config.yaml
var Config AppConfig

// La liste des répertoires et fichiers du contenu du site Hugo
var hugo []HugoFile

// HugoFile propriétés d'un fichier dans le dossier hugoDir
type HugoFile struct {
	Action       string
	Base         string
	Categories   string
	Content      string
	Date         string
	DateExpiry   string
	DatePublish  string
	Dir          string
	Draft        string
	Expired      bool   // page dont la date a expirée
	Ext          string // extension du fichier
	HugoPath     string // path de la page sur le site /exposant/expostants.md -> /exposant/exposant
	Inline       bool   // page en ligne et visible
	IsAudio      bool
	IsDir        bool
	IsDrawio     bool
	IsExcel      bool
	IsImage      bool
	IsMarkdown   bool
	IsPdf        bool
	IsPowerpoint bool
	IsSystem     bool
	IsText       bool
	IsWord       bool
	Mode         string // Codemirror mode yaml-frontmatter json conf
	Order        int
	Path         string
	PathAbsolu   string
	PathReal     string // path réel du fichier Ex. /data/config.yaml
	Planified    bool   // page qui sera en ligne bientôt
	Prefix       string
	Rang         int
	Tags         string
	Title        string
	URL          string
}

// HugoFileMeta meta données
type HugoFileMeta struct {
	Title       string   `yaml:"title"`
	Draft       bool     `yaml:"draft"`
	Date        string   `yaml:"date"`
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
	Appname     string
	Version     string
	Title       string
	Description string
	Background  string
	Favicon     string
	Icon        string
	HugoRacine  string // /volshare/foirexpo
	HugoTheme   string
	HelpEditor  string
	Github      string
	// Calculé dans main
	HugoContentDir string // /volshare/foirexpo/content
	HugoPrivateDir string // /volshare/foirexpo/private
	HugoPublicDir  string // /volshare/foirexpo/public
}

// Breadcrumb as
type Breadcrumb struct {
	Base   string
	Path   string
	IsLast bool
}

var err error

// GetFilesFolder retourne la liste des fichiers du <folder>
func GetFilesFolder(folder string) (hugoFiles []HugoFile, err error) {
	// suppression du / à la fin
	hugoFolder := strings.TrimSuffix(Config.HugoContentDir+folder, "/")
	var pis []HugoPathInfo
	err = readFolder(hugoFolder, &pis)
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
	lenPrefixe := len(Config.HugoContentDir)
	path := pathAbsolu[lenPrefixe:]
	if path == "" {
		return
	}
	var content []byte

	record.PathAbsolu = pathAbsolu
	record.Path = path // on enlève la partie hugoDirectory du chemin
	record.Dir = filepath.Dir(path)
	record.Base = filepath.Base(path)
	record.Rang = strings.Count(record.Path, "/")
	if info.IsDir() {
		record.IsDir = true
		if record.Dir == "/" {
			record.Dir += record.Base
		} else {
			record.Dir += "/" + record.Base
		}
	} else {
		record.IsDir = false
		// lecture du fichier
		content, err = ioutil.ReadFile(pathAbsolu)
		if err != nil {
			logs.Error(err)
		}
	}
	record.Ext = filepath.Ext(path)

	record.Order = 9
	if record.IsDir {
		record.Order = 0
	}
	if contains([]string{".md"}, record.Ext) {
		record.IsMarkdown = true
		record.Mode = "yaml-frontmatter"
		record.Order = 1
	}
	if contains([]string{".txt"}, record.Ext) {
		record.IsText = true
		record.Mode = ""
		record.Order = 1
	}
	if strings.Contains(strings.ToLower(record.Base), "dockerfile") {
		record.IsSystem = true
		record.Mode = "dockerfile"
		record.Order = 1
	}
	if contains([]string{".sh"}, record.Ext) {
		record.IsSystem = true
		record.Mode = "shell"
		record.Order = 1
	}
	if contains([]string{".json", ".js"}, record.Ext) {
		record.IsSystem = true
		record.Mode = "javascript"
		record.Order = 1
	}
	if contains([]string{".ini", ".conf", ".properties"}, record.Ext) {
		record.IsSystem = true
		record.Mode = "properties"
		record.Order = 1
	}
	if contains([]string{".yaml", ".toml"}, record.Ext) {
		record.IsSystem = true
		record.Mode = strings.ReplaceAll(record.Ext, ".", "")
		record.Order = 1
	}
	if contains([]string{".jpeg", ".jpg", ".png", ".svg", ".gif"}, record.Ext) {
		if strings.Contains(string(content[:]), "Cmxfile") {
			record.IsDrawio = true
		}
		record.IsImage = true
		record.Order = 1
	}
	if contains([]string{".drawio"}, record.Ext) {
		record.IsDrawio = true
		record.Order = 1
	}
	if contains([]string{".pdf"}, record.Ext) {
		record.IsPdf = true
		record.Order = 1
	}
	if contains([]string{".doc", ".dot", ".docx", ".dotx", ".odt"}, record.Ext) {
		record.IsWord = true
		record.Order = 1
	}
	if contains([]string{".xls", ".xlsx", ".ods"}, record.Ext) {
		record.IsExcel = true
		record.Order = 1
	}
	if contains([]string{".ppt", ".pps", ".pptx", ".ppsx", ".odp"}, record.Ext) {
		record.IsPowerpoint = true
		record.Order = 1
	}
	if contains([]string{".wav", ".mp3", ".ogg", ".wma"}, record.Ext) {
		record.IsAudio = true
		record.Order = 1
	}

	ext := filepath.Ext(path)
	// if record.Base == "config.yaml" && record.Dir == "/site" {
	// 	// le fichier a son clone dans /config/_default
	// 	// lecture du fichier yaml
	// 	content, err := ioutil.ReadFile(pathAbsolu)
	// 	if err != nil {
	// 		logs.Error(err)
	// 	}
	// 	record.Content = string(content[:])
	// 	// Recopie dans /config/_default/
	// 	record.PathReal = strings.Replace(record.PathAbsolu, "/content/site/", "/config/_default/", 1)
	// } else if record.Base == "default.md" && record.Dir == "/site" {
	// 	// le fichier a son clone dans /archetypes/
	// 	// lecture du fichier
	// 	content, err := ioutil.ReadFile(pathAbsolu)
	// 	if err != nil {
	// 		logs.Error(err)
	// 	}
	// 	record.Content = string(content[:])
	// 	// Recopie dans /archetypes/
	// 	record.PathReal = strings.Replace(record.PathAbsolu, "/content/site/", "/archetypes/", 1)
	// } else if ext == ".md" || ext == ".yaml" {
	if record.IsMarkdown || ext == ".yaml" {
		// Extraction des meta entre les --- meta ---
		var meta HugoFileMeta
		err = yaml.Unmarshal(content, &meta)
		if err != nil {
			logs.Error(err)
		}
		record.Title = meta.Title
		record.Date = meta.Date
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
		if record.Base == "_index.md" {
			record.HugoPath = strings.Replace(record.Path, "_index.md", "", 1)
		} else if record.Base == "index.md" {
			record.HugoPath = strings.Replace(record.Path, "index.md", "", 1)
		} else {
			record.HugoPath = strings.Replace(strings.ToLower(record.Path), ".md", "", 1) + "/"
		}
		// maj meta
		// for _, v := range meta.Categories {
		// 	metaCat[v] = append(metaCat[v], id)
		// }
		// for _, v := range meta.Tags {
		// 	metaTag[v] = append(metaTag[v], id)
		// }

	} else if record.IsSystem || record.IsText {
		record.Content = string(content[:])
	}
	return
}

// readFolder retourne la liste des fichiers dans HugoPathInfo
func readFolder(dirname string, info *[]HugoPathInfo) (err error) {
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
			readFolder(dirname+"/"+file.Name(), info)
		}
	}
	return
}

// loadHugo retourne la liste des HugoFile
func loadHugo() {
	hugoFolder := Config.HugoContentDir
	var pis []HugoPathInfo
	err := readFolder(hugoFolder, &pis)
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
	// tri des fichiers text,image,et le reste
	sort.SliceStable(hugoFiles, func(i, j int) bool {
		return hugoFiles[i].Order < hugoFiles[j].Order
	})

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
		if record.IsDir {
			hugoFiles = append(hugoFiles, record)
		}
	}
	// tri des dossiers
	sort.SliceStable(hugoFiles, func(i, j int) bool {
		return hugoFiles[i].Path < hugoFiles[j].Path
	})

	return
}

// HugoReload demande de rechargement de hugo
func HugoReload() {
	hugo = nil
	loadHugo()
	return
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
