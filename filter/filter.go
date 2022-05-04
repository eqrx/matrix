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

// Package filter allows creating sync filters.
package filter

import (
	"context"
	"net/http"

	"eqrx.net/matrix"
)

// Filter defines a filter used for syncing content.
type Filter struct {
	AccountData Event    `json:"account_data"`
	Fields      []string `json:"event_fields,omitempty"`
	Format      string   `json:"event_format,omitempty"`
	Presence    Event    `json:"presence"`
	Room        Room     `json:"room"`
}

// Event allows filering general events.
type Event struct {
	Limit      int      `json:"limit,omitempty"`
	NotSenders []string `json:"not_senders,omitempty"`
	Senders    []string `json:"senders,omitempty"`
	NotTypes   []string `json:"not_types,omitempty"`
	Types      []string `json:"types,omitempty"`
}

// Room allows to filter information specific to rooms.
type Room struct {
	AccountData  RoomEvent `json:"account_data"`
	Ephemeral    RoomEvent `json:"ephemeral"`
	IncludeLeave bool      `json:"include_leave,omitempty"`
	NotRooms     []string  `json:"not_rooms,omitempty"`
	Rooms        []string  `json:"rooms,omitempty"`
	State        RoomEvent `json:"state,omitempty"`
	Timeline     RoomEvent `json:"timeline,omitempty"`
}

// RoomEvent allows filering room events.
type RoomEvent struct {
	ContainsURL             *bool    `json:"contains_url,omitempty"`
	IncludeRedundantMembers bool     `json:"include_redundant_members,omitempty"`
	LazyLoadMembers         bool     `json:"lazy_load_members,omitempty"`
	Limit                   int      `json:"limit,omitempty"`
	NotSenders              []string `json:"not_senders,omitempty"`
	Senders                 []string `json:"senders,omitempty"`
	NotTypes                []string `json:"not_types,omitempty"`
	Types                   []string `json:"types,omitempty"`
	NotRooms                []string `json:"not_rooms,omitempty"`
	Rooms                   []string `json:"rooms,omitempty"`
}

type response struct {
	matrix.Response
	Filter string `json:"filter_id"`
}

// Register the filter with the given matrix client and return its ID.
func (f Filter) Register(ctx context.Context, cli matrix.Client) (string, error) {
	var response response

	path := "/_matrix/client/v3/user/" + cli.User() + "/filter"
	if err := cli.HTTP(ctx, http.MethodPost, path, f, &response); err != nil {
		return "", err
	}

	if err := response.AsError(); err != nil {
		return "", err
	}

	return response.Filter, nil
}
