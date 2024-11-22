package interaction

import "context"

type Repository interface {
	CreateInteraction(ctx context.Context, i Interaction) (*Interaction, error)
}
