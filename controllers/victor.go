package controllers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/pbillerot/victor/models"
)

// Main as get and Post
func (c *MainController) Main() {
	// Chargement de hugoFiles et des meta du répertoire courant
	hugoFiles, _ := models.GetFilesFolder("/")

	c.Data["Records"] = hugoFiles
	c.Data["Folder"] = "/"

	c.TplName = "index.html"
}

// Folder Demande de lister le répertoire
func (c *MainController) Folder() {
	path := "/" + c.Ctx.Input.Param(":path")
	// Chargement de hugoFiles et des meta du répertoire courant
	hugoFiles, _ := models.GetFilesFolder(path)

	c.Data["Records"] = hugoFiles
	c.Data["Folder"] = path

	c.TplName = "index.html"
}

// Image Visualiser Modifier une image
func (c *MainController) Image() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	// for _, rec := range hugoFiles {
	// 	if rec.Key == keyid {
	// 		record = rec
	// 		break
	// 	}
	// }
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
	c.Ctx.Output.Cookie("from", fmt.Sprintf("/image/%s", appid))
	c.TplName = "hugo_image.html"
}

// Pdf Visualiser Modifier une image
func (c *MainController) Pdf() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	// for _, rec := range hugoFiles {
	// 	if rec.Key == keyid {
	// 		record = rec
	// 		break
	// 	}
	// }
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
	c.Ctx.Output.Cookie("from", fmt.Sprintf("/image/%s", appid))
	c.TplName = "hugo_pdf.html"
}

// Document Visualiser Modifier un document
func (c *MainController) Document() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	// for _, rec := range hugoFiles {
	// 	if rec.Key == keyid {
	// 		record = rec
	// 		break
	// 	}
	// }
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
		// hugoFiles = nil
		// // hugoDirectoryRecord(c)
		// for _, rec := range hugoFiles {
		// 	if rec.Key == keyid {
		// 		record = rec
		// 		break
		// 	}
		// }
		// Demande d'actualisation de l'arborescence
		c.Ctx.Output.Cookie("hugo-refresh-"+appid, "true")
	}

	// Remplissage du contexte pour le template
	c.Data["Record"] = record
	c.Data["KeyID"] = keyid
	c.Ctx.Output.Cookie("hugo-"+appid, keyid)
	c.Ctx.Output.Cookie("from", fmt.Sprintf("/document/%s", appid))
	c.TplName = "hugo_document.html"
}

// File Gestion du fichier
func (c *MainController) File() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	// for _, rec := range hugoFiles {
	// 	if rec.Key == keyid {
	// 		record = rec
	// 		break
	// 	}
	// }
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	// Remplissage du contexte pour le template
	c.Data["Record"] = record
	c.Data["KeyID"] = keyid

	c.Ctx.Output.Cookie("from", fmt.Sprintf("/file/%s/%s", appid, keyid))
	c.TplName = "hugo_file.html"
}

// Directory Gestion du répertoire
func (c *MainController) Directory() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	// for _, rec := range hugoFiles {
	// 	if rec.Key == keyid {
	// 		record = rec
	// 		break
	// 	}
	// }
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	// Remplissage du contexte pour le template
	c.Data["Record"] = record
	c.Data["KeyID"] = keyid

	c.Ctx.Output.Cookie("from", fmt.Sprintf("/directory/%s/%s", appid, keyid))
	c.TplName = "hugo_directory.html"
}

// FileMv Renommer le fichier
func (c *MainController) FileMv() {
	appid := c.Ctx.Input.Param(":app")
	keyid := c.Ctx.Input.Param(":key")

	// Recherche du record
	var record models.HugoFile
	// for _, rec := range hugoFiles {
	// 	if rec.Key == keyid {
	// 		record = rec
	// 		break
	// 	}
	// }
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
		err = os.Rename(record.PathAbsolu, models.Config.HugoRacine+newFile)
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
		err = ioutil.WriteFile(models.Config.HugoRacine+newFile, data, 0644)
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
	// hugoFiles = nil

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
	// for _, rec := range hugoFiles {
	// 	if rec.Key == keyid {
	// 		record = rec
	// 		break
	// 	}
	// }
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
	err = ioutil.WriteFile(models.Config.HugoRacine+"/"+newFile, data, 0644)
	if err != nil {
		msg := fmt.Sprintf("HugoFileCp %s : %s", newFile, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	// Vidage de Hugo pour reconstruction
	// hugoFiles = nil

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
	// for _, rec := range hugoFiles {
	// 	if rec.Key == keyid {
	// 		record = rec
	// 		break
	// 	}
	// }
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	newFile := c.GetString("new_file")
	data, err := ioutil.ReadFile(models.Config.HugoRacine + "/site/modele.md")
	if err != nil {
		msg := fmt.Sprintf("HugoFileNew %s : %s", record.Path, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}
	err = ioutil.WriteFile(models.Config.HugoRacine+"/"+newFile, data, 0644)
	if err != nil {
		msg := fmt.Sprintf("HugoFileNew %s : %s", newFile, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	// Vidage de Hugo pour reconstruction
	// hugoFiles = nil

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
	// for _, rec := range hugoFiles {
	// 	if rec.Key == keyid {
	// 		record = rec
	// 		break
	// 	}
	// }
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
	// hugoFiles = nil

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
	// for _, rec := range hugoFiles {
	// 	if rec.Key == keyid {
	// 		record = rec
	// 		break
	// 	}
	// }
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
	// hugoFiles = nil

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
	// for _, rec := range hugoFiles {
	// 	if rec.Key == keyid {
	// 		record = rec
	// 		break
	// 	}
	// }
	flash := beego.ReadFromRequest(&c.Controller)
	if record.Key == "" {
		logs.Error("Fichier non trouvé", c.GetSession("Username").(string), appid)
		flash.Error("Fichier non trouvé : %s", keyid)
		flash.Store(&c.Controller)
	}

	newDir := c.GetString("new_dir")
	err = os.MkdirAll(models.Config.HugoRacine+"/"+newDir, 0744)
	if err != nil {
		msg := fmt.Sprintf("HugoFileMkdir %s : %s", newDir, err)
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/")
		return
	}

	// Vidage de Hugo pour reconstruction
	// hugoFiles = nil

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
