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

// Package matrix allows interfacing with a matrix server.
package matrix

import (
	"context"
	"net/http"
	"strings"
)

// Client to interface with a matrix server.
type Client struct {
	homeserver string
	token      string
	User       string
	Device     string
}

type whoamiResponse struct {
	Response
	User   string `json:"user_id"`
	Device string `json:"device_id"`
}

// New creates a new matrix client. It takes the homeserver url to contact
// and the client token to use as an argument.  It does a whoami request
// to get user ID and device ID of the token. Returns an error if that
// fails.
func New(ctx context.Context, homeserver, token string) (Client, error) {
	if homeserver == "" {
		panic("homeserver empty")
	}

	if token == "" {
		panic("token empty")
	}

	cli := Client{strings.TrimRight(homeserver, "/"), token, "", ""}

	var resp whoamiResponse

	err := HTTP(ctx, homeserver, token, http.MethodGet, "/_matrix/client/v3/account/whoami", nil, &resp)
	if err != nil {
		return cli, err
	}

	if err := resp.AsError(); err != nil {
		return cli, err
	}

	cli.User = resp.User
	cli.Device = resp.Device

	return cli, nil
}
