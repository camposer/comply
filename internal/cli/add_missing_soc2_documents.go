package cli

import (
	"errors"
	"github.com/strongdm/comply/internal/config"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

var addMissingSOC2DocumentsCommand = cli.Command{
	Name:   "add",
	Usage:  "add missing default SOC2 documents for a given document type",
	Action: addMissingSOC2DocumentsAction,
	Before: projectMustExist,
}

func addMissingSOC2DocumentsAction(c *cli.Context) error {
	validSOC2Folders := map[string]bool{"narratives": true, "policies": true, "procedures": true, "standards": true}
	documentType := c.Args().First()
	if !validSOC2Folders[documentType] {
		return errors.New("invalid SOC2 document type")
	}

	_, b, _, _ := runtime.Caller(0)
	complyPath := filepath.Dir(b)
	documentsRootDir := strings.Split(strings.ReplaceAll(complyPath, "\\", "/"), "/")
	documentsRootPath := strings.Join(documentsRootDir[0:len(documentsRootDir)-2], "/") + "/themes/comply-soc2/" + documentType
	documentFilesInfo, err := ioutil.ReadDir(documentsRootPath)
	if err != nil {
		log.Fatal(err)
	}

	orgDocumentFilesInfo, err := ioutil.ReadDir(config.ProjectRoot() + "/" + documentType)
	if err != nil {
		log.Fatal(err)
	}

	for _, documentFileInfo := range documentFilesInfo {
		orgAlreadyHasFile := false
		for _, orgDocumentFileInfo := range orgDocumentFilesInfo {
			if orgDocumentFileInfo.Name() == documentFileInfo.Name() {
				orgAlreadyHasFile = true
				break
			}
		}

		if !orgAlreadyHasFile {
			documentFile, err := ioutil.ReadFile(documentsRootPath + "/" + documentFileInfo.Name())
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(config.ProjectRoot()+"/"+documentType+"/"+documentFileInfo.Name(), documentFile, 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
