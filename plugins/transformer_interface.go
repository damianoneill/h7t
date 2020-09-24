package plugins

import (
	"net/rpc"

	"github.com/damianoneill/h7t/pkg/dsl"
	"github.com/hashicorp/go-plugin"
)

// Arguments - composite for passing data across net/rpc
type Arguments struct {
	InputDirectory string
	CmdLineArgs    []string
}

// Transformer is the interface that must be implemented by plugin authors.
type Transformer interface {
	Devices(args Arguments) (dsl.Devices, error)
}

// TransformerRPC - an implementation that talks over RPC
type TransformerRPC struct{ client *rpc.Client }

// Devices - interface implementation
func (g *TransformerRPC) Devices(args Arguments) (resp dsl.Devices, err error) {
	err = g.client.Call("Plugin.Devices", args, &resp)
	return
}

// TransformerRPCServer - is the RPC server that TransformerRPC talks to, conforming to
// the requirements of net/rpc
type TransformerRPCServer struct {
	Impl Transformer
}

// Devices - Server implementation
func (s *TransformerRPCServer) Devices(args Arguments, resp *dsl.Devices) (err error) {
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
