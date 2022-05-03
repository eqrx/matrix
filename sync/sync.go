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

package sync

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"eqrx.net/matrix"
)

// Sync state with the given client. Since indicate where to start the response, filter which filter to use and
// timeoutMilliSeconds tells the server how long to block if the return limit set by the filter is not reached yet.
// The sync request will time out 10 seconds after that limit.
func Sync(ctx context.Context, cli matrix.Client, since, filter string, timeoutMilliSeconds int) (Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutMilliSeconds)*time.Millisecond+10*time.Second)
	defer cancel()

	path := "/_matrix/client/v3/sync?timeout=" + strconv.Itoa(timeoutMilliSeconds)

	if since != "" {
		path += "&since=" + since
	}

	if filter != "" {
		path += "&filter=" + filter
	}

	var response Response

	if err := cli.HTTP(ctx, http.MethodGet, path, nil, &response); err != nil {
		return response, fmt.Errorf("sync: %w", err)
	}

	if err := response.AsError(); err != nil {
		return response, fmt.Errorf("sync: %w", err)
	}

	return response, nil
}
