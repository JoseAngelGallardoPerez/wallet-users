package models

import "time"

type Invite struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	Code      string    `gorm:"column:code;unique_index" json:"code"`
	To        string    `gorm:"column:to" json:"to"`
	Uses      uint64    `gorm:"column:uses" json:"uses"`
	MaxUsages uint64    `gorm:"column:max_usages" json:"maxUsages"`
	UserUID   string    `json:"userUID"`
	User      *User     `gorm:"foreignkey:UserUID;association_foreignkey:UID;association_autoupdate:false" json:"user"`
	ExpiresAt time.Time `gorm:"column:expires_at" json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CanBeRedeemed checks if an invite code can be used redeemed.
func (i *Invite) CanBeRedeemed() bool {
	return !i.IsExpired() && !i.IsSoldOut() && !i.HasRestrictedUsage()
}

// UsageRestrictedToEmail check if an invite is to the user who has the specified email.
func (i *Invite) UsageRestrictedToEmail(email string) bool {
	return i.To == email
}

// HasRestrictedUsage checks if an invite is usable for only one person.
func (i *Invite) HasRestrictedUsage() bool {
	return i.To != ""
}

// HasRestrictedUsage checks if an invite is usable for only one person.
func (i *Invite) IsExpired() bool {
	if i.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().Unix() <= i.ExpiresAt.Unix()
}

// IsSoldOut checks if the invite code has been sold out.
func (i *Invite) IsSoldOut() bool {
	if i.MaxUsages == 0 {
		return false
	}
	return i.Uses >= i.MaxUsages
}

func (i *Invite) CanBeUsedOnce() *Invite {
	i.MaxUsages = 1
	return i
}

// RestrictUsageTo set the user who can use this invite.
func (i *Invite) RestrictUsageTo(email string) *Invite {
	i.To = email
	return i
}

// Set the invite expiration date.
func (i *Invite) ExpiresDate(t time.Time) *Invite {
	i.ExpiresAt = t
	return i
}

// Set the expiration date to days from now.
func (i *Invite) ExpiresIn(days int) *Invite {
	i.ExpiresAt = time.Now().Local().AddDate(0, 0, days)
	return i
}
