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

package define

import (
	"net/http"

	"github.com/sacloud/iaas-api-go/internal/define/names"
	"github.com/sacloud/iaas-api-go/internal/define/ops"
	"github.com/sacloud/iaas-api-go/internal/dsl"
	"github.com/sacloud/iaas-api-go/internal/dsl/meta"
	"github.com/sacloud/iaas-api-go/naked"
)

const (
	databaseAPIName     = "Database"
	databaseAPIPathName = "appliance"
)

var databaseAPI = &dsl.Resource{
	Name:       databaseAPIName,
	PathName:   databaseAPIPathName,
	PathSuffix: dsl.CloudAPISuffix,
	Operations: dsl.Operations{
		// find
		ops.FindAppliance(databaseAPIName, databaseNakedType, findParameter, databaseView),

		// create
		ops.CreateAppliance(databaseAPIName, databaseNakedType, databaseCreateParam, databaseView),

		// read
		ops.ReadAppliance(databaseAPIName, databaseNakedType, databaseView),

		// update
		ops.UpdateAppliance(databaseAPIName, databaseNakedType, databaseUpdateParam, databaseView),

		// updateSettings
		ops.UpdateApplianceSettings(databaseAPIName, databaseUpdateSettingsNakedType, databaseUpdateSettingsParam, databaseView),

		// delete
		ops.Delete(databaseAPIName),

		// config
		ops.Config(databaseAPIName),

		// power management(boot/shutdown/reset)
		ops.Boot(databaseAPIName),
		ops.Shutdown(databaseAPIName),
		ops.Reset(databaseAPIName),

		// monitor
		ops.MonitorChild(databaseAPIName, "CPU", "cpu",
			monitorParameter, monitors.cpuTimeModel()),
		ops.MonitorChild(databaseAPIName, "Disk", "disk/0",
			monitorParameter, monitors.diskModel()),
		ops.MonitorChild(databaseAPIName, "Interface", "interface",
			monitorParameter, monitors.interfaceModel()),
		ops.MonitorChild(databaseAPIName, "Database", "database",
			monitorParameter, monitors.databaseModel()),

		// status
		{
			ResourceName: databaseAPIName,
			Name:         "Status",
			PathFormat:   dsl.IDAndSuffixPathFormat("status"),
			Method:       http.MethodGet,
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
			},
			ResponseEnvelope: dsl.ResponseEnvelope(&dsl.EnvelopePayloadDesc{
				Type: meta.Static(naked.DatabaseStatusResponse{}),
				Name: "Appliance",
			}),
			Results: dsl.Results{
				{
					SourceField: "Appliance",
					DestField:   databaseStatusView.Name,
					IsPlural:    false,
					Model:       databaseStatusView,
				},
			},
		},
		// Parameter
		{
			ResourceName: databaseAPIName,
			Name:         "GetParameter",
			PathFormat:   dsl.IDAndSuffixPathFormat("database/parameter"),
			Method:       http.MethodGet,
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
			},
			ResponseEnvelope: dsl.ResponseEnvelope(&dsl.EnvelopePayloadDesc{
				Type: databaseParameterNakedType,
				Name: "Database",
			}),
			Results: dsl.Results{
				{
					SourceField: "Database",
					DestField:   databaseParameterView.Name,
					IsPlural:    false,
					Model:       databaseParameterView,
				},
			},
		},
		{
			ResourceName: databaseAPIName,
			Name:         "SetParameter",
			PathFormat:   dsl.IDAndSuffixPathFormat("database/parameter"),
			Method:       http.MethodPut,
			RequestEnvelope: dsl.RequestEnvelope(&dsl.EnvelopePayloadDesc{
				Type: meta.Static(naked.DatabaseParameterSetting{}),
				Name: "Parameter",
			}),
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
				{
					Name:       "param",
					Type:       meta.Static(map[string]interface{}{}),
					MapConvTag: "Parameter.Attr",
				},
			},
		},
	},
}

var (
	databaseNakedType               = meta.Static(naked.Database{})
	databaseUpdateSettingsNakedType = meta.Static(naked.DatabaseSettingsUpdate{})
	databaseStatusNakedType         = meta.Static(naked.DatabaseStatus{})
	databaseParameterNakedType      = meta.Static(naked.DatabaseParameter{})

	databaseView = &dsl.Model{
		Name:      databaseAPIName,
		NakedType: databaseNakedType,
		Fields: []*dsl.FieldDesc{
			fields.ID(),
			fields.Class(),
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.Availability(),
			fields.IconID(),
			fields.CreatedAt(),
			fields.ModifiedAt(),
			// settings
			fields.DatabaseSettingsCommon(),
			fields.DatabaseSettingsBackup(),
			fields.DatabaseSettingsReplication(),
			fields.DatabaseSettingsInterfaces(),
			fields.SettingsHash(),

			// instance
			fields.InstanceHostName(),
			fields.InstanceHostInfoURL(),
			fields.InstanceStatus(),
			fields.InstanceStatusChangedAt(),
			// plan
			fields.AppliancePlanID(),
			// switch
			fields.ApplianceSwitchID(),
			// remark
			fields.RemarkDBConf(),
			fields.RemarkDefaultRoute(),
			fields.RemarkNetworkMaskLen(),
			fields.RemarkServerIPAddress(),
			fields.RemarkZoneID(),
			// interfaces
			fields.Interfaces(),
		},
	}

	databaseCreateParam = &dsl.Model{
		Name:      names.CreateParameterName(databaseAPIName),
		NakedType: databaseNakedType,
		ConstFields: []*dsl.ConstFieldDesc{
			{
				Name:  "Class",
				Type:  meta.TypeString,
				Value: `"database"`,
			},
		},
		Fields: []*dsl.FieldDesc{
			fields.AppliancePlanID(),
			fields.ApplianceSwitchID(),
			fields.ApplianceIPAddresses(),
			fields.RemarkNetworkMaskLen(),
			fields.RemarkDefaultRoute(),
			fields.RemarkDBConf(),
			fields.RemarkSourceAppliance(),

			fields.DatabaseSettingsCommon(),
			fields.DatabaseSettingsBackup(),
			fields.DatabaseSettingsReplication(),
			fields.DatabaseSettingsInterfaces(),

			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.IconID(),
		},
	}

	databaseUpdateParam = &dsl.Model{
		Name:      names.UpdateParameterName(databaseAPIName),
		NakedType: databaseNakedType,
		Fields: []*dsl.FieldDesc{
			// common fields
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.IconID(),

			// settings
			fields.DatabaseSettingsCommon(),
			fields.DatabaseSettingsBackup(),
			fields.DatabaseSettingsReplication(),
			fields.DatabaseSettingsInterfaces(),
			// settings hash
			fields.SettingsHash(),
		},
	}

	databaseUpdateSettingsParam = &dsl.Model{
		Name:      names.UpdateSettingsParameterName(databaseAPIName),
		NakedType: databaseNakedType,
		Fields: []*dsl.FieldDesc{
			// settings
			fields.DatabaseSettingsCommon(),
			fields.DatabaseSettingsBackup(),
			fields.DatabaseSettingsReplication(),
			fields.DatabaseSettingsInterfaces(),
			// settings hash
			fields.SettingsHash(),
		},
	}

	databaseStatusView = &dsl.Model{
		Name:      "DatabaseStatus",
		NakedType: databaseStatusNakedType,
		Fields: []*dsl.FieldDesc{
			{
				Name: "Status",
				Type: meta.TypeInstanceStatus,
				Tags: &dsl.FieldTags{
					MapConv: "SettingsResponse.Status",
				},
			},
			{
				Name: "MariaDBStatus",
				Type: meta.TypeString,
				Tags: &dsl.FieldTags{
					MapConv: "SettingsResponse.DBConf.MariaDB.Status",
				},
			},
			{
				Name: "PostgresStatus",
				Type: meta.TypeString,
				Tags: &dsl.FieldTags{
					MapConv: "SettingsResponse.DBConf.Postgres.Status",
				},
			},
			{
				Name: "IsFatal",
				Type: meta.TypeFlag,
				Tags: &dsl.FieldTags{
					MapConv: "SettingsResponse.IsFatal",
				},
			},
			{
				Name: "Version",
				Type: databaseStatusVersionView,
				Tags: &dsl.FieldTags{
					MapConv: "SettingsResponse.DBConf.Version,recursive",
				},
			},
			{
				Name: "Logs",
				Type: databaseStatusLogView,
				Tags: &dsl.FieldTags{
					MapConv: "SettingsResponse.DBConf.[]Log,recursive",
				},
			},
			{
				Name: "Backups",
				Type: databaseStatusBackupHistoryView,
				Tags: &dsl.FieldTags{
					MapConv: "SettingsResponse.DBConf.Backup.[]History,recursive",
				},
			},
		},
	}

	databaseStatusVersionView = &dsl.Model{
		Name:      "DatabaseVersionInfo",
		NakedType: meta.Static(naked.DatabaseStatusVersion{}),
		Fields: []*dsl.FieldDesc{
			fields.Def("LastModified", meta.TypeString),
			fields.Def("CommitHash", meta.TypeString),
			fields.Def("Status", meta.TypeString),
			fields.Def("Tag", meta.TypeString),
			fields.Def("Expire", meta.TypeString),
		},
	}

	databaseStatusLogView = &dsl.Model{
		Name:      "DatabaseLog",
		NakedType: meta.Static(naked.DatabaseLog{}),
		IsArray:   true,
		Fields: []*dsl.FieldDesc{
			fields.Def("Name", meta.TypeString),
			fields.Def("Data", meta.TypeString),
			fields.Def("Size", meta.TypeStringNumber),
		},
	}
	databaseStatusBackupHistoryView = &dsl.Model{
		Name:      "DatabaseBackupHistory",
		NakedType: meta.Static(naked.DatabaseBackupHistory{}),
		IsArray:   true,
		Fields: []*dsl.FieldDesc{
			fields.Def("CreatedAt", meta.TypeTime),
			fields.Def("Availability", meta.TypeString),
			fields.Def("RecoveredAt", meta.TypeTime),
			fields.Def("Size", meta.TypeInt64),
		},
	}

	databaseParameterView = &dsl.Model{
		Name:      "DatabaseParameter",
		NakedType: databaseParameterNakedType,
		Fields: []*dsl.FieldDesc{
			{
				Name: "Settings",
				Type: meta.Static(map[string]interface{}{}),
				Tags: &dsl.FieldTags{
					MapConv: "Parameter.Attr",
				},
			},
			{
				Name: "MetaInfo",
				Type: databaseParameterMeta,
				Tags: &dsl.FieldTags{
					MapConv: "Remark.[]Form,recursive",
				},
			},
		},
	}

	databaseParameterMeta = &dsl.Model{
		Name:      "DatabaseParameterMeta",
		NakedType: meta.Static(naked.DatabaseParameterFormMeta{}),
		IsArray:   true,
		Fields: []*dsl.FieldDesc{
			{
				Name: "Type",
				Type: meta.TypeString,
				Tags: mapConvTag("Options.Type"),
			},
			fields.Def("Name", meta.TypeString),
			fields.Def("Label", meta.TypeString),
			{
				Name: "Text",
				Type: meta.TypeString,
				Tags: mapConvTag("Options.Text"),
			},
			{
				Name: "Example",
				Type: meta.TypeString,
				Tags: mapConvTag("Options.Example"),
			},
			{
				Name: "Min",
				Type: meta.TypeFloat64,
				Tags: mapConvTag("Options.Min"),
			},
			{
				Name: "Max",
				Type: meta.TypeFloat64,
				Tags: mapConvTag("Options.Max"),
			},
			{
				Name: "MaxLen",
				Type: meta.TypeInt,
				Tags: mapConvTag("Options.MaxLen"),
			},
			{
				Name: "Reboot",
				Type: meta.TypeString,
				Tags: mapConvTag("Options.Reboot"),
			},
		},
	}
)
