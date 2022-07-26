// Copyright 2022 Mia-Platform

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package git

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
)

func TestUrlStringFromModule(t *testing.T) {
	expectedURL := "https://github.com/mia-platform/distribution"
	module := v1alpha1.Module{Version: "1.2.3", Weight: 1, Disable: false}

	url := urlForModule(module)
	if url != expectedURL {
		t.Fatalf("Unexpected url: %s", url)
	}

	module = v1alpha1.Module{}
	url = urlForModule(module)
	if url != expectedURL {
		t.Fatalf("Unexpected url for empty module: %s", url)
	}
}

func TestUrlStringFromAddOn(t *testing.T) {
	expectedURL := "https://github.com/mia-platform/distribution"
	addon := v1alpha1.AddOn{Version: "1.2.3", Disable: false}

	url := urlForAddon(addon)
	if url != expectedURL {
		t.Fatalf("Unexpected url: %s", url)
	}

	addon = v1alpha1.AddOn{}
	url = urlForAddon(addon)
	if url != expectedURL {
		t.Fatalf("Unexpected url for empty add-on: %s", url)
	}
}

func TestGetAuths(t *testing.T) {
	addonAuth := authForAddon(v1alpha1.AddOn{})
	if addonAuth != nil {
		t.Fatalf("Unexpected auth configuration %s", addonAuth)
	}

	moduleAuth := authForModule(v1alpha1.Module{})
	if moduleAuth != nil {
		t.Fatalf("Unexpected auth configuration %s", addonAuth)
	}
}

func TestTagReferences(t *testing.T) {
	addonName := "addon-name/with-slash"
	addonVersion := "1.0.0"
	expectedReference := "refs/tags/addon-" + addonName + "-" + addonVersion
	tag := tagReferenceForAddon(addonName, addonVersion)
	if tag != plumbing.ReferenceName(expectedReference) {
		t.Fatalf("Unexpected addon tag reference %s, expected %s", tag, expectedReference)
	}
	if !tag.IsTag() {
		t.Fatalf("The addon reference %s is not a tag reference", tag)
	}

	moduleName := "module-name/flavor"
	moduleVersion := "1.0.0"
	expectedReference = "refs/tags/module-module-name-" + addonVersion
	tag = tagReferenceForModule(moduleName, moduleVersion)
	if tag != plumbing.ReferenceName(expectedReference) {
		t.Fatalf("Unexpected module tag reference %s, expected %s", tag, expectedReference)
	}
	if !tag.IsTag() {
		t.Fatalf("The module reference %s is not a tag reference", tag)
	}
}

func TestCloneOptions(t *testing.T) {
	addon := v1alpha1.AddOn{Version: "1.0.0", Disable: false}
	addonName := "addon-name"
	options := cloneOptionsForAddon(addonName, addon)

	if options.URL != urlForAddon(addon) {
		t.Fatalf("Unexpected URL for addon %s: %s", addonName, options.URL)
	}
	if options.Auth != nil {
		t.Fatalf("Unexpected Auth for addon %s: %s", addonName, options.Auth)
	}
	if options.ReferenceName != tagReferenceForAddon(addonName, addon.Version) {
		t.Fatalf("Unexpected reference name for addon %s: %s", addonName, options.ReferenceName)
	}
	if !options.ReferenceName.IsTag() {
		t.Fatalf("Reference created for addon %s is not a branch: %s", addonName, options.ReferenceName)
	}

	module := v1alpha1.Module{Version: "1.0.0", Weight: 10, Disable: false}
	moduleName := "module-name/flavor-name"
	options = cloneOptionsForModule(moduleName, module)

	if options.URL != urlForModule(module) {
		t.Fatalf("Unexpected URL for module %s: %s", moduleName, options.URL)
	}
	if options.Auth != nil {
		t.Fatalf("Unexpected Auth for module %s: %s", moduleName, options.Auth)
	}
	if options.ReferenceName != tagReferenceForModule(moduleName, addon.Version) {
		t.Fatalf("Unexpected reference name for module %s: %s", moduleName, options.ReferenceName)
	}
	if !options.ReferenceName.IsTag() {
		t.Fatalf("Reference created for module %s is not a branch: %s", moduleName, options.ReferenceName)
	}
}
