package room

import (
	"context"
	"net/http"
	"strconv"

	"eqrx.net/matrix"
)

// ReactionContent is the payload of a reaction.
type ReactionContent struct {
	RelatesTo MessageReference `json:"m.relates_to"`
}

// SendReaction sends a reaction (emoticon) to a given source event id.
func (r *Room) SendReaction(ctx context.Context, source, key string) error {
	txID := strconv.FormatUint(r.server.TxID(), 10)
	path := "/_matrix/client/v3/rooms/" + r.id + "/send/m.reaction/" + txID
	reaction := ReactionContent{MessageReference{source, key, "m.annotation"}}

	var response matrix.Response

	status, err := r.server.HTTP(ctx, http.MethodPut, path, reaction, &response)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if err := response.AsError(status); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}
