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

package power

import (
	"context"
	"strings"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
)

// ServerAPI APIクライアント
type ServerAPI interface {
	Read(ctx context.Context, zone string, id types.ID) (*iaas.Server, error)
	Boot(ctx context.Context, zone string, id types.ID) error
	BootWithVariables(ctx context.Context, zone string, id types.ID, param *iaas.ServerBootVariables) error
	Shutdown(ctx context.Context, zone string, id types.ID, shutdownOption *iaas.ShutdownOption) error
}

type serverHandler struct {
	ctx       context.Context
	client    ServerAPI
	zone      string
	id        types.ID
	variables []string
}

func (h *serverHandler) boot() error {
	if len(h.variables) > 0 && strings.Join(h.variables, "") != "" {
		userData := strings.Join(h.variables, "\n")
		return h.client.BootWithVariables(h.ctx, h.zone, h.id, &iaas.ServerBootVariables{UserData: userData})
	}
	return h.client.Boot(h.ctx, h.zone, h.id)
}

func (h *serverHandler) shutdown(force bool) error {
	return h.client.Shutdown(h.ctx, h.zone, h.id, &iaas.ShutdownOption{Force: force})
}

func (h *serverHandler) read() (interface{}, error) {
	return h.client.Read(h.ctx, h.zone, h.id)
}

// LoadBalancerAPI APIクライアント
type LoadBalancerAPI interface {
	Read(ctx context.Context, zone string, id types.ID) (*iaas.LoadBalancer, error)
	Boot(ctx context.Context, zone string, id types.ID) error
	Shutdown(ctx context.Context, zone string, id types.ID, shutdownOption *iaas.ShutdownOption) error
}

type loadBalancerHandler struct {
	ctx    context.Context
	client LoadBalancerAPI
	zone   string
	id     types.ID
}

func (h *loadBalancerHandler) boot() error {
	return h.client.Boot(h.ctx, h.zone, h.id)
}

func (h *loadBalancerHandler) shutdown(force bool) error {
	return h.client.Shutdown(h.ctx, h.zone, h.id, &iaas.ShutdownOption{Force: force})
}

func (h *loadBalancerHandler) read() (interface{}, error) {
	return h.client.Read(h.ctx, h.zone, h.id)
}

// DatabaseAPI APIクライアント
type DatabaseAPI interface {
	Read(ctx context.Context, zone string, id types.ID) (*iaas.Database, error)
	Boot(ctx context.Context, zone string, id types.ID) error
	Shutdown(ctx context.Context, zone string, id types.ID, shutdownOption *iaas.ShutdownOption) error
}

type databaseHandler struct {
	ctx    context.Context
	client DatabaseAPI
	zone   string
	id     types.ID
}

func (h *databaseHandler) boot() error {
	return h.client.Boot(h.ctx, h.zone, h.id)
}

func (h *databaseHandler) shutdown(force bool) error {
	return h.client.Shutdown(h.ctx, h.zone, h.id, &iaas.ShutdownOption{Force: force})
}

func (h *databaseHandler) read() (interface{}, error) {
	return h.client.Read(h.ctx, h.zone, h.id)
}

// VPCRouterAPI APIクライアント
type VPCRouterAPI interface {
	Read(ctx context.Context, zone string, id types.ID) (*iaas.VPCRouter, error)
	Boot(ctx context.Context, zone string, id types.ID) error
	Shutdown(ctx context.Context, zone string, id types.ID, shutdownOption *iaas.ShutdownOption) error
}

type vpcRouterHandler struct {
	ctx    context.Context
	client VPCRouterAPI
	zone   string
	id     types.ID
}

func (h *vpcRouterHandler) boot() error {
	return h.client.Boot(h.ctx, h.zone, h.id)
}

func (h *vpcRouterHandler) shutdown(force bool) error {
	return h.client.Shutdown(h.ctx, h.zone, h.id, &iaas.ShutdownOption{Force: force})
}

func (h *vpcRouterHandler) read() (interface{}, error) {
	return h.client.Read(h.ctx, h.zone, h.id)
}

// NFSAPI APIクライアント
type NFSAPI interface {
	Read(ctx context.Context, zone string, id types.ID) (*iaas.NFS, error)
	Boot(ctx context.Context, zone string, id types.ID) error
	Shutdown(ctx context.Context, zone string, id types.ID, shutdownOption *iaas.ShutdownOption) error
}

type nfsHandler struct {
	ctx    context.Context
	client NFSAPI
	zone   string
	id     types.ID
}

func (h *nfsHandler) boot() error {
	return h.client.Boot(h.ctx, h.zone, h.id)
}

func (h *nfsHandler) shutdown(force bool) error {
	return h.client.Shutdown(h.ctx, h.zone, h.id, &iaas.ShutdownOption{Force: force})
}

func (h *nfsHandler) read() (interface{}, error) {
	return h.client.Read(h.ctx, h.zone, h.id)
}

// MobileGatewayAPI APIクライアント
type MobileGatewayAPI interface {
	Read(ctx context.Context, zone string, id types.ID) (*iaas.MobileGateway, error)
	Boot(ctx context.Context, zone string, id types.ID) error
	Shutdown(ctx context.Context, zone string, id types.ID, shutdownOption *iaas.ShutdownOption) error
}

type mobileGatewayHandler struct {
	ctx    context.Context
	client MobileGatewayAPI
	zone   string
	id     types.ID
}

func (h *mobileGatewayHandler) boot() error {
	return h.client.Boot(h.ctx, h.zone, h.id)
}

func (h *mobileGatewayHandler) shutdown(force bool) error {
	return h.client.Shutdown(h.ctx, h.zone, h.id, &iaas.ShutdownOption{Force: force})
}

func (h *mobileGatewayHandler) read() (interface{}, error) {
	return h.client.Read(h.ctx, h.zone, h.id)
}
