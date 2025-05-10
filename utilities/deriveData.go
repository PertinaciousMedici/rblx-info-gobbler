package utilities

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

const (
	EmployeeFlag                  = 1
	PartneredServerOwnerFlag      = 2
	HypeSquadEventsFlag           = 4
	BugHunterLevel1Flag           = 8
	HouseBraveryFlag              = 64
	HouseBrillianceFlag           = 128
	HouseBalanceFlag              = 256
	EarlySupporterFlag            = 512
	BugHunterLevel2Flag           = 1024
	EarlyVerifiedBotDeveloperFlag = 131072
)

func CalculateMaxPages(maxPerPage int, sizeOfCollection int) int {
	if maxPerPage <= 0 {
		return 1
	}
	if sizeOfCollection == 0 {
		return 1
	}
	return (sizeOfCollection + maxPerPage - 1) / maxPerPage
}

func deriveCreationDate(discID string) string {
	IntID, _ := strconv.ParseInt(discID, 10, 64)
	Timestamp := (IntID >> 22) + int64(1420070400000)
	UnixTime := time.Unix(0, Timestamp*int64(time.Millisecond))
	return UnixTime.Format(time.RFC3339)
}

func checkFlags(flagsNumber uint64) ([]string, bool) {
	var matchedFlags []string
	var hasFlags bool

	if (flagsNumber & EmployeeFlag) != 0 {
		matchedFlags = append(matchedFlags, "Discord Employee")
	}

	if (flagsNumber & PartneredServerOwnerFlag) != 0 {
		matchedFlags = append(matchedFlags, "Partnered Server Owner")
	}

	if (flagsNumber & HypeSquadEventsFlag) != 0 {
		matchedFlags = append(matchedFlags, "HypeSquad Events")
	}

	if (flagsNumber & BugHunterLevel1Flag) != 0 {
		matchedFlags = append(matchedFlags, "Bug H. Lvl1")
	}

	if (flagsNumber & HouseBraveryFlag) != 0 {
		matchedFlags = append(matchedFlags, "House Bravery")
	}

	if (flagsNumber & HouseBrillianceFlag) != 0 {
		matchedFlags = append(matchedFlags, "House Brilliance")
	}

	if (flagsNumber & HouseBalanceFlag) != 0 {
		matchedFlags = append(matchedFlags, "House Balance")
	}

	if (flagsNumber & EarlySupporterFlag) != 0 {
		matchedFlags = append(matchedFlags, "Early Supporter")
	}

	if (flagsNumber & BugHunterLevel2Flag) != 0 {
		matchedFlags = append(matchedFlags, "Bug H. Lvl2")
	}

	if (flagsNumber & EarlyVerifiedBotDeveloperFlag) != 0 {
		matchedFlags = append(matchedFlags, "E.V.B.D.")
	}

	hasFlags = len(matchedFlags) != 0
	return matchedFlags, hasFlags
}

func ExtractUsernames(description string) (string, string) {
	re := regexp.MustCompile(`\*\*(.*?)\*\* & \*\*(.*?)\*\*`)
	matches := re.FindStringSubmatch(description)

	if len(matches) == 3 {
		return matches[1], matches[2]
	}

	return "", ""
}

func BuildIterator(typeof string, currentPage int, embed *discord.Embed) *api.SendMessageData {
	IDFirst := fmt.Sprintf("%s_first", typeof)
	IDPRev := fmt.Sprintf("%s_prev_%d", typeof, currentPage)
	IDNext := fmt.Sprintf("%s_next_%d", typeof, currentPage)
	IDLast := fmt.Sprintf("%s_last", typeof)

	complexMessage := api.SendMessageData{
		Embeds: []discord.Embed{*embed},
		Components: []discord.ContainerComponent{
			&discord.ActionRowComponent{
				&discord.ButtonComponent{
					Style:    discord.SecondaryButtonStyle(),
					CustomID: discord.ComponentID(IDFirst),
					Label:    "<<",
				},
				&discord.ButtonComponent{
					Style:    discord.SecondaryButtonStyle(),
					CustomID: discord.ComponentID(IDPRev),
					Label:    "<",
				},
				&discord.ButtonComponent{
					Style:    discord.SecondaryButtonStyle(),
					CustomID: discord.ComponentID(IDNext),
					Label:    ">",
				},
				&discord.ButtonComponent{
					Style:    discord.SecondaryButtonStyle(),
					CustomID: discord.ComponentID(IDLast),
					Label:    ">>",
				},
			},
		},
	}

	return &complexMessage
}

func PaginateArray(arr []interface{}, maxPerPage int, currentPage int) []interface{} {
	totalPages := CalculateMaxPages(maxPerPage, len(arr))

	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}

	startIndex := (currentPage - 1) * maxPerPage
	endIndex := startIndex + maxPerPage

	if endIndex > len(arr) {
		endIndex = len(arr)
	}

	return arr[startIndex:endIndex]
}
