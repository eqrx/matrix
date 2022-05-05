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

const (
	// EventTypeMessage in a event type field indicates that the event is a room message.
	EventTypeMessage = "m.room.message"
	// MessageTypeText in a message content message type field indicates that the event is a text message.
	MessageTypeText = "m.text"
)

// IsMessageEvent returns true if the given event metadata indicates that the event is a room message.
func IsMessageEvent(evt event.Metadata) bool {
	return evt.Type == EventTypeMessage
}

// MessageEvent is a message sent to a room.
type MessageEvent struct {
	event.Metadata
	Content MessageContent `struct:"content"`
}

// MessageContent represents the content of room message.
type MessageContent struct {
	MessageType string          `json:"msgtype"`
	Body        string          `json:"body"`
	RelatesTo   *MessageRelates `json:"m.relates_to,omitempty"`
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

// NewTextMessage creates a new MessageEvent with the content of a text message with the given body
// and sets the room of the event.
func NewTextMessage(room, body string) MessageEvent {
	return MessageEvent{event.Metadata{Room: room}, MessageContent{"m.text", body, nil}}
}

// AsMessageEvent converts the given opaque event to a room message. Panics is metadata indicates that the event is not
// a room message and returns an error if unmarshalling of the content failed.
func AsMessageEvent(evt event.Opaque) (MessageEvent, error) {
	if evt.Type != EventTypeMessage {
		panic("not message")
	}

	mevt := MessageEvent{Metadata: evt.Metadata}

	if err := json.Unmarshal(evt.Content, &mevt.Content); err != nil {
		return mevt, nil
	}

	return mevt, nil
}

// AsReplyTo marks the message as a reply to the given event ID.
func (m MessageEvent) AsReplyTo(toID string) MessageEvent {
	m.Content.RelatesTo = &MessageRelates{MessageReference{ID: toID}}

	return m
}

// Send the event via the given matrix client with the given transaction ID.
// The room field of the even metadata must be set. Returns the event ID of the sent content.
func (m MessageEvent) Send(ctx context.Context, cli matrix.Client) (string, error) {
	return sendContent(ctx, cli, m.Room, EventTypeMessage, m.Content)
}
