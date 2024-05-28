package internal

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	spineapi "github.com/enbility/spine-go/api"
)

// return the current manufacturer data for a entity
//
// possible errors:
//   - ErrNoCompatibleEntity if entity is not compatible
//   - and others
func ManufacturerData(localEntity spineapi.EntityLocalInterface, entity spineapi.EntityRemoteInterface) (api.ManufacturerData, error) {
	if entity == nil {
		return api.ManufacturerData{}, api.ErrNoCompatibleEntity
	}

	deviceClassification, err := client.NewDeviceClassification(localEntity, entity)
	if err != nil {
		return api.ManufacturerData{}, err
	}

	data, err := deviceClassification.GetManufacturerDetails()
	if err != nil {
		return api.ManufacturerData{}, err
	}

	ret := api.ManufacturerData{
		DeviceName:                     Deref((*string)(data.DeviceName)),
		DeviceCode:                     Deref((*string)(data.DeviceCode)),
		SerialNumber:                   Deref((*string)(data.SerialNumber)),
		SoftwareRevision:               Deref((*string)(data.SoftwareRevision)),
		HardwareRevision:               Deref((*string)(data.HardwareRevision)),
		VendorName:                     Deref((*string)(data.VendorName)),
		VendorCode:                     Deref((*string)(data.VendorCode)),
		BrandName:                      Deref((*string)(data.BrandName)),
		PowerSource:                    Deref((*string)(data.PowerSource)),
		ManufacturerNodeIdentification: Deref((*string)(data.ManufacturerNodeIdentification)),
		ManufacturerLabel:              Deref((*string)(data.ManufacturerLabel)),
		ManufacturerDescription:        Deref((*string)(data.ManufacturerDescription)),
	}

	return ret, nil
}
