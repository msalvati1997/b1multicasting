package multicastapp

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	"net/http"
	"sync"
)

type HttpMethod int

const (
	HttpMethodGET HttpMethod = iota
	HttpMethodPOST
	HttpMethodPUT
	HttpMethodDELETE
	HttpMethodHEAD
)

var (
	HttpMethodName = map[HttpMethod]string{
		HttpMethodGET:    "GET",
		HttpMethodPOST:   "POST",
		HttpMethodPUT:    "PUT",
		HttpMethodDELETE: "DELETE",
		HttpMethodHEAD:   "HEAD",
	}
	MulticastTypeValue = map[string]HttpMethod{
		"GET":    HttpMethodGET,
		"POST":   HttpMethodPOST,
		"PUT":    HttpMethodPUT,
		"DELETE": HttpMethodDELETE,
		"HEAD":   HttpMethodHEAD,
	}
)

var restAPI = map[string]map[HttpMethod]func(ctx *gin.Context){

	"/groups": {
		HttpMethodGET: GetGroups,
	},
	"/groups/:multicastId": {
		HttpMethodPUT:  CreateGroup,
		HttpMethodPOST: StartGroup,
	},
}

func GetRestAPI() map[string]map[HttpMethod]func(ctx *gin.Context) {
	api := make(map[string]map[HttpMethod]func(ctx *gin.Context))

	for key, value := range restAPI {
		api[key] = value
	}

	return api
}

func GetRouterMethods(router *gin.RouterGroup) (map[HttpMethod]func(path string, handlers ...gin.HandlerFunc) gin.IRoutes, error) {
	routerMap := make(map[HttpMethod]func(path string, handlers ...gin.HandlerFunc) gin.IRoutes)

	routerMap[HttpMethodGET] = router.GET
	routerMap[HttpMethodPOST] = router.POST
	routerMap[HttpMethodPUT] = router.PUT
	routerMap[HttpMethodDELETE] = router.DELETE
	routerMap[HttpMethodHEAD] = router.HEAD

	return routerMap, nil

}

// GetGroups returns the groups to which the node is registered
func GetGroups(ctx *gin.Context) {

	groups := make([]*MulticastGroupInfo, 0)

	groupsMu.RLock()
	defer groupsMu.RUnlock()

	for _, group := range multicastGroups {
		group.groupMu.RLock()
		groups = append(groups, group.group)
		group.groupMu.RUnlock()
	}

	ctx.IndentedJSON(http.StatusOK, groups)
}

// CreateGroup initializes a new multicast group.
func CreateGroup(ctx *gin.Context) {
	var config GroupConfig

	multicastId := ctx.Param("multicastId")

	err := ctx.BindJSON(&config)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	multicastType, ok := registryservice.MulticastType[config.MulticastType]

	if !ok {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "multicast_type not supported", "supported": registryservice.MulticastTypes})
		return
	}

	groupsMu.Lock()
	defer groupsMu.Unlock()

	group, ok := multicastGroups[multicastId]

	if ok {
		ctx.IndentedJSON(http.StatusConflict, gin.H{"error": "group already exists"})
		return
	}

	registrationAns, err := registryClient.Register(context.Background(), &protoregistry.Rinfo{
		MulticastId:   multicastId,
		MulticastType: multicastType,
		ClientPort:    uint32(grpcPort),
	})

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	members := make(map[string]Member, 0)

	for memberId, member := range registrationAns.GroupInfo.Members {
		members[memberId] = Member{
			MemberId: member.Id,
			Address:  member.Address,
			Ready:    member.Ready,
		}
	}

	group = &MulticastGroup{
		clientId: registrationAns.ClientId,
		group: &MulticastGroupInfo{
			MulticastId:      registrationAns.GroupInfo.MulticastId,
			MulticastType:    protoregistry.MulticastType_name[int32(registrationAns.GroupInfo.MulticastType)],
			ReceivedMessages: 0,
			Status:           protoregistry.Status_name[int32(registrationAns.GroupInfo.Status)],
			Members:          members,
		},
		messages: make([]basic.Message, 0),
		groupMu:  sync.RWMutex{},
	}

	multicastGroups[registrationAns.GroupInfo.MulticastId] = group

	go initializeGroup(registrationAns.GroupInfo, group, len(registrationAns.GroupInfo.Members) == 1)

	ctx.Status(http.StatusOK)
}

// StartGroup completes the initialization phase of a multicast group.
// The aforementioned group can now be used to perform multicast operations
func StartGroup(ctx *gin.Context) {

	multicastId := ctx.Param("multicastId")

	groupsMu.RLock()
	defer groupsMu.RUnlock()

	group, ok := multicastGroups[multicastId]

	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	group.groupMu.RLock()
	defer group.groupMu.RUnlock()

	_, err := registryClient.StartGroup(context.Background(), &protoregistry.RequestData{MulticastId: group.group.MulticastId, MId: group.clientId})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
