package action

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func EnableRepoOverride(cmd *cobra.Command, f *cmdutils.Factory) { // TODO factory would cause circular dependency in cmdutils
	carapace.Gen(cmd).FlagCompletion(carapace.ActionMap{
		"repo": ActionRepo(cmd, f),
	})
}

func ActionRepo(cmd *cobra.Command, f *cmdutils.Factory) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		config, err := config.Init()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}
		configHosts, err := config.Hosts()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		if strings.Contains(c.CallbackValue, "/") {
			isKnownHost := false
			for _, host := range configHosts {
				if strings.HasPrefix(c.CallbackValue, host) {
					isKnownHost = true
					break
				}
			}
			if !isKnownHost {
				return carapace.ActionValues() // only complete full host style for simplicity
			}
		}

		return carapace.ActionMultiParts("/", func(c carapace.Context) carapace.Action {
			if len(c.Parts) > 0 {
				if err := f.RepoOverride(fmt.Sprintf("%v/fake/repo", c.Parts[0])); err != nil {
					return carapace.ActionMessage(err.Error())
				}
			}

			switch len(c.Parts) {
			case 0:
				return carapace.ActionValues(configHosts...).Invoke(c).Suffix("/").ToA()
			case 1:
				users := ActionUsers(cmd, f, &gitlab.ListUsersOptions{}).Invoke(c).Suffix("/")
				groups := ActionGroups(cmd, f, &gitlab.ListGroupsOptions{}).Invoke(c).Suffix("/")
				return users.Merge(groups).ToA()
			case 2:
				subgroups := ActionSubgroups(cmd, f, c.Parts[1], &gitlab.ListSubgroupsOptions{}).Invoke(c).Suffix("/")
				groupProjects := ActionGroupProjects(cmd, f, c.Parts[1], &gitlab.ListGroupProjectsOptions{}).Invoke(c)
				userProjects := ActionUserProjects(cmd, f, c.Parts[1], &gitlab.ListProjectsOptions{}).Invoke(c)
				return subgroups.Merge(groupProjects, userProjects).ToA()
			case 3:
				groupProjects := ActionGroupProjects(cmd, f, strings.Join(c.Parts[1:], "/"), &gitlab.ListGroupProjectsOptions{}).Invoke(c)
				return groupProjects.ToA()
			default:
				return carapace.ActionValues()
			}
		})
	})
}

func ActionApiCallback(cmd *cobra.Command, f *cmdutils.Factory, callback func(client *gitlab.Client, c carapace.Context) carapace.Action) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if repoFlag := cmd.Flag("repo"); repoFlag != nil && repoFlag.Changed {
			f.RepoOverride("https://" + repoFlag.Value.String())
		}

		client, err := f.HttpClient()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		return callback(client, c)
	})
}
