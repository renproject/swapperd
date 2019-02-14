package updater

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type updater struct {
	version   string
	homeDir   string
	frequency time.Duration
	logger    logrus.FieldLogger
}

func New(version, homeDir string, frequency time.Duration, logger logrus.FieldLogger) *updater {
	return &updater{
		version:   version,
		homeDir:   homeDir,
		frequency: frequency,
		logger:    logger,
	}
}

func (updater *updater) Run(done <-chan struct{}) {
	ticker := time.NewTicker(updater.frequency)
	for {
		select {
		case <-done:
		case <-ticker.C:
			if err := updater.Update(done); err != nil {
				updater.logger.Error(err)
			}
		}
	}
}

func (updater *updater) Update(done <-chan struct{}) error {
	updater.logger.Info("looking for latest version ...")
	ver, err := getLatestVersion()
	if err != nil {
		return err
	}

	updater.logger.Infof("latest version is %s", ver)
	res, err := compareVersions(updater.version, ver)
	if err != nil {
		return err
	}

	if res {
		updater.logger.Info("downloading swapperd ...")
		if err := downloadSwapperd(ver, updater.homeDir); err != nil {
			return err
		}

		if err := updater.unzipSwapperd(); err != nil {
			return err
		}

		if err := updater.updateLocalVersion(ver); err != nil {
			return err
		}
	}
	return nil
}

func getLatestVersion() (string, error) {
	release := struct {
		TagName string `json:"tag_name"`
	}{}

	resp, err := http.DefaultClient.Get("https://api.github.com/repos/renproject/swapperd/releases/latest")
	if err != nil {
		return release.TagName, err
	}

	if resp.StatusCode != 200 {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return release.TagName, err
		}
		return release.TagName, fmt.Errorf("Failed to get the latest version: %s", string(respBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return release.TagName, err
	}

	return release.TagName, nil
}

func (updater *updater) updateLocalVersion(newVer string) error {
	config := struct {
		Version   string        `json:"version"`
		Frequency time.Duration `json:"frequency"`
	}{
		Version:   newVer,
		Frequency: updater.frequency,
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fmt.Sprintf("%s/config.json", updater.homeDir), configBytes, 0644)
}

func (updater *updater) unzipSwapperd() error {
	src := fmt.Sprintf("%s/swapperd.zip", updater.homeDir)
	dest := fmt.Sprintf("%s", updater.homeDir)

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func downloadSwapperd(version, homeDir string) error {
	// Get the data
	resp, err := http.Get(fmt.Sprintf("https://github.com/renproject/swapperd/releases/download/%s/swapper_%s_amd64.zip", version, runtime.GOOS))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fmt.Sprintf("%s/swapperd.zip", homeDir))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func compareVersions(curr, latest string) (bool, error) {
	currVersion := strings.Split(strings.Trim(curr, "v"), "-")
	latestVersion := strings.Split(strings.Trim(latest, "v"), "-")

	currShares := strings.Split(currVersion[0], ".")
	latestShares := strings.Split(latestVersion[0], ".")

	currMajor, err := strconv.ParseInt(currShares[0], 10, 64)
	if err != nil {
		return false, err
	}
	currMinor, err := strconv.ParseInt(currShares[1], 10, 64)
	if err != nil {
		return false, err
	}
	currPatch, err := strconv.ParseInt(currShares[2], 10, 64)
	if err != nil {
		return false, err
	}

	latestMajor, err := strconv.ParseInt(latestShares[0], 10, 64)
	if err != nil {
		return false, err
	}

	latestMinor, err := strconv.ParseInt(latestShares[1], 10, 64)
	if err != nil {
		return false, err
	}

	latestPatch, err := strconv.ParseInt(latestShares[2], 10, 64)
	if err != nil {
		return false, err
	}

	if currMajor > latestMajor ||
		(currMajor == latestMajor && currMinor > latestMinor) ||
		(currMajor == latestMajor && currMinor == latestMinor && currPatch > latestPatch) {
		if currVersion[0] == latestVersion[0] && len(currVersion) == 2 {
			return compareTags(currVersion[1], latestVersion[1])
		}
		return false, nil
	}
	return true, nil
}

func compareTags(curr, latest string) (bool, error) {
	currTypeShares := strings.Split(curr, ".")
	latestTypeShares := strings.Split(latest, ".")

	currType, err := convertTypeToNumber(currTypeShares[0])
	if err != nil {
		return false, err
	}

	latestType, err := convertTypeToNumber(latestTypeShares[0])
	if err != nil {
		return false, err
	}

	currTypePatch, err := strconv.ParseInt(currTypeShares[1], 10, 64)
	if err != nil {
		return false, err
	}

	latestTypePatch, err := strconv.ParseInt(latestTypeShares[1], 10, 64)
	if err != nil {
		return false, err
	}

	if currType > latestType ||
		(currType == latestType && currTypePatch >= latestTypePatch) {
		return false, nil
	}
	return true, nil
}

func convertTypeToNumber(releaseType string) (int, error) {
	switch releaseType {
	case "alpha":
		return 1, nil
	case "beta":
		return 2, nil
	case "rc":
		return 3, nil
	default:
		return -1, fmt.Errorf("unknown release type: %v", releaseType)
	}
}
