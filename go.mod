go 1.18

module github.com/cosmos/ibc-go/v3

retract [v3.0.0, v3.3.0] // depends on SDK version without dragonberry fix

require (
	github.com/armon/go-metrics v0.4.0
	github.com/confio/ics23/go v0.7.0
	github.com/cosmos/cosmos-sdk v0.45.12-0.20221116140330-9c145c827001
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.6.0
	github.com/spf13/viper v1.13.0
	github.com/stretchr/testify v1.8.0
	github.com/tendermint/tendermint v0.34.23
	github.com/tendermint/tm-db v0.6.6
	google.golang.org/genproto v0.0.0-20221014213838-99cd37c6964a
	google.golang.org/grpc v1.50.1
	google.golang.org/protobuf v1.28.2-0.20220831092852-f930b1dc76e8
	gopkg.in/yaml.v2 v2.4.0
)

require github.com/gin-gonic/gin v1.7.0 // indirect

replace (
	// dragonberry replace for ics23
	github.com/confio/ics23/go => github.com/cosmos/cosmos-sdk/ics23/go v0.8.0

	// protocol buffers replace
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)
