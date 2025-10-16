# Gemini CLI DevOps Extension

The DevOps extension provides Gemini powered AI assisted [CI/CD](https://en.wikipedia.org/wiki/CI/CD). It supports deployment to Cloud Run and Cloud Storage as well as creation of a robust CI/CD pipeline. Extension also sets up infrastructure for this pipeline. The generated pipelines follow best practices for CI/CD including testing and security best practices.


## 📋 Features

- **Simple code deployment**: Use `/deploy` command to deploy your code base to Google Cloud. The extension will leverage Gemini's advance capabilities to determine if your code requires Cloud Run (fro dynamic apps) or Cloud Storage (for static websites). Before deployment, it will scan your workspace for secrets, keys and password to avoid unintentional leaks.
- **AI-powered CI/CD pipeline**: Design and implement a robust and secure CI/CD pipeline in seconds. Collaborate with Gemini to design a pipeline that suites your needs. This will also setup required Google Cloud infrastructure. 
- **Interact with Google Cloud from Gemini CLI**: The DevOps extension offers tools to and commands to interact with Google Cloud's CI/CD services. You can run builds, check CVEs, SBOM, pull build logs into Gemini to help you investigate build failures.
- **Build complex release flows**: DevOps extension helps you build complex Cloud Deploy release pipelines in seconds based on simple questions.
- **DevOps MCP Server**: The extension brings power of MCP and integrates Gemini CLI with Google Cloud services for CI/CD (Cloud Build, Artifact Registry, Artifact Analysis, Cloud Deploy, Developer Connect).


## ⚙️ Installation

Install the DevOps extension by running the following command from your terminal *(requires Gemini CLI v0.8.0 or newer)*:

```bash
gemini extensions install https://github.com/gemini-cli-extensions/devops
```

## ✅ Prerequisites

- Authenticate with Google Cloud: `gcloud auth application-default login`
- Have the `gcloud` CLI installed and in your PATH.

## ☕ Usage
[TODO]


## Resources

- [Gemini CLI extensions](https://geminicli.com/extensions/about/): Documentation about using extensions in Gemini CLI
- [GitHub issues](https://github.com/gemini-cli-extensions/devops/issues): Report bugs or request features

