package bot

import (
	"fmt"

	"PanoptisMouthNew/bot/commands"
	. "PanoptisMouthNew/structures/bot"
)

const (
	Utilities   = "Utilities"
	Information = "Information"
	Moderation  = "Moderation"
)

func SetHooks(botInstance *BotInstance) {
	botInstance.InstanceMutex.Lock()
	botInstance.CommandHooks["help"] = Command{
		Name:              "help",
		Usage:             fmt.Sprintf("%shelp [cmd]", botPrefix),
		Description:       "Get information on command usage",
		Category:          Information,
		ExpectedArguments: 0,
		Execute:           commands.Help,
	}
	botInstance.CommandHooks["softarchive"] = Command{
		Name:              "softarchive",
		Usage:             fmt.Sprintf("%ssoftarchive", botPrefix),
		Description:       "Delete all messages except for those pinned",
		Category:          Moderation,
		ExpectedArguments: 0,
		Execute:           commands.SoftArchive,
	}
	botInstance.CommandHooks["search"] = Command{
		Name:              "search",
		Usage:             fmt.Sprintf("%ssearch @user1|userID|self", botPrefix),
		Description:       "Search for a Discord user's Roblox account",
		Category:          Utilities,
		ExpectedArguments: 1,
		Execute:           commands.Search,
	}
	botInstance.CommandHooks["rsearch"] = Command{
		Name:              "rsearch",
		Usage:             fmt.Sprintf("%srsearch username|userID", botPrefix),
		Description:       "Search for a Roblox user's Discord account",
		Category:          Utilities,
		ExpectedArguments: 1,
		Execute:           commands.RSearch,
	}
	botInstance.CommandHooks["qsearch"] = Command{
		Name:              "qsearch",
		Usage:             fmt.Sprintf("%sqsearch username|userID", botPrefix),
		Description:       "Query Roblox for information on a user",
		Category:          Utilities,
		ExpectedArguments: 1,
		Execute:           commands.QSearch,
	}
	botInstance.CommandHooks["groups"] = Command{
		Name:              "groups",
		Usage:             fmt.Sprintf("%sgroups username|userID", botPrefix),
		Description:       "Query Roblox for a user's groups",
		Category:          Utilities,
		ExpectedArguments: 1,
		Execute:           commands.Groups,
	}
	botInstance.CommandHooks["friends"] = Command{
		Name:              "friends",
		Usage:             fmt.Sprintf("%sfriends username|userID", botPrefix),
		Description:       "Query Roblox for a user's friends",
		Category:          Utilities,
		ExpectedArguments: 1,
		Execute:           commands.Friends,
	}
	botInstance.CommandHooks["mgroups"] = Command{
		Name:              "mgroups",
		Usage:             fmt.Sprintf("%smgroups user1 user2", botPrefix),
		Description:       "Query Roblox for the groups shared by two users",
		Category:          Utilities,
		ExpectedArguments: 2,
		Execute:           commands.MGroups,
	}
	botInstance.CommandHooks["mfriends"] = Command{
		Name:              "mfriends",
		Usage:             fmt.Sprintf("%smfriends user1 user2", botPrefix),
		Description:       "Query Roblox for the friends shared by two users",
		Category:          Utilities,
		ExpectedArguments: 2,
		Execute:           commands.MFriends,
	}
	botInstance.InstanceMutex.Unlock()
}
