package controllers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/beego/beego/v2/core/logs"

	"github.com/beego/beego/v2/server/web"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/models"
	"github.com/pbillerot/victor/shutil"
)

// Main as get and Post
func (c *MainController) Main() {
	// Chargement de hugoFiles et des meta du dossier courant
	hugoFiles := models.HugoGetFolder("/")

	c.Data["Records"] = hugoFiles

	setSession(c, "Folder", "/")

	web.ReadFromRequest(&c.Controller)
	c.TplName = "index.html"
}

// Folder Demande de lister le dossier
func (c *MainController) Folder() {
	pathFolder := "/" + c.Ctx.Input.Param(":path")
	setSession(c, "Folder", pathFolder)

	web.ReadFromRequest(&c.Controller)
	// Chargement de hugoFiles et des meta du dossier courant
	hugoFiles := models.HugoGetFolder(pathFolder)

	c.Data["Records"] = hugoFiles

	c.TplName = "index.html"
}

// Image Visualiser Modifier une image
func (c *MainController) Image() {
	pathFile := "/" + c.Ctx.Input.Param(":path") + "." + c.Ctx.Input.Param(":ext")
	setSession(c, "File", pathFile)

	// Recherche du record
	record := models.HugoGetRecord(pathFile)

	flash := beego.ReadFromRequest(&c.Controller)
	if record.Path == "" {
		msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
	}

	if c.Ctx.Input.Method() == "POST" {
		// ENREGISTREMENT DE L'IMAGE
		simage := c.GetString("image")
		b64data := simage[strings.IndexByte(simage, ',')+1:]
		unbased, err := base64.StdEncoding.DecodeString(b64data)
		// img, _, err := image.Decode(bytes.NewReader([]byte(element.SQLout)))
		if err != nil {
			msg := fmt.Sprintf("HugoImage %s : %s", pathFile, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
		pathAbsolu := models.Config.HugoRacine + "/content" + pathFile
		err = ioutil.WriteFile(pathAbsolu, unbased, 0644)
		if err != nil {
			msg := fmt.Sprintf("HugoImage %s : %s", pathAbsolu, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
		models.HugoReload()
	}

	// Load Folder
	pathFolder := c.GetSession("Folder").(string)
	hugoFiles := models.HugoGetFolder(pathFolder)

	// Remplissage du contexte pour le template
	c.Data["Record"] = record
	c.Data["Records"] = hugoFiles

	c.TplName = "index.html"
}

// Pdf Visualiser Modifier une image
func (c *MainController) Pdf() {
	pathFile := "/" + c.Ctx.Input.Param(":path") + "." + c.Ctx.Input.Param(":ext")
	setSession(c, "File", pathFile)

	// Load Folder
	pathFolder := c.GetSession("Folder").(string)
	hugoFiles := models.HugoGetFolder(pathFolder)

	// Recherche du record
	record := models.HugoGetRecord(pathFile)

	flash := beego.ReadFromRequest(&c.Controller)
	if record.Path == "" {
		msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
	}

	// Remplissage du contexte pour le template
	c.Data["Record"] = record
	c.Data["Records"] = hugoFiles

	c.TplName = "index.html"
}

// Document Visualiser Modifier un document
func (c *MainController) Document() {
	pathFile := "/" + c.Ctx.Input.Param(":path") + "." + c.Ctx.Input.Param(":ext")
	setSession(c, "File", pathFile)

	// Recherche du record
	record := models.HugoGetRecord(pathFile)

	flash := beego.ReadFromRequest(&c.Controller)
	if record.Path == "" {
		msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
		logs.Error(msg)
		flash.Error(msg)
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
			msg := fmt.Sprintf("hugoFile %s : %s", record.PathAbsolu, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
		models.HugoReload()
	}
	// Load Folder
	pathFolder := c.GetSession("Folder").(string)
	hugoFiles := models.HugoGetFolder(pathFolder)

	// Remplissage du contexte pour le template
	c.Data["Record"] = record
	c.Data["Records"] = hugoFiles

	c.TplName = "index.html"
}

// FileRename Renommer le fichier
func (c *MainController) FileRename() {
	path := "/" + c.Ctx.Input.Param(":path")
	if c.Ctx.Input.Param(":ext") != "" {
		path += "." + c.Ctx.Input.Param(":ext")
	}
	pathFolder := c.GetSession("Folder").(string)
	flash := beego.ReadFromRequest(&c.Controller)
	// Recherche du record
	record := models.HugoGetRecord(path)
	if record.Path == "" {
		msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	if record.IsDir == 1 {
		newName := c.GetString("new_name")
		path := strings.Split(record.PathAbsolu, "/")
		path[len(path)-1] = newName
		newFile := strings.Join(path, "/")
		if _, err := os.Stat(newFile); err == nil {
			// path/to/whatever exists
			msg := fmt.Sprintf("Renommer en [%s] : %s", newFile, "existe déjà")
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/folder"+pathFolder)
			return
		}
		err = os.Rename(record.PathAbsolu, newFile)
		if err != nil {
			msg := fmt.Sprintf("Renommer en [%s] : %s", record.Path, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
	} else {
		// Copie du fichier sur la cible
		newFile := models.Config.HugoRacine + "/content" + record.Dir + "/" + c.GetString("new_name")
		if _, err := os.Stat(newFile); err == nil {
			// path/to/whatever exists
			msg := fmt.Sprintf("Renommer en [%s] : %s", newFile, "existe déjà")
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/folder"+pathFolder)
			return
		}

		data, err := ioutil.ReadFile(record.PathAbsolu)
		if err != nil {
			msg := fmt.Sprintf("Renommer en [%s] : %s", record.Path, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
		err = ioutil.WriteFile(newFile, data, 0644)
		if err != nil {
			msg := fmt.Sprintf("Renommer en [%s] : %s", newFile, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
		// Suppression du fichier source
		err = os.RemoveAll(record.PathAbsolu)
		if err != nil {
			msg := fmt.Sprintf("Renommer en [%s] : %s", record.Path, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
	}

	// reLoad Folder
	models.HugoReload()
	c.Ctx.Redirect(302, "/folder"+pathFolder)
	return
}

// FileMove Déplacer le fichier
func (c *MainController) FileMove() {
	path := "/" + c.Ctx.Input.Param(":path")
	if c.Ctx.Input.Param(":ext") != "" {
		path += "." + c.Ctx.Input.Param(":ext")
	}
	pathFolder := c.GetSession("Folder").(string)
	flash := beego.ReadFromRequest(&c.Controller)
	// Recherche du record
	record := models.HugoGetRecord(path)
	if record.Path == "" {
		msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	if record.IsDir == 1 {
		newFile := models.Config.HugoRacine + "/content" + c.GetString("new_path") + "/" + record.Base
		if _, err := os.Stat(newFile); err == nil {
			// path/to/whatever exists
			msg := fmt.Sprintf("Déplacer [%s] : %s", newFile, "existe déjà")
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/folder"+pathFolder)
			return
		}
		err = os.Rename(record.PathAbsolu, newFile)
		if err != nil {
			msg := fmt.Sprintf("Déplacer [%s] : %s", record.Path, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
	} else {
		// Copie du fichier sur la cible
		newFile := models.Config.HugoRacine + "/content" + c.GetString("new_path") + "/" + record.Base
		if _, err := os.Stat(newFile); err == nil {
			// path/to/whatever exists
			msg := fmt.Sprintf("Déplacer [%s] : %s", newFile, "existe déjà")
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/folder"+pathFolder)
			return
		}

		data, err := ioutil.ReadFile(record.PathAbsolu)
		if err != nil {
			msg := fmt.Sprintf("Déplacer [%s] : %s", record.Path, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
		err = ioutil.WriteFile(newFile, data, 0644)
		if err != nil {
			msg := fmt.Sprintf("Déplacer [%s] : %s", newFile, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
		// Suppression du fichier source
		err = os.RemoveAll(record.PathAbsolu)
		if err != nil {
			msg := fmt.Sprintf("Renommer en [%s] : %s", record.Path, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			return
		}
	}

	// reLoad Folder
	models.HugoReload()
	c.Ctx.Redirect(302, "/folder"+pathFolder)
	return
}

// FileCp Recopier le fichier ou dossier
func (c *MainController) FileCp() {
	path := "/" + c.Ctx.Input.Param(":path")
	if c.Ctx.Input.Param(":ext") != "" {
		path += "." + c.Ctx.Input.Param(":ext")
	}
	pathFolder := c.GetSession("Folder").(string)
	flash := beego.ReadFromRequest(&c.Controller)
	// Recherche du record
	record := models.HugoGetRecord(path)
	if record.Path == "" {
		msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	if record.IsDir == 1 {
		newPath := models.Config.HugoRacine + "/content" + c.GetString("new_path") + "/" + record.Base
		err = shutil.CopyTree(record.PathAbsolu, newPath, nil)
		// err = os.Rename(record.PathAbsolu, newPath)
		if err != nil {
			msg := fmt.Sprintf("Copie vers [%s] : %s", newPath, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			models.HugoReload()
			c.Ctx.Redirect(302, "/folder"+pathFolder)
			return
		}
	} else {
		newPath := models.Config.HugoRacine + "/content" + c.GetString("new_path") + "/" + record.Base
		if _, err := os.Stat(newPath); err == nil {
			// path/to/whatever exists
			msg := fmt.Sprintf("Copie vers [%s] : %s", newPath, "existe déjà")
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			models.HugoReload()
			c.Ctx.Redirect(302, "/folder"+pathFolder)
			return
		}
		data, err := ioutil.ReadFile(record.PathAbsolu)
		if err != nil {
			msg := fmt.Sprintf("Copie vers [%s] : %s", newPath, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			models.HugoReload()
			c.Ctx.Redirect(302, "/folder"+pathFolder)
			return
		}
		err = ioutil.WriteFile(newPath, data, 0644)
		if err != nil {
			msg := fmt.Sprintf("Copie vers [%s : %s", newPath, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			models.HugoReload()
			c.Ctx.Redirect(302, "/folder"+pathFolder)
			return
		}
	}

	// reLoad Folder
	models.HugoReload()
	c.Ctx.Redirect(302, "/folder"+pathFolder)
	return
}

// FileNew Nouveau document à partir du modele.md
func (c *MainController) FileNew() {
	path := "/" + c.Ctx.Input.Param(":path")
	if c.Ctx.Input.Param(":ext") != "" {
		path += "." + c.Ctx.Input.Param(":ext")
	}
	pathFolder := c.GetSession("Folder").(string)
	flash := beego.ReadFromRequest(&c.Controller)

	newFile := models.Config.HugoRacine + "/content" + pathFolder + "/" + c.GetString("new_name")
	modele := models.Config.HugoRacine + "/content/site/modele.md"
	data, err := ioutil.ReadFile(modele)
	if err != nil {
		msg := fmt.Sprintf("Modèle fichier %s : %s", modele, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/folder"+pathFolder)
		return
	}

	err = ioutil.WriteFile(newFile, data, 0644)
	if err != nil {
		msg := fmt.Sprintf("Nouveau fichier %s : %s", newFile, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/folder"+pathFolder)
		return
	}

	// reLoad Folder
	models.HugoReload()
	c.Ctx.Redirect(302, "/folder"+pathFolder)
	return
}

// FileRm Supprimer le fichier ou dossier
func (c *MainController) FileRm() {
	path := "/" + c.Ctx.Input.Param(":path")
	if c.Ctx.Input.Param(":ext") != "" {
		path += "." + c.Ctx.Input.Param(":ext")
	}
	pathFolder := c.GetSession("Folder").(string)
	flash := beego.ReadFromRequest(&c.Controller)
	// Recherche du record
	record := models.HugoGetRecord(path)
	if record.Path == "" {
		msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	err = os.RemoveAll(record.PathAbsolu)
	if err != nil {
		msg := fmt.Sprintf("Suppression de [%s] : %s", record.Path, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}
	if path == c.GetSession("File").(string) {
		c.SetSession("File", "")
	}
	// reLoad Folder
	models.HugoReload()
	c.Ctx.Redirect(302, "/folder"+pathFolder)
	return
}

// FileUpload Charger le fichier sur le serveur
func (c *MainController) FileUpload() {
	pathFolder := c.GetSession("Folder").(string)
	flash := beego.ReadFromRequest(&c.Controller)

	files, err := c.GetFiles("files")
	if err != nil {
		msg := fmt.Sprintf("Import : %s", err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/folder"+pathFolder)
		return
	}
	for _, mfile := range files {
		file, err := mfile.Open()
		defer file.Close()
		if err != nil {
			msg := fmt.Sprintf("Import : %s", err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/folder"+pathFolder)
		}
		fileContents, err := ioutil.ReadAll(file)
		err = ioutil.WriteFile(models.Config.HugoRacine+"/content"+pathFolder+mfile.Filename, fileContents, 0644)
		if err != nil {
			msg := fmt.Sprintf("Import : %s", err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/folder"+pathFolder)
		}
	}
	// reLoad Folder
	models.HugoReload()
	c.Ctx.Redirect(302, "/folder"+pathFolder)
	return
}

// FileMkdir Créer un dossier
func (c *MainController) FileMkdir() {
	path := "/" + c.Ctx.Input.Param(":path")
	if c.Ctx.Input.Param(":ext") != "" {
		path += "." + c.Ctx.Input.Param(":ext")
	}
	pathFolder := c.GetSession("Folder").(string)
	flash := beego.ReadFromRequest(&c.Controller)

	newDir := models.Config.HugoRacine + "/content" + pathFolder + "/" + c.GetString("new_name")

	err = os.MkdirAll(newDir, 0744)
	if err != nil {
		msg := fmt.Sprintf("Nouveau dossier %s : %s", newDir, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	// reLoad Folder
	models.HugoReload()
	c.Ctx.Redirect(302, "/folder"+pathFolder)
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

// APIFolders as /api/folders
func (c *MainController) APIFolders() {
	type myList struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	type myStruct struct {
		Success bool     `json:"success"`
		Message string   `json:"message"`
		Results []myList `json:"results"`
	}
	var jlist []myList
	for _, record := range models.HugoGetFolders() {
		jlist = append(jlist, myList{Name: record.Path, Value: record.Path})
	}
	var resp myStruct
	resp.Success = true
	resp.Message = "ok coral"
	resp.Results = jlist

	c.Data["json"] = &resp
	c.ServeJSON()
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
	// hugoFiles = nil
}
