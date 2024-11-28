package chatuser

type ChatUser struct {
	ID         *string `json:"id" db:"id"`
	Age        int     `json:"age" db:"age"`
	Gender     string  `json:"gender" db:"gender"`
	Occupation string  `json:"occupation" db:"occupation"`
}
