// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package artifactregistry

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	artifactregistrypb "google.golang.org/genproto/googleapis/devtools/artifactregistry/v1"

	"devops-mcp-server/artifactregistryiface/mocks"
)

func TestCreateRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGRPClient := mocks.NewMockGRPClient(ctrl)
	client, err := NewClient(mockGRPClient)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	ctx := context.Background()
	projectID := "test-project"
	location := "us-central1"
	repositoryID := "test-repo"
	format := "DOCKER"

	req := &artifactregistrypb.CreateRepositoryRequest{
		Parent:       "projects/test-project/locations/us-central1",
		RepositoryId: "test-repo",
		Repository: &artifactregistrypb.Repository{
			Format: artifactregistrypb.Repository_DOCKER,
		},
	}

	repo := &artifactregistrypb.Repository{
		Name:   "projects/test-project/locations/us-central1/repositories/test-repo",
		Format: artifactregistrypb.Repository_DOCKER,
	}

	mockOperation := mocks.NewMockCreateRepositoryOperation(ctrl)

	mockGRPClient.EXPECT().CreateRepository(ctx, req).Return(mockOperation, nil)
	mockOperation.EXPECT().Wait(ctx).Return(repo, nil)

	_, err = client.CreateRepository(ctx, projectID, location, repositoryID, format)
	if err != nil {
		t.Fatalf("CreateRepository() failed: %v", err)
	}
}