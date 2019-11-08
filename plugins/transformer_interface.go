package plugins

import (
	"net/rpc"

	"github.com/damianoneill/h7t/dsl"
	"github.com/hashicorp/go-plugin"
)

// Transformer is the interface that we're exposing as a plugin.
type Transformer interface {
	Devices(args []string) (dsl.Devices, error)
}

// TransformerRPC - an implementation that talks over RPC
type TransformerRPC struct{ client *rpc.Client }

// Devices - interface implementation
func (g *TransformerRPC) Devices(args []string) dsl.Devices {
	var resp dsl.Devices
	err := g.client.Call("Plugin.Devices", args, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// TransformerRPCServer - is the RPC server that TransformerRPC talks to, conforming to
// the requirements of net/rpc
type TransformerRPCServer struct {
	Impl Transformer
}

// Devices - Server implementation
func (s *TransformerRPCServer) Devices(args []string, resp *dsl.Devices) (err error) {
	*resp, err = s.Impl.Devices(args)
	return
}

// TransformerPlugin is the implementation of plugin.Devices so we can serve/consume this
type TransformerPlugin struct {
	// Impl Injection
	Impl Transformer
}

// Server - muxing
func (p *TransformerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &TransformerRPCServer{Impl: p.Impl}, nil
}

// Client - muxing
func (TransformerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &TransformerRPC{client: c}, nil
}
