// Copyright 2022-2025 The sacloud/iaas-api-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package naked

import "github.com/sacloud/iaas-api-go/types"

// Zone ゾーン情報
type Zone struct {
	ID           types.ID   `json:",omitempty" yaml:"id,omitempty" structs:",omitempty"`
	DisplayOrder int        `json:",omitempty" yaml:"display_order,omitempty" structs:",omitempty"`
	Name         string     `json:",omitempty" yaml:"name,omitempty" structs:",omitempty"`
	Description  string     `json:",omitempty" yaml:"description,omitempty" structs:",omitempty"`
	IsDummy      bool       `yaml:"is_dummy"`
	VNCProxy     *VNCProxy  `json:",omitempty" yaml:"vnc_proxy,omitempty" structs:",omitempty"`
	FTPServer    *FTPServer `json:",omitempty" yaml:"ftp_server,omitempty" structs:",omitempty"`
	Region       *Region    `json:",omitempty" yaml:"region,omitempty" structs:",omitempty"`
}
