package common

import "log"

const (
	PluginDBMain  = "mysql"
	PluginPubSub  = "pubsub"
	PluginAPIItem = "item-api"

	TopicUserLikeItem   = "TopicUserLikeItem"
	TopicUserUnlikeItem = "TopicUserUnlikeItem"
)

type DbType int

const (
	DbTypeItem DbType = 1
	DbTypeUser DbType = 2
)

const (
	CurrentUser = "current_user"
)

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

func Recovery() {
	if r := recover(); r != nil {
		log.Println("recovered from: ", r)
	}
}

type TokenPayload struct {
	UId   int    `json:"user_id"`
	URole string `json:"role"`
}

func (p TokenPayload) UserId() int {
	return p.UId
}

func (p TokenPayload) Role() string {
	return p.URole
}

func IsAdmin(requester Requester) bool {
	return requester.GetRole() == "admin" || requester.GetRole() == "mod"
}
