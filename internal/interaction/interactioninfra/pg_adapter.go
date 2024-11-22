package interactioninfra

import (
	"context"
	"fmt"

	"github.com/Abraxas-365/opd/internal/interaction"
	"github.com/Abraxas-365/toolkit/pkg/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgresStore struct {
	db *sqlx.DB
}

// NewInteractionStore creates a new PostgresStore for interaction repository
func NewInteractionStore(db *sqlx.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// CreateInteraction inserts a new interaction
func (s *PostgresStore) CreateInteraction(ctx context.Context, i interaction.Interaction) (*interaction.Interaction, error) {
	query := `
		INSERT INTO interactions (user_chat_id, context_interaction) 
		VALUES ($1, $2) 
		RETURNING id, user_chat_id, context_interaction`

	err := s.db.QueryRowContext(
		ctx,
		query,
		i.UserChatID,
		pq.Array(i.ContextInteraction),
	).Scan(
		&i.ID,
		&i.UserChatID,
		pq.Array(&i.ContextInteraction),
	)

	if err != nil {
		// Check if it's a foreign key violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, errors.ErrNotFound("Referenced chat user not found")
		}
		return nil, errors.ErrDatabase(fmt.Sprintf("Failed to create interaction: %v", err))
	}

	return &i, nil
}
