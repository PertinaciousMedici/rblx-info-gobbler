package botStructures

import (
	"regexp"
)

// RequestQueryByUsernames
/*
 * Structure to marshal before the API is queried.
 * Usernames is a list of usernames to fetch the accounts of
 * ExcludeBannedUsers is whether to exclude users banned from Roblox.
 */
type RequestQueryByUsernames struct {
	Usernames          []string `json:"usernames"`
	ExcludeBannedUsers bool     `json:"excludeBannedUsers"`
}

// ResponseQueryByUsernames
/*
 * Unmarshalled response from the Roblox API to a RequestQueryByUsernames
 * Users is a composite array of users associated with the usernames passed.
 */
type ResponseQueryByUsernames struct {
	Users []*FetchedUserByUsernames `json:"data"`
}

// FetchedUserByUsernames
/*
 * Building block of the unmarshalled response from the Roblox API to a RequestQueryByUsernames
 * UserID is the numeric identifier associated with the account.
 * Username is the @username of the account, a unique identifier.
 * DisplayName is a generic username associated with the user, chosen freely.
 */
type FetchedUserByUsernames struct {
	UserID      uint64 `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

// FetchedUserById
/*
 * Unmarshalled response from the Roblox API to a fetch accounts request
 * UserID is the numeric identifier associated with the account.
 * Username is the @username of the account, a unique identifier.
 * DisplayName is a generic username associated with the user, chosen freely.
 * CreationDate is the date and time the account was created.
 * Description is the account's about me field's content.
 * Premium is whether the account possesses a premium subscription.
 * FriendCount is the length of the FetchedRblxUserFriends of the user.
 * AvatarURL is the URL of the user's Roblox avatar, fetched separately.
 */
type FetchedUserById struct {
	UserID       string `json:"id"`
	Username     string `json:"name"`
	DisplayName  string `json:"displayName"`
	CreationDate string `json:"createTime"`
	Description  string `json:"about"`
	Premium      bool   `json:"premium"`
	FriendCount  int
	AvatarURL    string
}

// FetchedUserGroupMemberships
/*
 * FetchedUserGroupMemberships is an abstraction around a slice of FetchedUserMembership structs.
 * Groups is the slice around which the struct is built.
 */
type FetchedUserGroupMemberships struct {
	Groups []*FetchedUserMembership `json:"data"`
}

// FetchedUserMembership
/*
 * FetchedUserMembership establishes a connection between a user and a group and thus contains relational data.
 * Group is the data regarding the group itself, and is a struct of type FetchedRblxGroup.
 * Role is the data regarding the rank in said group, and is a struct of type FetchedRblxGroupRole.
 */
type FetchedUserMembership struct {
	Group *FetchedRblxGroup     `json:"group"`
	Role  *FetchedRblxGroupRole `json:"role"`
}

// FetchedFullRblxGroup
/*
 * FetchedRblxGroup with considerably more information, fetched through a direct query to group info.
 * GroupID is an uint64 identifier assigned to the Roblox group, through which one can fetch its data.
 * GroupOwner is the FetchedUserById.Username of the owner of the fetched Roblox group.
 * Description is the text description assigned to the Roblox group.
 * CreationDate is the date and time on which the group was created.
 * UpdateTime is the date and time on which the group was last updated (shouts, comments, et cetera).
 * MemberCount is the total amount of users in the Roblox group.
 * LockedGroup is whether the Roblox group is invite-only.
 */
type FetchedFullRblxGroup struct {
	GroupID      uint64 `json:"id"`
	GroupName    string `json:"displayName"`
	GroupOwner   string `json:"owner"`
	Description  string `json:"description"`
	CreationDate string `json:"creationDate"`
	UpdateTime   string `json:"updateTime"`
	MemberCount  uint64 `json:"memberCount"`
	LockedGroup  bool   `json:"lockedGroup"`
}

func (group *FetchedFullRblxGroup) GetOwnerID() string {
	regex := regexp.MustCompile(`\d+`)
	matches := regex.FindString(group.GroupOwner)
	return matches
}

// FetchedRblxGroup
/*
 * FetchedRblxGroup contains data fetched through a FetchedUserMembership query.
 * GroupID is an uint64 identifier assigned to the Roblox group, through which one can fetch its data.
 * GroupName is the unique alphabetic identifier of a Roblox group.
 * MemberCount is amount of members that a Roblox group possesses.
 */
type FetchedRblxGroup struct {
	GroupID     uint64 `json:"id"`
	GroupName   string `json:"name"`
	MemberCount uint64 `json:"memberCount"`
}

// FetchedRblxGroupRole
/*
 * FetchedRblxGroupRole contains data fetched through a FetchedUserMembership query.
 * RankName is the name of the role the user holds within the FetchedRblxGroup.
 * RankID is the numeric identifier (1-255) of the rank within the group hierarchy.
 */
type FetchedRblxGroupRole struct {
	RankName string `json:"name"`
	RankID   uint64 `json:"rank"`
}

// FetchedRblxUserFriends
/*
 * FetchedRblxUserFriends is an abstraction around a slice of FetchedRblxFriendship structs.
 * Friends is the slice around which the struct is built.
 */
type FetchedRblxUserFriends struct {
	Friends []*FetchedRblxFriendship `json:"data"`
}

// FetchedRblxFriendship
/*
 * FetchedRblxFriendship establishes a friendship relation between two users, used by FetchedRblxUserFriends structs.
 * UserID is the numeric identifier associated with the account.
 * Username is the @username of the account, a unique identifier.
 * DisplayName is a generic username associated with the user, chosen freely.
 */
type FetchedRblxFriendship struct {
	UserID      uint64 `json:"id"`
	Username    string `json:"name"`
	DisplayName string `json:"displayName"`
}

// FetchedRblxUserAvatar
/*
 * Simple abstraction around a look-up avatar request, using anonymous structs.
 * QueryResult is an abstraction around the image URL.
 * ImageURL is the avatar per se.
 */
type FetchedRblxUserAvatar struct {
	QueryResult []struct {
		ImageURL string `json:"imageUrl"`
	} `json:"data"`
}

// FetchedRblxUserBadges
/*
 * Wrapper struct for badge look-up requests and iteration.
 * CursorPrevious is a snowflake for cross-request pagination.
 * CursorNext is a snowflake for cross-request pagination.
 * Badges is the array to the pointers of the badge data structures.
 */
type FetchedRblxUserBadges struct {
	CursorPrevious string                  `json:"previousPageCursor"`
	CursorNext     string                  `json:"nextPageCursor"`
	Badges         []*FetchedRblxUserBadge `json:"data"`
}

// FetchedRblxUserBadge
/*
 * Simple abstraction around an awarded game badge.
 * BadgeID, its numeric identifier, used to fetch the award dates.
 * SecondaryID, is also the numeric identifier but with separate JSON keys, since endpoints differ.
 * AwardedDate, is the date the badge was awarded to the user it belongs to.
 */
type FetchedRblxUserBadge struct {
	BadgeID     uint64 `json:"id"`
	SecondaryID uint64 `json:"badgeId"`
	AwardedDate string `json:"awardedDate"`
}
