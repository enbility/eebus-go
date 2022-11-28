package spine

import (
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/ahmetb/go-linq/v3"
)

type BindingManager interface {
	AddBinding(localDevice *DeviceLocalImpl, remoteDevice *DeviceRemoteImpl, data model.BindingManagementRequestCallType) error
	RemoveBinding(data model.BindingManagementDeleteCallType, remoteDevice *DeviceRemoteImpl) error
	Bindings(remoteDevice *DeviceRemoteImpl) []*BindingEntry
	BindingsOnFeature(featureAddress model.FeatureAddressType) []*BindingEntry
}

type BindingEntry struct {
	id            uint64
	serverFeature FeatureLocal
	clientFeature *FeatureRemoteImpl
}

type BindingManagerImpl struct {
	bindingNum     uint64
	bindingEntries []*BindingEntry
	// TODO: add persistence
}

func NewBindingManager() BindingManager {
	c := &BindingManagerImpl{
		bindingNum: 0,
	}

	return c
}

// is sent from the client (remote device) to the server (local device)
func (c *BindingManagerImpl) AddBinding(localDevice *DeviceLocalImpl, remoteDevice *DeviceRemoteImpl, data model.BindingManagementRequestCallType) error {

	serverFeature := localDevice.FeatureByAddress(data.ServerAddress)
	if serverFeature == nil {
		return fmt.Errorf("server feature '%s' in local device '%s' not found", data.ServerAddress, *localDevice.Address())
	}
	if err := c.checkRoleAndType(serverFeature, model.RoleTypeServer, *data.ServerFeatureType); err != nil {
		return err
	}

	clientFeature := remoteDevice.FeatureByAddress(data.ClientAddress)
	if clientFeature == nil {
		return fmt.Errorf("client feature '%s' in remote device '%s' not found", data.ClientAddress, *remoteDevice.Address())
	}
	if err := c.checkRoleAndType(clientFeature, model.RoleTypeClient, *data.ServerFeatureType); err != nil {
		return err
	}

	bindingEntry := &BindingEntry{
		id:            c.bindingId(),
		serverFeature: serverFeature,
		clientFeature: clientFeature,
	}

	// TOV-TODO: check if binding already exists
	c.bindingEntries = append(c.bindingEntries, bindingEntry)

	payload := EventPayload{
		Ski:        remoteDevice.ski,
		EventType:  EventTypeBindingChange,
		ChangeType: ElementChangeAdd,
		Data:       data,
		Feature:    clientFeature,
	}
	Events.Publish(payload)

	// TOV-TODO: Send heartbeat to the feature which subscribed to DeviceDiagnostic

	return nil
}

func (c *BindingManagerImpl) RemoveBinding(data model.BindingManagementDeleteCallType, remoteDevice *DeviceRemoteImpl) error {
	// TODO: test this!!!

	var newBindingEntries []*BindingEntry

	// according to the spec 7.4.4
	// a. The absence of "bindingDelete. clientAddress. device" SHALL be treated as if it was
	//    present and set to the sender's "device" address part.
	// b. The absence of "bindingDelete. serverAddress. device" SHALL be treated as if it was
	//    present and set to the recipient's "device" address part.

	clientAddress := data.ClientAddress
	if data.ClientAddress.Device == nil {
		clientAddress.Device = remoteDevice.Address()
	}

	for _, item := range c.bindingEntries {
		if !reflect.DeepEqual(item.clientFeature.Address(), clientAddress) {
			newBindingEntries = append(newBindingEntries, item)
		}
	}

	if len(newBindingEntries) == len(c.bindingEntries) {
		return errors.New("could not find requested BindingId to be removed")
	}

	c.bindingEntries = newBindingEntries

	// TOV-TODO: stop heartbeat for remote device when it has no binding to DeviceDiagnostic anymore
	return nil
}

func (c *BindingManagerImpl) Bindings(remoteDevice *DeviceRemoteImpl) []*BindingEntry {
	var result []*BindingEntry

	linq.From(c.bindingEntries).WhereT(func(s *BindingEntry) bool {
		return s.clientFeature.Device().Ski() == remoteDevice.Ski()
	}).ToSlice(&result)

	return result
}

func (c *BindingManagerImpl) BindingsOnFeature(featureAddress model.FeatureAddressType) []*BindingEntry {
	var result []*BindingEntry

	linq.From(c.bindingEntries).WhereT(func(s *BindingEntry) bool {
		return reflect.DeepEqual(*s.serverFeature.Address(), featureAddress)
	}).ToSlice(&result)

	return result
}

func (c *BindingManagerImpl) bindingId() uint64 {
	i := atomic.AddUint64(&c.bindingNum, 1)
	return i
}

func (c *BindingManagerImpl) checkRoleAndType(feature Feature, role model.RoleType, featureType model.FeatureTypeType) error {
	if feature.Role() != model.RoleTypeSpecial && feature.Role() != role {
		return fmt.Errorf("found feature '%s' is not matching required role '%s'", feature.Type(), role)
	}

	if feature.Type() != featureType && feature.Type() != model.FeatureTypeTypeGeneric {
		return fmt.Errorf("found feature '%s' is not matching required type '%s'", feature.Type(), featureType)
	}

	return nil
}
