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
	"eqrx.net/matrix"
	"eqrx.net/matrix/event"
)

// Response of a sync request.
type Response struct {
	matrix.Response
	AccountData            EventContainer `json:"account_data"`
	DeviceLists            DeviceLists    `json:"device_lists"`
	DeviceOneTimeKeysCount map[string]int `json:"device_one_time_keys_count"`
	NextBatch              string         `json:"next_batch"`
	Presence               EventContainer `json:"presence"`
	Rooms                  Rooms          `json:"rooms"`
	ToDevice               EventContainer `json:"to_device"`
}

// EventContainer represents map that just contains an event field.
//
// There is a lot of those in Response...
type EventContainer struct {
	Events []event.Opaque `json:"events"`
}

// DeviceLists contains information about device changes.
type DeviceLists struct {
	Changed []string `json:"changed"`
	Left    []string `json:"left"`
}

// Rooms contains information about rooms the client has interacted with, grouped by state.
type Rooms struct {
	Invited map[string]InvitedRoom `json:"invite"`
	Joined  map[string]JoinedRoom  `json:"join"`
	Knocked map[string]KnockedRoom `json:"knock"`
	Left    map[string]LeftRoom    `json:"leave"`
}

// InvitedRoom is a room the client was invited to.
type InvitedRoom struct {
	State EventContainer `json:"invite_state"`
}

// JoinedRoom is a room the client has joined.
type JoinedRoom struct {
	AccountData               EventContainer            `json:"account_data"`
	Ephemeral                 EventContainer            `json:"ephemeral"`
	State                     EventContainer            `json:"state"`
	Summary                   RoomSummary               `json:"summary"`
	Timeline                  Timeline                  `json:"timeline"`
	UnreadNotificationsCounts UnreadNotificationsCounts `json:"unread_notifications"`
}

// UnreadNotificationsCounts contains notication counts of a joined room.
type UnreadNotificationsCounts struct {
	Highlighted int `json:"highlight_count"`
	Total       int `json:"total"`
}

// RoomSummary for joined rooms.
type RoomSummary struct {
	Heros              []string `json:"m.heroes"`
	InvitedMemberCount int      `json:"m.invited_member_count"`
	JoinedMemeberCount int      `json:"m.joined_member_count"`
}

// Timeline represents events of a room.
type Timeline struct {
	Events        []event.Opaque `json:"events"`
	Limited       bool           `json:"limited"`
	PreviousBatch string         `json:"prev_batch"`
}

// KnockedRoom is a room the client has knocked on.
type KnockedRoom struct {
	KnockState EventContainer `json:"knock_state"`
}

// LeftRoom is a room the client has left.
type LeftRoom struct {
	AccountData EventContainer `json:"account_data"`
	State       EventContainer `json:"ephemeral"`
	Timeline    Timeline       `json:"timeline"`
}
