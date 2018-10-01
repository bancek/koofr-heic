package koofrheictojpg

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	gopath "path"
	"path/filepath"
	"strings"

	"github.com/bancek/koofr-heic/app/models/heictojpg"
	koofrclient "github.com/koofr/go-koofrclient"
)

func Convert(koofr *koofrclient.KoofrClient, mountId string, path string, logger func(string)) (err error) {
	originalName := "_original"
	originalPath := gopath.Join(path, originalName)

	tmpDir, err := ioutil.TempDir("", "koofr-heic")
	if err != nil {
		return err
	}

	defer func() {
		os.RemoveAll(tmpDir)
	}()

	logger(fmt.Sprintf("Listing files for %s %s", mountId, path))

	files, err := koofr.FilesList(mountId, path)
	if err != nil {
		return fmt.Errorf("Convert koofr.FilesList failed: %s", err)
	}

	heicNames := []string{}
	originalDirExists := false

	for _, file := range files {
		name := file.Name

		if strings.HasSuffix(strings.ToLower(name), ".heic") {
			heicNames = append(heicNames, name)
		}
		if name == originalName {
			originalDirExists = true
		}
	}

	heicNamesCount := len(heicNames)

	if heicNamesCount == 0 {
		logger("No HEIC files found.")
		return nil
	}

	logger(fmt.Sprintf("Converting %d HEIC files.", heicNamesCount))

	if !originalDirExists {
		err = koofr.FilesNewFolder(mountId, path, originalName)
		if err != nil {
			return err
		}
	}

	for i, name := range heicNames {
		logger(fmt.Sprintf("Convert %d/%d", i+1, heicNamesCount))

		err = ConvertFile(koofr, mountId, path, name, originalPath, tmpDir, logger)
		if err != nil {
			logger(fmt.Sprintf("Convert %s failed: %s", name, err))
		}
	}

	logger("Done.")

	return nil
}

func ConvertFile(koofr *koofrclient.KoofrClient, mountId string, path string, name string, originalPath string, tmpDir string, logger func(string)) (err error) {
	heicRemotePath := gopath.Join(path, name)
	heicOriginalRemotePath := gopath.Join(originalPath, name)
	heicLocalPath := filepath.Join(tmpDir, name)

	jpgName := name[:len(name)-5] + ".jpg"
	jpgLocalPath := filepath.Join(tmpDir, jpgName)

	heicRemoteReader, err := koofr.FilesGet(mountId, heicRemotePath)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s koofr.FilesGet failed: %s", mountId, path, name, err)
	}
	defer heicRemoteReader.Close()

	heicLocalFile, err := os.Create(heicLocalPath)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s os.Create failed: %s", mountId, path, name, err)
	}
	defer heicLocalFile.Close()
	defer os.Remove(heicLocalPath)

	_, err = io.Copy(heicLocalFile, heicRemoteReader)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s io.Copy failed: %s", mountId, path, name, err)
	}

	err = heictojpg.HeicToJpg(heicLocalPath, jpgLocalPath)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s heictojpg.HeicToJpg failed: %s", mountId, path, name, err)
	}

	jpgLocalFile, err := os.Open(jpgLocalPath)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s os.Open failed: %s", mountId, path, name, err)
	}

	_, err = koofr.FilesPutOptions(mountId, path, jpgName, jpgLocalFile, nil)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s koofr.FilesPut failed: %s", mountId, path, name, err)
	}

	err = koofr.FilesMove(mountId, heicRemotePath, mountId, heicOriginalRemotePath)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s koofr.FilesMove failed: %s", mountId, path, name, err)
	}

	return nil
}
