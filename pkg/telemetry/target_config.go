// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package telemetry

import "github.com/daytonaio/daytona/pkg/models"

type TargetConfigEventName string

const (
	TargetConfigEventLifecycleCreated        TargetConfigEventName = "target_config_lifecycle_created"
	TargetConfigEventLifecycleCreationFailed TargetConfigEventName = "target_config_lifecycle_creation_failed"
	TargetConfigEventLifecycleDeleted        TargetConfigEventName = "target_config_lifecycle_deleted"
	TargetConfigEventLifecycleDeletionFailed TargetConfigEventName = "target_config_lifecycle_deletion_failed"
)

type targetConfigEvent struct {
	AbstractEvent
	targetConfig *models.TargetConfig
}

func NewTargetConfigEvent(name TargetConfigEventName, tc *models.TargetConfig, err error, extras map[string]interface{}) Event {
	return targetConfigEvent{
		targetConfig: tc,
		AbstractEvent: AbstractEvent{
			name:   string(name),
			extras: extras,
			err:    err,
		},
	}
}

func (e targetConfigEvent) Props() map[string]interface{} {
	props := e.AbstractEvent.Props()

	if e.targetConfig != nil {
		props["target_config_id"] = e.targetConfig.Id
		props["provider_name"] = e.targetConfig.ProviderInfo.Name
		props["provider_version"] = e.targetConfig.ProviderInfo.Version
		props["deleted"] = e.targetConfig.Deleted
		props["runner_id"] = e.targetConfig.ProviderInfo.RunnerId
	}

	return props
}
