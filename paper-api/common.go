package paperapi

import "github.com/sprpgmr/papermc-fetch/files"

const baseURL = "https://api.papermc.io/v2/projects/paper"

// GetPaperAPIService builds dependencies and passes them into the PaperApiServiceImpl for use
func GetPaperAPIService() Service {

	versionsListService := newVersionsListServiceImpl(baseURL)
	buildsListService := newBuildsListServiceImpl(baseURL)
	buildInfoService := newBuildInfoServiceImpl(baseURL)

	return newServiceImpl(buildInfoService, versionsListService, buildsListService, files.GetFileService(), baseURL)
}
