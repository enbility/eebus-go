package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/service"
	"golang.org/x/exp/jsonrpc2"

	shipapi "github.com/enbility/ship-go/api"
)

type Remote struct {
	rpc     *jsonrpc2.Server
	service *service.Service

	connections    []*jsonrpc2.Connection
	remoteServices []shipapi.RemoteService

	rpcServices map[string]rpcServiceFunc
}

func NewRemote(configuration *api.Configuration) (*Remote, error) {
	r := Remote{
		connections:    []*jsonrpc2.Connection{},
		remoteServices: []shipapi.RemoteService{},

		rpcServices: make(map[string]rpcServiceFunc),
	}
	r.service = service.NewService(configuration, &r)

	if err := r.service.Setup(); err != nil {
		return nil, err
	}
	r.service.SetLogging(&stdoutLogger{})

	err := r.registerStaticReceiverProxy("Service", r.service)
	if err != nil {
		return nil, err
	}

	err = r.registerStaticReceiverProxy("Remote", &r)
	if err != nil {
		return nil, err
	}

	err = r.registerStaticReceiverProxy("LocalDevice", r.service.LocalDevice())
	if err != nil {
		return nil, err
	}

	err = r.registerDynamicReceiverProxy("Call")
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *Remote) registerStaticReceiverProxy(name string, rcvr any) error {
	proxy, err := newStaticReceiverProxy(rcvr, name, true)
	if err != nil {
		return err
	}
	r.rpcServices[proxy.name] = proxy

	return nil
}

func (r *Remote) registerDynamicReceiverProxy(name string) error {
	r.rpcServices[strings.ToLower(name)] = dynamicReceiverProxy{}

	return nil
}

func (r Remote) RemoteServices() []shipapi.RemoteService {
	return r.remoteServices
}

func (r Remote) ConnectedDevices() []string {
	remoteDevices := r.service.LocalDevice().RemoteDevices()
	skiList := make([]string, len(remoteDevices))

	for i, dev := range remoteDevices {
		skiList[i] = dev.Ski()
	}

	return skiList
}

func (r Remote) LocalSKI() string {
	return r.service.LocalService().SKI()
}

func (r *Remote) Bind(context context.Context, conn *jsonrpc2.Connection) (jsonrpc2.ConnectionOptions, error) {
	connOpts := jsonrpc2.ConnectionOptions{
		Framer:    NewlineFramer{},
		Preempter: nil,
		Handler:   jsonrpc2.HandlerFunc(r.handleRPC),
	}

	r.connections = append(r.connections, conn)
	return connOpts, nil
}

func (r *Remote) Listen(context context.Context, network, address string) error {
	listener, err := jsonrpc2.NetListener(context, network, address, jsonrpc2.NetListenOptions{})
	if err != nil {
		return err
	}

	conn, err := jsonrpc2.Serve(context, listener, r)
	if err != nil {
		return err
	}
	r.rpc = conn

	r.service.Start()
	go func() {
		<-context.Done()
		r.service.Shutdown()
	}()

	return nil
}

func (r *Remote) handleRPC(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	if req.IsCall() {
		slash := strings.LastIndex(req.Method, "/")
		if slash < 0 {
			return nil, jsonrpc2.ErrMethodNotFound
		}
		serviceName := strings.ToLower(req.Method[:slash])
		methodName := req.Method[slash+1:]

		svc, found := r.rpcServices[serviceName]
		if !found {
			return nil, jsonrpc2.ErrMethodNotFound
		}

		output, err := svc.Call(r, methodName, req.Params)
		if err != nil {
			return nil, err
		}

		var resp interface{}
		numOut := len(output)
		switch numOut {
		case 0:
			resp = []interface{}{}
		default:
			resp = output
		}
		return resp, nil
	} else {
		// RPC Notification
		// TODO: implement
	}

	return nil, nil
}

// Implement api.ServiceReaderInterface
func (r Remote) RemoteSKIConnected(service api.ServiceInterface, ski string) {
	// necessary because RemoteSKIConnected is called before remote device actually exists
	go func() {
		params := make(map[string]interface{}, 1)
		params["ski"] = ski

		for {
			// wait until RemoteDevice available for SKI
			device := service.LocalDevice().RemoteDeviceForSki(ski)
			if device != nil && device.Address() != nil {
				params["device"] = *device.Address()
				break
			}
			time.Sleep(1 * time.Second)
		}

		for _, conn := range r.connections {
			_ = conn.Notify(context.Background(), "remote/RemoteSKIConnected", params)
		}
	}()
}

func (r Remote) RemoteSKIDisconnected(service api.ServiceInterface, ski string) {
	params := make(map[string]interface{}, 1)
	params["ski"] = ski
	for _, conn := range r.connections {
		_ = conn.Notify(context.Background(), "remote/RemoteSKIDisconnected", params)
	}
}

func (r *Remote) VisibleRemoteServicesUpdated(service api.ServiceInterface, entries []shipapi.RemoteService) {
	r.remoteServices = entries

	for _, conn := range r.connections {
		_ = conn.Notify(context.Background(), "remote/VisibleRemoteServicesUpdated", entries)
	}
}

func (r Remote) ServiceShipIDUpdate(ski string, shipID string) {
	params := make(map[string]interface{}, 2)
	params["ski"] = ski
	params["shipID"] = shipID

	for _, conn := range r.connections {
		_ = conn.Notify(context.Background(), "remote/ServiceShipIDUpdate", params)
	}
}

func (r Remote) ServicePairingDetailUpdate(ski string, detail *shipapi.ConnectionStateDetail) {
}

// Logging interface

type stdoutLogger struct{}

func (l *stdoutLogger) Trace(args ...interface{}) {
	// l.print("TRACE", args...)
}

func (l *stdoutLogger) Tracef(format string, args ...interface{}) {
	// l.printFormat("TRACE", format, args...)
}

func (l *stdoutLogger) Debug(args ...interface{}) {
	// l.print("DEBUG", args...)
}

func (l *stdoutLogger) Debugf(format string, args ...interface{}) {
	// l.printFormat("DEBUG", format, args...)
}

func (l *stdoutLogger) Info(args ...interface{}) {
	l.print("INFO ", args...)
}

func (l *stdoutLogger) Infof(format string, args ...interface{}) {
	l.printFormat("INFO ", format, args...)
}

func (l *stdoutLogger) Error(args ...interface{}) {
	l.print("ERROR", args...)
}

func (l *stdoutLogger) Errorf(format string, args ...interface{}) {
	l.printFormat("ERROR", format, args...)
}

func (l *stdoutLogger) currentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (l *stdoutLogger) print(msgType string, args ...interface{}) {
	value := fmt.Sprintln(args...)
	fmt.Printf("%s %s %s", l.currentTimestamp(), msgType, value)
}

func (l *stdoutLogger) printFormat(msgType, format string, args ...interface{}) {
	value := fmt.Sprintf(format, args...)
	fmt.Println(l.currentTimestamp(), msgType, value)
}
