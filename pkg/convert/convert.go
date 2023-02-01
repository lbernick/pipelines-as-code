package convert

import (
	"fmt"
	"strings"

	"github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/keys"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/v1alpha1"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/matcher"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

const gitRefPlaceholder = "gitRefPlaceholder"

func ConvertPipelineRunToWorkflow(pr *v1beta1.PipelineRun) (*v1alpha1.Workflow, error) {
	wf := v1alpha1.Workflow{Name: pr.Name}
	var events []string
	var targetBranches []string
	var taskRefs map[string]v1beta1.ResolverRef
	for k, v := range pr.Annotations {
		var err error
		switch k {
		case keys.OnEvent:
			if len(events) > 0 {
				return nil, fmt.Errorf("multiple events annotations")
			}
			events, err = matcher.GetAnnotationValues(v)
		case keys.OnTargetBranch:
			if len(targetBranches) > 0 {
				return nil, fmt.Errorf("multiple targetbranch annotations")
			}
			targetBranches, err = matcher.GetAnnotationValues(v)
		case keys.OnCelExpression:
			// TODO: Not supporting CEL expressions for now
			continue
		default:
			if !strings.HasPrefix(k, keys.Task) {
				// Not supporting any other annotation types
				continue
			}
			// You can have multiple TaskRefs by putting them into one "task" annotation
			// or by creating "task-1", "task-2", etc annotations.
			// HACK due to laziness: assume they are all in the same annotation and just overwrite taskRefs if it's already defined
			taskRefs, err = getTaskRefs(v)
		}
		if err != nil {
			return nil, err
		}
	}
	for _, event := range events {
		wf.Spec.Events = append(wf.Spec.Events, v1alpha1.Event{Type: event, Filters: []v1alpha1.Filter{{TargetBranches: targetBranches}}})
	}

	if pr.Spec.PipelineSpec != nil {
		ps, err := getReplacedPipelineSpec(*pr.Spec.PipelineSpec, taskRefs)
		if err != nil {
			return nil, err
		}
		pr.Spec.PipelineSpec = ps
	}

	return &wf, nil
}

func getTaskRefs(v string) (map[string]v1beta1.ResolverRef, error) {
	var tasks []string
	out := make(map[string]v1beta1.ResolverRef)
	if strings.HasPrefix(v, "[") {
		v = strings.TrimPrefix(v, "[")
		v = strings.TrimSuffix(v, "]")
		tasks = append(tasks, strings.Split(v, ", ")...)
	} else {
		tasks = append(tasks, v)
	}

	hasRemoteTask := false
	for _, task := range tasks {
		if strings.HasPrefix(task, "http://") || strings.HasPrefix(task, "https://") {
			// There's no way to know the name of the task being fetched,
			// so it's not clear which pipeline task should be replaced by the gitRef
			// unless there is only one remote pipeline task (or we assume that the file name is the task name).
			if hasRemoteTask {
				return nil, fmt.Errorf("cannot convert to workflow with multiple remote task refs")
			}
			hasRemoteTask = true

			useHttps := false
			if strings.HasPrefix(task, "http://") {
				v = strings.TrimPrefix(task, "http://")
			} else {
				useHttps = true
				v = strings.TrimPrefix(v, "https://")
			}
			components := strings.Split(v, "/")
			// HACK: Assuming a string of the form https://github.com/openshift-pipelines/pac/path/in/repo
			scm := components[0]
			owner := components[1]
			repo := components[2]
			pathInRepo := strings.Join(components[3:], "/")

			var url string
			if useHttps {
				url = "https://"
			} else {
				url = "http://"
			}
			url = strings.Join([]string{url, scm, owner, repo}, "/")

			out[gitRefPlaceholder] = v1beta1.ResolverRef{
				Resolver: "git",
				// Use anonymous cloning mode
				Params: []v1beta1.Param{{
					Name:  "url",
					Value: *v1beta1.NewArrayOrString(url),
				}, {
					Name:  "pathInRepo",
					Value: *v1beta1.NewArrayOrString(pathInRepo),
				}},
			}
		} else if strings.Contains(task, ".yaml") {
			// HACK: Assume the file name is the same as the task name
			components := strings.Split(task, "/")
			taskname := strings.TrimSuffix(components[len(components)-1], ".yaml")
			out[taskname] = v1beta1.ResolverRef{
				Resolver: "git",
				// Use anonymous cloning mode
				Params: []v1beta1.Param{{
					Name:  "url",
					Value: *v1beta1.NewArrayOrString("$(wf.params.repo_url)"),
				}, {
					Name:  "pathInRepo",
					Value: *v1beta1.NewArrayOrString(task),
				}},
			}
		} else {
			var taskname string
			version := "latest"
			if strings.Contains(task, ":") {
				components := strings.Split(task, ":")
				taskname = components[0]
				version = components[1]
			} else {
				taskname = task
			}
			out[taskname] = v1beta1.ResolverRef{
				Resolver: "hub",
				Params: []v1beta1.Param{{
					Name:  "name",
					Value: *v1beta1.NewArrayOrString(taskname),
				}, {
					Name:  "version",
					Value: *v1beta1.NewArrayOrString(version),
				}},
			}
		}
	}
	return out, nil
}

func getReplacedPipelineSpec(ps v1beta1.PipelineSpec, taskRefs map[string]v1beta1.ResolverRef) (*v1beta1.PipelineSpec, error) {
	newSpec := ps.DeepCopy()
	for i, task := range newSpec.Tasks {
		if task.TaskSpec != nil {
			continue
		}
		if task.TaskRef == nil {
			return nil, fmt.Errorf("must specify task spec or task ref for pipeline task %s", task.Name)
		}
		if task.TaskRef.Bundle != "" || task.TaskRef.Resolver != "" || task.TaskRef.Kind != "" {
			continue
		}
		resolverRef, ok := taskRefs[task.Name]
		if !ok {
			// Assume remote task
			resolverRef, ok = taskRefs[gitRefPlaceholder]
			if !ok {
				return nil, fmt.Errorf("error converting pipeline task %s", task.Name)
			}
		}
		newSpec.Tasks[i].Name = ""
		newSpec.Tasks[i].TaskRef = &v1beta1.TaskRef{ResolverRef: resolverRef}
	}
	return newSpec, nil
}
