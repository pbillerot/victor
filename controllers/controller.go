package controllers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/models"
	"gopkg.in/yaml.v2"
)

// MainController as
type MainController struct {
	beego.Controller
}

// Get as
func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.html"
}

// Prepare implements Prepare method for loggedRouter.
func (c *MainController) Prepare() {

	// Contexte lié à app.conf
	c.Data["Config"] = AppConfig

	// XSRF protection des formulaires
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	// Sera ajouté derrière les urls pour ne pas utiliser le cache des images dynamiques
	c.Data["Composter"] = time.Now().Unix()
}

// Liste des répertoires et fichiers du répertoire hugoFiles
var hugoFiles []models.HugoFile
var hugoRacine string
var metaTag map[string][]int
var metaCat map[string][]int
var metaDraft []int
var metaPlanified []int
var metaExpired []int

// AppConfig de config.yaml
var AppConfig models.AppConfig

var err error

// List Liste des fichiers et répertoires
func (c *MainController) List() {
	appid := c.Ctx.Input.Param(":app")
	dirid := c.Ctx.Input.Param(":dir")
	baseid := c.Ctx.Input.Param(":base")
	hugoRacine = c.Data["DataDir"].(string) + "/content"

	// RECHERCHE DANS LA VUE
	search := strings.ToLower(c.GetString("search"))
	ctxSearch := fmt.Sprintf("hugo-%s-search", appid)
	if strings.ToLower(c.GetString("searchstop")) != "" {
		c.DelSession(ctxSearch)
		search = ""
	}
	if search != "" {
		c.SetSession(ctxSearch, search)
	} else {
		if c.GetSession(ctxSearch) != nil {
			search = c.GetSession(ctxSearch).(string)
		}
	}

	var recordFiltered []models.HugoFile
	if search != "" {
		// Chargement des répertoires et fichiers
		if len(hugoFiles) == 0 {
			hugoDirectoryRecord(c, c.Data["DataDir"].(string))
		}

		for _, record := range hugoFiles {
			if strings.Contains(search, "draft") && containsInt(metaDraft, record.ID) {
				recordFiltered = append(recordFiltered, record)
			}
			if strings.Contains(search, "planif") && containsInt(metaPlanified, record.ID) {
				recordFiltered = append(recordFiltered, record)
			}
			if strings.Contains(search, "expir") && containsInt(metaExpired, record.ID) {
				recordFiltered = append(recordFiltered, record)
			}
			if meta, ok := metaTag[search]; ok {
				if containsInt(meta, record.ID) {
					recordFiltered = append(recordFiltered, record)
				}
			}
			if meta, ok := metaCat[search]; ok {
				if containsInt(meta, record.ID) {
					recordFiltered = append(recordFiltered, record)
				}
			}
		}
	} else {
		if len(hugoFiles) == 0 {
			hugoDirectoryRecord(c, c.Data["DataDir"].(string))
		}
	} // end if search

	// ********************************************************************
	// Remplissage du contexte pour le template
	c.Data["DirId"] = dirid
	c.Data["BaseId"] = baseid
	c.Data["Search"] = search
	if search == "" {
		c.Data["Records"] = hugoFiles
	} else {
		c.Data["Records"] = recordFiltered
	}
	c.Ctx.Output.Cookie("from", fmt.Sprintf("/bee/hugo/list/%s", appid))
	c.TplName = "hugo_list.html"
}

// Image Visualiser Modifier une image
func (c *MainController) Image() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	for _, rec := range hugoFiles {
		if rec.Key == keyid {
			record = rec
			break
		}
	}
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	if c.Ctx.Input.Method() == "POST" {
		// ENREGISTREMENT DE L'IMAGE
		simage := c.GetString("image")
		b64data := simage[strings.IndexByte(simage, ',')+1:]
		unbased, err := base64.StdEncoding.DecodeString(b64data)
		// img, _, err := image.Decode(bytes.NewReader([]byte(element.SQLout)))
		if err != nil {
			msg := fmt.Sprintf("HugoImage %s : %s", record.PathAbsolu, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}

		err = ioutil.WriteFile(record.PathAbsolu, unbased, 0644)
		// outputFile, err := os.Create(record.PathAbsolu)
		if err != nil {
			msg := fmt.Sprintf("HugoImage %s : %s", record.PathAbsolu, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
		// defer outputFile.Close()

		// outputFile.Write(unbased)
		// Fermeture de la fenêtre
		c.TplName = "bee_close.html"
		return
	}

	// Remplissage du contexte pour le template
	c.Data["Record"] = record
	c.Data["KeyID"] = keyid
	c.Ctx.Output.Cookie("hugo-"+appid, keyid)
	c.Ctx.Output.Cookie("from", fmt.Sprintf("/bee/hugo/image/%s", appid))
	c.TplName = "hugo_image.html"
}

// Pdf Visualiser Modifier une image
func (c *MainController) Pdf() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	for _, rec := range hugoFiles {
		if rec.Key == keyid {
			record = rec
			break
		}
	}
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	// Remplissage du contexte pour le template
	c.Data["Record"] = record
	c.Data["KeyID"] = keyid
	c.Ctx.Output.Cookie("hugo-"+appid, keyid)
	c.Ctx.Output.Cookie("from", fmt.Sprintf("/bee/hugo/image/%s", appid))
	c.TplName = "hugo_pdf.html"
}

// Document Visualiser Modifier un document
func (c *MainController) Document() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	for _, rec := range hugoFiles {
		if rec.Key == keyid {
			record = rec
			break
		}
	}
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	if c.Ctx.Input.Method() == "POST" {
		// ENREGISTREMENT DU DOCUMENT
		document := c.GetString("document")
		if record.PathReal != "" {
			err = ioutil.WriteFile(record.PathReal, []byte(document), 0644)
			if err != nil {
				msg := fmt.Sprintf("hugoFileument %s : %s", record.PathReal, err)
				logs.Error(msg)
				flash.Error(msg)
				flash.Store(&c.Controller)
				c.Ctx.Redirect(302, "/")
				return
			}
			err = ioutil.WriteFile(record.PathAbsolu, []byte(document), 0644)
		} else {
			err = ioutil.WriteFile(record.PathAbsolu, []byte(document), 0644)
		}
		if err != nil {
			msg := fmt.Sprintf("hugoFileument %s : %s", record.PathAbsolu, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}

		// Vidage de Hugo puis reconstruction
		hugoFiles = nil
		hugoDirectoryRecord(c, c.Data["DataDir"].(string))
		for _, rec := range hugoFiles {
			if rec.Key == keyid {
				record = rec
				break
			}
		}
		// Demande d'actualisation de l'arborescence
		c.Ctx.Output.Cookie("hugo-refresh-"+appid, "true")
	}

	// Remplissage du contexte pour le template
	c.Data["Record"] = record
	c.Data["KeyID"] = keyid
	c.Ctx.Output.Cookie("hugo-"+appid, keyid)
	c.Ctx.Output.Cookie("from", fmt.Sprintf("/bee/hugo/document/%s", appid))
	c.TplName = "hugo_document.html"
}

// File Gestion du fichier
func (c *MainController) File() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	for _, rec := range hugoFiles {
		if rec.Key == keyid {
			record = rec
			break
		}
	}
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	// Remplissage du contexte pour le template
	c.Data["Record"] = record
	c.Data["KeyID"] = keyid

	c.Ctx.Output.Cookie("from", fmt.Sprintf("/bee/hugo/file/%s/%s", appid, keyid))
	c.TplName = "hugo_file.html"
}

// Directory Gestion du répertoire
func (c *MainController) Directory() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	for _, rec := range hugoFiles {
		if rec.Key == keyid {
			record = rec
			break
		}
	}
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	// Remplissage du contexte pour le template
	c.Data["Record"] = record
	c.Data["KeyID"] = keyid

	c.Ctx.Output.Cookie("from", fmt.Sprintf("/bee/hugo/directory/%s/%s", appid, keyid))
	c.TplName = "hugo_directory.html"
}

// FileMv Renommer le fichier
func (c *MainController) FileMv() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	for _, rec := range hugoFiles {
		if rec.Key == keyid {
			record = rec
			break
		}
	}
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	newFile := c.GetString("new_name")
	if record.IsDir == 1 {
		err = os.Rename(record.PathAbsolu, hugoRacine+newFile)
		if err != nil {
			msg := fmt.Sprintf("HugoFileMv %s : %s", record.Path, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
	} else {
		// Copie du fichier sur la cible
		data, err := ioutil.ReadFile(record.PathAbsolu)
		if err != nil {
			msg := fmt.Sprintf("HugoFileMv %s : %s", record.Path, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
		err = ioutil.WriteFile(hugoRacine+newFile, data, 0644)
		if err != nil {
			msg := fmt.Sprintf("HugoFileCp %s : %s", newFile, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
		// Suppression du fichier source
		err = os.RemoveAll(record.PathAbsolu)
		if err != nil {
			msg := fmt.Sprintf("HugoFileMv %s : %s", record.Path, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
	}

	// Le cookie ancrage est déplacé sur le répertoire root
	c.Ctx.Output.Cookie("hugo-"+appid, record.Root)
	// Demande d'actualisation de l'arborescence
	c.Ctx.Output.Cookie("hugo-refresh-"+appid, "true")
	// Vidage de Hugo pour reconstruction
	hugoFiles = nil

	// Fermeture de la fenêtre
	c.TplName = "bee_close.html"
	return

}

// FileCp Recopier le fichier ou répertoire
func (c *MainController) FileCp() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	for _, rec := range hugoFiles {
		if rec.Key == keyid {
			record = rec
			break
		}
	}
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	newFile := c.GetString("copy_file")
	data, err := ioutil.ReadFile(record.PathAbsolu)
	if err != nil {
		msg := fmt.Sprintf("HugoFileCp %s : %s", record.Path, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}
	err = ioutil.WriteFile(hugoRacine+"/"+newFile, data, 0644)
	if err != nil {
		msg := fmt.Sprintf("HugoFileCp %s : %s", newFile, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	// Vidage de Hugo pour reconstruction
	hugoFiles = nil

	// Demande d'actualisation de l'arborescence
	c.Ctx.Output.Cookie("hugo-refresh-"+appid, "true")
	// Fermeture de la fenêtre
	c.TplName = "bee_close.html"
	return
}

// FileNew Nouveau document à partir du modele.md
func (c *MainController) FileNew() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	for _, rec := range hugoFiles {
		if rec.Key == keyid {
			record = rec
			break
		}
	}
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	newFile := c.GetString("new_file")
	data, err := ioutil.ReadFile(hugoRacine + "/site/modele.md")
	if err != nil {
		msg := fmt.Sprintf("HugoFileNew %s : %s", record.Path, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}
	err = ioutil.WriteFile(hugoRacine+"/"+newFile, data, 0644)
	if err != nil {
		msg := fmt.Sprintf("HugoFileNew %s : %s", newFile, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	// Vidage de Hugo pour reconstruction
	hugoFiles = nil

	// Demande d'actualisation de l'arborescence
	c.Ctx.Output.Cookie("hugo-refresh-"+appid, "true")
	// Fermeture de la fenêtre
	c.TplName = "bee_close.html"
	return
}

// FileRm Supprimer le fichier ou répertoire
func (c *MainController) FileRm() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	for _, rec := range hugoFiles {
		if rec.Key == keyid {
			record = rec
			break
		}
	}
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	err = os.RemoveAll(record.PathAbsolu)
	if err != nil {
		msg := fmt.Sprintf("HugoFileRm %s : %s", record.Path, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	// Le cookie ancrage est déplacé sur le répertoire root
	c.Ctx.Output.Cookie("hugo-"+appid, record.Root)

	// Vidage de Hugo pour reconstruction
	hugoFiles = nil

	// Demande d'actualisation de l'arborescence
	c.Ctx.Output.Cookie("hugo-refresh-"+appid, "true")
	// Fermeture de la fenêtre
	c.TplName = "bee_close.html"
	return
}

// FileUpload Charger le fichier sur le serveur
func (c *MainController) FileUpload() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	for _, rec := range hugoFiles {
		if rec.Key == keyid {
			record = rec
			break
		}
	}
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	file, handler, err := c.Ctx.Request.FormFile("new_file")
	if err != nil {
		msg := fmt.Sprintf("HugoFileUpload %s : %s", "new_file", err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		msg := fmt.Sprintf("HugoFileUpload %s : %s", handler.Filename, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}
	filepath := fmt.Sprintf("%s/%s", record.PathAbsolu, handler.Filename)
	err = ioutil.WriteFile(filepath, fileBytes, 0644)
	if err != nil {
		msg := fmt.Sprintf("HugoDirectory %s : %s", handler.Filename, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	// Le cookie ancrage est déplacé sur le répertoire root
	c.Ctx.Output.Cookie("hugo-"+appid, record.Root)

	// Vidage de Hugo pour reconstruction
	hugoFiles = nil

	// Demande d'actualisation de l'arborescence
	c.Ctx.Output.Cookie("hugo-refresh-"+appid, "true")
	// Fermeture de la fenêtre
	c.TplName = "bee_close.html"
	return
}

// FileMkdir Créer un répertoire
func (c *MainController) FileMkdir() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	for _, rec := range hugoFiles {
		if rec.Key == keyid {
			record = rec
			break
		}
	}
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	newDir := c.GetString("new_dir")
	err = os.MkdirAll(hugoRacine+"/"+newDir, 0744)
	if err != nil {
		msg := fmt.Sprintf("HugoFileMkdir %s : %s", newDir, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	// Vidage de Hugo pour reconstruction
	hugoFiles = nil

	// Le cookie ancrage est déplacé sur le répertoire root
	c.Ctx.Output.Cookie("hugo-"+appid, record.Root)

	// Demande d'actualisation de l'arborescence
	c.Ctx.Output.Cookie("hugo-refresh-"+appid, "true")
	// Fermeture de la fenêtre
	c.TplName = "bee_close.html"
	return
}

// Action Action
func (c *MainController) Action() {
	appid := c.Ctx.Input.Param(":app")
	action := c.Ctx.Input.Param(":action")

	switch action {
	case "refreshHugo":
		refreshHugo(c)
	case "publishDev":
		publishDev(c)
	case "pushProd":
		pushProd(c)
	}
	// Demande d'actualisation de l'arborescence
	c.Ctx.Output.Cookie("hugo-refresh-"+appid, "true")
	// Fermeture de la fenêtre
	c.TplName = "bee_close.html"
	return
}

var hugoLinks = []map[string]string{
	{
		"title":  "Site de test",
		"url":    "http://localhost:1313",
		"posx":   "left",
		"target": "hugo_test",
	},
}

func readDir(dirname string, info *[]models.HugoPathInfo) error {
	// ouverture du répertoire
	f, err := os.Open(dirname)
	if err != nil {
		return err
	}
	// lecture ds fichiers et répertoires du répertoire courant
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return err
	}
	// tri des fichiers sur le nom
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name() < list[j].Name()
	})
	// Rangement des fichiers au début
	for _, file := range list {
		if !file.IsDir() {
			var pi models.HugoPathInfo
			pi.Path = dirname + "/" + file.Name()
			pi.Info = file
			*info = append(*info, pi)
		}
	}
	// rangement des répertoires à la fin
	// tri des répertoires sur le nom inversé si numérique
	sort.Slice(list, func(i, j int) bool {
		if _, err := strconv.Atoi(list[i].Name()); err == nil {
			if _, err := strconv.Atoi(list[j].Name()); err == nil {
				return list[i].Name() > list[j].Name()
			}
			return list[i].Name() < list[j].Name()
		}
		return list[i].Name() < list[j].Name()
	})
	for _, file := range list {
		if file.IsDir() {
			var pi models.HugoPathInfo
			pi.Path = dirname + "/" + file.Name()
			pi.Info = file
			*info = append(*info, pi)
			// appel récursif des répertoires
			readDir(dirname+"/"+file.Name(), info)
		}
	}
	return nil
}

/**
 * hugoDirectoryRecord:
 * - lecture des répertoires /content et /data de foirexpo
 *
 **/
func hugoDirectoryRecord(c *MainController, hugoDirectory string) (err error) {

	// raz des meta
	metaTag = make(map[string][]int)
	metaCat = make(map[string][]int)
	metaDraft = []int{}
	metaPlanified = []int{}
	metaExpired = []int{}
	// Lecture des répertoires et insertion d'un record par document
	var id int

	var pis []models.HugoPathInfo
	err = readDir(hugoDirectory, &pis)
	if err != nil {
		return err
	}
	for _, pi := range pis {
		if strings.Contains(pi.Path, hugoRacine) {
			id++
			record := hugoFileRecord(hugoDirectory, pi.Path, pi.Info, id)
			record.ID = id
			if record.Level == 0 {
				continue
			}
			if record.Dir[1:] == record.Base {
				record.Key = record.Base
			} else {
				record.Key = strconv.Itoa(id)
			}
			record.URL = fmt.Sprintf("%s/%d", c.Data["DataUrl"].(string), id)
			record.SRC = fmt.Sprintf("%s/content%s", c.Data["DataUrl"].(string), record.Path)
			// ajout dans hugo
			hugoFiles = append(hugoFiles, record)
		}
	}
	// Update site public Hugo
	// publishDev(c)

	return
}

func hugoFileRecord(hugoDirectory string, pathAbsolu string, info os.FileInfo, id int) (record models.HugoFile) {

	// On elève le chemin absolu du path
	lenPrefixe := len(hugoDirectory + "/content")
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
	islash := strings.Index(record.Dir[1:], "/")
	if islash > 0 {
		record.Root = record.Dir[1 : islash+1]
	} else {
		record.Root = record.Dir[1:]
	}
	record.Level = strings.Count(record.Dir, "/")
	record.Ext = filepath.Ext(path)
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
		var meta models.HugoMeta
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
			metaPlanified = append(metaPlanified, id)
		}
		if record.DateExpiry != "" && record.DateExpiry <= time.Now().Format("2006-01-02") {
			record.Inline = false
			record.Expired = true
			metaExpired = append(metaExpired, id)
		}
		if meta.Draft {
			record.Draft = "1"
			record.Inline = false
			metaDraft = append(metaDraft, id)
		} else {
			record.Draft = "0"
		}
		record.Tags = strings.Join(meta.Tags, ",")
		record.Categories = strings.Join(meta.Categories, ",")
		record.Content = string(content[:])
		record.HugoPath = strings.Replace(record.Path, ".md", "", 1)
		// maj meta
		for _, v := range meta.Categories {
			metaCat[v] = append(metaCat[v], id)
		}
		for _, v := range meta.Tags {
			metaTag[v] = append(metaTag[v], id)
		}

	}
	return
}

// int contains in slice
func containsInt(sl []int, in int) bool {
	for _, v := range sl {
		if v == in {
			return true
		}
	}
	return false
}

// string contains in slice
func containsString(sl []string, in string) bool {
	for _, v := range sl {
		if v == in {
			return true
		}
	}
	return false
}

// publishDev : Exécution du moteur Hugo pour mettre à jour le site de développement
func publishDev(c *MainController) {
	logs.Info("publishDev", c.Data["HugoDev"].(string))
	cmd := exec.Command("hugo", "-d", c.Data["HugoDev"].(string),
		"--environment", "DEV", "--cleanDestinationDir")
	cmd.Dir = c.Data["HugoDir"].(string)
	out, err := cmd.CombinedOutput()
	flash := beego.ReadFromRequest(&c.Controller)
	if err != nil {
		logs.Error("publishDev", err)
		flash.Error("ERREURG Génération des pages : %v", err)
		flash.Store(&c.Controller)
	}
	logs.Info("publishDev", string(out))
}

// pushProd : Exécution du moteur Hugo pour mettre à jour le site de développement
func pushProd(c *MainController) {
	flash := beego.ReadFromRequest(&c.Controller)
	logs.Info("pushProd", c.Data["HugoProd"].(string))

	// Hugo Git push
	cmd := exec.Command("sh", "-c", "./project/git-push.sh")
	cmd.Dir = c.Data["HugoDir"].(string)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logs.Error("pushProd", err)
		flash.Error("ERREUR: pushProd : %v", err)
		flash.Store(&c.Controller)
	}
	logs.Info("pushProd", string(out))
}

// refreshHugo : Exécution du moteur Hugo pour mettre à jour le site de développement
func refreshHugo(c *MainController) {
	logs.Info("refreshHugo")
	hugoFiles = nil
}
