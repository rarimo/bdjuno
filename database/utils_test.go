package database_test

import (
	"gitlab.com/rarimo/bdjuno/v3/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_InsertEnableModules() {
	modules := []string{"auth", "bank", "consensus", "distribution", "gov", "mint", "pricefeed", "staking", "supply"}
	err := suite.database.InsertEnableModules(modules)
	suite.Require().NoError(err)

	var results types.ModuleRows
	err = suite.database.Sqlx.Select(&results, "SELECT * FROM modules")
	suite.Require().NoError(err)

	expected := types.NewModuleRows(modules)
	suite.Require().True(results.Equal(&expected))

}
