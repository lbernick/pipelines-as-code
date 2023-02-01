package v1alpha1

import (
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// At the moment, Workflow is just a configuration file in the repo.
// In the future, we will design a way for it to be applied to a cluster
// or stored in a separate repo, and turn it into a CRD.

type Workflow struct {
	metav1.TypeMeta
	Name string
	Spec WorkflowSpec
}

type WorkflowSpec struct {
	PipelineRun v1beta1.PipelineRun
	Params      []v1beta1.Param
	Events      []Event
}

type Event struct {
	Type    string
	Filters []Filter
}

type Filter struct {
	TargetBranches []string
	SourceBranches []string
	PathsChanged   []string
	CEL            []string
}
