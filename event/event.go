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

// Package event defines the basic event type for interacting with matrix servers.
package event

import "encoding/json"

// Opaque is an event with Metadata and OpaqueContent as content. To get a more concrete type out of this check the
// value of the type field in the metadata and unmarshal the event into concrete types.
type Opaque struct {
	Metadata
	Content OpaqueContent `json:"content"`
}

// Metadata of an event. Depending of the event type fields may be set or not.
//
// Reason for this is that there are multiple groups of event types in matrix.
// One could have created types for each of them but I do not see the benefit.
type Metadata struct {
	Type      string        `json:"type"`
	ID        string        `json:"event_id"`
	Sender    string        `json:"sender"`
	Room      string        `json:"room_id"`
	StateKey  string        `json:"state_key"`
	Timestamp int           `json:"origin_server_ts"`
	Unsigned  *UnsignedData `json:"unsigned"`
}

// UnsignedData is the portion of Metadata that is not set by sender but
// servers on the way and it thus unsigned.
type UnsignedData struct {
	Age             int           `json:"age"`
	TXID            string        `json:"transaction_id"`
	PreviousContent OpaqueContent `json:"prev_content"`
	RedactedBecause *Opaque       `json:"redacted_because"`
}

// OpaqueContent is the content field of an event that should not be interpreted (yet).
//
// Normally one would give it the type interface{} when unmarshalling via json but that
// would result in a nested map[string]interface{} structure with I find suboptimal.
// It implements json.Unmarshaler and simply causes the json unmarshaller to store
// the content field as bytes until it gets unmarshalled into a concrete content type.
type OpaqueContent []byte

// UnmarshalJSON tells the json unmarshaller to leave the content field as is.
func (o *OpaqueContent) UnmarshalJSON(b []byte) error {
	*o = b

	return nil
}

var _ json.Unmarshaler = &OpaqueContent{}
