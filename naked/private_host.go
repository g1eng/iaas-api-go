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

import (
	"time"

	"github.com/sacloud/iaas-api-go/types"
)

// PrivateHost 専有ホスト
type PrivateHost struct {
	ID               types.ID         `yaml:"id"`
	Name             string           `json:",omitempty" yaml:"name,omitempty" structs:",omitempty"`
	Description      string           `yaml:"description"`
	Tags             types.Tags       `yaml:"tags"`
	Icon             *Icon            `json:",omitempty" yaml:"icon,omitempty" structs:",omitempty"`
	CreatedAt        *time.Time       `json:",omitempty" yaml:"created_at,omitempty" structs:",omitempty"`
	Plan             *PrivateHostPlan `json:",omitempty" yaml:"plan,omitempty" structs:",omitempty"`
	Host             *Host            `json:",omitempty" yaml:"host,omitempty" structs:",omitempty"`
	AssignedCPU      int              `json:",omitempty" yaml:"assigned_cpu,omitempty" structs:",omitempty"`
	AssignedMemoryMB int              `json:",omitempty" yaml:"assigned_memory_mb,omitempty" structs:",omitempty"`
}

// PrivateHostPlan 専有ホストプラン
type PrivateHostPlan struct {
	ID           types.ID            `json:",omitempty" yaml:"id,omitempty" structs:",omitempty"`
	Name         string              `json:",omitempty" yaml:"name,omitempty" structs:",omitempty"`
	Class        string              `json:",omitempty" yaml:"class,omitempty" structs:",omitempty"`
	CPU          int                 `json:",omitempty" yaml:"cpu,omitempty" structs:",omitempty"`
	MemoryMB     int                 `json:",omitempty" yaml:"memory_mb,omitempty" structs:",omitempty"`
	ServiceClass string              `json:",omitempty" yaml:"service_class,omitempty" structs:",omitempty"`
	Availability types.EAvailability `json:",omitempty" yaml:"availability,omitempty" structs:",omitempty"`
}
