package action

import (
	"github.com/profclems/glab/internal/config"
	"github.com/rsteube/carapace"
)

func ActionConfigAliases() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		config, err := config.Init()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}
		configAliases, err := config.Aliases()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		vals := make([]string, 0)
		for alias, desc := range configAliases.All() {
			vals = append(vals, alias, desc)
		}
		return carapace.ActionValuesDescribed(vals...)
	})
}

func ActionConfigHosts() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		config, err := config.Init()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}
		configHosts, err := config.Hosts()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}
		return carapace.ActionValues(configHosts...)
	})
}

func ActionConfigKeys() carapace.Action {
	return carapace.ActionValuesDescribed(
		"token", "Your gitlab access token, defaults to environment variables",
		"gitlab_uri", "if unset, defaults to https://gitlab.com",
		"browser", "if unset, defaults to environment variables",
		"editor", "if unset, defaults to environment variables.",
		"visual", "alternative for editor. if unset, defaults to environment variables.",
		"glamour_style", "Your desired markdown renderer style.",
	)
}

func ActionConfigValues(key string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		actions := map[string]carapace.Action{
			"token":         carapace.ActionValues(),
			"gitlab_uri":    carapace.ActionValues(),
			"browser":       carapace.ActionFiles(),
			"editor":        carapace.ActionFiles(),
			"visual":        carapace.ActionFiles(),
			"glamour_style": carapace.ActionValues("dark", "light", "notty"),
		}

		return actions[key]
	})
}
