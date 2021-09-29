//go:build e2e
// +build e2e

package test

import (
	"context"
	"fmt"
	"testing"

	tgithub "github.com/openshift-pipelines/pipelines-as-code/test/pkg/github"
	twait "github.com/openshift-pipelines/pipelines-as-code/test/pkg/wait"
	"github.com/tektoncd/pipeline/pkg/names"
	"gotest.tools/v3/assert"
)

func TestPullRequestPrivateRepository(t *testing.T) {
	for _, onWebhook := range []bool{false, true} {
		targetNS := names.SimpleNameGenerator.RestrictLengthWithRandomSuffix("pac-e2e-ns")
		ctx := context.Background()
		runcnx, opts, ghcnx, err := setup(ctx, onWebhook)
		assert.NilError(t, err)
		if onWebhook {
			runcnx.Clients.Log.Info("Testing with Direct Webhook integration")
		} else {
			runcnx.Clients.Log.Info("Testing with Github APPS integration")
		}
		entries, err := getEntries("testdata/pipelinerun_git_clone_private.yaml", targetNS, mainBranch, pullRequestEvent)
		assert.NilError(t, err)

		repoinfo, err := createRepoCRD(ctx, t, ghcnx, runcnx, opts, targetNS, pullRequestEvent, mainBranch, runcnx)
		assert.NilError(t, err)

		targetRefName := fmt.Sprintf("refs/heads/%s",
			names.SimpleNameGenerator.RestrictLengthWithRandomSuffix("pac-e2e-test"))

		title := "TestPullRequestPrivateRepository "
		if onWebhook {
			title += "OnWebhook"
		}
		title += "- " + targetRefName

		sha, err := tgithub.PushFilesToRef(ctx, ghcnx.Client, title, repoinfo.GetDefaultBranch(), targetRefName,
			opts.Owner, opts.Repo, entries)
		assert.NilError(t, err)
		runcnx.Clients.Log.Infof("Commit %s has been created and pushed to %s", sha, targetRefName)

		number, err := tgithub.PRCreate(ctx, runcnx, ghcnx, opts.Owner, opts.Repo, targetRefName, repoinfo.GetDefaultBranch(), title)
		assert.NilError(t, err)

		defer tearDown(ctx, t, runcnx, ghcnx, number, targetRefName, targetNS, opts)

		runcnx.Clients.Log.Infof("Waiting for Repository to be updated")
		waitOpts := twait.Opts{
			RepoName:        targetNS,
			Namespace:       targetNS,
			MinNumberStatus: 0,
			PollTimeout:     defaultTimeout,
			TargetSHA:       sha,
		}
		err = twait.UntilRepositoryUpdated(ctx, runcnx.Clients, waitOpts)
		assert.NilError(t, err)

		checkSuccess(ctx, t, runcnx, opts, pullRequestEvent, targetNS, sha, title)
	}
}

// Local Variables:
// compile-command: "go test -tags=e2e -v -info TestPullRequestPrivateRepository$ ."
// End:
