package update

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/profclems/glab/internal/request"
)

type ReleaseInfo struct {
	Version     string    `json:"tag_name"`
	PreRelease  bool      `json:"prerelease"`
	URL         string    `json:"html_url"`
	PublishedAt time.Time `json:"published_at"`
}

// GetUpdateInfo checks for latest glab release and returns the ReleaseInfo
func GetUpdateInfo() (ReleaseInfo, error) {
	releasesUrl := "https://api.github.com/repos/profclems/glab/releases/latest"
	resp, err := request.MakeRequest("{}", releasesUrl, "GET")
	var releaseInfo ReleaseInfo
	if err != nil {
		return ReleaseInfo{}, err
	}
	err = json.Unmarshal([]byte(resp), &releaseInfo)
	if err != nil {
		return ReleaseInfo{}, err
	}
	return releaseInfo, nil
}

func isOlderVersion(latestVersion, appVersion string) bool {
	latestVersion = strings.TrimSpace(latestVersion)
	appVersion = strings.TrimSpace(appVersion)

	vv, ve := version.NewVersion(latestVersion)
	vw, we := version.NewVersion(appVersion)

	return ve == nil && we == nil && vv.GreaterThan(vw)
}
