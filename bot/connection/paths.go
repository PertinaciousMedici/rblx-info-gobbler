package connection

const GET = "GET"
const POST = "POST"

const RwRoot = "https://api.rowifi.xyz/v3/"
const RwReverseSearchGet = RwRoot + "guilds/%s/members/roblox/%s"
const RwRegularSearchGet = RwRoot + "guilds/%s/members/%s"

const DsRoot = "https://discord.com/api/v10/"
const DsUserSearchGet = DsRoot + "users/%s"

const DsMediaRoot = "https://cdn.discordapp.com/"
const DsAvatarSearchGet = DsMediaRoot + "avatars/%s/%s.png"

const RblxUserSearchGet = "https://apis.roblox.com/cloud/v2/users/%s"
const RblxGroupSearchGet = "https://apis.roblox.com/cloud/v2/groups%s"
const RblxGroupView = "https://www.roblox.com/communities/%d"
const RblxUserView = "https://www.roblox.com/users/%d/profile"
const RblxUserGroupsSearchGet = "https://groups.roblox.com/v2/users/%s/groups/roles"
const RblxUserFriendsSearchGet = "https://friends.roblox.com/v1/users/%s/friends"
const RblxUsernameSearchPost = "https://users.roblox.com/v1/usernames/users"
const RblxPastUsernameSearchGet = "https://users.roblox.com/v1/users/%s/username-history?limit=25"
const RblxAvatarSearchGet = "https://thumbnails.roblox.com/v1/users/avatar?userIds=%s&size=420x420&format=Png&isCircular=true"
const RblxUserBadgeSearchGet = "https://badges.roblox.com/v1/users/%s/badges?limit=%d&cursor=%s"
const RblxBadgeTimestampSearchGet = "https://badges.roblox.com/v1/users/%s/badges/awarded-dates?badgeIds="
