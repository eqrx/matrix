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

package matrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTP performs a http exchange with the matrix server. It takes the http method, the path including the query
// and pointers to request and response payload as arguments. Request payload may be nil.
func (c Client) HTTP(ctx context.Context, method, path string, request, response interface{}) error {
	return HTTP(ctx, c.http, c.homeserver, c.token, method, path, request, response)
}

func httpRequest(
	ctx context.Context, cli http.Client, homeserver, token, method, path string, request interface{},
) (*http.Response, error) {
	requestBody := &bytes.Buffer{}

	if request != nil {
		if err := json.NewEncoder(requestBody).Encode(request); err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, homeserver+path, requestBody)
	if err != nil {
		panic(fmt.Sprintf("new http req: %v", err))
	}

	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}

	return cli.Do(httpReq)
}

func httpResponse(httpResp *http.Response, response interface{}) (err error) {
	defer func() {
		cErr := httpResp.Body.Close()
		if err == nil {
			err = cErr
		}
	}()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %v", httpResp.Status)
	}

	return json.NewDecoder(httpResp.Body).Decode(response)
}

// HTTP performs a http exchange via the matrix client. It takes the homeserver url, the client token, the http method,
// the path including the query and pointers to request and response payload as arguments. Request payload may be nil.
// Token may be "" to send an unauthenticated request.
func HTTP(
	ctx context.Context, cli http.Client, homeserver, token, method, path string, request, response interface{},
) error {
	httpResp, err := httpRequest(ctx, cli, homeserver, token, method, path, request)
	if err != nil {
		return err
	}

	if err := httpResponse(httpResp, response); err != nil {
		return err
	}

	return nil
}
