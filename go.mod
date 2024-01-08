module github.com/enbility/eebus-go

go 1.18

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-cmp v0.6.0
	github.com/miekg/dns v1.1.57 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rickb777/plural v1.4.1 // indirect
	github.com/stretchr/objx v0.5.1 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/tools v0.16.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/DerAndereAndi/zeroconf/v2 v2.0.0-20231028092313-1ae0ab54a2df
	github.com/ahmetb/go-linq/v3 v3.2.0
	github.com/godbus/dbus/v5 v5.1.0
	github.com/gorilla/websocket v1.5.1
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/holoplot/go-avahi v0.0.0-20231130181535-fd8f1d0307ed
	github.com/rickb777/date v1.20.5
	github.com/stretchr/testify v1.8.4
	gitlab.com/c0b/go-ordered-json v0.0.0-20201030195603-febf46534d5a
	go.uber.org/mock v0.4.0
)

retract (
	v0.2.2 // Contains retractions only.
	v0.2.1 // Published accidentally.
)
