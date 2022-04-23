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
