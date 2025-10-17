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

package artifactregistryiface_test

import (
	"context"
	"testing"

	"devops-mcp-server/artifactregistryiface/mocks"
	"github.com/golang/mock/gomock"
	artifactregistrypb "google.golang.org/genproto/googleapis/devtools/artifactregistry/v1"
)

func TestMockGRPClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockGRPClient(ctrl)
	mockOperation := mocks.NewMockCreateRepositoryOperation(ctrl)

	req := &artifactregistrypb.CreateRepositoryRequest{
		Parent: "projects/my-project/locations/us-central1",
		Repository: &artifactregistrypb.Repository{
			Format: artifactregistrypb.Repository_DOCKER,
		},
		RepositoryId: "my-repo",
	}

	repo := &artifactregistrypb.Repository{
		Name: "projects/my-project/locations/us-central1/repositories/my-repo",
	}

	mockClient.EXPECT().CreateRepository(gomock.Any(), req).Return(mockOperation, nil)
	mockOperation.EXPECT().Wait(gomock.Any()).Return(repo, nil)

	op, err := mockClient.CreateRepository(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateRepository() err = %v, want nil", err)
	}

	gotRepo, err := op.Wait(context.Background())
	if err != nil {
		t.Fatalf("Wait() err = %v, want nil", err)
	}

	if gotRepo.Name != repo.Name {
		t.Errorf("Wait() = %v, want %v", gotRepo, repo)
	}
}
