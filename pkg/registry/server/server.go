package server

import (
	"b1multicasting/pkg/registry/proto"
	"context"
	"sync"
)

type server struct {
	proto.UnimplementedRegistryServer
}

var (
	groups   map[string]*Mgroup
	Mugroups sync.RWMutex
)

type Mgroup struct {
	mu        sync.RWMutex
	groupInfo *proto.MGroup
}

// Register registers the calling node to a multicast group
func (s *server) Register(ctx context.Context, in *proto.Rinfo) (*proto.Ranswer, error) {

	Mugroups.Lock()
	defer Mugroups.Unlock()

	//Check if the group already exists, if not start group

	return nil, nil
}
