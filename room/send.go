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

	"eqrx.net/matrix"
)

type sendResponse struct {
	matrix.Response
	ID string `json:"event_id"`
}

func sendContent(
	ctx context.Context, cli matrix.Client, roomID, eventType, txID string, content interface{},
) (string, error) {
	if roomID == "" || eventType == "" || txID == "" || content == nil {
		panic("parameter empty")
	}

	path := "/_matrix/client/v3/rooms/" + roomID + "/send/" + eventType + "/" + txID

	var response sendResponse

	if err := cli.HTTP(ctx, http.MethodPut, path, content, &response); err != nil {
		return "", err
	}

	if err := response.AsError(); err != nil {
		return "", err
	}

	return response.ID, nil
}
