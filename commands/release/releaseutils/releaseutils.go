package releaseutils

import (
	"fmt"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/tableprinter"

	"github.com/xanzy/go-gitlab"
)

func DisplayAllReleases(releases []*gitlab.Release, repoName string) string {
	table := tableprinter.NewTablePrinter()
	for _, r := range releases {
		table.AddRow(r.Name, r.TagName, iostreams.Gray(utils.TimeToPrettyTimeAgo(*r.CreatedAt)))
	}

	return table.Render()
}

func RenderReleaseAssertLinks(assets []*gitlab.ReleaseLink) string {
	var assetsPrint string
	if len(assets) == 0 {
		return "There are no assets for this release"
	}
	for _, asset := range assets {
		assetsPrint += asset.URL + "\n"
	}
	return assetsPrint
}

func DisplayRelease(r *gitlab.Release, glamourStyle string) string {
	duration := utils.TimeToPrettyTimeAgo(*r.CreatedAt)
	description, err := utils.RenderMarkdown(r.Description, glamourStyle)
	if err != nil {
		description = r.Description

	}

	var assetsSources string
	for _, asset := range r.Assets.Sources {
		assetsSources += asset.URL + "\n"
	}
	return fmt.Sprintf("%s\n%s released this %s \n%s - %s \n%s \n%s \n%s \n%s \n%s", // whoops
		iostreams.Bold(r.Name), r.Author.Name, duration, r.Commit.ShortID, r.TagName, description, iostreams.Bold("ASSETS"),
		RenderReleaseAssertLinks(r.Assets.Links), iostreams.Bold("SOURCES"), assetsSources,
	)
}
