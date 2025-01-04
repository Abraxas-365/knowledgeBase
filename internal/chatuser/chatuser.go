package chatuser

type ChatUser struct {
	ID        *string `json:"id" db:"id"`
	Age       int     `json:"age" db:"age"`
	Gender    string  `json:"gender" db:"gender"`
	Ocupation string  `json:"ocupation" db:"occupation"`
	Location  string  `json:"location" db:"location"`
}
