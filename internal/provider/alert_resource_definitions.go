package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/alerts"
)

func NewThresholdAlertResource() resource.Resource {
	return &AlertResource[ThresholdAlertModel]{
		typeName:          ALERT_TYPE_THRESHOLD,
		fromModelFunc:     ThresholdAlertFromModel,
		toModelFunc:       ThresholdAlertToModel,
		getIdFunc:         func(m *ThresholdAlertModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ThresholdAlertModel) basetypes.StringValue { return m.PipelineId },
		schema:            ThresholdAlertResourceSchema,
	}
}
func NewChangeAlertResource() resource.Resource {
	return &AlertResource[ChangeAlertModel]{
		typeName:          ALERT_TYPE_CHANGE,
		fromModelFunc:     ChangeAlertFromModel,
		toModelFunc:       ChangeAlertToModel,
		getIdFunc:         func(m *ChangeAlertModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ChangeAlertModel) basetypes.StringValue { return m.PipelineId },
		schema:            ChangeAlertResourceSchema,
	}
}
func NewAbsenceAlertResource() resource.Resource {
	return &AlertResource[AbsenceAlertModel]{
		typeName:          ALERT_TYPE_ABSENCE,
		fromModelFunc:     AbsenceAlertFromModel,
		toModelFunc:       AbsenceAlertToModel,
		getIdFunc:         func(m *AbsenceAlertModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *AbsenceAlertModel) basetypes.StringValue { return m.PipelineId },
		schema:            AbsenceAlertResourceSchema,
	}
}
