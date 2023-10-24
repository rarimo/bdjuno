package rarimocore

import (
	"github.com/rarimo/bdjuno/types"
	rarimocoretypes "github.com/rarimo/rarimo-core/x/rarimocore/types"
)

func (m *Module) saveViolationReports(slice []rarimocoretypes.ViolationReport) error {
	reports := make([]types.ViolationReport, len(slice))
	for index, report := range slice {
		reports[index] = types.ViolationReportFromCore(report)
	}
	return m.db.SaveViolationReports(reports)
}
