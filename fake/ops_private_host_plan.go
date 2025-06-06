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

package fake

import (
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
)

// Find is fake implementation
func (o *PrivateHostPlanOp) Find(ctx context.Context, zone string, conditions *iaas.FindCondition) (*iaas.PrivateHostPlanFindResult, error) {
	results, _ := find(o.key, zone, conditions)
	var values []*iaas.PrivateHostPlan
	for _, res := range results {
		dest := &iaas.PrivateHostPlan{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &iaas.PrivateHostPlanFindResult{
		Total:            len(results),
		Count:            len(results),
		From:             0,
		PrivateHostPlans: values,
	}, nil
}

// Read is fake implementation
func (o *PrivateHostPlanOp) Read(ctx context.Context, zone string, id types.ID) (*iaas.PrivateHostPlan, error) {
	value := getPrivateHostPlanByID(zone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &iaas.PrivateHostPlan{}
	copySameNameField(value, dest)
	return dest, nil
}
