package user

type User struct {
	ID         string `json:"id" db:"id"`
	Email      string `json:"email" db:"email"`
	IsAdmin    bool   `json:"isAdmin" db:"is_admin"`
	ProviderID string `json:"-" db:"provider_id"`
	Provider   string `json:"-" db:"provider"`
}

func (au *User) GetID() string {
	return au.ID
}
