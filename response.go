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

import "fmt"

// Response is the base for all HTTP responses returned by a matrix server.
type Response struct {
	ErrCode string `json:"errcode"`
	ErrMsg  string `json:"error"`
}

// AsError returns a descriptive error instance if the matrix server has included an error
// in the response. Returns nil otherwise.
func (r Response) AsError() error {
	if r.ErrCode != "" || r.ErrMsg != "" {
		return fmt.Errorf("matrix server response: %s: %s", r.ErrCode, r.ErrMsg)
	}

	return nil
}
