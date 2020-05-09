/*
Copyright 2016 The Kubernetes Authors.

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

package selinux

import (
	corev1 "k8s.io/api/core/v1"
	policy "k8s.io/api/policy/v1beta1"
	"testing"
)

func TestRunAsAnyOptions(t *testing.T) {
	_, err := NewRunAsAny(nil)
	if err != nil {
		t.Fatalf("unexpected error initializing NewRunAsAny %v", err)
	}
	_, err = NewRunAsAny(&policy.SELinuxStrategyOptions{})
	if err != nil {
		t.Errorf("unexpected error initializing NewRunAsAny %v", err)
	}
}

func TestRunAsAnyGenerate(t *testing.T) {
	s, err := NewRunAsAny(&policy.SELinuxStrategyOptions{})
	if err != nil {
		t.Fatalf("unexpected error initializing NewRunAsAny %v", err)
	}
	uid, err := s.Generate(nil, nil)
	if uid != nil {
		t.Errorf("expected nil uid but got %v", *uid)
	}
	if err != nil {
		t.Errorf("unexpected error generating uid %v", err)
	}
}

func TestRunAsAnyValidate(t *testing.T) {
	s, err := NewRunAsAny(&policy.SELinuxStrategyOptions{
		SELinuxOptions: &corev1.SELinuxOptions{
			Level: "foo",
		},
	},
	)
	if err != nil {
		t.Fatalf("unexpected error initializing NewRunAsAny %v", err)
	}
	errs := s.Validate(nil, nil, nil, nil)
	if len(errs) != 0 {
		t.Errorf("unexpected errors validating with ")
	}
	s, err = NewRunAsAny(&policy.SELinuxStrategyOptions{})
	if err != nil {
		t.Fatalf("unexpected error initializing NewRunAsAny %v", err)
	}
	errs = s.Validate(nil, nil, nil, nil)
	if len(errs) != 0 {
		t.Errorf("unexpected errors validating %v", errs)
	}
}
