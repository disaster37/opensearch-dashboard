// A generated module for OpensearchDashboard functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/opensearch-dashboard/internal/dagger"
	"fmt"

	"emperror.dev/errors"
	"github.com/disaster37/dagger-library-go/lib/helper"
)

const (
	OpensearchVersion string = "3.4.0"
	username          string = "admin"
	password          string = "vLPeJYa8.3RqtZCcAK6jNz"
	mockgenVersion           = "v0.3.0"
	gitUsername       string = "ci"
	gitEmail          string = "ci@localhost"
	defaultGitBranch  string = "3.x"
)

type OpensearchDashboard struct {
	// Src is a directory that contains the projects source code
	// +private
	Src *dagger.Directory

	// +private
	GolangModule *dagger.Golang
}

func New(
	ctx context.Context,
	// a path to a directory containing the source code
	// +required
	src *dagger.Directory,
) (*OpensearchDashboard, error) {

	return &OpensearchDashboard{
		Src:          src,
		GolangModule: dag.Golang(src.WithoutDirectory("ci")),
	}, nil
}

func (h *OpensearchDashboard) Ci(
	ctx context.Context,

	// Set tru if you are on CI
	// +default=false
	ci bool,

	// The codeCov token
	// +optional
	codeCoveToken *dagger.Secret,

	// The git branch where you should to push
	// You need to provide it when you are on PullRequest or on Tag
	// +optional
	gitBranch string,

	// Set true if current build is a tag
	// It will use the stable and alpha channel
	// alpha channel only instead
	// +optional
	isTag bool,

	// The git token
	// +optional
	gitToken *dagger.Secret,
) (dir *dagger.Directory, err error) {
	var stdout string

	h.GolangModule = dag.Golang(h.Src.WithoutDirectory("ci"), dagger.GolangOpts{Base: h.GolangModule.Container().WithExec([]string{"go", "mod", "tidy"})})

	// Build
	if _, err = h.Build(ctx).Sync(ctx); err != nil {
		return nil, errors.Wrap(err, "Error when build project")
	}

	// Lint code
	stdout, err = h.Lint(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Error when lint project: %s", stdout)
	}

	// Format code
	dir = h.Format(ctx)

	// Test code
	reportFile := h.Test(ctx)
	dir = dir.WithFile("coverage.out", reportFile)

	// Generate mock
	dir = dir.WithDirectory(".", h.GenerateMock(ctx))

	if ci {
		if codeCoveToken == nil {
			return nil, errors.New("You need to provide CodeCov token")
		}
		stdout, err = h.CodeCov(ctx, dir, codeCoveToken)
		if err != nil {
			return nil, errors.Wrapf(err, "Error when upload report on CodeCov: %s", stdout)
		}

		git := dag.GitModule(dir.WithDirectory("ci", h.Src.Directory("ci")), dagger.GitModuleOpts{Ci: "github"}).
			SetConfig(dagger.GitModuleSetConfigOpts{
				Username: gitUsername,
				Email:    gitEmail,
			})

		if isTag {
			gitBranch = defaultGitBranch
		}

		if _, err = git.CommitAndPush(
			ctx,
			gitToken,
			dagger.GitModuleCommitAndPushOpts{
				BranchName: gitBranch,
				GitRepoURL: "https://github.com/disaster37/opensearch.git",
				Message:    "Commit from CI",
			},
		); err != nil {
			return nil, errors.Wrap(err, "Error when commit and push files change")
		}
	}

	return dir, nil

}

// Lint permit to lint code
func (h *OpensearchDashboard) Lint(
	ctx context.Context,
) (string, error) {
	return h.GolangModule.Lint(ctx)
}

// Format permit to format the golang code
func (h *OpensearchDashboard) Format(
	ctx context.Context,
) *dagger.Directory {
	return h.GolangModule.Format()
}

// Test permit to run tests
func (h *OpensearchDashboard) Test(
	ctx context.Context,
) *dagger.File {
	// Run Opensearch
	opensearchService := dag.Container().
		From(fmt.Sprintf("opensearchproject/opensearch:%s", OpensearchVersion)).
		WithEnvVariable("cluster.name", "test").
		WithEnvVariable("node.name", "opensearch-node1").
		WithEnvVariable("bootstrap.memory_lock", "true").
		WithEnvVariable("discovery.type", "single-node").
		WithEnvVariable("network.publish_host", "127.0.0.1").
		WithEnvVariable("logger.org.opensearchsearch", "warn").
		WithEnvVariable("OPENSEARCH_JAVA_OPTS", "-Xms1g -Xmx1g").
		WithEnvVariable("plugins.security.nodes_dn_dynamic_config_enabled", "true").
		WithEnvVariable("plugins.security.unsupported.restapi.allow_securityconfig_modification", "true").
		WithEnvVariable("OPENSEARCH_INITIAL_ADMIN_PASSWORD", password).
		WithExposedPort(9200).
		AsService()

	dashboardService := dag.Container().
		From(fmt.Sprintf("opensearchproject/opensearch-dashboards:%s", OpensearchVersion)).
		WithEnvVariable("OPENSEARCH_HOSTS", "https://opensearch.svc:9200").
		WithEnvVariable("OPENSEARCH_USERNAME", username).
		WithEnvVariable("OPENSEARCH_PASSWORD", password).
		WithServiceBinding("opensearch.svc", opensearchService).
		WithExposedPort(5601).
		AsService()

	return h.GolangModule.Container().
		WithServiceBinding("opensearch.svc", opensearchService).
		WithServiceBinding("dashboard.svc", dashboardService).
		WithExec(helper.ForgeScript("DASHBOARD_USERNAME=%s DASHBOARD_PASSWORD=%s go test ./... -v -count 1 -parallel 1 -race -coverprofile=coverage.out -covermode=atomic -timeout 120m", username, password)).
		File("coverage.out")
}

// Build permit to build project
func (h *OpensearchDashboard) Build(
	ctx context.Context,
) *dagger.Directory {
	return h.GolangModule.Build()
}

func (h *OpensearchDashboard) CodeCov(
	ctx context.Context,

	// Optional directory
	// +optional
	src *dagger.Directory,

	// The Codecov token
	// +required
	token *dagger.Secret,
) (stdout string, err error) {

	if src == nil {
		src = h.Src
	}

	return dag.Codecov().Upload(
		ctx,
		src,
		token,
		dagger.CodecovUploadOpts{
			Files: []string{"coverage.out"},
		},
	)
}

// Test permit to run tests
func (h *OpensearchDashboard) GenerateMock(
	ctx context.Context,
) *dagger.Directory {
	return h.GolangModule.Container().WithExec(helper.ForgeScript(`
go install go.uber.org/mock/mockgen@%s
mockgen --build_flags=--mod=mod -destination=mocks/client.go -package=mocks github.com/disaster37/opensearch-dashboard/v3 Client
mockgen --build_flags=--mod=mod -destination=mocks/api.go -package=mocks github.com/disaster37/opensearch-dashboard/v3/api Api,SavedObjectApi,ShortenUrlApi,StatusApi
	`, mockgenVersion)).Directory(".")
}
