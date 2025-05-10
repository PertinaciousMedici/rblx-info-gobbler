package botStructures

// RwFetchedUserByReverseSearch
/*
 * Unmarshalled response for a Discord user look-up through the RoWifi API (reverse-search).
 * DiscordID is an uint-64 snowflake, from which one can fetch the user or derive the account's creation date.
 */
type RwFetchedUserByReverseSearch struct {
	DiscordID string `json:"discord_id"`
}

// RwFetchedUserByRegularSearch
/*
 * Unmarshalled response for a Roblox user look-up through the RoWifi API (regular-search).
 * RobloxID is an uint-64 snowflake, from which one can fetch the user.
 */
type RwFetchedUserByRegularSearch struct {
	RobloxID uint64 `json:"roblox_id"`
}
