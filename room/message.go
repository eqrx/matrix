// Copyright (C) 2022 Alexander Sowitzki
//
// This program is free software: you can redistribute it and/or modify it under the terms of the
// GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
// warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
// details.
//
// You should have received a copy of the GNU Affero General Public License along with this program.
// If not, see <https://www.gnu.org/licenses/>.

package room

import (
	"context"
	"net/http"
	"strconv"

	"eqrx.net/matrix"
)

// MessageEventType indicates a room event is a message.
const MessageEventType = "m.room.message"

// MessageEvent is a message sent to a room.
type MessageEvent struct {
	EventWithoutContent
	Content MessageContent `struct:"content"`
}

// MessageContent represents the content of room message.
type MessageContent struct {
	MessageType   string          `json:"msgtype"`
	Body          string          `json:"body"`
	Format        string          `json:"format,omitempty"`
	FormattedBody string          `json:"formatted_body,omitempty"`
	Relates       *MessageRelates `json:"m.relates_to,omitempty"`
}

// MessageRelates indicates that a message relates so another one.
type MessageRelates struct {
	RepliesTo MessageReference `json:"m.in_reply_to"`
}

// MessageReference references another message.
type MessageReference struct {
	ID   string `json:"event_id"`
	Key  string `json:"key,omitempty"`
	Type string `json:"rel_type,omitempty"`
}

// SendMessage with the given text content to the room.
func (r *Room) SendMessage(ctx context.Context, text string) error {
	txID := strconv.FormatUint(r.server.TxID(), 10)
	path := "/_matrix/client/v3/rooms/" + r.id + "/send/m.room.message/" + txID
	msg := MessageContent{"m.text", text, "", "", nil}

	var response matrix.Response

	status, err := r.server.HTTP(ctx, http.MethodPut, path, msg, &response)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if err := response.AsError(status); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}

// SendReply sends a text reply to the given source event ID.
func (r *Room) SendReply(ctx context.Context, source, text string) error {
	txID := strconv.FormatUint(r.server.TxID(), 10)
	path := "/_matrix/client/v3/rooms/" + r.id + "/send/m.room.message/" + txID
	msg := MessageContent{"m.text", text, "", "", &MessageRelates{MessageReference{source, "", ""}}}

	var response matrix.Response

	status, err := r.server.HTTP(ctx, http.MethodPut, path, msg, &response)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if err := response.AsError(status); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}
