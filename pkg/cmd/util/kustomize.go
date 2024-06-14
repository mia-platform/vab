// Copyright Mia srl
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func KustomizeData(path string) ([]byte, error) {
	kOpts := krusty.MakeDefaultOptions()
	kOpts.Reorder = krusty.ReorderOptionLegacy
	k := krusty.MakeKustomizer(kOpts)
	m, err := k.Run(filesys.MakeFsOnDisk(), path)
	if err != nil {
		return nil, err
	}

	data, err := m.AsYaml()
	if err != nil {
		return nil, err
	}

	return data, err
}
