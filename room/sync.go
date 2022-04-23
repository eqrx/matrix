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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"eqrx.net/matrix"
	"github.com/go-logr/logr"
)

type syncResponse struct {
	matrix.Response
	NextBatch string `json:"next_batch"`
	Rooms     struct {
		Join map[string]struct {
			Timeline struct {
				Events        []OpaqueEvent `json:"events"`
				PreviousBatch string        `json:"prev_batch"`
			} `json:"timeline"`
		} `json:"join"`
	} `json:"rooms"`
}

// SyncEvents of this room and let all messages be handled by handler.
func (r *Room) SyncEvents(ctx context.Context, log logr.Logger, handler func(MessageEvent)) {
	since := ""

	for {
		nextBatch, events, err := r.eventBatch(ctx, since, 60)
		since = nextBatch

		switch {
		case err == nil:
		case errors.Is(err, ctx.Err()):
			return
		default:
			log.Error(err, "sync events")
			select {
			case <-ctx.Done():
				return
			case <-time.NewTimer(10 * time.Second).C:
				continue
			}
		}

		for _, e := range events {
			if e.IsRoomMessageEvent() {
				me, err := e.AsRoomMessageEvent()
				if err != nil {
					log.Error(err, "parsing message content")
				}

				handler(me)
			}
		}
	}
}

// eventBatch gets the batch of events for the room since the given value. Matrix server blocks
// for the given timeout in seconds.
func (r *Room) eventBatch(ctx context.Context, since string, timeout int) (string, []OpaqueEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, (5+time.Duration(timeout))*time.Second)
	defer cancel()

	path := "/_matrix/client/v3/sync?timeout=" + strconv.Itoa(timeout)

	if since != "" {
		path += "&since=" + since
	}

	var response syncResponse

	httpStatus, err := r.server.HTTP(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return "", nil, fmt.Errorf("sync: %w", err)
	}

	if err := response.AsError(httpStatus); err != nil {
		return "", nil, fmt.Errorf("sync: %w", err)
	}

	roomData, ok := response.Rooms.Join[r.id]
	if !ok {
		return response.NextBatch, []OpaqueEvent{}, nil
	}

	return response.NextBatch, roomData.Timeline.Events, nil
}
