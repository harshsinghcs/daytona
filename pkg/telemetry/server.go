// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package telemetry

type ServerEventName string

const (
	ServerEventPurgeStarted   ServerEventName = "server_purge_started"
	ServerEventPurgeCompleted ServerEventName = "server_purge_completed"
	ServerEventPurgeError     ServerEventName = "server_purge_error"
)

type serverEvent struct {
	AbstractEvent
	serverId string
}

func NewServerEvent(name ServerEventName, serverId string, err error, extras map[string]interface{}) Event {
	return serverEvent{
		AbstractEvent: AbstractEvent{
			name:   string(name),
			extras: extras,
			err:    err,
		},
		serverId: serverId,
	}
}

func (e serverEvent) Props() map[string]interface{} {
	props := e.AbstractEvent.Props()

	props["server_id"] = e.serverId

	return props
}
