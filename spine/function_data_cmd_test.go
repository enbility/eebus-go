package spine

import (
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestFunctionDataCmdSuite(t *testing.T) {
	suite.Run(t, new(FunctionDataCmdTestSuite))
}

type FunctionDataCmdTestSuite struct {
	suite.Suite
	function model.FunctionEnumType
	data     *model.DeviceClassificationManufacturerDataType
	sut      *FunctionDataCmdImpl[model.DeviceClassificationManufacturerDataType]
}

func (suite *FunctionDataCmdTestSuite) SetupSuite() {
	suite.function = model.FunctionEnumTypeDeviceClassificationManufacturerData
	suite.data = &model.DeviceClassificationManufacturerDataType{
		DeviceName: util.Ptr(model.DeviceClassificationStringType("device name")),
	}
	suite.sut = NewFunctionDataCmd[model.DeviceClassificationManufacturerDataType](suite.function)
	suite.sut.SetData(suite.data)
}

func (suite *FunctionDataCmdTestSuite) TestFunctionDataLocal_ReadCmd() {
	readCmd := suite.sut.ReadCmdType()
	assert.NotNil(suite.T(), readCmd.DeviceClassificationManufacturerData)
	assert.Nil(suite.T(), readCmd.DeviceClassificationManufacturerData.DeviceName)
	// TODO: assert on json
}

func (suite *FunctionDataCmdTestSuite) TestFunctionDataLocal_ReplyCmd() {
	readCmd := suite.sut.ReplyCmdType()
	assert.NotNil(suite.T(), readCmd.DeviceClassificationManufacturerData)
	assert.Equal(suite.T(), suite.data.DeviceName, readCmd.DeviceClassificationManufacturerData.DeviceName)
	// TODO: assert on json
}

// TODO: test NotifyCmdType
