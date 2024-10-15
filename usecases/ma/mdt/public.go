package mdt

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// Scenario 1

// return the momentary temperature of the domestic hot water circuit
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *MDT) Temperature(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil {
		return 0, err
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeTemperature),
		CommodityType:   util.Ptr(model.CommodityTypeTypeDomestichotwater),
		ScopeType:       util.Ptr(model.ScopeTypeTypeDhwTemperature),
	}
	data, err := measurement.GetDataForFilter(filter)
	if err != nil || len(data) == 0 || data[0].Value == nil {
		return 0, api.ErrDataNotAvailable
	}

	// take the first item
	value := data[0].Value

	return value.GetValue(), nil
}
