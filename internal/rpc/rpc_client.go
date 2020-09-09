package rpc

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/TeaOSLab/EdgeNode/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeNode/internal/const"
	"github.com/TeaOSLab/EdgeNode/internal/encrypt"
	"github.com/TeaOSLab/EdgeNode/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeNode/internal/utils"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

type RPCClient struct {
	apiConfig   *configs.APIConfig
	nodeClients []pb.NodeServiceClient
}

func NewRPCClient(apiConfig *configs.APIConfig) (*RPCClient, error) {
	if apiConfig == nil {
		return nil, errors.New("api config should not be nil")
	}

	nodeClients := []pb.NodeServiceClient{}

	conns := []*grpc.ClientConn{}
	for _, endpoint := range apiConfig.RPC.Endpoints {
		conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		conns = append(conns, conn)
	}
	if len(conns) == 0 {
		return nil, errors.New("[RPC]no available endpoints")
	}

	// node clients
	for _, conn := range conns {
		nodeClients = append(nodeClients, pb.NewNodeServiceClient(conn))
	}

	return &RPCClient{
		apiConfig:   apiConfig,
		nodeClients: nodeClients,
	}, nil
}

func (this *RPCClient) NodeRPC() pb.NodeServiceClient {
	if len(this.nodeClients) > 0 {
		return this.nodeClients[rands.Int(0, len(this.nodeClients)-1)]
	}
	return nil
}

func (this *RPCClient) Context() context.Context {
	ctx := context.Background()
	m := maps.Map{
		"timestamp": time.Now().Unix(),
		"type":      "node",
		"userId":    0,
	}
	method, err := encrypt.NewMethodInstance(teaconst.EncryptMethod, this.apiConfig.Secret, this.apiConfig.NodeId)
	if err != nil {
		utils.PrintError(err)
		return context.Background()
	}
	data, err := method.Encrypt(m.AsJSON())
	if err != nil {
		utils.PrintError(err)
		return context.Background()
	}
	token := base64.StdEncoding.EncodeToString(data)
	ctx = metadata.AppendToOutgoingContext(ctx, "nodeId", this.apiConfig.NodeId, "token", token)
	return ctx
}
