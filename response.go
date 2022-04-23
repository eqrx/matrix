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
	"errors"
	"fmt"
	"net/http"
)

// Response is the base for all responses returned by the matrix server.
type Response struct {
	ErrCode string `json:"errcode"`
	ErrMsg  string `json:"error"`
}

var errServer = errors.New("matrix server response")

// AsError returns a descriptive error if the matrix server has returned an error for the corresponding
// request of nil.
func (r Response) AsError(httpStatus int) error {
	if r.ErrCode != "" || r.ErrMsg != "" || httpStatus != http.StatusOK {
		return fmt.Errorf("%w: %s: %s: %s", errServer, http.StatusText(httpStatus), r.ErrCode, r.ErrMsg)
	}

	return nil
}
