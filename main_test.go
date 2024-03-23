package main

import (
	"fmt"
	"testing"

	paperapi "github.com/sprpgmr/papermc-fetch/paper-api"
)

type paperServiceMock struct {
	getLatestBuildHandler  func(s *paperServiceMock, unstable bool, versionPrefix string) (*paperapi.BuildInfo, error)
	isValidDownloadHandler func(s *paperServiceMock, filePath string, hash string) (bool, error)
	downloadJarHandler     func(s *paperServiceMock, buildInfo *paperapi.BuildInfo, filepath string) error
	downloadExistsHandler  func(s *paperServiceMock, filepath string, buildInfo *paperapi.BuildInfo) (bool, error)
	ranDownload            bool
}

func (s *paperServiceMock) GetLatestBuild(unstable bool, versionPrefix string) (*paperapi.BuildInfo, error) {
	if s.getLatestBuildHandler != nil {
		return s.getLatestBuildHandler(s, unstable, versionPrefix)
	}

	return nil, nil
}

func (s *paperServiceMock) IsValidDownload(filePath string, hash string) (bool, error) {
	if s.isValidDownloadHandler != nil {
		return s.isValidDownloadHandler(s, filePath, hash)
	}

	return true, nil
}

func (s *paperServiceMock) DownloadJar(buildInfo *paperapi.BuildInfo, filepath string) error {
	if s.downloadJarHandler != nil {
		return s.downloadJarHandler(s, buildInfo, filepath)
	}

	s.ranDownload = true

	return nil
}

func (s *paperServiceMock) DownloadExists(filepath string, buildInfo *paperapi.BuildInfo) (bool, error) {
	if s.downloadExistsHandler != nil {
		return s.downloadExistsHandler(s, filepath, buildInfo)
	}

	return false, nil
}

type fileServiceMock struct {
	fileExistsHandler     func(s *fileServiceMock, filepath string) bool
	deleteIfExistsHandler func(s *fileServiceMock, filepath string) error
	deleteIfExistsCalled  int
}

func (s *fileServiceMock) FileExists(filepath string) bool {
	if s.fileExistsHandler != nil {
		return s.fileExistsHandler(s, filepath)
	}

	return false
}

func (s *fileServiceMock) DeleteIfExists(filepath string) error {
	if s.deleteIfExistsHandler != nil {
		return s.deleteIfExistsHandler(s, filepath)
	}

	s.deleteIfExistsCalled++

	return nil
}

func TestRunMainProgram(t *testing.T) {
	args := []string{"--skip-download"}

	serviceMock := &paperServiceMock{}
	fileService := &fileServiceMock{}

	err := runMainProgram(serviceMock, fileService, args)
	if err != nil && err.Error() != "no builds found" {
		t.Error(err)
	}
}

func TestValidDownloadExistsExitsEarly(t *testing.T) {
	args := []string{""}

	serviceMock := &paperServiceMock{}
	fileService := &fileServiceMock{}

	serviceMock.getLatestBuildHandler = func(s *paperServiceMock, unstable bool, versionPrefix string) (*paperapi.BuildInfo, error) {
		buildInfo := &paperapi.BuildInfo{
			Version: "1.20.2",
			Build:   118,
			Downloads: &paperapi.DownloadInfo{
				Application: &paperapi.ApplicationInfo{
					Name:   "paper.jar",
					Sha256: "asdf",
				},
			},
			Channel: "default",
		}

		return buildInfo, nil
	}

	serviceMock.downloadExistsHandler = func(s *paperServiceMock, filepath string, buildInfo *paperapi.BuildInfo) (bool, error) {
		return true, nil
	}

	err := runMainProgram(serviceMock, fileService, args)
	if err != nil {
		t.Error(err)
	}

	if serviceMock.ranDownload {
		t.Error("Shouldn't have run download method!")
	}
}

func TestSkipDownloadWorks(t *testing.T) {
	args := []string{"--skip-download"}

	serviceMock := &paperServiceMock{}
	fileService := &fileServiceMock{}

	serviceMock.getLatestBuildHandler = func(s *paperServiceMock, unstable bool, versionPrefix string) (*paperapi.BuildInfo, error) {
		buildInfo := &paperapi.BuildInfo{
			Version: "1.20.2",
			Build:   118,
			Downloads: &paperapi.DownloadInfo{
				Application: &paperapi.ApplicationInfo{
					Name:   "paper.jar",
					Sha256: "asdf",
				},
			},
			Channel: "default",
		}

		return buildInfo, nil
	}

	err := runMainProgram(serviceMock, fileService, args)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", serviceMock)
	if serviceMock.ranDownload {
		t.Error("Shouldn't have run download")
	}
}

func TestDeleteIfInvalidDownloadExists(t *testing.T) {
	args := []string{}

	serviceMock := &paperServiceMock{}
	fileService := &fileServiceMock{}

	serviceMock.getLatestBuildHandler = func(s *paperServiceMock, unstable bool, versionPrefix string) (*paperapi.BuildInfo, error) {
		buildInfo := &paperapi.BuildInfo{
			Version: "1.20.2",
			Build:   118,
			Downloads: &paperapi.DownloadInfo{
				Application: &paperapi.ApplicationInfo{
					Name:   "paper.jar",
					Sha256: "asdf",
				},
			},
			Channel: "default",
		}

		return buildInfo, nil
	}

	err := runMainProgram(serviceMock, fileService, args)
	if err != nil {
		t.Error(err)
	}

	if fileService.deleteIfExistsCalled != 1 {
		t.Error("Expected delete if exists to be called once")
	}
}
