package controllers

import (
	"context"
	"encoding/json"
	"errors"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/msalvati1997/b1multicasting/internal/app"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"github.com/msalvati1997/b1multicasting/pkg/registry/proto"
	_ "github.com/msalvati1997/b1multicasting/pkg/registry/server"
	utils2 "github.com/msalvati1997/b1multicasting/pkg/utils"
	"github.com/prometheus/common/log"
	"time"
)

type RegistryController struct {
	beego.Controller
}

type MulticastController struct {
	beego.Controller
}

func (c *RegistryController) GetAllGroups() {

}

func (c *MulticastController) GetGroups() {

}

func (c *MulticastController) PutMessage() {
	var msg basic.Message
	multicastId := c.Ctx.Input.Param("multicastId")
	mtype := c.Ctx.Input.Param("type")
	msg.MessageHeader["type"] = mtype
	msg.MessageHeader["multicastId"] = multicastId
	json.Unmarshal(c.Ctx.Input.RequestBody, &msg.Payload)

	app.GMu.RLock()
	defer app.GMu.RUnlock()

	group, ok := app.MulticastGroups[multicastId]
	if !ok {
		c.Data["json"] = map[string]string{"Error": "GroupNotFound"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}

	group.GroupMu.RLock()
	defer group.GroupMu.RUnlock()

	if proto.Status(proto.Status_value[group.Group.Status]) != proto.Status_ACTIVE {
		c.Data["json"] = map[string]string{"Error": "group not ready"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
	utils2.GoPool.MessageCh <- msg
}

func (c *MulticastController) GetMessage() {

}

// CreateGroup initializes a new multicast group.
func (c *RegistryController) CreateGroup() {

	multicastId := c.Ctx.Input.Param("multicastId")
	app.GMu.RLock()
	defer app.GMu.RUnlock()

	group, ok := app.MulticastGroups[multicastId]

	if ok {
		c.Data["json"] = map[string]string{"Error": "group already exist"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}

	register, err := app.RegistryClient.Register(context.Background(), &proto.Rinfo{
		MulticastId: multicastId,
		ClientPort:  uint32(app.GrpcPort),
	})
	if err != nil {
		c.Data["json"] = map[string]string{"Error": "Error in creating group"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}

	members := make(map[string]app.Member, 0)

	for memberId, member := range register.GroupInfo.Members {
		members[memberId] = app.Member{
			MemberId: member.Id,
			Address:  member.Address,
			Ready:    member.Ready,
		}
	}

	group = &app.MulticastGroup{
		ClientId: register.ClientId,
		Group: &app.MulticastInfo{
			MulticastId:      register.GroupInfo.MulticastId,
			ReceivedMessages: 0,
			Status:           proto.Status_name[int32(register.GroupInfo.Status)],
			Members:          members,
		},
		Messages: make([]basic.Message, 0),
	}

	app.MulticastGroups[register.GroupInfo.MulticastId] = group
	go c.InitGroup(register.GroupInfo, group, len(register.GroupInfo.Members) == 1)
	c.Data["stato"] = "ok"
	err = c.ServeJSON()
	if err != nil {
		return
	}
}

func (c *RegistryController) InitGroup(info *proto.MGroup, group *app.MulticastGroup, b bool) {

	// Waiting that the group is ready
	updateMulticastGroup(info, group)
	groupInfo, err := waitStatusChange(info, group, proto.Status_OPENING)

	if err != nil {
		group.GroupMu.Lock()
		group.Group.Error = err
		group.GroupMu.Unlock()
		return
	}
	// Communicating to the registry that the node is ready
	groupInfo, err = app.RegistryClient.Ready(context.Background(), &proto.RequestData{
		MulticastId: group.Group.MulticastId,
		MId:         group.ClientId,
	})

	if err != nil {
		group.GroupMu.Lock()
		group.Group.Error = err
		group.GroupMu.Unlock()
		return
	}

	// Waiting tha all other nodes are ready
	updateMulticastGroup(groupInfo, group)
	groupInfo, err = waitStatusChange(groupInfo, group, proto.Status_STARTING)

	if err != nil {
		group.GroupMu.Lock()
		group.Group.Error = err
		group.GroupMu.Unlock()
		return
	}

	log.Info("Ready to multicast")
}

// StartGroup completes the initialization phase of a multicast group.
// The aforementioned group can now be used to perform multicast operations
func (c *RegistryController) StartGroup() {

	multicastId := c.Ctx.Input.Param("multicastId")

	app.GMu.RLock()
	defer app.GMu.RUnlock()

	group, ok := app.MulticastGroups[multicastId]

	if !ok {

		return
	}

	group.GroupMu.RLock()
	defer group.GroupMu.RUnlock()

	_, err := app.RegistryClient.StartGroup(context.Background(), &proto.RequestData{
		MulticastId: group.Group.MulticastId, MId: group.ClientId})
	if err != nil {
		c.Data["json"] = "Error " + err.Error()
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}

}

// updateMulticastGroup updates local group infos with ones received from the registry
func updateMulticastGroup(groupInfo *proto.MGroup, multicastGroup *app.MulticastGroup) {
	multicastGroup.GroupMu.Lock()
	defer multicastGroup.GroupMu.Unlock()

	multicastGroup.Group.Status = proto.Status_name[int32(groupInfo.Status)]

	for clientId, member := range groupInfo.Members {
		m, ok := multicastGroup.Group.Members[clientId]

		if !ok {
			m = app.Member{
				MemberId: member.Id,
				Address:  member.Address,
				Ready:    member.Ready,
			}

			multicastGroup.Group.Members[clientId] = m
		}

		if m.Ready != member.Ready {
			m.Ready = member.Ready
			multicastGroup.Group.Members[clientId] = m
		}

	}
}

// waitStatusChange waits until the status of the group in the registry changes
func waitStatusChange(groupInfo *proto.MGroup, multicastGroup *app.MulticastGroup, status proto.Status) (*proto.MGroup, error) {
	var err error

	for groupInfo.Status == status {
		time.Sleep(time.Second * 5)
		groupInfo, err = app.RegistryClient.GetStatus(context.Background(), &proto.MulticastId{MulticastId: groupInfo.MulticastId})
		if err != nil {
			return nil, err
		}

		updateMulticastGroup(groupInfo, multicastGroup)
	}

	if groupInfo.Status == proto.Status_CLOSED || groupInfo.Status == proto.Status_CLOSING {
		return nil, errors.New("multicast group is closed")
	}

	return groupInfo, nil
}
