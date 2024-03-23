package paperapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"testing"
)

func TestFilterVersions(t *testing.T) {
	input := &VersionsList{
		Versions: []string{"1.20", "1.2", "1.23.2", "1.19.0", "1.2.4"},
	}

	output := filterVersions("1.2", input)

	failed := false

	if len(output.Versions) != 2 {
		t.Logf("Expected two versions to be in the outputt but there was %d", len(output.Versions))
		failed = true
	}

	if !slices.Contains[[]string](output.Versions, "1.2") || !slices.Contains[[]string](output.Versions, "1.2.4") {
		t.Log("Expected versions weren't in the output.")
		failed = true
	}

	if failed {
		t.Fail()
	}
}

func TestIsValidDownload(t *testing.T) {
	expected := "d1bc8d3ba4afc7e109612cb73acbdddac052c93025aa1f82942edabb7deb82a1"
	fileName := ".sha256test"

	err := setupTestFile(fileName)
	if err != nil {
		t.Error(err)
	}

	defer cleanupTestFile(fileName)

	valid, err := newServiceImpl(nil, nil, nil, nil, "").IsValidDownload(fileName, expected)
	if err != nil {
		t.Error(err)
	}

	if !valid {
		t.Errorf("hash didn't match expected result.")
	}
}

func cleanupTestFile(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		os.Remove(fileName)
	}
}

func setupTestFile(fileName string) error {
	fileContents := "asdf\n"

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write([]byte(fileContents))

	if err != nil {
		return err
	}

	err = file.Sync()

	return err
}

type buildsListServiceMock struct {
	getBuildsListHandler func(s buildsListServiceMock, version string) (*BuildsList, error)
}

func (s buildsListServiceMock) GetBuildsList(version string) (*BuildsList, error) {
	return s.getBuildsListHandler(s, version)
}

type buildInfoServiceMock struct {
	getBuildInfoHandler func(s buildInfoServiceMock, version string, build int) (*BuildInfo, error)
}

func (s buildInfoServiceMock) GetBuildInfo(version string, build int) (*BuildInfo, error) {
	return s.getBuildInfoHandler(s, version, build)
}

type versionsListServiceMock struct {
	getVersionsListHandler func(s versionsListServiceMock) (*VersionsList, error)
}

func (s versionsListServiceMock) GetVersionsList() (*VersionsList, error) {
	return s.getVersionsListHandler(s)
}

func TestGetLatestBuildInfo(t *testing.T) {
	buildsListMock := buildsListServiceMock{}

	buildsListMock.getBuildsListHandler = func(s buildsListServiceMock, version string) (*BuildsList, error) {
		if version == "1.20.2" {
			buildsList := *&BuildsList{
				Builds:  []int{2, 3, 4, 5, 7, 8, 9},
				Version: version,
			}

			return &buildsList, nil
		}

		return nil, nil
	}

	buildInfoMock := buildInfoServiceMock{}

	buildInfoMock.getBuildInfoHandler = func(s buildInfoServiceMock, version string, build int) (*BuildInfo, error) {
		if version == "1.20.2" && build == 9 {
			applicationInfo := &ApplicationInfo{
				Name:   "asdf",
				Sha256: "1234",
			}

			downloadInfo := &DownloadInfo{
				Application: applicationInfo,
			}

			buildInfo := &BuildInfo{
				Version:   version,
				Build:     build,
				Channel:   "default",
				Downloads: downloadInfo,
			}

			return buildInfo, nil
		}

		return nil, nil
	}

	versionsListMock := versionsListServiceMock{}
	versionsListMock.getVersionsListHandler = func(s versionsListServiceMock) (*VersionsList, error) {
		return &VersionsList{Versions: []string{}}, nil
	}

	service := newServiceImpl(buildInfoMock, nil, buildsListMock, nil, "")

	buildInfo, err := service.getLatestBuildInfo("1.20.2")
	if err != nil {
		t.Error(err)
	}

	if buildInfo.Version != "1.20.2" {
		t.Error("expected version to be 1.20.2")
	}

	if buildInfo.Build != 9 {
		t.Error("expected build number to be 9")
	}
}

func TestGetLatestBuild(t *testing.T) {
	buildsListMock := buildsListServiceMock{
		getBuildsListHandler: handleGetBuildsListForTestGetLatestBuild,
	}

	buildsInfoMock := buildInfoServiceMock{
		getBuildInfoHandler: handleGetBuildInfoForTestGetLatestBuild,
	}

	versionsListMock := versionsListServiceMock{
		getVersionsListHandler: handleGetVersionsForTestGetLatestBuild,
	}

	service := newServiceImpl(buildsInfoMock, versionsListMock, buildsListMock, nil, "")

	buildInfo, err := service.GetLatestBuild(false, "")
	if err != nil {
		t.Error(err)
	}

	expectedVersion := "1.20.1"
	if buildInfo.Version != expectedVersion {
		t.Errorf("Expected latest stable version to be %s but it was %s", expectedVersion, buildInfo.Version)
	}

	if buildInfo.Build != 3 {
		t.Errorf("Expected build number to be 3, but it was %d", buildInfo.Build)
	}

	buildInfo, err = service.GetLatestBuild(true, "")
	if err != nil {
		t.Error(err)
	}

	expectedVersion = "1.20.3"
	if buildInfo.Version != expectedVersion {
		t.Errorf("Expected latest stable version to be %s but it was %s", expectedVersion, buildInfo.Version)
	}

	if buildInfo.Build != 3 {
		t.Errorf("Expected build number to be 3, but it was %d", buildInfo.Build)
	}

	buildInfo, err = service.GetLatestBuild(true, "1.19")
	if err != nil {
		t.Error(err)
	}

	expectedVersion = "1.19.0"
	if buildInfo.Version != expectedVersion {
		t.Errorf("Expected latest stable version to be %s but it was %s", expectedVersion, buildInfo.Version)
	}

	if buildInfo.Build != 3 {
		t.Errorf("Expected build number to be 3, but it was %d", buildInfo.Build)
	}
}

func handleGetBuildsListForTestGetLatestBuild(s buildsListServiceMock, version string) (*BuildsList, error) {
	buildsList := &BuildsList{
		Version: version,
		Builds: []int{
			1,
			2,
			3,
		},
	}

	return buildsList, nil
}

func handleGetBuildInfoForTestGetLatestBuild(s buildInfoServiceMock, version string, build int) (*BuildInfo, error) {
	if version == "1.20.3" {
		return &BuildInfo{
			Version: version,
			Build:   build,
			Channel: "experimental",
		}, nil
	}

	return &BuildInfo{
		Version: version,
		Build:   build,
		Channel: "default",
	}, nil
}

func handleGetVersionsForTestGetLatestBuild(s versionsListServiceMock) (*VersionsList, error) {
	versionsList := &VersionsList{
		Versions: []string{
			"1.19.0",
			"1.20.1",
			"1.20.3",
		},
	}

	return versionsList, nil
}

func TestDownloadFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := "asdf\n"

		fmt.Fprint(w, data)
	}))

	service := newServiceImpl(nil, nil, nil, nil, ts.URL)

	buildInfo := &BuildInfo{
		Version: "1.2.3",
		Build:   123,
		Downloads: &DownloadInfo{
			Application: &ApplicationInfo{
				Name:   "paper.jar",
				Sha256: "asdf",
			},
		},
	}

	filename := ".testfile"

	err := service.DownloadJar(buildInfo, filename)
	if err != nil {
		t.Error(err)
	}

	defer cleanupTestFile(filename)

	expected := "d1bc8d3ba4afc7e109612cb73acbdddac052c93025aa1f82942edabb7deb82a1"

	valid, err := service.IsValidDownload(filename, expected)
	if err != nil {
		t.Error(err)
	}

	if !valid {
		t.Error("Download isn't valid")
	}
}
