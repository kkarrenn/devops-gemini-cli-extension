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

package clouddeployclient

import (
	"context"
	"fmt"

	deploy "cloud.google.com/go/deploy/apiv1"
	deploypb "cloud.google.com/go/deploy/apiv1/deploypb"
	"google.golang.org/api/iterator"
)

// contextKey is a private type to use as a key for context values.
type contextKey string

const (
	cloudDeployClientContextKey contextKey = "cloudDeployClient"
)

// ClientFrom returns the CloudDeployClient stored in the context, if any.
func ClientFrom(ctx context.Context) (CloudDeployClient, bool) {
	client, ok := ctx.Value(cloudDeployClientContextKey).(CloudDeployClient)
	return client, ok
}

// ContextWithClient returns a new context with the provided CloudDeployClient.
func ContextWithClient(ctx context.Context, client CloudDeployClient) context.Context {
	return context.WithValue(ctx, cloudDeployClientContextKey, client)
}

// CloudDeployClient is an interface for interacting with the Cloud Deploy API.
type CloudDeployClient interface {
	ListDeliveryPipelines(ctx context.Context, projectID, location string) ([]*deploypb.DeliveryPipeline, error)
	ListTargets(ctx context.Context, projectID, location string) ([]*deploypb.Target, error)
	ListReleases(ctx context.Context, projectID, location, pipelineID string) ([]*deploypb.Release, error)
	ListRollouts(ctx context.Context, projectID, location, pipelineID, releaseID string) ([]*deploypb.Rollout, error)
	CreateRelease(ctx context.Context, projectID, location, pipelineID, releaseID string) (*deploy.CreateReleaseOperation, error)
}

// NewCloudDeployClient creates a new Cloud Deploy client.
func NewCloudDeployClient(ctx context.Context) (CloudDeployClient, error) {
	c, err := deploy.NewCloudDeployClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Deploy client: %v", err)
	}

	return &CloudDeployClientImpl{
		client: c,
	}, nil
}

// CloudDeployClientImpl is an implementation of the CloudDeployClient interface.
type CloudDeployClientImpl struct {
	client *deploy.CloudDeployClient
}

// ListDeliveryPipelines lists Delivery Pipelines in a given project and location
func (c *CloudDeployClientImpl) ListDeliveryPipelines(ctx context.Context, projectID, location string) ([]*deploypb.DeliveryPipeline, error) {
	req := &deploypb.ListDeliveryPipelinesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}
	it := c.client.ListDeliveryPipelines(ctx, req)
	var pipelines []*deploypb.DeliveryPipeline
	for {
		pipeline, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list delivery pipelines: %w", err)
		}
		pipelines = append(pipelines, pipeline)
	}
	return pipelines, nil
}

// ListTargets lists Targets in a given project and location
func (c *CloudDeployClientImpl) ListTargets(ctx context.Context, projectID, location string) ([]*deploypb.Target, error) {
	req := &deploypb.ListTargetsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}
	it := c.client.ListTargets(ctx, req)
	var targets []*deploypb.Target
	for {
		target, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list targets: %w", err)
		}
		targets = append(targets, target)
	}
	return targets, nil
}

// ListReleases lists Releases for a specific Delivery Pipeline
func (c *CloudDeployClientImpl) ListReleases(ctx context.Context, projectID, location, pipelineID string) ([]*deploypb.Release, error) {
	req := &deploypb.ListReleasesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/deliveryPipelines/%s", projectID, location, pipelineID),
	}
	it := c.client.ListReleases(ctx, req)
	var releases []*deploypb.Release
	for {
		release, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list releases: %w", err)
		}
		releases = append(releases, release)
	}
	return releases, nil
}

// ListRollouts lists Rollouts for a specific Release
func (c *CloudDeployClientImpl) ListRollouts(ctx context.Context, projectID, location, pipelineID, releaseID string) ([]*deploypb.Rollout, error) {
	req := &deploypb.ListRolloutsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/deliveryPipelines/%s/releases/%s", projectID, location, pipelineID, releaseID),
	}
	it := c.client.ListRollouts(ctx, req)
	var rollouts []*deploypb.Rollout
	for {
		rollout, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list rollouts: %w", err)
		}
		rollouts = append(rollouts, rollout)
	}
	return rollouts, nil
}

// CreateRelease creates a new Release to trigger a deployment
func (c *CloudDeployClientImpl) CreateRelease(ctx context.Context, projectID, location, pipelineID, releaseID string) (*deploy.CreateReleaseOperation, error) {
	req := &deploypb.CreateReleaseRequest{
		Parent:    fmt.Sprintf("projects/%s/locations/%s/deliveryPipelines/%s", projectID, location, pipelineID),
		ReleaseId: releaseID,
		Release:   &deploypb.Release{},
	}
	op, err := c.client.CreateRelease(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create release: %w", err)
	}
	return op, nil
}
