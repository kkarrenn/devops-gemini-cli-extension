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

package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"

	chromem "github.com/philippgille/chromem-go"
	"golang.org/x/oauth2/google"
)

// Source represents a data source to be fetched.
type Source struct {
	Name           string   `json:"name"`
	Extract        string   `json:"extract"`
	Type           string   `json:"type"`
	ExcludePattern string   `json:"exclude_pattern,omitempty"`
	Dir            string   `json:"dir,omitempty"`
	URLs           []string `json:"urls"`
	URLPattern     string   `json:"url_pattern,omitempty"`
}

var KNOWLEDGE_RAG_SOURCES = []Source{
	{
		Name:           "GCP_DOCS",
		Extract:        "devsite-content",
		Type:           "webpage",
		ExcludePattern: ".*\\?hl=.+$",
		Dir:            "GCP_DOCS",
		URLs: []string{
			"https://cloud.google.com/developer-connect/docs/api/reference/rest",
			"https://cloud.google.com/developer-connect/docs/authentication",
			"https://cloud.google.com/build/docs/api/reference/rest",
			"https://cloud.google.com/deploy/docs/api/reference/rest",
			"https://cloud.google.com/artifact-analysis/docs/reference/rest",
			"https://cloud.google.com/artifact-registry/docs/reference/rest",
			"https://cloud.google.com/infrastructure-manager/docs/reference/rest",
			"https://cloud.google.com/docs/buildpacks/stacks",
			"https://cloud.google.com/docs/buildpacks/base-images",
			"https://cloud.google.com/docs/buildpacks/build-application",
			"https://cloud.google.com/docs/buildpacks/python",
			"https://cloud.google.com/docs/buildpacks/nodejs",
			"https://cloud.google.com/docs/buildpacks/java",
			"https://cloud.google.com/docs/buildpacks/go",
			"https://cloud.google.com/docs/buildpacks/ruby",
			"https://cloud.google.com/docs/buildpacks/php",
			"https://cloud.google.com/build/docs/build-config-file-schema",
			"https://cloud.google.com/build/docs/private-pools/use-in-private-network",
			"https://cloud.google.com/deploy/docs/config-files",
			"https://cloud.google.com/deploy/docs/deploy-app-gke",
			"https://cloud.google.com/deploy/docs/deploy-app-run",
			"https://cloud.google.com/deploy/docs/overview",
			"https://cloud.google.com/build/docs/build-push-docker-image",
			"https://cloud.google.com/build/docs/deploy-containerized-application-cloud-run",
			"https://cloud.google.com/build/docs/automate-builds",
			"https://cloud.google.com/build/docs/configuring-builds/create-basic-configuration",
			"https://cloud.google.com/build/docs/automating-builds/create-manage-triggers",
			"https://cloud.google.com/build/docs/building/build-containers",
			"https://cloud.google.com/build/docs/building/build-nodejs",
			"https://cloud.google.com/build/docs/building/build-java",
			"https://cloud.google.com/build/docs/deploying-builds/deploy-cloud-run",
			"https://cloud.google.com/build/docs/deploying-builds/deploy-gke",
		},
	},
	{
		Name:       "Python_Specific_Docs",
		Extract:    "article",
		Type:       "webpage",
		URLPattern: ".*(?<!\\?hl=..)$",
		Dir:        "Python_Specific_Docs",
		URLs: []string{
			"https://packaging.python.org/en/latest/guides/tool-recommendations/",
			"https://packaging.python.org/en/latest/guides/section-build-and-publish/",
			"https://packaging.python.org/en/latest/tutorials/managing-dependencies/",
			"https://packaging.python.org/en/latest/tutorials/installing-packages/",
			"https://packaging.python.org/en/latest/tutorials/packaging-projects/",
			"https://packaging.python.org/en/latest/overview/",
			"https://packaging.python.org/en/latest/guides/",
			"https://packaging.python.org/en/latest/guides/tool-recommendations",
			"https://packaging.python.org/en/latest/glossary/",
			"https://packaging.python.org/en/latest/key_projects/",
			"https://py-pkgs.org/08-ci-cd.html",
			"https://switowski.com/blog/ci-101/",
		},
	},
	{
		Name:           "cloud_builder_docs",
		Extract:        "section",
		Type:           "git_repo",
		URLPattern:     "\\.md$",
		ExcludePattern: ".*(vendor|third_party|.github).*$",
		URLs: []string{
			"https://github.com/GoogleCloudPlatform/cloud-builders/archive/refs/heads/master.zip",
			"https://github.com/GoogleCloudPlatform/cloud-builders-community/archive/refs/heads/master.zip",
		},
	},
	{
		Name:           "GCP_Terraform_Docs",
		Extract:        "section",
		Type:           "git_repo",
		URLPattern:     "website/docs/.*\\.markdown$",
		ExcludePattern: ".*(vendor|third_party|.github).*$",
		URLs: []string{
			"https://github.com/hashicorp/terraform-provider-google/archive/refs/heads/main.zip",
		},
	},
}

func uploadDirectoryToRag(collection *chromem.Collection, dir string) {
	var docs []chromem.Document
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file %s: %v", path, err)
				return nil
			}
			doc := chromem.Document{
				ID:       path,
				Content:  string(content),
				Metadata: map[string]string{"source": path},
			}
			docs = append(docs, doc)
		}
		return nil
	})

	if len(docs) > 0 {
		err := collection.AddDocuments(context.Background(), docs, runtime.NumCPU())
		if err != nil {
			log.Printf("Error adding documents from %s to collection: %v", dir, err)
		}
	}
}

func processSource(source Source, tmpDir string) {
	sourceType := source.Type

	if sourceType == "webpage" {
		err := downloadWebsites(&source, tmpDir)
		if err != nil {
			log.Printf("Error downloading websites from source %s: %v", source.Name, err)
		}
	} else if sourceType == "git_repo" {
		for _, url := range source.URLs {
			repoDir := filepath.Join(tmpDir, source.Dir)
			err := fetchRepository(url, repoDir)
			if err != nil {
				log.Printf("Error downloading git repo %s: %v", url, err)
			}
		}
	} else {
		log.Printf("RAG Source type [%s] is not supported", sourceType)
	}
}

func main() {
	// Initialize the chromem database
	ctx := context.Background()

	// Use Application Default Credentials to get a TokenSource
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Fatalf("Failed to find default credentials: %v", err)
	}

	// Get an access token from the TokenSource
	token, err := creds.TokenSource.Token()
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	vertexEmbeddingFunc := chromem.NewEmbeddingFuncVertex(
		token.AccessToken,
		"haroonc-exp",
		chromem.EmbeddingModelVertexEnglishV4)
	db := chromem.NewDB()
	if len(os.Getenv("RAG_DB_PATH")) > 0 {
		//check if file exists, only import if it does
		if _, err := os.Stat(os.Getenv("RAG_DB_PATH")); os.IsNotExist(err) {
			log.Printf("RAG_DB_PATH file does not exist, skipping import: %v", os.Getenv("RAG_DB_PATH"))
		}else{
			err := db.ImportFromFile(os.Getenv("RAG_DB_PATH"), "")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	collection, err := db.GetOrCreateCollection("devops-rag", nil, vertexEmbeddingFunc)
	if err != nil {
		log.Fatal(err)
	}

	// Upload local directories
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	patternsDir := filepath.Join(pwd, "patterns")
	uploadDirectoryToRag(collection, patternsDir)

	knowledgeDir := filepath.Join(pwd, "knowledge")
	uploadDirectoryToRag(collection, knowledgeDir)

	// Create a temporary directory for downloads
	tmpDir, err := os.MkdirTemp("", "rag-data-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Process each data source
	for _, source := range KNOWLEDGE_RAG_SOURCES {
		processSource(source, tmpDir)
	}

	// Upload all files in the temporary directory to RAG
	uploadDirectoryToRag(collection, tmpDir)

	// Export the database to a file
	if len(os.Getenv("RAG_DB_PATH")) > 0 {
		err = db.ExportToFile(os.Getenv("RAG_DB_PATH"), true, "")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Database exported to %s", os.Getenv("RAG_DB_PATH"))
	}
}