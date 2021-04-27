/*
Copyright 2019 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeEventListeners implements EventListenerInterface
type FakeEventListeners struct {
	Fake *FakeTriggersV1alpha1
	ns   string
}

var eventlistenersResource = schema.GroupVersionResource{Group: "triggers.tekton.dev", Version: "v1alpha1", Resource: "eventlisteners"}

var eventlistenersKind = schema.GroupVersionKind{Group: "triggers.tekton.dev", Version: "v1alpha1", Kind: "EventListener"}

// Get takes name of the eventListener, and returns the corresponding eventListener object, and an error if there is any.
func (c *FakeEventListeners) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.EventListener, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(eventlistenersResource, c.ns, name), &v1alpha1.EventListener{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventListener), err
}

// List takes label and field selectors, and returns the list of EventListeners that match those selectors.
func (c *FakeEventListeners) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.EventListenerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(eventlistenersResource, eventlistenersKind, c.ns, opts), &v1alpha1.EventListenerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.EventListenerList{ListMeta: obj.(*v1alpha1.EventListenerList).ListMeta}
	for _, item := range obj.(*v1alpha1.EventListenerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested eventListeners.
func (c *FakeEventListeners) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(eventlistenersResource, c.ns, opts))

}

// Create takes the representation of a eventListener and creates it.  Returns the server's representation of the eventListener, and an error, if there is any.
func (c *FakeEventListeners) Create(ctx context.Context, eventListener *v1alpha1.EventListener, opts v1.CreateOptions) (result *v1alpha1.EventListener, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(eventlistenersResource, c.ns, eventListener), &v1alpha1.EventListener{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventListener), err
}

// Update takes the representation of a eventListener and updates it. Returns the server's representation of the eventListener, and an error, if there is any.
func (c *FakeEventListeners) Update(ctx context.Context, eventListener *v1alpha1.EventListener, opts v1.UpdateOptions) (result *v1alpha1.EventListener, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(eventlistenersResource, c.ns, eventListener), &v1alpha1.EventListener{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventListener), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeEventListeners) UpdateStatus(ctx context.Context, eventListener *v1alpha1.EventListener, opts v1.UpdateOptions) (*v1alpha1.EventListener, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(eventlistenersResource, "status", c.ns, eventListener), &v1alpha1.EventListener{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventListener), err
}

// Delete takes name of the eventListener and deletes it. Returns an error if one occurs.
func (c *FakeEventListeners) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(eventlistenersResource, c.ns, name), &v1alpha1.EventListener{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeEventListeners) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(eventlistenersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.EventListenerList{})
	return err
}

// Patch applies the patch and returns the patched eventListener.
func (c *FakeEventListeners) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.EventListener, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(eventlistenersResource, c.ns, name, pt, data, subresources...), &v1alpha1.EventListener{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventListener), err
}
