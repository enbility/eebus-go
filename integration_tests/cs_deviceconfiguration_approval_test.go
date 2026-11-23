package integrationtests

import (
	"slices"
	"sync"
	"testing"
	"time"

	eebusapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	cslpc "github.com/enbility/eebus-go/usecases/cs/lpc"
	cslpp "github.com/enbility/eebus-go/usecases/cs/lpp"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type deviceConfigurationApprovalEvents struct {
	mux    sync.Mutex
	events []eebusapi.EventType
}

func (e *deviceConfigurationApprovalEvents) add(_ string, _ spineapi.DeviceRemoteInterface, _ spineapi.EntityRemoteInterface, event eebusapi.EventType) {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.events = append(e.events, event)
}

func (e *deviceConfigurationApprovalEvents) contains(event eebusapi.EventType) bool {
	e.mux.Lock()
	defer e.mux.Unlock()

	return slices.Contains(e.events, event)
}

type csDeviceConfigurationApprovalSetup struct {
	remoteDevice      spineapi.DeviceRemoteInterface
	remoteEntity      spineapi.EntityRemoteInterface
	remoteFeature     spineapi.FeatureRemoteInterface
	localFeature      spineapi.FeatureLocalInterface
	deviceConfig      *server.DeviceConfiguration
	failsafeDuration  model.DeviceConfigurationKeyValueDescriptionDataType
	lpc               *cslpc.LPC
	lpp               *cslpp.LPP
	events            *deviceConfigurationApprovalEvents
	writeMessageCount model.MsgCounterType
}

func TestCSDeviceConfigurationApprovalRequiresLPCAndLPPForFailsafeDurationMinimum(t *testing.T) {
	t.Run("notifies both use cases and keeps write pending until both approve", func(t *testing.T) {
		setup := newCSDeviceConfigurationApprovalSetup(t)
		newDuration := 4 * time.Hour

		setup.writeFailsafeDurationMinimum(t, newDuration)
		setup.requireBothUsecasesPendingFailsafeDurationMinimum(t)

		setup.lpc.ApproveOrDenyDeviceConfiguration(setup.writeMessageCount, true, "")

		setup.requireFailsafeDurationMinimumValue(t, 0)
	})

	t.Run("does not write when one use case denies", func(t *testing.T) {
		setup := newCSDeviceConfigurationApprovalSetup(t)
		newDuration := 6 * time.Hour

		setup.writeFailsafeDurationMinimum(t, newDuration)
		setup.requireBothUsecasesPendingFailsafeDurationMinimum(t)

		setup.lpc.ApproveOrDenyDeviceConfiguration(setup.writeMessageCount, true, "")
		setup.lpp.ApproveOrDenyDeviceConfiguration(setup.writeMessageCount, false, "denied")

		setup.requireFailsafeDurationMinimumValue(t, 0)
	})

	t.Run("writes only after both use cases approve", func(t *testing.T) {
		setup := newCSDeviceConfigurationApprovalSetup(t)
		newDuration := 8 * time.Hour

		setup.writeFailsafeDurationMinimum(t, newDuration)
		setup.requireBothUsecasesPendingFailsafeDurationMinimum(t)

		setup.lpc.ApproveOrDenyDeviceConfiguration(setup.writeMessageCount, true, "")
		setup.requireFailsafeDurationMinimumValue(t, 0)

		setup.lpp.ApproveOrDenyDeviceConfiguration(setup.writeMessageCount, true, "")

		require.Eventually(t, func() bool {
			return setup.failsafeDurationMinimumValue(t) == newDuration
		}, time.Second, 10*time.Millisecond)
	})
}

func newCSDeviceConfigurationApprovalSetup(t *testing.T) *csDeviceConfigurationApprovalSetup {
	t.Helper()

	localDevice := spine.NewDeviceLocal(
		"TestBrandName",
		"TestDeviceModel",
		"TestSerialNumber",
		"TestDeviceCode",
		"TestDeviceAddress",
		model.DeviceTypeTypeEnergyManagementSystem,
		model.NetworkManagementFeatureSetTypeSmart,
	)
	localEntity := spine.NewEntityLocal(localDevice, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}), time.Second*60)
	localDevice.AddEntity(localEntity)

	events := &deviceConfigurationApprovalEvents{}
	lpc := cslpc.NewLPC(localEntity, events.add)
	lpc.AddFeatures()
	lpp := cslpp.NewLPP(localEntity, events.add)
	lpp.AddFeatures()

	localFeature := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	require.NotNil(t, localFeature)

	deviceConfig, err := server.NewDeviceConfiguration(localEntity)
	require.NoError(t, err)

	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
	}
	descriptions, err := deviceConfig.GetKeyValueDescriptionsForFilter(filter)
	require.NoError(t, err)
	require.Len(t, descriptions, 1)
	require.NotNil(t, descriptions[0].KeyId)

	writeHandler := &WriteMessageHandler{}
	_ = localDevice.SetupRemoteDevice("TestRemoteSki", writeHandler)
	remoteDevice := localDevice.RemoteDeviceForSki("TestRemoteSki")
	require.NotNil(t, remoteDevice)

	remoteEntity := spine.NewEntityRemote(remoteDevice, model.EntityTypeTypeGridGuard, spine.NewAddressEntityType([]uint{1}))
	remoteDevice.AddEntity(remoteEntity)
	remoteFeature := spine.NewFeatureRemote(1, remoteEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient)
	remoteEntity.AddFeature(remoteFeature)

	return &csDeviceConfigurationApprovalSetup{
		remoteDevice:      remoteDevice,
		remoteEntity:      remoteEntity,
		remoteFeature:     remoteFeature,
		localFeature:      localFeature,
		deviceConfig:      deviceConfig,
		failsafeDuration:  descriptions[0],
		lpc:               lpc,
		lpp:               lpp,
		events:            events,
		writeMessageCount: model.MsgCounterType(1),
	}
}

func (s *csDeviceConfigurationApprovalSetup) writeFailsafeDurationMinimum(t *testing.T, duration time.Duration) {
	t.Helper()

	msg := &spineapi.Message{
		RequestHeader: &model.HeaderType{
			AddressSource: &model.FeatureAddressType{
				Device:  util.Ptr(model.AddressDeviceType("remote")),
				Entity:  []model.AddressEntityType{1},
				Feature: util.Ptr(model.AddressFeatureType(1)),
			},
			AddressDestination: s.localFeature.Address(),
			MsgCounter:         util.Ptr(s.writeMessageCount),
			AckRequest:         util.Ptr(true),
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteFeature,
		DeviceRemote:  s.remoteDevice,
		EntityRemote:  s.remoteEntity,
		Cmd: model.CmdType{
			DeviceConfigurationKeyValueListData: &model.DeviceConfigurationKeyValueListDataType{
				DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
					{
						KeyId: s.failsafeDuration.KeyId,
						Value: &model.DeviceConfigurationKeyValueValueType{
							Duration: model.NewDurationType(duration),
						},
					},
				},
			},
		},
	}

	err := s.localFeature.HandleMessage(msg)
	require.Nil(t, err)
}

func (s *csDeviceConfigurationApprovalSetup) requireBothUsecasesPendingFailsafeDurationMinimum(t *testing.T) {
	t.Helper()

	require.Eventually(t, func() bool {
		return s.events.contains(cslpc.ConfigurationWriteApprovalRequired) &&
			s.events.contains(cslpp.ConfigurationWriteApprovalRequired)
	}, time.Second, 10*time.Millisecond)

	require.Eventually(t, func() bool {
		return pendingDeviceConfigurationHasKeyName(
			s.lpc.PendingDeviceConfigurations(),
			s.writeMessageCount,
			model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum,
		) && pendingDeviceConfigurationHasKeyName(
			s.lpp.PendingDeviceConfigurations(),
			s.writeMessageCount,
			model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum,
		)
	}, time.Second, 10*time.Millisecond)
}

func (s *csDeviceConfigurationApprovalSetup) failsafeDurationMinimumValue(t *testing.T) time.Duration {
	t.Helper()

	lpcDuration, lpcChangeable, err := s.lpc.FailsafeDurationMinimum()
	require.NoError(t, err)
	lppDuration, lppChangeable, err := s.lpp.FailsafeDurationMinimum()
	require.NoError(t, err)

	assert.Equal(t, lpcDuration, lppDuration)
	assert.Equal(t, lpcChangeable, lppChangeable)

	return lpcDuration
}

func (s *csDeviceConfigurationApprovalSetup) requireFailsafeDurationMinimumValue(t *testing.T, expected time.Duration) {
	t.Helper()

	assert.Equal(t, expected, s.failsafeDurationMinimumValue(t))
}

func pendingDeviceConfigurationHasKeyName(
	pending map[model.MsgCounterType][]ucapi.PendingDeviceConfiguration,
	msgCounter model.MsgCounterType,
	keyName model.DeviceConfigurationKeyNameType,
) bool {
	configs, ok := pending[msgCounter]
	if !ok {
		return false
	}

	for _, config := range configs {
		if config.KeyName == keyName {
			return true
		}
	}

	return false
}
