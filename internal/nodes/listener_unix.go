package nodes

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"net"
)

type UnixListener struct {
	BaseListener

	Group    *serverconfigs.ServerGroup
	Listener net.Listener
}

func (this *UnixListener) Serve() error {
	// TODO
	return nil
}

func (this *UnixListener) Close() error {
	// TODO
	return nil
}
