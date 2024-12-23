// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"encoding/json"

	"github.com/daytonaio/daytona/pkg/models"
)

type WorkspaceTemplateEventName string

var (
	WorkspaceTemplateEventLifecycleCreated        WorkspaceTemplateEventName = "workspace_template_lifecycle_created"
	WorkspaceTemplateEventLifecycleCreationFailed WorkspaceTemplateEventName = "workspace_template_lifecycle_creation_failed"
	WorkspaceTemplateEventLifecycleDeleted        WorkspaceTemplateEventName = "workspace_template_lifecycle_deleted"
	WorkspaceTemplateEventLifecycleDeletionFailed WorkspaceTemplateEventName = "workspace_template_lifecycle_deletion_failed"
	WorkspaceTemplateEventPrebuildCreated         WorkspaceTemplateEventName = "workspace_template_prebuild_created"
	WorkspaceTemplateEventPrebuildCreationFailed  WorkspaceTemplateEventName = "workspace_template_prebuild_creation_failed"
	WorkspaceTemplateEventPrebuildDeleted         WorkspaceTemplateEventName = "workspace_template_prebuild_deleted"
	WorkspaceTemplateEventPrebuildDeletionFailed  WorkspaceTemplateEventName = "workspace_template_prebuild_deletion_failed"
)

type workspaceTemplateEvent struct {
	AbstractEvent
	workspaceTemplate *models.WorkspaceTemplate
}

func NewWorkspaceTemplateEvent(name WorkspaceTemplateEventName, wt *models.WorkspaceTemplate, err error, extras map[string]interface{}) Event {
	return workspaceTemplateEvent{
		workspaceTemplate: wt,
		AbstractEvent: AbstractEvent{
			name:   string(name),
			extras: extras,
			err:    err,
		},
	}
}

func (e workspaceTemplateEvent) Props() map[string]interface{} {
	props := e.AbstractEvent.Props()

	if e.workspaceTemplate != nil {
		props["workspace_template_name"] = e.workspaceTemplate.Name
		prebuilds, err := json.Marshal(e.workspaceTemplate.Prebuilds)
		if err == nil {
			props["prebuilds"] = string(prebuilds)
		}
		if isImagePublic(e.workspaceTemplate.Image) {
			props["image"] = e.workspaceTemplate.Image
		}

		if isPublic(e.workspaceTemplate.RepositoryUrl) {
			props["repository_url"] = e.workspaceTemplate.RepositoryUrl
		}

		props["builder"] = getBuilder(e.workspaceTemplate.BuildConfig)
	}

	return props
}
