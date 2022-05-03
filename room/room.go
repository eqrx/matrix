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

// Package room handles matrix rooms.
package room

import (
	"context"
	"fmt"
	"net/http"

	"eqrx.net/matrix"
)

// Join the given room ID with the client.
func Join(ctx context.Context, cli matrix.Client, id string) error {
	path := "/_matrix/client/v3/join/" + id

	var joinRoomResponse matrix.Response

	if err := cli.HTTP(ctx, http.MethodPost, path, nil, &joinRoomResponse); err != nil {
		return fmt.Errorf("join rooms: %w", err)
	}

	if err := joinRoomResponse.AsError(); err != nil {
		return err
	}

	return nil
}

// Joined returns all rooms this client is part of.
func Joined(ctx context.Context, cli matrix.Client) ([]string, error) {
	path := "/_matrix/client/v3/joined_rooms"

	var listRoomsResponse struct {
		matrix.Response
		Rooms []string `json:"joined_rooms"`
	}

	if err := cli.HTTP(ctx, http.MethodGet, path, nil, &listRoomsResponse); err != nil {
		return nil, fmt.Errorf("list joined rooms: %w", err)
	}

	if err := listRoomsResponse.AsError(); err != nil {
		return nil, err
	}

	return listRoomsResponse.Rooms, nil
}
