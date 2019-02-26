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
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/renproject/swapperd/adapter/server"
	"github.com/renproject/swapperd/driver/service"
	"github.com/sirupsen/logrus"
)

type Updater struct {
	frequency time.Duration
	homeDir   string
	logger    logrus.FieldLogger
}

func New() (*Updater, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	homeDir := filepath.Dir(filepath.Dir(ex))
	logFile, err := os.OpenFile(fmt.Sprintf("%s/swapperd-updater.log", homeDir), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetOutput(logFile)
	frequency := time.Minute
	if configData, err := ioutil.ReadFile(fmt.Sprintf("%s/config.json", homeDir)); err == nil {
		config := struct {
			Frequency time.Duration `json:"frequency"`
		}{}
		if err := json.Unmarshal(configData, &config); err == nil {
			frequency = config.Frequency * time.Second
		}
	}
	return &Updater{
		frequency: frequency,
		homeDir:   homeDir,
		logger:    logger,
	}, nil
}

func (updater *Updater) Update() error {
	updater.logger.Info("looking for latest version ...")
	latVer, err := getLatestVersion()
	if err != nil {
		return err
	}
	updater.logger.Infof("latest version is %s", latVer)
	currVer, err := updater.getCurrentVersion()
	if err != nil {
		return updater.updateSwapperd(latVer)
	}
	updater.logger.Info("current version is %s", currVer)
	if res, err := compareVersions(currVer, latVer); err != nil || !res {
		return err
	}
	return updater.updateSwapperd(latVer)
}

func (updater *Updater) getCurrentVersion() (string, error) {
	info := server.GetInfoResponse{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:7927/info", nil)
	if err != nil {
		return info.Version, err
	}

	req.SetBasicAuth("", "")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return info.Version, err
	}

	if resp.StatusCode != 200 {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return info.Version, err
		}
		return info.Version, fmt.Errorf("Failed to get the current version: %s", string(respBytes))
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return info.Version, err
	}
	return info.Version, nil
}

func (updater *Updater) updateSwapperd(ver string) error {
	if err := updater.downloadSwapperd(ver); err != nil {
		return err
	}
	service.Stop("swapperd")
	service.Stop("swapperd-updater")
	if err := updater.unzipSwapperd(); err != nil {
		return err
	}
	service.Start("swapperd-updater")
	service.Start("swapperd")
	return updater.updateConfig(ver)
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

func (updater *Updater) updateConfig(version string) error {
	// Get the data
	resp, err := http.Get(fmt.Sprintf("https://github.com/renproject/swapperd/releases/download/%s/config.json", version))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get the latest config file (%d): %s", resp.StatusCode, respBytes)
	}
	return ioutil.WriteFile(fmt.Sprintf("%s/config.json", updater.homeDir), respBytes, 0644)
}

func (updater *Updater) unzipSwapperd() error {
	src := path.Join(updater.homeDir, "swapperd.zip")
	dest := updater.homeDir

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
				updater.logger.Errorf("could not update %s: %v", f.Name, err)
				continue
			}
			defer f.Close()
			if _, err := io.Copy(f, rc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (updater *Updater) downloadSwapperd(version string) error {
	// Get the data
	resp, err := http.Get(fmt.Sprintf("https://github.com/renproject/swapperd/releases/download/%s/swapper_%s_amd64.zip", version, runtime.GOOS))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to download the swapperd.zip file: %s", respBytes)
	}

	// Create the file
	out, err := os.Create(fmt.Sprintf("%s/swapperd.zip", updater.homeDir))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func compareVersions(curr, latest string) (bool, error) {
	rgx, err := regexp.Compile("v([0-9]+).([0-9]+).([0-9]+)(-(nightly|beta|stable)(.([0-9]+))?)?")
	if err != nil {
		return false, err
	}

	currShares := rgx.FindStringSubmatch(curr)
	latestShares := rgx.FindStringSubmatch(latest)

	currMajor, err := strconv.ParseInt(currShares[1], 10, 64)
	if err != nil {
		return false, err
	}
	currMinor, err := strconv.ParseInt(currShares[2], 10, 64)
	if err != nil {
		return false, err
	}
	currPatch, err := strconv.ParseInt(currShares[3], 10, 64)
	if err != nil {
		return false, err
	}

	latestMajor, err := strconv.ParseInt(latestShares[1], 10, 64)
	if err != nil {
		return false, err
	}

	latestMinor, err := strconv.ParseInt(latestShares[2], 10, 64)
	if err != nil {
		return false, err
	}

	latestPatch, err := strconv.ParseInt(latestShares[3], 10, 64)
	if err != nil {
		return false, err
	}

	if currShares[0] == latestShares[0] ||
		currMajor > latestMajor ||
		(currMajor == latestMajor && currMinor > latestMinor) ||
		(currMajor == latestMajor && currMinor == latestMinor && currPatch > latestPatch) {
		return false, nil
	}

	if currMajor == latestMajor && currMinor == latestMinor && currPatch == latestPatch {
		return compareTags(currShares[5:], latestShares[5:])
	}
	return true, nil
}

func compareTags(curr, latest []string) (bool, error) {
	currType, err := convertTypeToNumber(curr[0])
	if err != nil {
		return false, err
	}

	latestType, err := convertTypeToNumber(latest[0])
	if err != nil {
		return false, err
	}

	if currType > latestType {
		return false, nil
	} else if currType < latestType {
		return true, nil
	}

	if curr[1] != "" && latest[1] == "" {
		return false, nil
	} else if curr[1] == "" && latest[1] != "" {
		return true, nil
	}

	currTypePatch, err := strconv.ParseInt(curr[2], 10, 64)
	if err != nil {
		return false, err
	}

	latestTypePatch, err := strconv.ParseInt(latest[2], 10, 64)
	if err != nil {
		return false, err
	}

	if currTypePatch > latestTypePatch {
		return false, nil
	}
	return true, nil
}

func convertTypeToNumber(releaseType string) (int, error) {
	switch releaseType {
	case "nightly":
		return 1, nil
	case "beta":
		return 2, nil
	case "stable", "":
		return 3, nil
	default:
		return -1, fmt.Errorf("unknown release type: %v", releaseType)
	}
}
