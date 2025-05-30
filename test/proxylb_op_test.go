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

package test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/stretchr/testify/assert"
)

func TestProxyLBOp_CRUD(t *testing.T) {
	initProxyLBVariables()

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel: true,

		PreCheck: testutil.PreCheckEnvsFunc("SAKURACLOUD_PROXYLB_SERVER0", "SAKURACLOUD_PROXYLB_SERVER1", "SAKURACLOUD_PROXYLB_SERVER2"),

		SetupAPICallerFunc: singletonAPICaller,

		Create: &testutil.CRUDTestFunc{
			Func: testProxyLBCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createProxyLBExpected,
				IgnoreFields: ignoreProxyLBFields,
			}),
		},

		Read: &testutil.CRUDTestFunc{
			Func: testProxyLBRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createProxyLBExpected,
				IgnoreFields: ignoreProxyLBFields,
			}),
		},

		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testProxyLBUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateProxyLBExpected,
					IgnoreFields: ignoreProxyLBFields,
				}),
			},
			{
				Func: testProxyLBUpdatePlan,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateProxyLBPlanExpected,
					IgnoreFields: ignoreProxyLBFields,
				}),
			},
			{
				Func: testProxyLBUpdateSettings,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateProxyLBSettingsExpected,
					IgnoreFields: ignoreProxyLBFields,
				}),
			},
			{
				Func: testProxyLBUpdateToMin,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateProxyLBToMinExpected,
					IgnoreFields: ignoreProxyLBFields,
				}),
			},
		},

		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testProxyLBDelete,
		},
	})
}

var (
	ignoreProxyLBFields           []string
	createProxyLBParam            *iaas.ProxyLBCreateRequest
	createProxyLBExpected         *iaas.ProxyLB
	updateProxyLBParam            *iaas.ProxyLBUpdateRequest
	updateProxyLBExpected         *iaas.ProxyLB
	updateProxyLBPlanExpected     *iaas.ProxyLB
	updateProxyLBSettingsParam    *iaas.ProxyLBUpdateSettingsRequest
	updateProxyLBSettingsExpected *iaas.ProxyLB
	updateProxyLBToMinParam       *iaas.ProxyLBUpdateRequest
	updateProxyLBToMinExpected    *iaas.ProxyLB

	createProxyLBForACMEParam *iaas.ProxyLBCreateRequest
	updateProxyLBForACMEParam *iaas.ProxyLBUpdateRequest
)

func initProxyLBVariables() {
	ignoreProxyLBFields = []string{
		"ID",
		"CreatedAt",
		"ModifiedAt",
		"Class",
		"SettingsHash",
		"Region",
		"ProxyNetworks",
		"FQDN",
	}

	createProxyLBParam = &iaas.ProxyLBCreateRequest{
		Name:        testutil.ResourceName("proxylb"),
		Description: "desc",
		Tags:        []string{"tag1", "tag2"},
		Plan:        types.ProxyLBPlans.CPS100,
		HealthCheck: &iaas.ProxyLBHealthCheck{
			Protocol:  types.ProxyLBProtocols.HTTP,
			Path:      "/",
			DelayLoop: 10,
		},
		SorryServer: &iaas.ProxyLBSorryServer{
			IPAddress: os.Getenv("SAKURACLOUD_PROXYLB_SERVER0"),
			Port:      80,
		},
		BindPorts: []*iaas.ProxyLBBindPort{
			{
				ProxyMode:       types.ProxyLBProxyModes.HTTP,
				Port:            80,
				RedirectToHTTPS: true,
				AddResponseHeader: []*iaas.ProxyLBResponseHeader{
					{
						Header: "Cache-Control",
						Value:  "public, max-age=60",
					},
				},
			},
			{
				ProxyMode:    types.ProxyLBProxyModes.HTTPS,
				Port:         443,
				SupportHTTP2: true,
				SSLPolicy:    "TLS-1-3-2021-06",
			},
		},
		Servers: []*iaas.ProxyLBServer{
			{
				IPAddress:   os.Getenv("SAKURACLOUD_PROXYLB_SERVER1"),
				Port:        80,
				ServerGroup: "group1",
				Enabled:     true,
			},
			{
				IPAddress:   os.Getenv("SAKURACLOUD_PROXYLB_SERVER2"),
				Port:        80,
				ServerGroup: "group2",
				Enabled:     true,
			},
		},
		Rules: []*iaas.ProxyLBRule{
			{
				Host:        "www.usacloud.jp",
				Path:        "/path1",
				SourceIPs:   "192.168.0.1,192.168.1.0/24",
				ServerGroup: "group1",
				Action:      types.ProxyLBRuleActions.Forward,
			},
			{
				Host:        "www.usacloud.jp",
				Path:        "/path2",
				ServerGroup: "group2",
				Action:      types.ProxyLBRuleActions.Forward,
			},
			{
				RequestHeaderName:            "User-Agent",
				RequestHeaderValue:           "*badbot",
				RequestHeaderValueIgnoreCase: false,
				RequestHeaderValueNotMatch:   false,
				ServerGroup:                  "group1",
				Action:                       types.ProxyLBRuleActions.Forward,
			},
			{
				Host:             "www.usacloud.jp",
				Path:             "/fixed-response",
				Action:           types.ProxyLBRuleActions.Fixed,
				FixedStatusCode:  types.ProxyLBFixedStatusCodes.OK,
				FixedContentType: types.ProxyLBFixedContentTypes.Plain,
				FixedMessageBody: "foobar",
			},
			{
				Host:               "www.usacloud.jp",
				Path:               "/redirect",
				Action:             types.ProxyLBRuleActions.Redirect,
				RedirectLocation:   "https://redirect.usacloud.jp",
				RedirectStatusCode: types.ProxyLBRedirectStatusCodes.Found,
			},
		},
		LetsEncrypt: &iaas.ProxyLBACMESetting{
			Enabled: false,
		},
		StickySession: &iaas.ProxyLBStickySession{
			Method:  "cookie",
			Enabled: true,
		},
		Gzip: &iaas.ProxyLBGzip{
			Enabled: true,
		},
		BackendHttpKeepAlive: &iaas.ProxyLBBackendHttpKeepAlive{
			Mode: types.ProxyLBBackendHttpKeepAlive.Aggressive,
		},
		ProxyProtocol: &iaas.ProxyLBProxyProtocol{
			Enabled: true,
		},
		Syslog: &iaas.ProxyLBSyslog{
			Server: "133.242.0.1",
			Port:   514,
		},
		Timeout: &iaas.ProxyLBTimeout{
			InactiveSec: 30,
		},
		UseVIPFailover: true,
		Region:         types.ProxyLBRegions.Anycast,
	}
	createProxyLBExpected = &iaas.ProxyLB{
		Name:         createProxyLBParam.Name,
		Description:  createProxyLBParam.Description,
		Tags:         createProxyLBParam.Tags,
		Availability: types.Availabilities.Available,

		Plan:                 createProxyLBParam.Plan,
		HealthCheck:          createProxyLBParam.HealthCheck,
		SorryServer:          createProxyLBParam.SorryServer,
		BindPorts:            createProxyLBParam.BindPorts,
		Servers:              createProxyLBParam.Servers,
		Rules:                createProxyLBParam.Rules,
		LetsEncrypt:          createProxyLBParam.LetsEncrypt,
		StickySession:        createProxyLBParam.StickySession,
		Gzip:                 createProxyLBParam.Gzip,
		BackendHttpKeepAlive: createProxyLBParam.BackendHttpKeepAlive,
		ProxyProtocol:        createProxyLBParam.ProxyProtocol,
		Syslog:               createProxyLBParam.Syslog,
		Timeout:              createProxyLBParam.Timeout,
		UseVIPFailover:       createProxyLBParam.UseVIPFailover,
		Region:               createProxyLBParam.Region,
	}
	updateProxyLBParam = &iaas.ProxyLBUpdateRequest{
		Name:        testutil.ResourceName("proxylb-upd"),
		Description: "desc-upd",
		Tags:        []string{"tag1-upd", "tag2-upd"},
		IconID:      testIconID,
		HealthCheck: &iaas.ProxyLBHealthCheck{
			Protocol:  types.ProxyLBProtocols.HTTP,
			Path:      "/index.html",
			DelayLoop: 20,
		},
		SorryServer: &iaas.ProxyLBSorryServer{
			IPAddress: os.Getenv("SAKURACLOUD_PROXYLB_SERVER0"),
			Port:      8080,
		},
		BindPorts: []*iaas.ProxyLBBindPort{
			{
				ProxyMode:       types.ProxyLBProxyModes.HTTP,
				Port:            8080,
				RedirectToHTTPS: true,
			},
			{
				ProxyMode:    types.ProxyLBProxyModes.HTTPS,
				Port:         8443,
				SupportHTTP2: true,
				SSLPolicy:    "TLS-1-3-2021-06",
			},
		},
		Servers: []*iaas.ProxyLBServer{
			{
				IPAddress:   os.Getenv("SAKURACLOUD_PROXYLB_SERVER1"),
				Port:        8080,
				ServerGroup: "group1upd",
				Enabled:     true,
			},
			{
				IPAddress:   os.Getenv("SAKURACLOUD_PROXYLB_SERVER2"),
				Port:        8080,
				ServerGroup: "group2upd",
				Enabled:     true,
			},
		},
		Rules: []*iaas.ProxyLBRule{
			{
				Host:        "www-upd.usacloud.jp",
				Path:        "/path1-upd",
				ServerGroup: "group1upd",
				Action:      types.ProxyLBRuleActions.Forward,
			},
			{
				Host:        "www-upd.usacloud.jp",
				Path:        "/path2-upd",
				ServerGroup: "group2upd",
				Action:      types.ProxyLBRuleActions.Forward,
			},
			{
				RequestHeaderName:            "User-Agent",
				RequestHeaderValue:           "*badbot",
				RequestHeaderValueIgnoreCase: false,
				RequestHeaderValueNotMatch:   false,
				ServerGroup:                  "group1upd",
				Action:                       types.ProxyLBRuleActions.Forward,
			},
			{
				Host:             "www-upd.usacloud.jp",
				Path:             "/fixed-response-upd",
				Action:           types.ProxyLBRuleActions.Fixed,
				FixedStatusCode:  types.ProxyLBFixedStatusCodes.Forbidden,
				FixedContentType: types.ProxyLBFixedContentTypes.HTML,
				FixedMessageBody: "foobar-upd",
			},
			{
				Host:               "www-upd.usacloud.jp",
				Path:               "/redirect-upd",
				Action:             types.ProxyLBRuleActions.Redirect,
				RedirectLocation:   "https://redirect.usacloud.jp/upd",
				RedirectStatusCode: types.ProxyLBRedirectStatusCodes.MovedPermanently,
			},
		},
		// LetsEncryptのテストはA or CNAMEレコードの登録が必要なため別ケースで行う
		LetsEncrypt: &iaas.ProxyLBACMESetting{
			Enabled: false,
		},
		StickySession: &iaas.ProxyLBStickySession{
			Enabled: false,
		},
		Gzip: &iaas.ProxyLBGzip{
			Enabled: false,
		},
		BackendHttpKeepAlive: &iaas.ProxyLBBackendHttpKeepAlive{
			Mode: types.ProxyLBBackendHttpKeepAlive.Safe,
		},
		ProxyProtocol: &iaas.ProxyLBProxyProtocol{
			Enabled: false,
		},
		Syslog: &iaas.ProxyLBSyslog{
			Server: "",
			Port:   514,
		},
		Timeout: &iaas.ProxyLBTimeout{
			InactiveSec: 10,
		},
	}
	updateProxyLBExpected = &iaas.ProxyLB{
		Name:          updateProxyLBParam.Name,
		Description:   updateProxyLBParam.Description,
		Tags:          updateProxyLBParam.Tags,
		IconID:        testIconID,
		Availability:  types.Availabilities.Available,
		Plan:          createProxyLBParam.Plan,
		HealthCheck:   updateProxyLBParam.HealthCheck,
		SorryServer:   updateProxyLBParam.SorryServer,
		BindPorts:     updateProxyLBParam.BindPorts,
		Servers:       updateProxyLBParam.Servers,
		Rules:         updateProxyLBParam.Rules,
		LetsEncrypt:   updateProxyLBParam.LetsEncrypt,
		StickySession: updateProxyLBParam.StickySession,
		Timeout: &iaas.ProxyLBTimeout{
			InactiveSec: 10,
		},
		Gzip:                 updateProxyLBParam.Gzip,
		BackendHttpKeepAlive: updateProxyLBParam.BackendHttpKeepAlive,
		ProxyProtocol:        updateProxyLBParam.ProxyProtocol,
		Syslog:               updateProxyLBParam.Syslog,
		UseVIPFailover:       createProxyLBParam.UseVIPFailover,
		Region:               createProxyLBParam.Region,
	}
	updateProxyLBPlanExpected = &iaas.ProxyLB{
		Name:          updateProxyLBParam.Name,
		Description:   updateProxyLBParam.Description,
		Tags:          updateProxyLBParam.Tags,
		IconID:        testIconID,
		Availability:  types.Availabilities.Available,
		Plan:          types.ProxyLBPlans.CPS500,
		HealthCheck:   updateProxyLBParam.HealthCheck,
		SorryServer:   updateProxyLBParam.SorryServer,
		BindPorts:     updateProxyLBParam.BindPorts,
		Servers:       updateProxyLBParam.Servers,
		Rules:         updateProxyLBParam.Rules,
		LetsEncrypt:   updateProxyLBParam.LetsEncrypt,
		StickySession: updateProxyLBParam.StickySession,
		Timeout: &iaas.ProxyLBTimeout{
			InactiveSec: 10,
		},
		Gzip:                 updateProxyLBParam.Gzip,
		BackendHttpKeepAlive: updateProxyLBParam.BackendHttpKeepAlive,
		ProxyProtocol:        updateProxyLBParam.ProxyProtocol,
		Syslog:               updateProxyLBParam.Syslog,
		UseVIPFailover:       createProxyLBParam.UseVIPFailover,
		Region:               createProxyLBParam.Region,
	}
	updateProxyLBSettingsParam = &iaas.ProxyLBUpdateSettingsRequest{
		HealthCheck: &iaas.ProxyLBHealthCheck{
			Protocol:  types.ProxyLBProtocols.TCP,
			DelayLoop: 20,
		},
		SorryServer: &iaas.ProxyLBSorryServer{
			IPAddress: os.Getenv("SAKURACLOUD_PROXYLB_SERVER0"),
			Port:      8080,
		},
		BindPorts: []*iaas.ProxyLBBindPort{
			{
				ProxyMode: types.ProxyLBProxyModes.HTTP, Port: 8081,
				RedirectToHTTPS: true,
			},
			{
				ProxyMode:    types.ProxyLBProxyModes.HTTPS,
				Port:         8443,
				SupportHTTP2: true,
				SSLPolicy:    "TLS-1-3-2021-06",
			},
		},
		Servers: []*iaas.ProxyLBServer{
			{
				IPAddress:   os.Getenv("SAKURACLOUD_PROXYLB_SERVER1"),
				Port:        8081,
				ServerGroup: "group1upd2",
				Enabled:     true,
			},
			{
				IPAddress:   os.Getenv("SAKURACLOUD_PROXYLB_SERVER2"),
				Port:        8081,
				ServerGroup: "group2upd2",
				Enabled:     true,
			},
		},
		Rules: []*iaas.ProxyLBRule{
			{
				Host:        "www-upd2.usacloud.jp",
				Path:        "/path1-upd2",
				ServerGroup: "group1upd2",
				Action:      types.ProxyLBRuleActions.Forward,
			},
			{
				Host:        "www-upd2.usacloud.jp",
				Path:        "/path2-upd2",
				ServerGroup: "group2upd2",
				Action:      types.ProxyLBRuleActions.Forward,
			},
		},
		// LetsEncryptのテストはA or CNAMEレコードの登録が必要なため別ケースで行う
		LetsEncrypt: &iaas.ProxyLBACMESetting{
			Enabled: false,
		},
		StickySession: &iaas.ProxyLBStickySession{
			Method:  "cookie",
			Enabled: true,
		},
		Timeout: &iaas.ProxyLBTimeout{
			InactiveSec: 10,
		},
		Gzip: &iaas.ProxyLBGzip{
			Enabled: false,
		},
		ProxyProtocol: &iaas.ProxyLBProxyProtocol{
			Enabled: false,
		},
		Syslog: &iaas.ProxyLBSyslog{
			Server: "",
			Port:   514,
		},
	}
	updateProxyLBSettingsExpected = &iaas.ProxyLB{
		Name:          updateProxyLBParam.Name,
		Description:   updateProxyLBParam.Description,
		Tags:          updateProxyLBParam.Tags,
		IconID:        testIconID,
		Availability:  types.Availabilities.Available,
		Plan:          updateProxyLBPlanExpected.Plan,
		HealthCheck:   updateProxyLBSettingsParam.HealthCheck,
		SorryServer:   updateProxyLBSettingsParam.SorryServer,
		BindPorts:     updateProxyLBSettingsParam.BindPorts,
		Servers:       updateProxyLBSettingsParam.Servers,
		Rules:         updateProxyLBSettingsParam.Rules,
		LetsEncrypt:   updateProxyLBSettingsParam.LetsEncrypt,
		StickySession: updateProxyLBSettingsParam.StickySession,
		Timeout: &iaas.ProxyLBTimeout{
			InactiveSec: 10,
		},
		Gzip:                 updateProxyLBSettingsParam.Gzip,
		BackendHttpKeepAlive: updateProxyLBSettingsParam.BackendHttpKeepAlive,
		ProxyProtocol:        updateProxyLBSettingsParam.ProxyProtocol,
		Syslog:               updateProxyLBSettingsParam.Syslog,
		UseVIPFailover:       createProxyLBParam.UseVIPFailover,
		Region:               createProxyLBParam.Region,
	}

	updateProxyLBToMinParam = &iaas.ProxyLBUpdateRequest{
		Name: testutil.ResourceName("proxylb-to-min"),
		HealthCheck: &iaas.ProxyLBHealthCheck{
			Protocol:  types.ProxyLBProtocols.TCP,
			DelayLoop: 10,
		},
		LetsEncrypt: &iaas.ProxyLBACMESetting{
			Enabled: false,
		},
		StickySession: &iaas.ProxyLBStickySession{
			Enabled: false,
		},
		Timeout: &iaas.ProxyLBTimeout{
			InactiveSec: 10,
		},
		Gzip: &iaas.ProxyLBGzip{
			Enabled: false,
		},
		BackendHttpKeepAlive: &iaas.ProxyLBBackendHttpKeepAlive{
			Mode: types.ProxyLBBackendHttpKeepAlive.Safe,
		},
		ProxyProtocol: &iaas.ProxyLBProxyProtocol{
			Enabled: false,
		},
		Syslog: &iaas.ProxyLBSyslog{
			Server: "",
			Port:   514,
		},
		BindPorts: []*iaas.ProxyLBBindPort{},
		Rules:     []*iaas.ProxyLBRule{},
		Servers:   []*iaas.ProxyLBServer{},
	}
	updateProxyLBToMinExpected = &iaas.ProxyLB{
		Name:         updateProxyLBToMinParam.Name,
		Availability: types.Availabilities.Available,
		Plan:         updateProxyLBPlanExpected.Plan,
		HealthCheck:  updateProxyLBToMinParam.HealthCheck,
		SorryServer:  &iaas.ProxyLBSorryServer{},
		LetsEncrypt: &iaas.ProxyLBACMESetting{
			Enabled: false,
		},
		StickySession: &iaas.ProxyLBStickySession{
			Enabled: false,
		},
		Timeout: &iaas.ProxyLBTimeout{
			InactiveSec: 10,
		},
		BindPorts:            updateProxyLBToMinParam.BindPorts,
		Rules:                updateProxyLBToMinParam.Rules,
		Servers:              updateProxyLBToMinParam.Servers,
		Gzip:                 updateProxyLBToMinParam.Gzip,
		BackendHttpKeepAlive: updateProxyLBToMinParam.BackendHttpKeepAlive,
		ProxyProtocol:        updateProxyLBToMinParam.ProxyProtocol,
		Syslog:               updateProxyLBToMinParam.Syslog,
		UseVIPFailover:       createProxyLBParam.UseVIPFailover,
		Region:               createProxyLBParam.Region,
	}

	createProxyLBForACMEParam = &iaas.ProxyLBCreateRequest{
		Name: testutil.ResourceName("proxylb-acme"),
		Plan: types.ProxyLBPlans.CPS100,
		HealthCheck: &iaas.ProxyLBHealthCheck{
			Protocol:  types.ProxyLBProtocols.HTTP,
			Path:      "/",
			DelayLoop: 20,
		},
		BindPorts: []*iaas.ProxyLBBindPort{
			{
				ProxyMode:       types.ProxyLBProxyModes.HTTP,
				Port:            80,
				RedirectToHTTPS: true,
			},
			{
				ProxyMode:    types.ProxyLBProxyModes.HTTPS,
				Port:         443,
				SupportHTTP2: true,
			},
		},
		Servers: []*iaas.ProxyLBServer{
			{
				IPAddress: os.Getenv("SAKURACLOUD_PROXYLB_SERVER0"),
				Port:      80,
				Enabled:   true,
			},
			{
				IPAddress: os.Getenv("SAKURACLOUD_PROXYLB_SERVER1"),
				Port:      80,
				Enabled:   true,
			},
		},
		LetsEncrypt: &iaas.ProxyLBACMESetting{
			Enabled: false,
		},
		Timeout: &iaas.ProxyLBTimeout{
			InactiveSec: 10,
		},
		UseVIPFailover: true,
	}

	updateProxyLBForACMEParam = &iaas.ProxyLBUpdateRequest{
		Name: testutil.ResourceName("proxylb-acme"),
		HealthCheck: &iaas.ProxyLBHealthCheck{
			Protocol:  types.ProxyLBProtocols.HTTP,
			Path:      "/",
			DelayLoop: 20,
		},
		BindPorts: []*iaas.ProxyLBBindPort{
			{
				ProxyMode:       types.ProxyLBProxyModes.HTTP,
				Port:            80,
				RedirectToHTTPS: true,
			},
			{
				ProxyMode:    types.ProxyLBProxyModes.HTTPS,
				Port:         443,
				SupportHTTP2: true,
			},
		},
		Servers: []*iaas.ProxyLBServer{
			{
				IPAddress: os.Getenv("SAKURACLOUD_PROXYLB_SERVER0"),
				Port:      80,
				Enabled:   true,
			},
			{
				IPAddress: os.Getenv("SAKURACLOUD_PROXYLB_SERVER1"),
				Port:      80,
				Enabled:   true,
			},
		},
		LetsEncrypt: &iaas.ProxyLBACMESetting{
			CommonName:      os.Getenv("SAKURACLOUD_PROXYLB_COMMON_NAME"),
			Enabled:         true,
			SubjectAltNames: []string{os.Getenv("SAKURACLOUD_PROXYLB_ALT_NAME")},
		},
		Timeout: &iaas.ProxyLBTimeout{
			InactiveSec: 10,
		},
	}
}

func testProxyLBCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewProxyLBOp(caller)
	return client.Create(ctx, createProxyLBParam)
}

func testProxyLBRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewProxyLBOp(caller)
	return client.Read(ctx, ctx.ID)
}

func testProxyLBUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewProxyLBOp(caller)
	return client.Update(ctx, ctx.ID, updateProxyLBParam)
}

func testProxyLBUpdatePlan(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewProxyLBOp(caller)
	return client.ChangePlan(ctx, ctx.ID, &iaas.ProxyLBChangePlanRequest{
		ServiceClass: types.ProxyLBServiceClass(types.ProxyLBPlans.CPS500, types.ProxyLBRegions.Anycast),
	})
}

func testProxyLBUpdateSettings(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewProxyLBOp(caller)
	return client.UpdateSettings(ctx, ctx.ID, updateProxyLBSettingsParam)
}

func testProxyLBUpdateToMin(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewProxyLBOp(caller)
	return client.Update(ctx, ctx.ID, updateProxyLBToMinParam)
}

func testProxyLBDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewProxyLBOp(caller)
	return client.Delete(ctx, ctx.ID)
}

func TestProxyLBOpLetsEncryptAndHealth(t *testing.T) {
	if !isAccTest() {
		t.Skip("TestProxyLBOpLetsEncrypt only exec at Acceptance Test")
	}

	t.Parallel()
	initProxyLBVariables()
	testutil.PreCheckEnvsFunc(
		"SAKURACLOUD_PROXYLB_SERVER0",
		"SAKURACLOUD_PROXYLB_SERVER1",
		"SAKURACLOUD_PROXYLB_COMMON_NAME",
		"SAKURACLOUD_PROXYLB_ALT_NAME",
		"SAKURACLOUD_PROXYLB_ZONE_NAME",
	)(t)

	// prepare variables
	commonName := os.Getenv("SAKURACLOUD_PROXYLB_COMMON_NAME")
	altName := os.Getenv("SAKURACLOUD_PROXYLB_ALT_NAME")
	zoneName := os.Getenv("SAKURACLOUD_PROXYLB_ZONE_NAME")
	if !strings.HasSuffix(commonName, zoneName) {
		t.Fatal("$SAKURACLOUD_PROXYLB_COMMON_NAME does not have suffix $SAKURACLOUD_PROXYLB_ZONE_NAME")
	}
	if !strings.HasSuffix(altName, zoneName) {
		t.Fatal("$SAKURACLOUD_PROXYLB_ALT_NAME does not have suffix $SAKURACLOUD_PROXYLB_ZONE_NAME")
	}
	recordName1 := strings.ReplaceAll(commonName, "."+zoneName, "")
	recordName2 := strings.ReplaceAll(altName, "."+zoneName, "")

	ctx := context.Background()
	proxyLBOp := iaas.NewProxyLBOp(singletonAPICaller())

	// create proxyLB
	proxyLB, err := proxyLBOp.Create(ctx, createProxyLBForACMEParam)
	if !assert.NoError(t, err) {
		return
	}
	defer func() {
		proxyLBOp.Delete(ctx, proxyLB.ID) //nolint:errcheck
	}()

	// read DNS
	dns, err := lookupDNSByName(singletonAPICaller(), os.Getenv("SAKURACLOUD_PROXYLB_ZONE_NAME"))
	if !assert.NoError(t, err) {
		return
	}

	dns.Records = append(dns.Records, &iaas.DNSRecord{
		Name:  recordName1,
		Type:  types.DNSRecordTypes.CNAME,
		RData: fmt.Sprintf("%s.", proxyLB.FQDN),
		TTL:   10,
	})
	dns.Records = append(dns.Records, &iaas.DNSRecord{
		Name:  recordName2,
		Type:  types.DNSRecordTypes.CNAME,
		RData: fmt.Sprintf("%s.", proxyLB.FQDN),
		TTL:   10,
	})

	// update DNS record
	dnsOp := iaas.NewDNSOp(singletonAPICaller())
	dns, err = dnsOp.Update(ctx, dns.ID, &iaas.DNSUpdateRequest{
		Records: dns.Records,
	})
	if !assert.NoError(t, err) {
		return
	}
	defer func() {
		var records []*iaas.DNSRecord
		for i, r := range dns.Records {
			if r.Name != recordName1 && r.Name != recordName2 {
				records = append(records, dns.Records[i])
			}
		}
		dnsOp.Update(ctx, dns.ID, &iaas.DNSUpdateRequest{Records: records}) //nolint:errcheck
	}()

	time.Sleep(time.Minute)

	// update proxyLB
	retryMax := 10
	done := false

	for retryMax >= 0 {
		proxyLB, err = proxyLBOp.Update(ctx, proxyLB.ID, updateProxyLBForACMEParam)
		if err != nil {
			t.Log("Update Let's encrypt setting is failed. retry after 10 sec.")
			time.Sleep(10 * time.Second)
			retryMax--
			continue
		}
		done = true
		break
	}
	if !done {
		t.Error("Update Let's encrypt settings was failed: given up after 10 retries")
		return
	}

	// renew certs
	err = proxyLBOp.RenewLetsEncryptCert(ctx, proxyLB.ID)
	if !assert.NoError(t, err) {
		return
	}

	time.Sleep(time.Minute)

	// get cert
	certs, err := proxyLBOp.GetCertificates(ctx, proxyLB.ID)

	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, certs.PrimaryCert)
	assert.NotEmpty(t, certs.PrimaryCert.ServerCertificate)
	assert.NotEmpty(t, certs.PrimaryCert.IntermediateCertificate)
	assert.NotEmpty(t, certs.PrimaryCert.PrivateKey)
	assert.NotEmpty(t, certs.PrimaryCert.CertificateCommonName)
	assert.NotEmpty(t, certs.PrimaryCert.CertificateAltNames)
	assert.NotEmpty(t, certs.PrimaryCert.CertificateEndDate)

	// check health status
	status, err := proxyLBOp.HealthStatus(ctx, proxyLB.ID)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotEmpty(t, status.CurrentVIP)
}
