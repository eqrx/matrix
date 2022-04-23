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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
)

// Matrix represents a matrix server.
type Matrix struct {
	homeserver string
	token      string
	lastTxID   uint64
}

// New creates a new matrix server handle. It takes the homeserver url and the client token
// as an argument.
func New(homeserver, token string) *Matrix {
	if homeserver == "" {
		panic("homeserver empty")
	}

	if token == "" {
		panic("token empty")
	}

	return &Matrix{strings.TrimRight(homeserver, "/"), token, 0}
}

// TxID generates a unique transaction ID for requests.
func (m *Matrix) TxID() uint64 { return atomic.AddUint64(&m.lastTxID, 1) }

// HTTP performs a http exchange with the matrix server. It takes the http method, the path including queries
// and pointers to request and response payload as arguments. It returns the http status and io errors.
func (m *Matrix) HTTP(ctx context.Context, method, path string, request, response interface{}) (int, error) {
	requestBody := &bytes.Buffer{}

	if request != nil {
		if err := json.NewEncoder(requestBody).Encode(request); err != nil {
			return -1, fmt.Errorf("marshal request: %w", err)
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, m.homeserver+path, requestBody)
	if err != nil {
		panic(fmt.Sprintf("new http req: %v", err))
	}

	httpReq.Header.Set("Authorization", "Bearer "+m.token)

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return -1, fmt.Errorf("http: %w", err)
	}

	err = json.NewDecoder(httpResp.Body).Decode(response)
	cErr := httpResp.Body.Close()

	switch {
	case err != nil && cErr != nil:
		return httpResp.StatusCode, fmt.Errorf("resp unmarshal: %w, resp close: %v", err, cErr)
	case err != nil:
		return httpResp.StatusCode, fmt.Errorf("resp unmarshal: %w", err)
	case cErr != nil:
		return httpResp.StatusCode, fmt.Errorf("resp close: %w", cErr)
	default:
		return httpResp.StatusCode, nil
	}
}
