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

package clouddeploy

import (
	"context"
	"fmt"

	deployclient "devops-mcp-server/clouddeploy/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Handler holds the clients for the clouddeploy service.
type Handler struct {
	CdClient deployclient.CloudDeployClient
}

// Register registers the clouddeploy tools with the MCP server.
func (h *Handler) Register(server *mcp.Server) {
	addListDeliveryPipelinesTool(server, h.CdClient)
	addListTargetsTool(server, h.CdClient)
	addListReleasesTool(server, h.CdClient)
	addListRolloutsTool(server, h.CdClient)
	addCreateReleaseTool(server, h.CdClient)
}

// ListDeliveryPipelinesArgs arguments for listing pipelines
type ListDeliveryPipelinesArgs struct {
	ProjectID string `json:"project_id" jsonschema:"The Google Cloud project ID."`
	Location  string `json:"location" jsonschema:"The Google Cloud location."`
}

var listDeliveryPipelinesToolFunc func(ctx context.Context, req *mcp.CallToolRequest, args ListDeliveryPipelinesArgs) (*mcp.CallToolResult, any, error)

func addListDeliveryPipelinesTool(server *mcp.Server, cdClient deployclient.CloudDeployClient) {
	listDeliveryPipelinesToolFunc = func(ctx context.Context, req *mcp.CallToolRequest, args ListDeliveryPipelinesArgs) (*mcp.CallToolResult, any, error) {
		pipelines, err := cdClient.ListDeliveryPipelines(ctx, args.ProjectID, args.Location)
		if err != nil {
			return &mcp.CallToolResult{}, nil, fmt.Errorf("failed to list delivery pipelines: %w", err)
		}
		return &mcp.CallToolResult{}, map[string]any{"delivery_pipelines": pipelines}, nil
	}
	mcp.AddTool(server, &mcp.Tool{Name: "clouddeploy.list_delivery_pipelines", Description: "Lists the Cloud Deploy delivery pipelines in a specified GCP project and location."}, listDeliveryPipelinesToolFunc)
}

// ListTargetsArgs arguments for listing targets
type ListTargetsArgs struct {
	ProjectID string `json:"project_id" jsonschema:"The Google Cloud project ID."`
	Location  string `json:"location" jsonschema:"The Google Cloud location."`
}

var listTargetsToolFunc func(ctx context.Context, req *mcp.CallToolRequest, args ListTargetsArgs) (*mcp.CallToolResult, any, error)

func addListTargetsTool(server *mcp.Server, cdClient deployclient.CloudDeployClient) {
	listTargetsToolFunc = func(ctx context.Context, req *mcp.CallToolRequest, args ListTargetsArgs) (*mcp.CallToolResult, any, error) {
		targets, err := cdClient.ListTargets(ctx, args.ProjectID, args.Location)
		if err != nil {
			return &mcp.CallToolResult{}, nil, fmt.Errorf("failed to list targets: %w", err)
		}
		return &mcp.CallToolResult{}, map[string]any{"targets": targets}, nil
	}
	mcp.AddTool(server, &mcp.Tool{Name: "clouddeploy.list_targets", Description: "Lists the Cloud Deploy targets in a specified GCP project and location."}, listTargetsToolFunc)
}

// ListReleasesArgs arguments for listing releases
type ListReleasesArgs struct {
	ProjectID  string `json:"project_id" jsonschema:"The Google Cloud project ID."`
	Location   string `json:"location" jsonschema:"The Google Cloud location."`
	PipelineID string `json:"pipeline_id" jsonschema:"The Delivery Pipeline ID."`
}

var listReleasesToolFunc func(ctx context.Context, req *mcp.CallToolRequest, args ListReleasesArgs) (*mcp.CallToolResult, any, error)

func addListReleasesTool(server *mcp.Server, cdClient deployclient.CloudDeployClient) {
	listReleasesToolFunc = func(ctx context.Context, req *mcp.CallToolRequest, args ListReleasesArgs) (*mcp.CallToolResult, any, error) {
		releases, err := cdClient.ListReleases(ctx, args.ProjectID, args.Location, args.PipelineID)
		if err != nil {
			return &mcp.CallToolResult{}, nil, fmt.Errorf("failed to list releases: %w", err)
		}
		return &mcp.CallToolResult{}, map[string]any{"releases": releases}, nil
	}
	mcp.AddTool(server, &mcp.Tool{Name: "clouddeploy.list_releases", Description: "Lists the Cloud Deploy releases for a specified pipeline."}, listReleasesToolFunc)
}

// ListRolloutsArgs arguments for listing rollouts
type ListRolloutsArgs struct {
	ProjectID  string `json:"project_id" jsonschema:"The Google Cloud project ID."`
	Location   string `json:"location" jsonschema:"The Google Cloud location."`
	PipelineID string `json:"pipeline_id" jsonschema:"The Delivery Pipeline ID."`
	ReleaseID  string `json:"release_id" jsonschema:"The Release ID."`
}

var listRolloutsToolFunc func(ctx context.Context, req *mcp.CallToolRequest, args ListRolloutsArgs) (*mcp.CallToolResult, any, error)

func addListRolloutsTool(server *mcp.Server, cdClient deployclient.CloudDeployClient) {
	listRolloutsToolFunc = func(ctx context.Context, req *mcp.CallToolRequest, args ListRolloutsArgs) (*mcp.CallToolResult, any, error) {
		rollouts, err := cdClient.ListRollouts(ctx, args.ProjectID, args.Location, args.PipelineID, args.ReleaseID)
		if err != nil {
			return &mcp.CallToolResult{}, nil, fmt.Errorf("failed to list rollouts: %w", err)
		}
		return &mcp.CallToolResult{}, map[string]any{"rollouts": rollouts}, nil
	}
	mcp.AddTool(server, &mcp.Tool{Name: "clouddeploy.list_rollouts", Description: "Lists the Cloud Deploy rollouts for a specified release."}, listRolloutsToolFunc)
}

// CreateReleaseArgs arguments for creating a release
type CreateReleaseArgs struct {
	ProjectID  string `json:"project_id" jsonschema:"The Google Cloud project ID."`
	Location   string `json:"location" jsonschema:"The Google Cloud location."`
	PipelineID string `json:"pipeline_id" jsonschema:"The Delivery Pipeline ID."`
	ReleaseID  string `json:"release_id" jsonschema:"The ID of the release to create."`
}

var createReleaseToolFunc func(ctx context.Context, req *mcp.CallToolRequest, args CreateReleaseArgs) (*mcp.CallToolResult, any, error)

func addCreateReleaseTool(server *mcp.Server, cdClient deployclient.CloudDeployClient) {
	createReleaseToolFunc = func(ctx context.Context, req *mcp.CallToolRequest, args CreateReleaseArgs) (*mcp.CallToolResult, any, error) {
		op, err := cdClient.CreateRelease(ctx, args.ProjectID, args.Location, args.PipelineID, args.ReleaseID)
		if err != nil {
			return &mcp.CallToolResult{}, nil, fmt.Errorf("failed to create release: %w", err)
		}
		return &mcp.CallToolResult{}, map[string]any{"operation": op.Name()}, nil
	}
	mcp.AddTool(server, &mcp.Tool{Name: "clouddeploy.create_release", Description: "Creates a new Cloud Deploy release for a specified delivery pipeline."}, createReleaseToolFunc)
}
