package controllers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/kennygrant/sanitize"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/models"
	"github.com/pbillerot/victor/shutil"
)

// Hugo redir /hugo en /hugo/
// func (c *MainController) Hugo() {
// 	c.Ctx.Redirect(302, "/hugo/")
// }

// Main as get and Post
func (c *MainController) Main() {
	setSession(c, "Folder", "/")

	beego.ReadFromRequest(&c.Controller)

	if c.GetSession("Hugo").(models.Hugo).Name == "" {
		hugoApp := models.GetFirstHugoApp()
		if hugoApp.Name != "" {
			SetHugoApp(c, hugoApp)
			c.Ctx.Redirect(302, "/victor/folder/")
		}
	}

	// Remplissage du contexte pour le template
	c.Data["Record"] = models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), c.GetSession("File").(string))
	c.Data["Records"] = models.HugoGetFolder(c.GetSession("Hugo").(models.Hugo), c.GetSession("Folder").(string))

	c.TplName = "index.html"
}

// SelectHugoApp Sélection d'un folder à administrer
func (c *MainController) SelectHugoApp() {

	hugoApp := models.GetHugoApp(c.Ctx.Input.Param(":app"))

	if hugoApp.Name != "" {
		SetHugoApp(c, hugoApp)
	}
	c.Ctx.Redirect(302, "/victor/folder/")
}

// Folder Demande de lister le dossier
func (c *MainController) Folder() {
	pathFolder := "/" + c.Ctx.Input.Param(":path")
	setSession(c, "Folder", pathFolder)

	beego.ReadFromRequest(&c.Controller)

	// Remplissage du contexte pour le template
	c.Data["Record"] = models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), c.GetSession("File").(string))
	c.Data["Records"] = models.HugoGetFolder(c.GetSession("Hugo").(models.Hugo), c.GetSession("Folder").(string))

	c.TplName = "index.html"
}

// Image Visualiser Modifier une image
func (c *MainController) Image() {
	pathFile := "/" + c.Ctx.Input.Param(":path") + "." + c.Ctx.Input.Param(":ext")
	setSession(c, "File", pathFile)
	flash := beego.ReadFromRequest(&c.Controller)

	// Recherche du record
	record := models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), pathFile)
	if record.Path == "" {
		msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
	}
	// c.SetSession("Folder", record.Dir)

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
			c.Ctx.Redirect(302, c.Ctx.Request.URL.String())
			return
		}
		// pathAbsolu := c.Data["Hugo"].ContentDir(models.Hugo) + pathFile
		pathAbsolu := c.GetSession("Hugo").(models.Hugo).ContentDir + pathFile
		err = ioutil.WriteFile(pathAbsolu, unbased, 0644)
		if err != nil {
			msg := fmt.Sprintf("HugoImage %s : %s", pathAbsolu, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, c.Ctx.Request.URL.String())
			return
		}
		hugoReload(c)
		pushDev(c)
		c.Data["Refresh"] = true
	}

	// Remplissage du contexte pour le template
	c.Data["Record"] = models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), c.GetSession("File").(string))
	c.Data["Records"] = models.HugoGetFolder(c.GetSession("Hugo").(models.Hugo), c.GetSession("Folder").(string))
	c.Data["Folder"] = c.GetSession("Folder").(string)
	c.Data["File"] = c.GetSession("File").(string)
	if record.IsDrawio {
		dataURL, err := shutil.EncodePngToDataURL(strings.Replace(record.PathAbsolu, ".drawio", ".png", -1))
		if err != nil {
			c.Data["DataURL"] = template.URL("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVQYV2NgYAAAAAMAAWgmWQ0AAAAASUVORK5CYII=")
		} else {
			c.Data["DataURL"] = template.URL(dataURL)
		}
		// logs.Info("dataURL=[%s]", dataURL)
	}
	c.TplName = "file.html"
}

// Pdf Visualiser Modifier une image
func (c *MainController) Pdf() {
	pathFile := "/" + c.Ctx.Input.Param(":path") + "." + c.Ctx.Input.Param(":ext")
	setSession(c, "File", pathFile)
	flash := beego.ReadFromRequest(&c.Controller)

	// Recherche du record
	record := models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), pathFile)
	if record.Path == "" {
		msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
	}
	// c.SetSession("Folder", record.Dir)

	// Remplissage du contexte pour le template
	c.Data["Record"] = models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), c.GetSession("File").(string))
	c.Data["Records"] = models.HugoGetFolder(c.GetSession("Hugo").(models.Hugo), c.GetSession("Folder").(string))
	c.Data["Folder"] = c.GetSession("Folder").(string)
	c.Data["File"] = c.GetSession("File").(string)
	c.TplName = "file.html"
}

// Document Visualiser Modifier un document
func (c *MainController) Document() {
	// Mémorisation de la position du curseur
	cursorCh := c.GetString("cursor_ch")
	cursorLine := c.GetString("cursor_line")

	pathFile := "/" + c.Ctx.Input.Param(":path")
	if c.Ctx.Input.Param(":ext") != "" {
		pathFile += "." + c.Ctx.Input.Param(":ext")
	}
	setSession(c, "File", pathFile)
	flash := beego.ReadFromRequest(&c.Controller)

	// Recherche du record
	record := models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), pathFile)
	if record.Path == "" {
		msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
	}
	// c.SetSession("Folder", record.Dir)

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
				c.Ctx.Redirect(302, c.Ctx.Request.URL.String())
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
			c.Ctx.Redirect(302, c.Ctx.Request.URL.String())
			return
		}
		hugoReload(c)
		pushDev(c)
		c.Data["Refresh"] = true
	}
	// Remplissage du contexte pour le template
	c.Data["Record"] = models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), c.GetSession("File").(string))
	c.Data["Records"] = models.HugoGetFolder(c.GetSession("Hugo").(models.Hugo), c.GetSession("Folder").(string))
	c.Data["Folder"] = c.GetSession("Folder").(string)
	c.Data["File"] = c.GetSession("File").(string)
	// Mémorisation de la position du curseur
	c.Data["CursorCh"] = cursorCh
	c.Data["CursorLine"] = cursorLine

	c.TplName = "file.html"
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
	record := models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), path)
	if record.Path == "" {
		msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
		return
	}

	if record.IsDir {
		newName := sanitize.Name(c.GetString("new_name"))
		path := strings.Split(record.PathAbsolu, "/")
		path[len(path)-1] = newName
		newFile := strings.Join(path, "/")
		if _, err := os.Stat(newFile); err == nil {
			// path/to/whatever exists
			msg := fmt.Sprintf("Renommer en [%s] : %s", newFile, "existe déjà")
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}
		err = os.Rename(record.PathAbsolu, newFile)
		if err != nil {
			msg := fmt.Sprintf("Renommer en [%s] : %s", record.Path, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}
	} else {
		// Copie du fichier sur la cible
		newFile := c.GetSession("Hugo").(models.Hugo).ContentDir + record.Dir + "/" + c.GetString("new_name")
		if _, err := os.Stat(newFile); err == nil {
			// path/to/whatever exists
			msg := fmt.Sprintf("Renommer en [%s] : %s", newFile, "existe déjà")
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}

		data, err := ioutil.ReadFile(record.PathAbsolu)
		if err != nil {
			msg := fmt.Sprintf("Renommer en [%s] : %s", record.Path, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}
		err = ioutil.WriteFile(newFile, data, 0644)
		if err != nil {
			msg := fmt.Sprintf("Renommer en [%s] : %s", newFile, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}
		// Suppression du fichier source
		err = os.RemoveAll(record.PathAbsolu)
		if err != nil {
			msg := fmt.Sprintf("Renommer en [%s] : %s", record.Path, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}
	}
	if path == c.GetSession("File").(string) {
		c.SetSession("File", "")
	}
	// reLoad Folder
	hugoReload(c)
	pushDev(c)
	c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
}

// FileMove Déplacer le fichier
func (c *MainController) FileMove() {
	// liste des fichiers à déplacer sépârés avec des ,
	paths := strings.Split(c.GetString("paths"), ",")
	// Répertoire destination
	dest := c.GetString("dest")

	flash := beego.ReadFromRequest(&c.Controller)
	pathFolder := c.GetSession("Folder").(string) // répertoire source des fichiers
	// Traitement unitaire des fichiers ou répertoires
	for _, path := range paths {
		record := models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), path)
		if record.Path == "" {
			msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}
		if record.IsDir {
			newPath := c.GetSession("Hugo").(models.Hugo).ContentDir + dest + "/" + record.Base
			if _, err := os.Stat(newPath); err == nil {
				msg := fmt.Sprintf("Déplacer [%s] vers [%s] : Existe déjà", record.PathAbsolu, newPath)
				logs.Error(msg)
				flash.Error(msg)
				flash.Store(&c.Controller)
				hugoReload(c)
				pushDev(c)
				c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
				return
			}
			err = os.Rename(record.PathAbsolu, newPath)
			if err != nil {
				msg := fmt.Sprintf("Déplacer [%s] vers [%s] : %v", record.PathAbsolu, newPath, err)
				logs.Error(msg)
				flash.Error(msg)
				flash.Store(&c.Controller)
				hugoReload(c)
				pushDev(c)
				c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
				return
			}
		} else {
			newPath := c.GetSession("Hugo").(models.Hugo).ContentDir + dest + "/" + record.Base
			if _, err := os.Stat(newPath); err == nil {
				msg := fmt.Sprintf("Déplacer [%s] vers [%s] : Existe déjà", record.PathAbsolu, newPath)
				logs.Error(msg)
				flash.Error(msg)
				flash.Store(&c.Controller)
				hugoReload(c)
				pushDev(c)
				c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
				return
			}
			err = shutil.CopyFile(record.PathAbsolu, newPath, false)
			if err != nil {
				msg := fmt.Sprintf("Déplacer [%s] vers %s : %v", record.PathAbsolu, newPath, err)
				logs.Error(msg)
				flash.Error(msg)
				flash.Store(&c.Controller)
				hugoReload(c)
				pushDev(c)
				c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
				return
			}
			// Suppression du fichier source
			err = os.RemoveAll(record.PathAbsolu)
			if err != nil {
				msg := fmt.Sprintf("Déplacer [%s] vers %s : %v", record.PathAbsolu, newPath, err)
				logs.Error(msg)
				flash.Error(msg)
				flash.Store(&c.Controller)
				hugoReload(c)
				pushDev(c)
				c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
				return
			}
			if path == c.GetSession("File").(string) {
				c.SetSession("File", "")
			}
		}
	}

	// reLoad Folder
	hugoReload(c)
	pushDev(c)
	c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
}

// FileCp Recopier le fichier ou dossier
func (c *MainController) FileCp() {
	// liste des fichiers à déplacer sépârés avec des ,
	paths := strings.Split(c.GetString("paths"), ",")
	// Répertoire destination
	dest := c.GetString("dest")

	flash := beego.ReadFromRequest(&c.Controller)
	pathFolder := c.GetSession("Folder").(string) // répertoire source des fichiers
	// Traitement unitaire des fichiers ou répertoires
	for _, path := range paths {
		record := models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), path)
		if record.Path == "" {
			msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}
		if record.IsDir {
			newPath := c.GetSession("Hugo").(models.Hugo).ContentDir + dest + "/" + record.Base
			err = shutil.CopyTree(record.PathAbsolu, newPath, nil)
			if err != nil {
				msg := fmt.Sprintf("Copie [%s] vers [%s] : %v", record.PathAbsolu, newPath, err)
				logs.Error(msg)
				flash.Error(msg)
				flash.Store(&c.Controller)
				hugoReload(c)
				pushDev(c)
				c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
				return
			}
		} else {
			newPath := c.GetSession("Hugo").(models.Hugo).ContentDir + dest + "/" + record.Base
			if _, err := os.Stat(newPath); err == nil {
				// path/to/whatever exists
				newPath = c.GetSession("Hugo").(models.Hugo).ContentDir + dest + "/" + "Copy " + record.Base
			}
			err = shutil.CopyFile(record.PathAbsolu, newPath, false)
			if err != nil {
				msg := fmt.Sprintf("Copie [%s] vers [%s] : %v", record.PathAbsolu, newPath, err)
				logs.Error(msg)
				flash.Error(msg)
				flash.Store(&c.Controller)
				hugoReload(c)
				pushDev(c)
				c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
				return
			}
		}
	}

	// reLoad Folder
	hugoReload(c)
	pushDev(c)
	c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
}

// FileNew Nouveau document à partir du modele.md
func (c *MainController) FileNew() {
	// path := "/" + c.Ctx.Input.Param(":path")
	pathFolder := c.GetSession("Folder").(string)
	if pathFolder == "/" {
		pathFolder = ""
	}
	flash := beego.ReadFromRequest(&c.Controller)

	newName := c.GetString("new_name")
	newFile := c.GetSession("Hugo").(models.Hugo).ContentDir + pathFolder + "/" + newName
	if strings.Contains(newName, ".md") {
		err = newDocument(c, pathFolder, newName)
		if err != nil {
			msg := fmt.Sprintf("Modèle fichier %s : %s", newName, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}
	} else if strings.Contains(newName, ".drawio") {
		modele := "./static/img/new.png"
		data, err := ioutil.ReadFile(modele)
		if err != nil {
			msg := fmt.Sprintf("Modèle fichier %s : %s", modele, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}
		newFile = strings.Replace(newFile, ".drawio", ".png", 1)
		err = ioutil.WriteFile(newFile, data, 0644)
		if err != nil {
			msg := fmt.Sprintf("Nouveau fichier %s : %s", newFile, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}
	} else {
		err = ioutil.WriteFile(newFile, []byte("Created by Victor"), 0644)
	}
	if err != nil {
		msg := fmt.Sprintf("Nouveau fichier %s : %s", newFile, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
		return
	}

	// reLoad Folder
	hugoReload(c)
	pushDev(c)
	c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
}

// FileRm Supprimer le fichier ou dossier
func (c *MainController) FileRm() {
	// liste des fichiers à supprimer séparés avec des ,
	paths := strings.Split(c.GetString("paths"), ",")

	flash := beego.ReadFromRequest(&c.Controller)
	pathFolder := c.GetSession("Folder").(string) // répertoire source des fichiers
	// Traitement unitaire des fichiers ou répertoires
	for _, path := range paths {
		record := models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), path)
		if record.Path == "" {
			msg := fmt.Sprintf("[%s] : non trouvé", record.Path)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}
		err = os.RemoveAll(record.PathAbsolu)
		if err != nil {
			msg := fmt.Sprintf("Suppression de [%s] : %v", record.PathAbsolu, err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			hugoReload(c)
			pushDev(c)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
			return
		}
		if path == c.GetSession("File").(string) {
			c.SetSession("File", "")
		}
	}

	// reLoad Folder
	hugoReload(c)
	pushDev(c)
	c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
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
		c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
		return
	}
	for _, mfile := range files {
		file, err := mfile.Open()
		if err != nil {
			msg := fmt.Sprintf("Import : %s", err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
		}
		defer file.Close()
		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			msg := fmt.Sprintf("Import : %s", err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
		}
		path := c.GetSession("Hugo").(models.Hugo).ContentDir + pathFolder + "/" + mfile.Filename
		err = ioutil.WriteFile(path, fileContents, 0644)
		if err != nil {
			msg := fmt.Sprintf("Import : %s", err)
			logs.Error(msg)
			flash.Error(msg)
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
		}
	}
	// reLoad Folder
	hugoReload(c)
	pushDev(c)
	c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
}

// FileMkdir Créer un dossier
func (c *MainController) FileMkdir() {
	path := "/" + c.Ctx.Input.Param(":path")
	if c.Ctx.Input.Param(":ext") != "" {
		path += "." + c.Ctx.Input.Param(":ext")
	}
	pathFolder := c.GetSession("Folder").(string)
	flash := beego.ReadFromRequest(&c.Controller)

	newName := sanitize.Name(c.GetString("new_name"))
	newDir := c.GetSession("Hugo").(models.Hugo).ContentDir + pathFolder + "/" + newName

	err = os.MkdirAll(newDir, 0744)
	if err != nil {
		msg := fmt.Sprintf("Nouveau dossier %s : %s", newDir, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
		return
	}

	// reLoad Folder
	hugoReload(c)
	pushDev(c)
	c.Ctx.Redirect(302, "/victor/folder"+pathFolder)
}

// APIFolders as /api/folders
func (c *MainController) APIFolders() {
	type myList struct {
		Base     string `json:"base"`
		Path     string `json:"path"`
		Rang     int    `json:"rang"`
		Selected bool   `json:"selected"`
	}
	type myStruct struct {
		Success bool     `json:"success"`
		Message string   `json:"message"`
		Results []myList `json:"results"`
	}
	var list []myList
	for _, record := range models.HugoGetFolders(c.GetSession("Hugo").(models.Hugo)) {
		if record.Path == c.GetSession("Folder").(string) {
			list = append(list, myList{
				Base:     record.Base,
				Path:     record.Path,
				Rang:     record.Rang,
				Selected: true,
			})
		} else {
			list = append(list, myList{
				Base:     record.Base,
				Path:     record.Path,
				Rang:     record.Rang,
				Selected: false,
			})
		}
	}

	var resp myStruct
	resp.Success = true
	resp.Message = "ok coral"
	resp.Results = list

	c.Data["json"] = &resp
	c.ServeJSON()
}

// SetHugoApp positionnement du contaxte hugo
func SetHugoApp(c *MainController, hugoApp models.HugoApp) {
	// Changement de contexte webapp hugo
	hugo := models.Hugo{}
	hugo.Name = hugoApp.Name
	hugo.Title = hugoApp.Title
	hugo.BaseURL = hugoApp.BaseURL
	hugo.Racine = hugoApp.Folder
	hugo.Theme = hugoApp.Theme
	hugo.ThemeHelp = hugoApp.ThemeHelp
	hugo.Deploy = hugoApp.Deploy
	hugo.DeployLabel = hugoApp.DeployLabel
	hugo.ContentDir = hugoApp.Folder + "/content"
	hugo.PrivateDir = hugoApp.Folder + "/private"
	hugo.PublicDir = hugoApp.Folder + "/public"
	logs.Info("SetHugoApp:", hugo.Name, hugo.Racine)

	c.SetSession("Hugo", hugo)
	hugoReload(c)

	beego.SetStaticPath("/content", hugo.ContentDir)
	beego.SetStaticPath("/hugo", hugo.PrivateDir)

}

// newDocument : Exécution du moteur Hugo pour créer un nouveau document
func newDocument(c *MainController, folder string, filename string) error {
	folder = strings.TrimPrefix(folder, "/")
	var cmd *exec.Cmd
	if folder == "" {
		cmd = exec.Command("hugo", "new", filename)
	} else {
		cmd = exec.Command("hugo", "new", folder+"/"+filename)
	}
	// cmd := exec.Command("ls", "-l")
	logs.Info("newDocument", cmd)
	cmd.Dir = c.GetSession("Hugo").(models.Hugo).Racine
	out, err := cmd.CombinedOutput()
	if err != nil {
		logs.Error("newDocument", err)
	}
	logs.Info("newDocument", string(out))
	return err
}

// APIFile as /api/file
func (c *MainController) APIFile() {
	type myStruct struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Results models.HugoFile `json:"results"`
	}
	record := models.HugoGetRecord(c.GetSession("Hugo").(models.Hugo), c.GetSession("File").(string))

	var resp myStruct
	resp.Success = true
	resp.Message = "ok coral"
	resp.Results = record

	c.Data["json"] = &resp
	c.ServeJSON()
}

// Action Action
func (c *MainController) Action() {
	action := c.Ctx.Input.Param(":action")

	switch action {
	case "refresh":
		hugoReload(c)
		pushDev(c)
	case "publishDev":
		pushDev(c)
	case "pushProd":
		pushProd(c)
	case "deploy":
		deploy(c)
	case "gitUpdateTheme":
		gitUpdateTheme(c)
	}
	c.Ctx.Redirect(302, "/victor/folder"+c.GetSession("Folder").(string))
}

// pushDev : Exécution du moteur Hugo pour mettre à jour le site de développement
func pushDev(c *MainController) {
	flash := beego.ReadFromRequest(&c.Controller)

	cmd := exec.Command("hugo", "-b", "/hugo/", "-d", c.GetSession("Hugo").(models.Hugo).PrivateDir,
		"--cleanDestinationDir", "--cacheDir", "/tmp")
	// cmd := exec.Command("ls", "-l")
	logs.Info("pushDev", cmd)
	cmd.Dir = c.GetSession("Hugo").(models.Hugo).Racine
	out, err := cmd.CombinedOutput()
	if err != nil {
		logs.Error("publishDev", err)
		flash.Error("ERREURG Génération des pages : %v", err)
		flash.Error(string(out))
		flash.Store(&c.Controller)
	} else {
		flash.Success(string(out))
		flash.Store(&c.Controller)
	}
	logs.Info("pushDev", string(out))
}

// pushProd : Exécution du moteur Hugo pour mettre à jour le site de production
func pushProd(c *MainController) {
	flash := beego.ReadFromRequest(&c.Controller)

	cmd := exec.Command("hugo", "-b", c.GetSession("Hugo").(models.Hugo).BaseURL, "--cleanDestinationDir", "--cacheDir", "/tmp")
	logs.Info("pushProd", cmd)
	cmd.Dir = c.GetSession("Hugo").(models.Hugo).Racine
	out, err := cmd.CombinedOutput()
	if err != nil {
		logs.Error("pushProd", err)
		flash.Error("pushProd : %v", err)
		flash.Error(string(out))
		flash.Store(&c.Controller)
	} else {
		flash.Success(string(out))
		flash.Store(&c.Controller)
	}
	logs.Info("pushProd", string(out))
}

// deploy : déployer la production sur un site
func deploy(c *MainController) {
	flash := beego.ReadFromRequest(&c.Controller)

	cmd := exec.Command("/bin/sh", c.GetSession("Hugo").(models.Hugo).Deploy)
	logs.Info("deploy", cmd)
	cmd.Dir = c.GetSession("Hugo").(models.Hugo).Racine
	out, err := cmd.CombinedOutput()
	if err != nil {
		logs.Error("deploy", err)
		flash.Error("deploy : %v", err)
		flash.Error(string(out))
		flash.Store(&c.Controller)
		return
	}
	flash.Success(string(out))
	flash.Store(&c.Controller)
	logs.Info("deploy", string(out))
}

// gitUpdateTheme : Mise à jour du thème à partir du référentiel git
func gitUpdateTheme(c *MainController) {
	flash := beego.ReadFromRequest(&c.Controller)

	// git or submodule ? .gitmodules
	_, err = os.Open(c.GetSession("Hugo").(models.Hugo).Racine + "/.gitmodules")
	if !os.IsNotExist(err) {
		// submodule
		cmd := exec.Command("git", "submodule", "update", "--remote")
		logs.Info("gitUpdateTheme submodule", cmd)
		cmd.Dir = c.GetSession("Hugo").(models.Hugo).Racine
		out, err := cmd.CombinedOutput()
		if err != nil {
			logs.Error("gitUpdateTheme", err)
			flash.Error("gitUpdateTheme : %v", err)
			flash.Error(string(out))
			flash.Store(&c.Controller)
			return
		}
		flash.Success(string(out))
		flash.Store(&c.Controller)
		logs.Info("gitUpdateTheme", string(out))
	} else {
		// git
		cmd := exec.Command("git", "reset", "--hard")
		logs.Info("gitUpdateTheme git reset", cmd)
		cmd.Dir = c.GetSession("Hugo").(models.Hugo).Racine + "/themes/" + c.GetSession("Hugo").(models.Hugo).Theme
		out, err := cmd.CombinedOutput()
		if err != nil {
			logs.Error("gitUpdateTheme", err)
			flash.Error("gitUpdateTheme : %v", err)
			flash.Error(string(out))
			flash.Store(&c.Controller)
			return
		}
		flash.Success(string(out))
		logs.Info("gitUpdateTheme", string(out))
		cmd = exec.Command("git", "pull")
		logs.Info("gitUpdateTheme git reset", cmd)
		cmd.Dir = c.GetSession("Hugo").(models.Hugo).Racine + "/themes/" + c.GetSession("Hugo").(models.Hugo).Theme
		out, err = cmd.CombinedOutput()
		if err != nil {
			logs.Error("gitUpdateTheme", err)
			flash.Error("gitUpdateTheme : %v", err)
			flash.Error(string(out))
			flash.Store(&c.Controller)
			return
		}
		flash.Success(string(out))
		flash.Store(&c.Controller)
		logs.Info("gitUpdateTheme", string(out))
	}

}

func hugoReload(c *MainController) {
	hugo := c.GetSession("Hugo").(models.Hugo)
	hugo.LoadFile()
	c.SetSession("Hugo", hugo)
	c.Data["Hugo"] = hugo
}
