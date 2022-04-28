package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

var FeatureOperationsMap = map[model.FeatureTypeType]map[model.FunctionType]*Operations{
	model.FeatureTypeTypeDeviceClassification: {
		model.FunctionTypeDeviceClassificationManufacturerData: NewOperations(true, false),
	},
	model.FeatureTypeTypeNodeManagement: {
		model.FunctionTypeNodeManagementDetailedDiscoveryData:   NewOperations(true, false),
		model.FunctionTypeNodeManagementUseCaseData:             NewOperations(true, false),
		model.FunctionTypeNodeManagementSubscriptionData:        NewOperations(true, false),
		model.FunctionTypeNodeManagementSubscriptionRequestCall: NewOperations(false, false),
		model.FunctionTypeNodeManagementSubscriptionDeleteCall:  NewOperations(false, false),
		model.FunctionTypeNodeManagementBindingData:             NewOperations(true, false),
		model.FunctionTypeNodeManagementBindingRequestCall:      NewOperations(false, false),
		model.FunctionTypeNodeManagementBindingDeleteCall:       NewOperations(false, false),
	},
	model.FeatureTypeTypeDeviceDiagnosis: {
		model.FunctionTypeDeviceDiagnosisStateData: NewOperations(true, false),
	},
	// add more features here
}
