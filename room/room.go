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

// Package room interfaces with matrix rooms.
package room

import (
	"context"
	"fmt"
	"net/http"

	"eqrx.net/matrix"
)

// Room represents a matrix room.
type Room struct {
	server *matrix.Matrix
	id     string
}

// New creates a new matrix room handle on the given server for the given room ID.
func New(server *matrix.Matrix, id string) *Room {
	if id == "" {
		panic("id empty")
	}

	return &Room{server, id}
}

// Join the room.
func (r *Room) Join(ctx context.Context) error {
	path := "/_matrix/client/v3/joined_rooms"

	var listRoomsResponse struct {
		matrix.Response
		Rooms []string `json:"joined_rooms"`
	}

	status, err := r.server.HTTP(ctx, http.MethodGet, path, nil, &listRoomsResponse)
	if err != nil {
		return fmt.Errorf("list joined rooms: %w", err)
	}

	if err := listRoomsResponse.AsError(status); err != nil {
		return err //nolint:wrapcheck
	}

	for _, room := range listRoomsResponse.Rooms {
		if room == r.id {
			return nil
		}
	}

	path = "/_matrix/client/v3/join/" + r.id

	var joinRoomResponse matrix.Response

	status, err = r.server.HTTP(ctx, http.MethodPost, path, nil, &joinRoomResponse)
	if err != nil {
		return fmt.Errorf("join rooms: %w", err)
	}

	if err := joinRoomResponse.AsError(status); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}
