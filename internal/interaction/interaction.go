package interaction

type Interaction struct {
	ID                 int      `json:"id" db:"id"`
	UserChatID         string   `json:"user_chat_id" db:"user_chat_id"`
	ContextInteraction []string `json:"context_interaction" db:"context_interaction"`
}
