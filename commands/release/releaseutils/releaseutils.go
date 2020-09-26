package releaseutils

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/profclems/glab/internal/utils"
	"github.com/xanzy/go-gitlab"
)

func DisplayAllReleases(rs []*gitlab.Release, repo string) *uitable.Table {
	return utils.DisplayList(utils.ListInfo{
		Name:    "releases",
		Columns: []string{"Name", "Tag", "CreatedAt"},
		Total:   len(rs),
		GetCellValue: func(ri int, ci int) interface{} {
			row := rs[ri]
			switch ci {
			case 0:
				return row.Name
			case 1:
				return row.TagName
			case 2:
				return utils.Gray(utils.TimeToPrettyTimeAgo(*row.CreatedAt))
			default:
				return ""
			}
		},
	}, repo)
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
		utils.Bold(r.Name), r.Author.Name, duration, r.Commit.ShortID, r.TagName, description, utils.Bold("ASSETS"),
		RenderReleaseAssertLinks(r.Assets.Links), utils.Bold("SOURCES"), assetsSources,
	)
}
