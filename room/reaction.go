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
	"encoding/json"

	"eqrx.net/matrix"
	"eqrx.net/matrix/event"
)

// EventTypeReaction in a event type field indicates that the event is a room reaction.
const EventTypeReaction = "m.reaction"

// IsReactionEvent returns true if the given event metadata indicates that the event is a room reaction.
func IsReactionEvent(evt event.Metadata) bool {
	return evt.Type == EventTypeReaction
}

// ReactionEvent is a reaction sent to a room.
type ReactionEvent struct {
	event.Metadata
	Content ReactionContent `struct:"content"`
}

// ReactionContent represents the content of room reaction.
type ReactionContent struct {
	RelatesTo MessageReference `json:"m.relates_to"`
}

// NewReaction creates a new ReactionEvent with the content of a reaction with the given key and reference to the given
// event ID and sets the room of the event.
func NewReaction(room, key, toID string) ReactionEvent {
	if room == "" {
		panic("room id empty")
	}

	return ReactionEvent{
		event.Metadata{Room: room},
		ReactionContent{MessageReference{toID, key, "m.annotation"}},
	}
}

// AsReactionEvent converts the given opaque event to a room reaction. Panics is metadata indicates that the event is
// not a room message and returns an error if unmarshalling of the content failed.
func AsReactionEvent(evt event.Opaque) (ReactionEvent, error) {
	if evt.Type != EventTypeReaction {
		panic("not reaction")
	}

	rEvt := ReactionEvent{Metadata: evt.Metadata}

	if err := json.Unmarshal(evt.Content, &rEvt.Content); err != nil {
		return rEvt, nil
	}

	return rEvt, nil
}

// Send the event via the given matrix client with the given transaction ID.
// The room field of the even metadata must be set. Returns the event ID of the sent content.
func (r ReactionEvent) Send(ctx context.Context, cli matrix.Client) (string, error) {
	return sendContent(ctx, cli, r.Room, EventTypeReaction, r.Content)
}
