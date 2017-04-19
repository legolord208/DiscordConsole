package main

import (
	"github.com/chzyer/readline"
)

func setCompleter(rl *readline.Instance) {
	rl.Config.AutoComplete = readline.NewPrefixCompleter(
		readline.PcItem("guild", readline.PcItemDynamic(func(name string) []string {
			names := make([]string, len(cacheGuilds))
			for i, g := range cacheGuilds {
				names[i] = g.Name
			}
			return names
		})),
		readline.PcItem("channel", readline.PcItemDynamic(func(name string) []string {
			names := make([]string, len(cacheChannels))
			for i, c := range cacheChannels {
				names[i] = c.Name
			}
			return names
		})),

		readline.PcItem("edit", readline.PcItemDynamic(singleValue(&lastUsedMsg))),
		readline.PcItem("del", readline.PcItemDynamic(singleValue(&lastUsedMsg))),
		readline.PcItem("quote", readline.PcItemDynamic(singleValue(&lastUsedMsg))),
		readline.PcItem("reactadd", readline.PcItemDynamic(singleValue(&lastUsedMsg))),
		readline.PcItem("reactdel", readline.PcItemDynamic(singleValue(&lastUsedMsg))),

		readline.PcItem("roleedit", readline.PcItemDynamic(singleValue(&lastUsedRole))),
		readline.PcItem("roledelete", readline.PcItemDynamic(singleValue(&lastUsedRole))),

		readline.PcItem("roleadd", readline.PcItem(userID,
			readline.PcItemDynamic(singleValue(&lastUsedRole)),
		)),
		readline.PcItem("roledel", readline.PcItem(userID,
			readline.PcItemDynamic(singleValue(&lastUsedRole)),
		)),

		readline.PcItem("bookmark", readline.PcItemDynamic(bookmarkTab)),
		readline.PcItem("go", readline.PcItemDynamic(bookmarkTab)),

		readline.PcItem("read", readline.PcItemDynamic(singleValue(&lastUsedMsg))),
		readline.PcItem("cinfo",
			readline.PcItem("guild"),
			readline.PcItem("name"),
			readline.PcItem("topic"),
			readline.PcItem("type"),
		),
		readline.PcItem("ginfo",
			readline.PcItem("name"),
			readline.PcItem("icon"),
			readline.PcItem("region"),
			readline.PcItem("owner"),
			readline.PcItem("splash"),
			readline.PcItem("members"),
			readline.PcItem("level"),
		),
		readline.PcItem("messages",
			readline.PcItem("all"),
			readline.PcItem("mentions"),
			readline.PcItem("private"),
			readline.PcItem("current"),
			readline.PcItem("none"),
		),
		readline.PcItem("intercept",
			readline.PcItem("true"),
			readline.PcItem("false"),
		),
		readline.PcItem("output",
			readline.PcItem("true"),
			readline.PcItem("false"),
		),
	)
}
func bookmarkTab(line string) []string {
	items := make([]string, len(bookmarks))

	i := 0
	for key := range bookmarks {
		items[i] = key
		i++
	}
	return items
}
func singleValue(val *string) func(string) []string {
	return func(line string) []string {
		return []string{*val}
	}
}
