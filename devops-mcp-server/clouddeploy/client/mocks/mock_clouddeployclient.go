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

package mocks

import (
	"context"

	deploy "cloud.google.com/go/deploy/apiv1"
	deploypb "cloud.google.com/go/deploy/apiv1/deploypb"
)

// MockCloudDeployClient is a mock implementation of the CloudDeployClient interface.
type MockCloudDeployClient struct {
	ListDeliveryPipelinesFunc func(ctx context.Context, projectID, location string) ([]*deploypb.DeliveryPipeline, error)
	ListTargetsFunc           func(ctx context.Context, projectID, location string) ([]*deploypb.Target, error)
	ListReleasesFunc          func(ctx context.Context, projectID, location, pipelineID string) ([]*deploypb.Release, error)
	ListRolloutsFunc          func(ctx context.Context, projectID, location, pipelineID, releaseID string) ([]*deploypb.Rollout, error)
	CreateReleaseFunc         func(ctx context.Context, projectID, location, pipelineID, releaseID string) (*deploy.CreateReleaseOperation, error)
}

func (m *MockCloudDeployClient) ListDeliveryPipelines(ctx context.Context, projectID, location string) ([]*deploypb.DeliveryPipeline, error) {
	if m.ListDeliveryPipelinesFunc != nil {
		return m.ListDeliveryPipelinesFunc(ctx, projectID, location)
	}
	return nil, nil
}

func (m *MockCloudDeployClient) ListTargets(ctx context.Context, projectID, location string) ([]*deploypb.Target, error) {
	if m.ListTargetsFunc != nil {
		return m.ListTargetsFunc(ctx, projectID, location)
	}
	return nil, nil
}

func (m *MockCloudDeployClient) ListReleases(ctx context.Context, projectID, location, pipelineID string) ([]*deploypb.Release, error) {
	if m.ListReleasesFunc != nil {
		return m.ListReleasesFunc(ctx, projectID, location, pipelineID)
	}
	return nil, nil
}

func (m *MockCloudDeployClient) ListRollouts(ctx context.Context, projectID, location, pipelineID, releaseID string) ([]*deploypb.Rollout, error) {
	if m.ListRolloutsFunc != nil {
		return m.ListRolloutsFunc(ctx, projectID, location, pipelineID, releaseID)
	}
	return nil, nil
}

func (m *MockCloudDeployClient) CreateRelease(ctx context.Context, projectID, location, pipelineID, releaseID string) (*deploy.CreateReleaseOperation, error) {
	if m.CreateReleaseFunc != nil {
		return m.CreateReleaseFunc(ctx, projectID, location, pipelineID, releaseID)
	}
	return nil, nil
}
