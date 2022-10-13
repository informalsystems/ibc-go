package ica

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/client/cli"
	controllerkeeper "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/keeper"
	controllertypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/types"
	"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/host"
	hostkeeper "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/host/keeper"
	hosttypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/host/types"
	"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/simulation"
	"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	porttypes "github.com/cosmos/ibc-go/v5/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v5/modules/core/24-host"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}

	_ porttypes.IBCModule = host.IBCModule{}
)

// AppModuleBasic is the IBC interchain accounts AppModuleBasic
type AppModuleBasic struct{}

// Name implements AppModuleBasic interface
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec implements AppModuleBasic.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// RegisterInterfaces registers module concrete types into protobuf Any
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the IBC
// interchain accounts module
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis performs genesis state validation for the IBC interchain acounts module
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return gs.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the interchain accounts module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := controllertypes.RegisterQueryHandlerClient(context.Background(), mux, controllertypes.NewQueryClient(clientCtx))
	if err != nil {
		panic(err)
	}

	err = hosttypes.RegisterQueryHandlerClient(context.Background(), mux, hosttypes.NewQueryClient(clientCtx))
	if err != nil {
		panic(err)
	}
}

// GetTxCmd implements AppModuleBasic interface
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd implements AppModuleBasic interface
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// AppModule is the application module for the IBC interchain accounts module
type AppModule struct {
	AppModuleBasic
	controllerKeeper *controllerkeeper.Keeper
	hostKeeper       *hostkeeper.Keeper
}

// NewAppModule creates a new IBC interchain accounts module
func NewAppModule(controllerKeeper *controllerkeeper.Keeper, hostKeeper *hostkeeper.Keeper) AppModule {
	return AppModule{
		controllerKeeper: controllerKeeper,
		hostKeeper:       hostKeeper,
	}
}

// InitModule will initialize the interchain accounts moudule. It should only be
// called once and as an alternative to InitGenesis.
func (am AppModule) InitModule(ctx sdk.Context, controllerParams controllertypes.Params, hostParams hosttypes.Params) {
	if am.controllerKeeper != nil {
		am.controllerKeeper.SetParams(ctx, controllerParams)
	}

	if am.hostKeeper != nil {
		am.hostKeeper.SetParams(ctx, hostParams)

		cap := am.hostKeeper.BindPort(ctx, types.PortID)
		if err := am.hostKeeper.ClaimCapability(ctx, cap, ibchost.PortPath(types.PortID)); err != nil {
			panic(fmt.Sprintf("could not claim port capability: %v", err))
		}
	}
}

// RegisterInvariants implements the AppModule interface
func (AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

// NewHandler implements the AppModule interface
func (AppModule) NewHandler() sdk.Handler {
	return nil
}

// RegisterServices registers module services
func (am AppModule) RegisterServices(cfg module.Configurator) {
	if am.controllerKeeper != nil {
		controllertypes.RegisterQueryServer(cfg.QueryServer(), am.controllerKeeper)
	}

	if am.hostKeeper != nil {
		hosttypes.RegisterQueryServer(cfg.QueryServer(), am.hostKeeper)
	}
}

// InitGenesis performs genesis initialization for the interchain accounts module.
// It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	if am.controllerKeeper != nil {
		controllerkeeper.InitGenesis(ctx, *am.controllerKeeper, genesisState.ControllerGenesisState)
	}

	if am.hostKeeper != nil {
		hostkeeper.InitGenesis(ctx, *am.hostKeeper, genesisState.HostGenesisState)
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the interchain accounts module
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	var (
		controllerGenesisState = types.DefaultControllerGenesis()
		hostGenesisState       = types.DefaultHostGenesis()
	)

	if am.controllerKeeper != nil {
		controllerGenesisState = controllerkeeper.ExportGenesis(ctx, *am.controllerKeeper)
	}

	if am.hostKeeper != nil {
		hostGenesisState = hostkeeper.ExportGenesis(ctx, *am.hostKeeper)
	}

	gs := types.NewGenesisState(controllerGenesisState, hostGenesisState)

	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock implements the AppModule interface
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
}

// EndBlock implements the AppModule interface
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the ics27 module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations is unimplemented.
func (am AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}

// RandomizedParams creates randomized ibc-transfer param changes for the simulator.
func (am AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return simulation.ParamChanges(r, am.controllerKeeper, am.hostKeeper)
}

// RegisterStoreDecoder registers a decoder for interchain accounts module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[types.StoreKey] = simulation.NewDecodeStore()
}
