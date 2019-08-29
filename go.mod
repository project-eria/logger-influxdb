module github.com/project-eria/logger-influxdb

go 1.12

require (
	github.com/influxdata/influxdb1-client v0.0.0-20190809212627-fc22c7df067e
	github.com/project-eria/eria-base/config-manager v0.0.2
	github.com/project-eria/eria-base/helpers v0.0.2
	github.com/project-eria/logger v0.0.1
	github.com/project-eria/xaal-go/device v0.5.0
	github.com/project-eria/xaal-go/engine v0.5.0
	github.com/project-eria/xaal-go/message v0.5.0
)

replace github.com/project-eria/xaal-go/device => ../xaal-go/device

replace github.com/project-eria/xaal-go/messagefactory => ../xaal-go/messagefactory

replace github.com/project-eria/xaal-go/message => ../xaal-go/message

replace github.com/project-eria/xaal-go/network => ../xaal-go/network

replace github.com/project-eria/xaal-go/utils => ../xaal-go/utils

replace github.com/project-eria/xaal-go/engine => ../xaal-go/engine

replace github.com/project-eria/xaal-go/schemas => ../xaal-go/schemas

replace github.com/project-eria/eria-base/helpers => ../eria-base/helpers

replace github.com/project-eria/eria-base/config-manager => ../eria-base/config-manager
