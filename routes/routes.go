package routes

import (
	"example/hivemind-be/account"
	"example/hivemind-be/comment"
	"example/hivemind-be/content"
	"example/hivemind-be/hive"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	// Content
	router.GET("/content", content.GetContent)
	router.GET("/content/id/:id", content.GetContentById)
	router.GET("/content/uuid/:uuid", content.GetContentByUuid)
	router.POST("/content", content.CreateContent)
	router.PATCH("/content/uuid/:uuid/add-upvote", content.AddContentUpvoteByUuid)
	router.PATCH("/content/uuid/:uuid/remove-upvote", content.RemoveContentUpvoteByUuid)
	router.PATCH("/content/uuid/:uuid/add-downvote", content.AddContentDownvoteByUuid)
	router.PATCH("/content/uuid/:uuid/remove-downvote", content.RemoveContentDownvoteByUuid)
	router.PATCH("/content/uuid/:uuid/delete", content.DeleteContentByUuid)
	router.PATCH("/content/uuid/:uuid/undelete", content.UndeleteContentByUuid)
	router.PATCH("/content/uuid/:uuid/update", content.UpdateContentByUuid)

	// Comment via Content
	router.GET("/content/uuid/:uuid/comment", comment.GetCommentsByContentUuid)
	router.POST("/content/uuid/:uuid/comment", comment.CreateComment)
	router.POST("/content/uuid/:uuid/comment/:parentuuid/reply", comment.CreateCommentReply)

	// Comment
	router.GET("/comment/uuid/:uuid", comment.GetCommentByUuid)
	router.GET("/comment/uuid/:uuid/replies", comment.GetCommentByUuidWithReplies)
	router.PATCH("/comment/uuid/:uuid/delete", comment.DeleteCommentByUuid)
	router.PATCH("/comment/uuid/:uuid/undelete", comment.UndeleteCommentByUuid)
	router.PATCH("/comment/uuid/:uuid/update", comment.UpdateCommentByUuid)
	router.PATCH("/comment/uuid/:uuid/add-upvote", comment.AddCommentUpvoteByUuid)
	router.PATCH("/comment/uuid/:uuid/remove-upvote", comment.RemoveCommentUpvoteByUuid)
	router.PATCH("/comment/uuid/:uuid/add-downvote", comment.AddCommentDownvoteByUuid)
	router.PATCH("/comment/uuid/:uuid/remove-downvote", comment.RemoveCommentDownvoteByUuid)

	// Hive
	router.GET("/hive", hive.GetHive)
	router.GET("/hive/uuid/:uuid/content", content.GetContentByHiveUuid)
	router.POST("/hive", hive.CreateHive)
	router.PATCH("/hive/uuid/:uuid/ban", hive.BanHiveByUuid)
	router.PATCH("/hive/uuid/:uuid/unban", hive.UnBanHiveByUuid)
	router.PATCH("/hive/uuid/:uuid/archive", hive.ArchiveHiveByUuid)
	router.PATCH("/hive/uuid/:uuid/unarchive", hive.UnArchiveHiveByUuid)
	router.PATCH("/hive/uuid/:uuid/update", hive.UpdateHiveByUuid)

	// Account
	router.POST("/account/create", account.CreateAccount)
	router.POST("/account/login", account.AccountLogin)
	router.POST("/account/token/refresh", account.RefreshAuthToken)
	router.GET("/account/token/validate", account.ValidateAccountToken)
	router.GET("/account", account.GetAccount)
	router.PATCH("/account/change-password", account.ChangePassword)
}
