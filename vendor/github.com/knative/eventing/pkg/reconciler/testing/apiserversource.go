/*
Copyright 2019 The Knative Authors

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

package testing

import (
	"fmt"
	"time"

	"github.com/knative/eventing/pkg/apis/sources/v1alpha1"
	"github.com/knative/eventing/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ApiServerSourceOption enables further configuration of a ApiServer.
type ApiServerSourceOption func(*v1alpha1.ApiServerSource)

// NewApiServerSource creates a ApiServer with ApiServerOptions
func NewApiServerSource(name, namespace, uid string, o ...ApiServerSourceOption) *v1alpha1.ApiServerSource {
	c := &v1alpha1.ApiServerSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			UID:       types.UID(uid),
		},
	}
	for _, opt := range o {
		opt(c)
	}
	//c.SetDefaults(context.Background()) // TODO: We should add defaults and validation.
	return c
}

// WithInitApiServerSourceConditions initializes the ApiServerSource's conditions.
func WithInitApiServerSourceConditions(s *v1alpha1.ApiServerSource) {
	s.Status.InitializeConditions()
}

func WithApiServerSourceSinkNotFound(s *v1alpha1.ApiServerSource) {
	s.Status.MarkNoSink("NotFound", "")
}

func WithApiServerSourceSink(uri string) ApiServerSourceOption {
	return func(s *v1alpha1.ApiServerSource) {
		s.Status.MarkSink(uri)
	}
}

func WithApiServerSourceDeploymentUnavailable(s *v1alpha1.ApiServerSource) {
	// The Deployment uses GenerateName, so its name is empty.
	name := utils.GenerateFixedName(s, fmt.Sprintf("apiserversource-%s", s.Name))
	s.Status.PropagateDeploymentAvailability(NewDeployment(name, "any"))
}

func WithApiServerSourceDeployed(s *v1alpha1.ApiServerSource) {
	s.Status.PropagateDeploymentAvailability(NewDeployment("any", "any", WithDeploymentAvailable()))
}

func WithApiServerSourceEventTypes(s *v1alpha1.ApiServerSource) {
	s.Status.MarkEventTypes()
}

func WithApiServerSourceDeleted(c *v1alpha1.ApiServerSource) {
	t := metav1.NewTime(time.Unix(1e9, 0))
	c.ObjectMeta.SetDeletionTimestamp(&t)
}

func WithApiServerSourceSpec(spec v1alpha1.ApiServerSourceSpec) ApiServerSourceOption {
	return func(c *v1alpha1.ApiServerSource) {
		c.Spec = spec
	}
}
