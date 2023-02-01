package convert_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/v1alpha1"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/convert"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/pipeline/test/parse"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConvert(t *testing.T) {
	pr := parse.MustParseV1beta1PipelineRun(t, `
metadata:
  name: release-pipeline
  annotations:
    pipelinesascode.tekton.dev/on-event: "[push]"
    pipelinesascode.tekton.dev/on-target-branch: "[refs/tags/*]"
    pipelinesascode.tekton.dev/task: "[.tekton/tasks/goreleaser.yaml, git-clone:0.5, https://remote.url/owner/repo/my-task.yaml]"
spec:
  params:
  - name: repo_url
    value: "{{repo_url}}"
  - name: revision
    value: "{{revision}}"
  pipelineSpec:
    params:
    - name: repo_url
    - name: revision
    tasks:
    - name: fetch-repository
      taskRef:
        name: git-clone
    - name: gorelease
      taskRef:
        name: goreleaser
    - name: other-random-task
      taskRef:
        name: my-task
`)
	want := v1alpha1.Workflow{
		Name: "release-pipeline",
		Spec: v1alpha1.WorkflowSpec{
			Events: []v1alpha1.Event{{
				Type: "push",
				Filters: []v1alpha1.Filter{{
					TargetBranches: []string{"refs/tags/*"},
				}},
			}},
			Params: []v1beta1.Param{{
				Name:  "repo_url",
				Value: *v1beta1.NewArrayOrString("$(context.repo_url)"),
			}, {
				Name:  "revision",
				Value: *v1beta1.NewArrayOrString("$(context.revision)"),
			}},
			PipelineRun: v1beta1.PipelineRun{ObjectMeta: v1.ObjectMeta{Name: "release-pipeline"},
				Spec: v1beta1.PipelineRunSpec{
					Params: []v1beta1.Param{{
						Name:  "repo_url",
						Value: *v1beta1.NewArrayOrString("$(wf.params.repo_url)"),
					}, {
						Name:  "revision",
						Value: *v1beta1.NewArrayOrString("$(wf.params.revision)"),
					}},
					PipelineSpec: &v1beta1.PipelineSpec{
						Tasks: []v1beta1.PipelineTask{{
							Name: "fetch-repository",
							TaskRef: &v1beta1.TaskRef{
								ResolverRef: v1beta1.ResolverRef{
									Resolver: "hub",
									Params: []v1beta1.Param{{
										Name:  "name",
										Value: *v1beta1.NewArrayOrString("git-clone"),
									}, {
										Name:  "version",
										Value: *v1beta1.NewArrayOrString("0.5"),
									}},
								},
							},
						}, {
							Name: "gorelease",
							TaskRef: &v1beta1.TaskRef{
								ResolverRef: v1beta1.ResolverRef{
									Resolver: "git",
									Params: []v1beta1.Param{{
										Name:  "url",
										Value: *v1beta1.NewArrayOrString("$(wf.params.repo_url)"),
									}, {
										Name:  "pathInRepo",
										Value: *v1beta1.NewArrayOrString(".tekton/tasks/goreleaser.yaml"),
									}},
								},
							},
						}, {
							Name: "other-random-task",
							TaskRef: &v1beta1.TaskRef{
								ResolverRef: v1beta1.ResolverRef{
									Resolver: "git",
									Params: []v1beta1.Param{{
										Name:  "url",
										Value: *v1beta1.NewArrayOrString("https://remote.url/owner/repo"),
									}, {
										Name:  "pathInRepo",
										Value: *v1beta1.NewArrayOrString("my-task.yaml"),
									}},
								},
							},
						}},
					},
				}},
		},
	}
	got, err := convert.ConvertPipelineRunToWorkflow(pr)
	if err != nil {
		t.Error(err)
	}
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf(d)
	}
}
