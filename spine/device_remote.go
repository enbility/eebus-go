package spine

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type DeviceRemoteImpl struct {
	*DeviceImpl
	ski string

	entities []*EntityRemoteImpl
	sender   Sender

	localDevice *DeviceLocalImpl

	// The read channel for incoming messages
	readChannel <-chan []byte

	// Handles closing of the connection
	closeChannel chan bool

	// Heartbeat Sender
	heartbeatSender *HeartbeatSender
}

func NewDeviceRemoteImpl(localDevice *DeviceLocalImpl, ski string, readC <-chan []byte, writeC chan<- []byte) *DeviceRemoteImpl {
	sender := NewSender(writeC)
	res := DeviceRemoteImpl{
		DeviceImpl:      NewDeviceImpl(nil, nil),
		ski:             ski,
		localDevice:     localDevice,
		readChannel:     readC,
		closeChannel:    make(chan bool),
		sender:          sender,
		heartbeatSender: NewHeartbeatSender(sender),
	}
	go res.readPump()

	return &res
}

// return the device SKI
func (d *DeviceRemoteImpl) Ski() string {
	return d.ski
}

// Needs to be called by the CEM implementation once a subscription for the local DeviceDiagnosis server feature is received
func (d *DeviceRemoteImpl) StartHeartbeatSend(senderAddr, destinationAddr *model.FeatureAddressType) {
	d.heartbeatSender.StartHeartbeatSend(senderAddr, destinationAddr)
}

// this connection is closed
func (d *DeviceRemoteImpl) CloseConnection() {
	d.heartbeatSender.StopHeartbeat()
	d.closeChannel <- true
}

// read all incoming spine messages from the associated SHIP connection
func (d *DeviceRemoteImpl) readPump() {
	for {
		select {
		case <-d.closeChannel:
			return
		case data, ok := <-d.readChannel:
			if !ok {
				return
			}

			datagram := model.Datagram{}
			if err := json.Unmarshal([]byte(data), &datagram); err != nil {
				fmt.Println(err)
				continue
			}
			_ = d.localDevice.ProcessCmd(datagram.Datagram, d)
		}
	}
}

func (d *DeviceRemoteImpl) Sender() Sender {
	return d.sender
}

func (d *DeviceRemoteImpl) Entity(id []model.AddressEntityType) *EntityRemoteImpl {
	for _, e := range d.entities {
		if reflect.DeepEqual(id, e.Address().Entity) {
			return e
		}
	}
	return nil
}

func (d *DeviceRemoteImpl) FeatureByAddress(address *model.FeatureAddressType) *FeatureRemoteImpl {
	entity := d.Entity(address.Entity)
	if entity != nil {
		return entity.Feature(address.Feature)
	}
	return nil
}

func (d *DeviceRemoteImpl) RemoveByAddress(addr []model.AddressEntityType) *EntityRemoteImpl {
	entityForRemoval := d.Entity(addr)
	if entityForRemoval == nil {
		return nil
	}

	var newEntities []*EntityRemoteImpl
	for _, item := range d.entities {
		if !reflect.DeepEqual(item, entityForRemoval) {
			newEntities = append(newEntities, item)
		}
	}
	d.entities = newEntities

	// if removedEntity != nil {
	// 	events.EntityChanged.Publish(events.EntityChangedPayload{Entity: removedEntity, Mode: events.EntityRemoved})
	// }

	return entityForRemoval
}

func (r *DeviceRemoteImpl) FeatureByTypeAndRole(featureType model.FeatureTypeType, role model.RoleType) *FeatureRemoteImpl {
	if len(r.entities) < 1 {
		return nil
	}

	for _, entity := range r.entities {
		for _, feature := range entity.Features() {
			if feature.Type() == featureType && feature.Role() == role {
				return feature
			}
		}
	}

	return nil
}

func (d *DeviceRemoteImpl) UpdateDevice(description *model.NetworkManagementDeviceDescriptionDataType) {
	if description != nil {
		if description.DeviceAddress != nil && description.DeviceAddress.Device != nil {
			d.address = description.DeviceAddress.Device
		}
		if description.DeviceType != nil {
			d.dType = description.DeviceType
		}
	}
}

func (d *DeviceRemoteImpl) AddEntityAndFeatures(initialData bool, data *model.NodeManagementDetailedDiscoveryDataType) error {
	for _, ei := range data.EntityInformation {
		if err := d.CheckEntityInformation(initialData, ei); err != nil {
			return err
		}

		entityAddress := ei.Description.EntityAddress.Entity

		entity := d.Entity(entityAddress)
		if entity == nil {
			entity = d.addNewEntity(*ei.Description.EntityType, entityAddress)
		}

		entity.SetDescription(ei.Description.Description)

		for _, fi := range data.FeatureInformation {
			if reflect.DeepEqual(fi.Description.FeatureAddress.Entity, entityAddress) {
				if f := unmarshalFeature(entity, fi); f != nil {
					entity.AddFeature(f)
				}
			}
		}

		// TOV-TODO: check this approach
		// if err := f.announceFeatureDiscovery(entity); err != nil {
		// 	return err
		// }
	}

	return nil
}

// check if the provided entity information is correct
// provide initialData to check if the entity is new and not an update
func (d *DeviceRemoteImpl) CheckEntityInformation(initialData bool, entity model.NodeManagementDetailedDiscoveryEntityInformationType) error {
	description := entity.Description
	if description == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description")
	}

	if description.EntityAddress == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description.EntityAddress")
	}

	if description.EntityAddress.Entity == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description.EntityAddress.Entity")
	}

	// Consider on initial NodeManagement Detailed Discovery, the device being empty as it is not yet known
	if initialData {
		return nil
	}

	address := d.Address()
	if description.EntityAddress.Device != nil && address != nil && *description.EntityAddress.Device != *address {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: device address mismatch")
	}

	return nil
}

func (d *DeviceRemoteImpl) addNewEntity(eType model.EntityTypeType, address []model.AddressEntityType) *EntityRemoteImpl {
	newEntity := NewEntityRemoteImpl(d, eType, address)
	return d.addEntity(newEntity)
}

func (d *DeviceRemoteImpl) addEntity(entity *EntityRemoteImpl) *EntityRemoteImpl {
	d.entities = append(d.entities, entity)

	return entity
}

func unmarshalFeature(entity *EntityRemoteImpl,
	featureData model.NodeManagementDetailedDiscoveryFeatureInformationType,
) *FeatureRemoteImpl {
	var result *FeatureRemoteImpl

	if fid := featureData.Description; fid != nil {

		result = NewFeatureRemoteImpl(uint(*fid.FeatureAddress.Feature), entity, *fid.FeatureType, *fid.Role)

		result.SetDescription(fid.Description)
		result.SetMaxResponseDelay(fid.MaxResponseDelay)
		result.SetOperations(fid.SupportedFunction)
	}

	return result
}
