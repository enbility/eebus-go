package main

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"go/token"
	"os"
	"reflect"

	"github.com/enbility/spine-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return token.IsExported(t.Name()) || t.PkgPath() == ""
}

// json.Marshall won't marshall error types, marshall as string
func errorAsJson(v reflect.Value) interface{} {
	if v.IsNil() {
		// passthrough nil as nil, otherwise nil.(error) will panic
		return nil
	} else {
		return v.Interface().(error).Error()
	}
}

func transformReturnValues(values []reflect.Value) []interface{} {
	result := make([]interface{}, len(values))

	for i, e := range values {
		valueType := e.Type()

		switch valueType {
		case reflect.TypeFor[spineapi.DeviceRemoteInterface]():
			result[i] = e.Interface().(spineapi.DeviceRemoteInterface).Address()
		case reflect.TypeFor[[]spineapi.DeviceRemoteInterface]():
			rawValues := e.Interface().([]api.DeviceRemoteInterface)
			transformedValues := make([]model.AddressDeviceType, len(rawValues))

			for j, r := range rawValues {
				transformedValues[j] = *r.Address()
			}
			result[i] = transformedValues
		case reflect.TypeFor[spineapi.EntityRemoteInterface]():
			result[i] = e.Interface().(spineapi.EntityRemoteInterface).Address()
		case reflect.TypeFor[[]spineapi.EntityRemoteInterface]():
			rawValues := e.Interface().([]api.EntityRemoteInterface)
			transformedValues := make([]model.EntityAddressType, len(rawValues))

			for j, r := range rawValues {
				transformedValues[j] = *r.Address()
			}
			result[i] = transformedValues
		default:
			switch {
			case valueType.Implements(reflect.TypeFor[error]()):
				result[i] = errorAsJson(e)
			default:
				result[i] = e.Interface()
			}
		}
	}

	return result
}

func WriteKey(cert tls.Certificate, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch v := cert.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		bytes, err := x509.MarshalECPrivateKey(v)
		if err != nil {
			return err
		}

		err = pem.Encode(file, &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: bytes,
		})
	default:
		return fmt.Errorf("Unable to serialize private key of type %T", v)
	}

	return nil
}

func WriteCertificate(cert tls.Certificate, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, leaf := range cert.Certificate {
		err = pem.Encode(file, &pem.Block{Type: "CERTIFICATE", Bytes: leaf})
		if err != nil {
			return err
		}
	}

	return nil
}
