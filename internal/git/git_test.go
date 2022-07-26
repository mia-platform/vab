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

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
)

func TestUrlStringFromModule(t *testing.T) {
	module := v1alpha1.Module{Version: "1.2.3", Weight: 1, Disable: false}

	url := urlForModule(module)
	if url != defaultGitUrl {
		t.Fatalf("Unexpected url: %s", url)
	}

	module = v1alpha1.Module{}
	url = urlForModule(module)
	if url != defaultGitUrl {
		t.Fatalf("Unexpected url for empty module: %s", url)
	}
}

func TestUrlStringFromAddOn(t *testing.T) {
	addon := v1alpha1.AddOn{Version: "1.2.3", Disable: false}

	url := urlForAddon(addon)
	if url != defaultGitUrl {
		t.Fatalf("Unexpected url: %s", url)
	}

	addon = v1alpha1.AddOn{}
	url = urlForAddon(addon)
	if url != defaultGitUrl {
		t.Fatalf("Unexpected url for empty add-on: %s", url)
	}
}
