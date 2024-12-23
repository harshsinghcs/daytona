// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"github.com/daytonaio/daytona/pkg/common"
	"github.com/daytonaio/daytona/pkg/models"
)

type RunnerEventName string

const (
	RunnerEventLifecycleCreated        RunnerEventName = "runner_lifecycle_created"
	RunnerEventLifecycleCreationFailed RunnerEventName = "runner_lifecycle_creation_failed"
	RunnerEventLifecycleStarted        RunnerEventName = "runner_lifecycle_started"
	RunnerEventLifecycleStartFailed    RunnerEventName = "runner_lifecycle_start_failed"
	RunnerEventLifecycleStopped        RunnerEventName = "runner_lifecycle_stopped"
	RunnerEventLifecycleStopFailed     RunnerEventName = "runner_lifecycle_stop_failed"
	RunnerEventLifecycleDeleted        RunnerEventName = "runner_lifecycle_deleted"
	RunnerEventLifecycleDeletionFailed RunnerEventName = "runner_lifecycle_deletion_failed"
	RunnerEventProviderInstalled       RunnerEventName = "runner_provider_installed"
	RunnerEventProviderInstallFailed   RunnerEventName = "runner_provider_install_failed"
	RunnerEventProviderUninstalled     RunnerEventName = "runner_provider_uninstalled"
	RunnerEventProviderUninstallFailed RunnerEventName = "runer_provider_uninstall_failed"
	RunnerEventProviderUpdated         RunnerEventName = "runner_provider_updated"
	RunnerEventProviderUpdateFailed    RunnerEventName = "runner_provider_update_failed"
)

type runnerEvent struct {
	AbstractEvent
	runner *models.Runner
}

func NewRunnerEvent(name RunnerEventName, r *models.Runner, err error, extras map[string]interface{}) Event {
	return runnerEvent{
		runner: r,
		AbstractEvent: AbstractEvent{
			name:   string(name),
			extras: extras,
			err:    err,
		},
	}
}

func (e runnerEvent) Props() map[string]interface{} {
	props := e.AbstractEvent.Props()

	if e.runner != nil {
		props["runner_id"] = e.runner.Id
		props["is_local_runner"] = e.runner.Id == common.LOCAL_RUNNER_ID
	}

	return props
}
