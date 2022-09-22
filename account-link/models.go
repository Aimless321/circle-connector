package account_link

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	DiscordId uint64 `gorm:"column:discordid;unique;uniqueIndex:idx_discordid_steamid" json:"discordId,omitempty"`
	SteamID64 uint64 `gorm:"column:steamid;unique;uniqueIndex:idx_discordid_steamid" json:"steamId,omitempty"`
}
