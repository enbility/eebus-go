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

type SubscriptionManager interface {
	AddSubscription(remoteDevice *DeviceRemoteImpl, data model.SubscriptionManagementRequestCallType) error
	RemoveSubscription(data model.SubscriptionManagementDeleteCallType, remoteDevice *DeviceRemoteImpl) error
	RemoveSubscriptionsForDevice(remoteDevice *DeviceRemoteImpl)
	RemoveSubscriptionsForEntity(remoteEntity *EntityRemoteImpl)
	Subscriptions(remoteDevice *DeviceRemoteImpl) []*SubscriptionEntry
	SubscriptionsOnFeature(featureAddress model.FeatureAddressType) []*SubscriptionEntry
}

type SubscriptionEntry struct {
	id            uint64
	serverFeature FeatureLocal
	clientFeature *FeatureRemoteImpl
}

type SubscriptionManagerImpl struct {
	localDevice *DeviceLocalImpl

	subscriptionNum     uint64
	subscriptionEntries []*SubscriptionEntry

	mux sync.Mutex
	// TODO: add persistence
}

func NewSubscriptionManager(localDevice *DeviceLocalImpl) SubscriptionManager {
	c := &SubscriptionManagerImpl{
		subscriptionNum: 0,
		localDevice:     localDevice,
	}

	return c
}

// is sent from the client (remote device) to the server (local device)
func (c *SubscriptionManagerImpl) AddSubscription(remoteDevice *DeviceRemoteImpl, data model.SubscriptionManagementRequestCallType) error {

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

	subscriptionEntry := &SubscriptionEntry{
		id:            c.subscriptionId(),
		serverFeature: serverFeature,
		clientFeature: clientFeature,
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	for _, item := range c.subscriptionEntries {
		if reflect.DeepEqual(item.serverFeature, serverFeature) && reflect.DeepEqual(item.clientFeature, clientFeature) {
			return fmt.Errorf("requested subscription is already present")
		}
	}

	c.subscriptionEntries = append(c.subscriptionEntries, subscriptionEntry)

	payload := EventPayload{
		Ski:        remoteDevice.ski,
		EventType:  EventTypeSubscriptionChange,
		ChangeType: ElementChangeAdd,
		Data:       data,
		Feature:    clientFeature,
	}
	Events.Publish(payload)

	return nil
}

// Remove a specific subscription that is provided by a delete message from a remote device
func (c *SubscriptionManagerImpl) RemoveSubscription(data model.SubscriptionManagementDeleteCallType, remoteDevice *DeviceRemoteImpl) error {
	var newSubscriptionEntries []*SubscriptionEntry

	// according to the spec 7.4.4
	// a. The absence of "subscriptionDelete. clientAddress. device" SHALL be treated as if it was
	//    present and set to the sender's "device" address part.
	// b. The absence of "subscriptionDelete. serverAddress. device" SHALL be treated as if it was
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

	for _, item := range c.subscriptionEntries {
		itemAddress := item.clientFeature.Address()

		if !reflect.DeepEqual(*itemAddress, clientAddress) &&
			!reflect.DeepEqual(item.serverFeature, serverFeature) {
			newSubscriptionEntries = append(newSubscriptionEntries, item)
		}
	}

	if len(newSubscriptionEntries) == len(c.subscriptionEntries) {
		return errors.New("could not find requested SubscriptionId to be removed")
	}

	c.subscriptionEntries = newSubscriptionEntries

	payload := EventPayload{
		Ski:        remoteDevice.ski,
		EventType:  EventTypeSubscriptionChange,
		ChangeType: ElementChangeRemove,
		Data:       data,
		Device:     remoteDevice,
		Feature:    clientFeature,
	}
	Events.Publish(payload)

	return nil
}

// Remove all existing subscriptions for a given remote device
func (c *SubscriptionManagerImpl) RemoveSubscriptionsForDevice(remoteDevice *DeviceRemoteImpl) {
	if remoteDevice == nil {
		return
	}

	for _, entity := range remoteDevice.Entities() {
		c.RemoveSubscriptionsForEntity(entity)
	}
}

// Remove all existing subscriptions for a given remote device entity
func (c *SubscriptionManagerImpl) RemoveSubscriptionsForEntity(remoteEntity *EntityRemoteImpl) {
	if remoteEntity == nil {
		return
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	var newSubscriptionEntries []*SubscriptionEntry
	for _, item := range c.subscriptionEntries {
		if !reflect.DeepEqual(item.clientFeature.Address().Entity, remoteEntity.Address().Entity) {
			newSubscriptionEntries = append(newSubscriptionEntries, item)
			continue
		}

		clientFeature := remoteEntity.Feature(item.clientFeature.address.Feature)
		payload := EventPayload{
			Ski:        remoteEntity.Device().ski,
			EventType:  EventTypeSubscriptionChange,
			ChangeType: ElementChangeRemove,
			Entity:     remoteEntity,
			Feature:    clientFeature,
		}
		Events.Publish(payload)
	}

	c.subscriptionEntries = newSubscriptionEntries
}

func (c *SubscriptionManagerImpl) Subscriptions(remoteDevice *DeviceRemoteImpl) []*SubscriptionEntry {
	var result []*SubscriptionEntry

	c.mux.Lock()
	defer c.mux.Unlock()

	linq.From(c.subscriptionEntries).WhereT(func(s *SubscriptionEntry) bool {
		return s.clientFeature.Device().Ski() == remoteDevice.Ski()
	}).ToSlice(&result)

	return result
}

func (c *SubscriptionManagerImpl) SubscriptionsOnFeature(featureAddress model.FeatureAddressType) []*SubscriptionEntry {
	var result []*SubscriptionEntry

	c.mux.Lock()
	defer c.mux.Unlock()

	linq.From(c.subscriptionEntries).WhereT(func(s *SubscriptionEntry) bool {
		return reflect.DeepEqual(*s.serverFeature.Address(), featureAddress)
	}).ToSlice(&result)

	return result
}

func (c *SubscriptionManagerImpl) subscriptionId() uint64 {
	i := atomic.AddUint64(&c.subscriptionNum, 1)
	return i
}

func (c *SubscriptionManagerImpl) checkRoleAndType(feature Feature, role model.RoleType, featureType model.FeatureTypeType) error {
	if feature.Role() != model.RoleTypeSpecial && feature.Role() != role {
		return fmt.Errorf("found feature %s is not matching required role %s", feature.Type(), role)
	}

	if feature.Type() != featureType && feature.Type() != model.FeatureTypeTypeGeneric {
		return fmt.Errorf("found feature %s is not matching required type %s", feature.Type(), featureType)
	}

	return nil
}
