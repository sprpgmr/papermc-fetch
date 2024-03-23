package paperapi

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/sprpgmr/papermc-fetch/files"
)

// Service contains methods to get paper api info conveniently.
type Service interface {
	GetLatestBuild(unstable bool, versionPrefix string) (*BuildInfo, error)
	IsValidDownload(filePath string, hash string) (bool, error)
	DownloadJar(buildInfo *BuildInfo, filepath string) error
	DownloadExists(filePath string, buildInfo *BuildInfo) (bool, error)
}

type serviceImpl struct {
	buildInfoService    BuildInfoService
	versionsListService VersionsListService
	buildsListService   BuildsListService
	fileService         files.Service
	baseURL             string
}

func newServiceImpl(buildInfoService BuildInfoService, versionsListService VersionsListService, buildsListService BuildsListService, fileService files.Service, baseURL string) *serviceImpl {
	return &serviceImpl{
		buildInfoService:    buildInfoService,
		versionsListService: versionsListService,
		buildsListService:   buildsListService,
		fileService:         fileService,
		baseURL:             baseURL,
	}
}

// GetLatestBuild will look for and return the BuildInfo of the latest stable version available, or latest unstable version available if unstable is true.
func (s *serviceImpl) GetLatestBuild(unstable bool, versionPrefix string) (*BuildInfo, error) {
	if !unstable {
		return s.getLatestStableVersion(versionPrefix)
	}

	versions, err := s.getFilteredVersionsList(versionPrefix)
	if err != nil {
		return nil, err
	}

	if len(versions.Versions) == 0 {
		return nil, errors.New("no versions found")
	}

	latestVersion := versions.Versions[len(versions.Versions)-1]

	return s.getLatestBuildInfo(latestVersion)
}

func (s *serviceImpl) getFilteredVersionsList(versionPrefix string) (*VersionsList, error) {
	versions, err := s.versionsListService.GetVersionsList()
	if err != nil {
		return nil, err
	}

	return filterVersions(versionPrefix, versions), nil
}

func filterVersions(versionPrefix string, versions *VersionsList) *VersionsList {
	if len(versionPrefix) == 0 {
		return versions
	}

	filteredVersions := &VersionsList{}
	filteredVersions.Versions = make([]string, 0)

	for i := 0; i < len(versions.Versions); i++ {
		if strings.Index(versions.Versions[i], versionPrefix) == 0 {
			if len(versions.Versions[i]) == len(versionPrefix) || versions.Versions[i][len(versionPrefix)] == '.' {
				filteredVersions.Versions = append(filteredVersions.Versions, versions.Versions[i])
			}
		}
	}

	return filteredVersions
}

func (s *serviceImpl) getLatestStableVersion(versionPrefix string) (*BuildInfo, error) {
	versions, err := s.getFilteredVersionsList(versionPrefix)
	if err != nil {
		return nil, err
	}

	for i := len(versions.Versions) - 1; i >= 0; i-- {
		buildInfo, err := s.getLatestBuildInfo(versions.Versions[i])
		if buildInfo.Channel == "default" {
			return buildInfo, err
		}
	}

	return nil, errors.New("no stable versions found")
}

func (s *serviceImpl) getLatestBuildInfo(version string) (*BuildInfo, error) {
	builds, err := s.buildsListService.GetBuildsList(version)
	if err != nil {
		return nil, err
	}

	if len(builds.Builds) == 0 {
		return nil, errors.New("no builds for this version exist")
	}

	latestBuild := builds.Builds[len(builds.Builds)-1]

	buildInfo, err := s.buildInfoService.GetBuildInfo(version, latestBuild)
	return buildInfo, nil
}

// IsValidDownload checks the sha256 sum of the filepath and compares it with the provided hash, returns true if they match
func (s *serviceImpl) IsValidDownload(filePath string, hash string) (bool, error) {
	h := sha256.New()

	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}

	defer file.Close()

	io.Copy(h, file)

	output := fmt.Sprintf("%x", h.Sum(nil))

	return output == hash, nil
}

// DownloadJar will download the paper jar file for the specific version and build number provided, to filepath
func (s *serviceImpl) DownloadJar(info *BuildInfo, filepath string) error {
	url := fmt.Sprint(s.baseURL, "/versions/", info.Version, "/builds/", info.Build, "/downloads/", info.Downloads.Application.Name)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func (s *serviceImpl) DownloadExists(filepath string, buildInfo *BuildInfo) (bool, error) {

	if s.fileService.FileExists(filepath) {
		valid, err := s.IsValidDownload(filepath, buildInfo.Downloads.Application.Sha256)
		if err != nil {
			return false, err
		}

		return valid, nil
	}

	return false, nil
}
