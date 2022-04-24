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
	"encoding/json"
	"fmt"
)

// EventWithoutContent represents a room event without the content field.
// This struct should be extended by more concrete types.
type EventWithoutContent struct {
	EventID               string `json:"event_id"`
	OriginSenderTimestamp int    `json:"origin_server_ts"`
	Sender                string `json:"sender"`
	StateKey              string `json:"state_key"`
	EventType             string `json:"type"`
}

// IsRoomMessageEvent returns true if the event represents a message sent
// to a matrix room.
func (e EventWithoutContent) IsRoomMessageEvent() bool {
	return e.EventType == MessageEventType
}

// OpaqueContent is the marshalled form of a room event content.
type OpaqueContent []byte

// UnmarshalJSON tells the json unmarshaller to leave the content field as is.
func (o *OpaqueContent) UnmarshalJSON(b []byte) error {
	*o = b

	return nil
}

// ensure OpaqueContent implements json.Unmarshaller.
var _ json.Unmarshaler = &OpaqueContent{}

// OpaqueEvent is a room event with opaque content.
type OpaqueEvent struct {
	EventWithoutContent
	Content OpaqueContent `json:"content,string"`
}

// AsRoomMessageEvent converts this event to a MessageEvent. Panics if wrong type.
func (r OpaqueEvent) AsRoomMessageEvent() (MessageEvent, error) {
	if !r.IsRoomMessageEvent() {
		panic("not a message event")
	}

	var content MessageContent

	if err := json.Unmarshal(r.Content, &content); err != nil {
		return MessageEvent{}, fmt.Errorf("unmarshal room message event content: %w", err)
	}

	return MessageEvent{r.EventWithoutContent, content}, nil
}
