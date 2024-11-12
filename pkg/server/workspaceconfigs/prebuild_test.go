// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspaceconfigs_test

import (
	"time"

	"github.com/daytonaio/daytona/internal/util"
	"github.com/daytonaio/daytona/pkg/gitprovider"
	"github.com/daytonaio/daytona/pkg/models"
	"github.com/daytonaio/daytona/pkg/server/builds"
	build_dto "github.com/daytonaio/daytona/pkg/server/builds/dto"
	"github.com/daytonaio/daytona/pkg/server/workspaceconfigs"
	"github.com/daytonaio/daytona/pkg/server/workspaceconfigs/dto"
)

var prebuild1 = &models.PrebuildConfig{
	Id:             "1",
	Branch:         "feat",
	CommitInterval: util.Pointer(3),
	Retention:      3,
	TriggerFiles:   []string{"file1", "file2"},
}

var prebuild2 = &models.PrebuildConfig{
	Id:             "2",
	Branch:         "dev",
	CommitInterval: util.Pointer(1),
	Retention:      3,
	TriggerFiles:   []string{"file1", "file2"},
}

var prebuild3 = &models.PrebuildConfig{
	Id:             "3",
	Branch:         "new",
	CommitInterval: util.Pointer(1),
	Retention:      3,
	TriggerFiles:   []string{"file1", "file2"},
}

var prebuild1Dto = &dto.PrebuildDTO{
	WorkspaceConfigName: workspaceConfig1.Name,
	Id:                  prebuild1.Id,
	Branch:              prebuild1.Branch,
	CommitInterval:      prebuild1.CommitInterval,
	Retention:           prebuild1.Retention,
	TriggerFiles:        prebuild1.TriggerFiles,
}

var repository1 *gitprovider.GitRepository = &gitprovider.GitRepository{
	Url:    "https://github.com/daytonaio/daytona.git",
	Branch: "main",
	Sha:    "sha1",
}

var expectedPrebuilds []*models.PrebuildConfig
var expectedFilteredPrebuilds []*models.PrebuildConfig

var expectedPrebuildsMap map[string]*models.PrebuildConfig
var expectedFilteredPrebuildsMap map[string]*models.PrebuildConfig

func (s *WorkspaceConfigServiceTestSuite) TestSetPrebuild() {
	require := s.Require()

	s.gitProviderService.On("GetGitProviderForUrl", repository1.Url).Return(&s.gitProvider, "github", nil)
	s.gitProvider.On("GetRepositoryContext", gitprovider.GetRepositoryContext{
		Url: repository1.Url,
	}).Return(repository1, nil)
	s.gitProviderService.On("GetPrebuildWebhook", "github", repository1, "").Return(util.Pointer("webhook-id"), nil)

	newPrebuildDto, err := s.workspaceConfigService.SetPrebuild(workspaceConfig1.Name, dto.CreatePrebuildDTO{
		Id:             &prebuild3.Id,
		Branch:         prebuild3.Branch,
		CommitInterval: prebuild3.CommitInterval,
		Retention:      prebuild3.Retention,
		TriggerFiles:   prebuild3.TriggerFiles,
	})
	require.Nil(err)

	prebuildDtos, err := s.workspaceConfigService.ListPrebuilds(&workspaceconfigs.WorkspaceConfigFilter{
		Name: &workspaceConfig1.Name,
	}, nil)
	require.Nil(err)
	require.Contains(prebuildDtos, newPrebuildDto)
}

func (s *WorkspaceConfigServiceTestSuite) TestFindPrebuild() {
	require := s.Require()

	prebuild, err := s.workspaceConfigService.FindPrebuild(&workspaceconfigs.WorkspaceConfigFilter{
		Name: &workspaceConfig1.Name,
	}, &workspaceconfigs.PrebuildFilter{
		Id: &prebuild1.Id,
	})
	require.Nil(err)
	require.Equal(prebuild1Dto, prebuild)
}
func (s *WorkspaceConfigServiceTestSuite) TestListPrebuilds() {
	require := s.Require()

	prebuildDtos, err := s.workspaceConfigService.ListPrebuilds(&workspaceconfigs.WorkspaceConfigFilter{
		Name: &workspaceConfig1.Name,
	}, nil)
	require.Nil(err)

	require.Contains(prebuildDtos, prebuild1Dto)
}

func (s *WorkspaceConfigServiceTestSuite) TestDeletePrebuild() {
	expectedPrebuilds = expectedPrebuilds[:1]

	require := s.Require()

	s.buildService.On("MarkForDeletion", &builds.BuildFilter{
		PrebuildIds: &[]string{prebuild2.Id},
	}, false).Return([]error{})

	err := s.workspaceConfigService.DeletePrebuild(workspaceConfig1.Name, prebuild2.Id, false)
	require.Nil(err)

	prebuildDtos, errs := s.workspaceConfigService.ListPrebuilds(&workspaceconfigs.WorkspaceConfigFilter{
		Name: &workspaceConfig1.Name,
	}, nil)
	require.Nil(errs)
	require.ElementsMatch([]*dto.PrebuildDTO{
		prebuild1Dto,
	}, prebuildDtos)
}

func (s *WorkspaceConfigServiceTestSuite) TestProcessGitEventCommitInterval() {
	require := s.Require()

	s.gitProviderService.On("GetGitProviderForUrl", repository1.Url).Return(&s.gitProvider, "github", nil)
	s.gitProvider.On("GetRepositoryContext", gitprovider.GetRepositoryContext{
		Url: repository1.Url,
	}).Return(repository1, nil)

	s.buildService.On("Create", build_dto.BuildCreationData{
		PrebuildId: prebuild1.Id,
		Repository: repository1,
		User:       workspaceConfig1.User,
		Image:      workspaceConfig1.Image,
	}).Return("", nil)

	s.buildService.On("Find", &builds.BuildFilter{
		PrebuildIds: &[]string{prebuild1.Id},
		GetNewest:   util.Pointer(true),
	}).Return(&models.Build{
		Id:         "1",
		PrebuildId: prebuild1.Id,
		Repository: repository1,
	}, nil)

	data := gitprovider.GitEventData{
		Url:           repository1.Url,
		Branch:        "feat",
		Sha:           "sha4",
		Owner:         repository1.Owner,
		AffectedFiles: []string{},
	}

	s.gitProvider.On("GetCommitsRange", repository1, repository1.Sha, data.Sha).Return(3, nil)

	err := s.workspaceConfigService.ProcessGitEvent(data)
	require.Nil(err)
}

func (s *WorkspaceConfigServiceTestSuite) TestProcessGitEventTriggerFiles() {
	require := s.Require()

	s.gitProviderService.On("GetGitProviderForUrl", repository1.Url).Return(&s.gitProvider, "github", nil)
	s.gitProvider.On("GetRepositoryContext", gitprovider.GetRepositoryContext{
		Url: repository1.Url,
	}).Return(repository1, nil)

	s.buildService.On("Create", build_dto.BuildCreationData{
		PrebuildId: prebuild1.Id,
		Repository: repository1,
		User:       workspaceConfig1.User,
		Image:      workspaceConfig1.Image,
	}).Return("", nil)

	data := gitprovider.GitEventData{
		Url:    repository1.Url,
		Branch: "feat",
		Sha:    "sha4",
		Owner:  repository1.Owner,
		AffectedFiles: []string{
			"file1",
		},
	}

	err := s.workspaceConfigService.ProcessGitEvent(data)
	require.Nil(err)
}

func (s *WorkspaceConfigServiceTestSuite) TestEnforceRetentionPolicy() {
	require := s.Require()

	s.buildService.On("List", &builds.BuildFilter{
		States: &[]models.BuildState{models.BuildStatePublished},
	}).Return([]*models.Build{
		{
			Id:         "1",
			PrebuildId: "1",
			State:      models.BuildStatePublished,
			CreatedAt:  time.Now().Add(time.Hour * -4),
		},
		{
			Id:         "2",
			PrebuildId: "1",
			State:      models.BuildStatePublished,
			CreatedAt:  time.Now().Add(time.Hour * -3),
		},
		{
			Id:         "3",
			PrebuildId: "1",
			State:      models.BuildStatePublished,
			CreatedAt:  time.Now().Add(time.Hour * -2),
		},
		{
			Id:         "4",
			PrebuildId: "1",
			State:      models.BuildStatePublished,
			CreatedAt:  time.Now().Add(time.Hour * -1),
		},
	}, nil)

	s.buildService.On("MarkForDeletion", &builds.BuildFilter{
		Id: util.Pointer("1"),
	}, false).Return([]error{})

	err := s.workspaceConfigService.EnforceRetentionPolicy()
	require.Nil(err)
}