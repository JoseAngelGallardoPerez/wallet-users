package validators

type PatchUser struct {
	FirstName *string `json:"firstName" binding:"omitempty,max=255"`
	LastName  *string `json:"lastName" binding:"omitempty,max=255"`
	Nickname  *string `json:"nickname" binding:"omitempty,max=100"`
}
