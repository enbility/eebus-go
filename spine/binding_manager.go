package spine

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/ahmetb/go-linq/v3"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
)

type BindingEntry struct {
	id            uint64
	serverFeature FeatureLocal
	clientFeature FeatureRemote
}

type BindingManagerImpl struct {
	localDevice *DeviceLocalImpl

	bindingNum     uint64
	bindingEntries []*BindingEntry

	mux sync.Mutex
	// TODO: add persistence
}

func NewBindingManager(localDevice *DeviceLocalImpl) BindingManager {
	c := &BindingManagerImpl{
		bindingNum:  0,
		localDevice: localDevice,
	}

	return c
}

// is sent from the client (remote device) to the server (local device)
func (c *BindingManagerImpl) AddBinding(remoteDevice *DeviceRemoteImpl, data model.BindingManagementRequestCallType) error {

	serverFeature := c.localDevice.FeatureByAddress(data.ServerAddress)
	if serverFeature == nil {
		return fmt.Errorf("server feature '%s' in local device '%s' not found", data.ServerAddress, *c.localDevice.Address())
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

	c.mux.Lock()
	defer c.mux.Unlock()

	for _, item := range c.bindingEntries {
		if reflect.DeepEqual(item.serverFeature, serverFeature) && reflect.DeepEqual(item.clientFeature, clientFeature) {
			return fmt.Errorf("requested binding is already present")
		}
	}

	c.bindingEntries = append(c.bindingEntries, bindingEntry)

	payload := EventPayload{
		Ski:        remoteDevice.ski,
		EventType:  EventTypeBindingChange,
		ChangeType: ElementChangeAdd,
		Data:       data,
		Feature:    clientFeature,
	}
	Events.Publish(payload)

	return nil
}

func (c *BindingManagerImpl) RemoveBinding(data model.BindingManagementDeleteCallType, remoteDevice *DeviceRemoteImpl) error {
	var newBindingEntries []*BindingEntry

	// according to the spec 7.4.4
	// a. The absence of "bindingDelete. clientAddress. device" SHALL be treated as if it was
	//    present and set to the sender's "device" address part.
	// b. The absence of "bindingDelete. serverAddress. device" SHALL be treated as if it was
	//    present and set to the recipient's "device" address part.

	var clientAddress model.FeatureAddressType
	util.DeepCopy(data.ClientAddress, &clientAddress)
	if data.ClientAddress.Device == nil {
		clientAddress.Device = remoteDevice.Address()
	}

	clientFeature := remoteDevice.FeatureByAddress(data.ClientAddress)
	if clientFeature == nil {
		return fmt.Errorf("client feature '%s' in remote device '%s' not found", data.ClientAddress, *remoteDevice.Address())
	}

	serverFeature := c.localDevice.FeatureByAddress(data.ServerAddress)
	if serverFeature == nil {
		return fmt.Errorf("server feature '%s' in local device '%s' not found", data.ServerAddress, *c.localDevice.Address())
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	for _, item := range c.bindingEntries {
		itemAddress := item.clientFeature.Address()

		if !reflect.DeepEqual(*itemAddress, clientAddress) &&
			!reflect.DeepEqual(item.serverFeature, serverFeature) {
			newBindingEntries = append(newBindingEntries, item)
		}
	}

	if len(newBindingEntries) == len(c.bindingEntries) {
		return errors.New("could not find requested BindingId to be removed")
	}

	c.bindingEntries = newBindingEntries

	payload := EventPayload{
		Ski:        remoteDevice.ski,
		EventType:  EventTypeBindingChange,
		ChangeType: ElementChangeRemove,
		Data:       data,
		Device:     remoteDevice,
		Feature:    clientFeature,
	}
	Events.Publish(payload)

	return nil
}

// Remove all existing bindings for a given remote device
func (c *BindingManagerImpl) RemoveBindingsForDevice(remoteDevice *DeviceRemoteImpl) {
	if remoteDevice == nil {
		return
	}

	for _, entity := range remoteDevice.Entities() {
		c.RemoveBindingsForEntity(entity)
	}
}

// Remove all existing bindings for a given remote device entity
func (c *BindingManagerImpl) RemoveBindingsForEntity(remoteEntity *EntityRemoteImpl) {
	if remoteEntity == nil {
		return
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	var newBindingEntries []*BindingEntry
	for _, item := range c.bindingEntries {
		if !reflect.DeepEqual(item.clientFeature.Address().Entity, remoteEntity.Address().Entity) {
			newBindingEntries = append(newBindingEntries, item)
			continue
		}

		clientFeature := remoteEntity.Feature(item.clientFeature.Address().Feature)
		payload := EventPayload{
			Ski:        remoteEntity.Device().ski,
			EventType:  EventTypeBindingChange,
			ChangeType: ElementChangeRemove,
			Entity:     remoteEntity,
			Feature:    clientFeature,
		}
		Events.Publish(payload)
	}

	c.bindingEntries = newBindingEntries
}

func (c *BindingManagerImpl) Bindings(remoteDevice *DeviceRemoteImpl) []*BindingEntry {
	var result []*BindingEntry

	c.mux.Lock()
	defer c.mux.Unlock()

	linq.From(c.bindingEntries).WhereT(func(s *BindingEntry) bool {
		return s.clientFeature.Device().Ski() == remoteDevice.Ski()
	}).ToSlice(&result)

	return result
}

func (c *BindingManagerImpl) BindingsOnFeature(featureAddress model.FeatureAddressType) []*BindingEntry {
	var result []*BindingEntry

	c.mux.Lock()
	defer c.mux.Unlock()

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
		return fmt.Errorf("found feature %s is not matching required role %s", feature.Type(), role)
	}

	if feature.Type() != featureType && feature.Type() != model.FeatureTypeTypeGeneric {
		return fmt.Errorf("found feature %s is not matching required type %s", feature.Type(), featureType)
	}

	return nil
}
