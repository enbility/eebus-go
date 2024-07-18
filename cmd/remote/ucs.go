package main

import (
	"context"
	"fmt"

	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type UseCaseId string
type UseCaseTypeType string

const (
	UseCaseTypeLPC UseCaseTypeType = "LPC"
)

type UseCaseBuilder func(spineapi.EntityLocalInterface, api.EntityEventCallback) api.UseCaseInterface

func (r *Remote) RegisterUseCase(entityType model.EntityTypeType, usecaseId string, builder UseCaseBuilder) error {
	// entityType/uc
	var identifier UseCaseId = UseCaseId(fmt.Sprintf("%s/%s", entityType, usecaseId))

	localInterface := r.service.LocalDevice().EntityForType(entityType)
	uc := builder(localInterface, func(
		ski string,
		device spineapi.DeviceRemoteInterface,
		entity spineapi.EntityRemoteInterface,
		event api.EventType,
	) {
		r.PropagateEvent(identifier, ski, device, entity, event)
	})
	r.service.AddUseCase(uc)

	return r.registerStaticReceiverProxy(usecaseId, uc)
}

func (r *Remote) PropagateEvent(
	id UseCaseId,
	ski string,
	device spineapi.DeviceRemoteInterface,
	entity spineapi.EntityRemoteInterface,
	event api.EventType,
) {
	params := make(map[string]interface{}, 2)
	params["ski"] = ski
	params["device"] = device.Address()
	params["entity"] = entity.Address()
	for _, conn := range r.connections {
		_ = conn.Notify(context.Background(), string(event), params)
	}
}
