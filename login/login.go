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

// Package login allows logging in with a matrix server.
package login

import (
	"context"
	"net/http"

	"eqrx.net/matrix"
)

type loginRequest struct {
	ID       id     `json:"identifier"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

type id struct {
	Type string `json:"type"`
	User string `json:"user"`
}

type loginResponse struct {
	matrix.Response
	Token  string `json:"access_token"`
	Device string `json:"device_id"`
}

// Login to the the matrix server behind the homeserver URL using the given http.Client, username and password.
// Returns the device ID and token.
func Login(ctx context.Context, cli http.Client, homeserver, user, password string) (string, string, error) {
	request := loginRequest{id{"m.id.user", user}, password, "m.login.password"}

	var response loginResponse

	path := "/_matrix/client/v3/login"
	if err := matrix.HTTP(ctx, cli, homeserver, "", http.MethodPost, path, request, &response); err != nil {
		return "", "", err
	}

	if err := response.AsError(); err != nil {
		return "", "", err
	}

	return response.Device, response.Token, nil
}
