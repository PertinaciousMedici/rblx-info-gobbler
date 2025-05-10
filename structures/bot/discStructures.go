package botStructures

// DiscFetchedUserById
/*
 * Unmarshalled response from the Discord API to a user query request.
 * UserID is an uint64 snowflake, from which one can fetch an account or derive its creation date.
 * Username is a unique alphabetic identifier chosen by each user.
 * Avatar is a hash of the account's current avatar attachment, necessary to fetch the image from a CDN.
 * Flags is a bitmask of a user's badges and achievements.
 */
type DiscFetchedUserById struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Flags    uint64 `json:"public_flags"`
}
