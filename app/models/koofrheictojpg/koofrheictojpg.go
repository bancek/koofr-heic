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
	"github.com/bancek/koofr-heic/app/models/movtomp4"
	koofrclient "github.com/koofr/go-koofrclient"
)

func Convert(koofr *koofrclient.KoofrClient, mountId string, path string, convertMovToMp4 bool, logger func(string)) (err error) {
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
	movNames := []string{}
	originalDirExists := false

	for _, file := range files {
		name := file.Name
		nameLower := strings.ToLower(name)

		if strings.HasSuffix(nameLower, ".heic") {
			heicNames = append(heicNames, name)
		}
		if convertMovToMp4 && strings.HasSuffix(nameLower, ".mov") {
			movNames = append(movNames, name)
		}
		if name == originalName {
			originalDirExists = true
		}
	}

	heicNamesCount := len(heicNames)
	movNamesCount := len(movNames)

	if heicNamesCount == 0 && movNamesCount == 0 {
		logger("No HEIC or MOV files found.")
		return nil
	}

	logger(fmt.Sprintf("Converting %d HEIC files.", heicNamesCount))
	if convertMovToMp4 {
		logger(fmt.Sprintf("Converting %d MOV files.", movNamesCount))
	}

	if !originalDirExists {
		err = koofr.FilesNewFolder(mountId, path, originalName)
		if err != nil {
			return err
		}
	}

	for i, name := range heicNames {
		logger(fmt.Sprintf("Convert HEIC %d/%d", i+1, heicNamesCount))

		err = ConvertHeicFile(koofr, mountId, path, name, originalPath, tmpDir, logger)
		if err != nil {
			logger(fmt.Sprintf("Convert HEIC %s failed: %s", name, err))
		}
	}
	for i, name := range movNames {
		logger(fmt.Sprintf("Convert MOV %d/%d", i+1, movNamesCount))

		err = ConvertMovFile(koofr, mountId, path, name, originalPath, tmpDir, logger)
		if err != nil {
			logger(fmt.Sprintf("Convert MOV %s failed: %s", name, err))
		}
	}

	logger("Done.")

	return nil
}

func ConvertHeicFile(koofr *koofrclient.KoofrClient, mountId string, path string, name string, originalPath string, tmpDir string, logger func(string)) (err error) {
	heicRemotePath := gopath.Join(path, name)
	heicOriginalRemotePath := gopath.Join(originalPath, name)
	heicLocalPath := filepath.Join(tmpDir, name)

	jpgName := name[:len(name)-5] + ".jpg"
	jpgLocalPath := filepath.Join(tmpDir, jpgName)

	defer os.Remove(heicLocalPath)
	defer os.Remove(jpgLocalPath)

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

func ConvertMovFile(koofr *koofrclient.KoofrClient, mountId string, path string, name string, originalPath string, tmpDir string, logger func(string)) (err error) {
	movRemotePath := gopath.Join(path, name)
	movOriginalRemotePath := gopath.Join(originalPath, name)
	movLocalPath := filepath.Join(tmpDir, name)

	mp4Name := name[:len(name)-4] + ".mp4"
	mp4LocalPath := filepath.Join(tmpDir, mp4Name)

	defer os.Remove(movLocalPath)
	defer os.Remove(mp4LocalPath)

	movRemoteReader, err := koofr.FilesGet(mountId, movRemotePath)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s koofr.FilesGet failed: %s", mountId, path, name, err)
	}
	defer movRemoteReader.Close()

	movLocalFile, err := os.Create(movLocalPath)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s os.Create failed: %s", mountId, path, name, err)
	}
	defer movLocalFile.Close()

	_, err = io.Copy(movLocalFile, movRemoteReader)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s io.Copy failed: %s", mountId, path, name, err)
	}

	err = movtomp4.MovToMp4(movLocalPath, mp4LocalPath)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s movtomp4.MovToMp4 failed: %s", mountId, path, name, err)
	}

	mp4LocalFile, err := os.Open(mp4LocalPath)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s os.Open failed: %s", mountId, path, name, err)
	}

	_, err = koofr.FilesPutOptions(mountId, path, mp4Name, mp4LocalFile, nil)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s koofr.FilesPut failed: %s", mountId, path, name, err)
	}

	err = koofr.FilesMove(mountId, movRemotePath, mountId, movOriginalRemotePath)
	if err != nil {
		return fmt.Errorf("ConvertFile mountId=%s path=%s name=%s koofr.FilesMove failed: %s", mountId, path, name, err)
	}

	return nil
}
