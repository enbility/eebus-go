package features

import (
	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type DeviceDiagnosisType struct {
	OperatingState       model.DeviceDiagnosisOperatingStateType
	PowerSupplyCondition model.PowerSupplyConditionType
}

type DeviceDiagnosis struct {
	*FeatureImpl
}

func NewDeviceDiagnosis(localRole, remoteRole model.RoleType, spineLocalDevice *spine.DeviceLocalImpl, entity *spine.EntityRemoteImpl) (*DeviceDiagnosis, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeDeviceDiagnosis, localRole, remoteRole, spineLocalDevice, entity)
	if err != nil {
		return nil, err
	}

	dd := &DeviceDiagnosis{
		FeatureImpl: feature,
	}

	return dd, nil
}

// request DeviceDiagnosisStateData from a remote entity
func (d *DeviceDiagnosis) RequestStateForEntity() (*model.MsgCounterType, error) {
	// request FunctionTypeDeviceDiagnosisStateData from a remote entity
	msgCounter, err := d.requestData(model.FunctionTypeDeviceDiagnosisStateData, nil, nil)
	if err != nil {
		logging.Log.Error(err)
		return nil, err
	}

	return msgCounter, nil
}

// get the current diagnosis state for an device entity
func (d *DeviceDiagnosis) GetState() (*DeviceDiagnosisType, error) {
	if d.featureRemote == nil {
		return nil, ErrDataNotAvailable
	}

	rData := d.featureRemote.Data(model.FunctionTypeDeviceDiagnosisStateData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.DeviceDiagnosisStateDataType)

	details := &DeviceDiagnosisType{}
	if data.OperatingState != nil {
		details.OperatingState = *data.OperatingState
	}
	if data.PowerSupplyCondition != nil {
		details.PowerSupplyCondition = *data.PowerSupplyCondition
	}

	return details, nil
}

func (d *DeviceDiagnosis) SendDeviceDiagnosisState(operatingState *model.DeviceDiagnosisStateDataType) {
	d.featureLocal.SetData(model.FunctionTypeDeviceDiagnosisStateData, operatingState)

	_, _ = d.featureLocal.NotifyData(model.FunctionTypeDeviceDiagnosisStateData, nil, nil, false, nil, d.featureRemote)
}
