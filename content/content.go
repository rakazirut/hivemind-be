package content

import (
	"example/hivemind-be/db"
	"example/hivemind-be/hive"
	"example/hivemind-be/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// GORM uses the name of your type as the DB table to query. Here the type is Message so gorm will use the messages table by default.
type Content struct {
	ID           int32       `json:"Id" gorm:"primaryKey:type:int32"` //cannot be updated
	Hive         string      `json:"Hive"`                            //cannot be updated
	Title        string      `json:"Title"`                           //can be updated
	Author       string      `json:"Author"`                          //cannot be updated
	Message      string      `json:"Message"`                         //can be updated
	UUID         string      `json:"Uuid"`                            //cannot be update
	HiveUUID     string      `json:"HiveUuid"`                        //cannot be update
	AccountUUID  string      `json:"AccountUuid"`                     //cannot be update
	Link         string      `json:"Link" gorm:"default:null"`        //can be updated
	ImageLink    string      `json:"ImageLink" gorm:"default:null"`   //can be updated
	Upvote       int32       `json:"Upvote"`                          //cannot be updated
	Downvote     int32       `json:"Downvote"`                        //cannot be updated
	CommentCount int32       `json:"CommentCount"`                    //cannot be updated
	Deleted      bool        `json:"Deleted"`                         //can be updated
	Created      pq.NullTime `json:"Created"`                         //cannot be updated
	LastEdited   pq.NullTime `json:"LastEdited"`                      //updated when an update occurs
}

type ContentVote struct {
	ID          int32
	AccountUUID string
	ContentUUID string
	Upvote      bool
	Downvote    bool
	LastEdited  pq.NullTime
}

func GetContent(c *gin.Context) {
	var content []Content

	cookie, err := c.Request.Cookie("Token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	if result := db.Db.Order("id asc").Find(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &content)
}

func GetContentById(c *gin.Context) {
	var content Content

	cookie, err := c.Request.Cookie("Token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	id := c.Param("id")
	if result := db.Db.First(&content, id); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, content)
}

func GetContentByUuid(c *gin.Context) {
	var content Content

	cookie, err := c.Request.Cookie("Token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uuid := c.Param("uuid")
	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, content)
}

func GetContentByHiveUuid(c *gin.Context) {
	var content []Content

	cookie, err := c.Request.Cookie("Token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uuid := c.Param("uuid")
	if result := db.Db.Where("hive_uuid = ?", uuid).Find(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, content)
}

func CreateContent(c *gin.Context) {
	var content Content
	var hive hive.Hive

	cookie, err := c.Request.Cookie("Token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	claims, err := token.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	if err := c.BindJSON(&content); err != nil {
		return
	}

	if result := db.Db.Where("name = ?", content.Hive).First(&hive); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": "Hive not found! Please use an existing hive or create a new hive first.",
		})
		return
	}

	content.Author = claims.Username
	content.UUID = uuid.NewString()
	content.HiveUUID = hive.UUID
	content.AccountUUID = claims.AccountUUID
	content.Upvote = 0
	content.Downvote = 0
	content.CommentCount = 0
	content.Deleted = false
	content.LastEdited = pq.NullTime{Valid: false}
	content.Created = pq.NullTime{Time: time.Now(), Valid: true}

	if result := db.Db.Create(&content); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "There was an error creating this content. Please try again.",
		})
		return
	}

	hive.TotalContent += 1
	db.Db.Save(&hive)
	c.JSON(http.StatusCreated, content)
}

func AddContentUpvoteByUuid(c *gin.Context) {
	var content Content
	var hive hive.Hive
	var contentVote ContentVote

	cookie, err := c.Request.Cookie("Token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	claims, err := token.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uuid := c.Param("uuid")

	//check content exsit
	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//check hive exist
	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	voteQuery := map[string]interface{}{
		"account_uuid": claims.AccountUUID,
		"content_uuid": uuid,
	}

	//check account vote
	if result := db.Db.Where(voteQuery).First(&contentVote); result.Error != nil {
		//user has no record
		content.Upvote += 1
		hive.TotalUpvotes += 1
		contentVote.AccountUUID = claims.AccountUUID
		contentVote.ContentUUID = uuid
		contentVote.Upvote = true
		contentVote.Downvote = false
		contentVote.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
		db.Db.Save(&content)
		db.Db.Save(&hive)
		db.Db.Save(&contentVote)
		c.JSON(http.StatusOK, gin.H{
			"Message": "User successfully upvoted!",
		})
		return
	}

	//error if user has already voted
	if contentVote.Upvote || contentVote.Downvote {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User has already voted on this content!",
		})
		return
	}

	//user has false for both upvote and downvote
	if !contentVote.Upvote && !contentVote.Downvote {
		content.Upvote += 1
		hive.TotalUpvotes += 1
		contentVote.Upvote = true
		contentVote.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
		db.Db.Save(&content)
		db.Db.Save(&hive)
		db.Db.Save(&contentVote)
		c.JSON(http.StatusOK, gin.H{
			"Message": "User successfully upvoted!",
		})
		return
	}
}

func RemoveContentUpvoteByUuid(c *gin.Context) {
	var content Content
	var hive hive.Hive
	var contentVote ContentVote

	cookie, err := c.Request.Cookie("Token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	claims, err := token.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uuid := c.Param("uuid")

	//check content exist
	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//check hive exist
	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	voteQuery := map[string]interface{}{
		"account_uuid": claims.AccountUUID,
		"content_uuid": uuid,
	}

	//check account vote
	if result := db.Db.Where(voteQuery).First(&contentVote); result.Error != nil {
		//user has not voted at all
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User has not voted on this content!",
		})
		return
	}

	//error if user has not already upvoted or has downvoted
	if !contentVote.Upvote || contentVote.Downvote {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User has not upvoted on this content!",
		})
		return
	}

	content.Upvote -= 1
	hive.TotalUpvotes -= 1
	contentVote.Upvote = false
	contentVote.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
	db.Db.Save(&content)
	db.Db.Save(&hive)
	db.Db.Save(&contentVote)
	c.JSON(http.StatusOK, gin.H{
		"Message": "User upvote removed sucessfully!",
	})
}

func AddContentDownvoteByUuid(c *gin.Context) {
	var content Content
	var hive hive.Hive
	var contentVote ContentVote

	cookie, err := c.Request.Cookie("Token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	claims, err := token.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uuid := c.Param("uuid")

	//check content exist
	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//check hive exist
	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	voteQuery := map[string]interface{}{
		"account_uuid": claims.AccountUUID,
		"content_uuid": uuid,
	}

	//check account vote
	if result := db.Db.Where(voteQuery).First(&contentVote); result.Error != nil {
		//user has no record
		content.Downvote += 1
		hive.TotalDownvotes += 1
		contentVote.AccountUUID = claims.AccountUUID
		contentVote.ContentUUID = uuid
		contentVote.Upvote = false
		contentVote.Downvote = true
		contentVote.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
		db.Db.Save(&content)
		db.Db.Save(&hive)
		db.Db.Save(&contentVote)
		c.JSON(http.StatusOK, gin.H{
			"Message": "User successfully downvoted!",
		})
		return
	}

	//error if user has already voted
	if contentVote.Upvote || contentVote.Downvote {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User has already voted on this content!",
		})
		return
	}

	//user has false for both upvote and downvote
	if !contentVote.Upvote && !contentVote.Downvote {
		content.Downvote += 1
		hive.TotalDownvotes += 1
		contentVote.Downvote = true
		contentVote.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
		db.Db.Save(&content)
		db.Db.Save(&hive)
		db.Db.Save(&contentVote)
		c.JSON(http.StatusOK, gin.H{
			"Message": "User successfully downvoted!",
		})
		return
	}
}

func RemoveContentDownvoteByUuid(c *gin.Context) {
	var content Content
	var hive hive.Hive
	var contentVote ContentVote

	cookie, err := c.Request.Cookie("Token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	claims, err := token.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uuid := c.Param("uuid")

	//check content exist
	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//check hive exist
	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	voteQuery := map[string]interface{}{
		"account_uuid": claims.AccountUUID,
		"content_uuid": uuid,
	}

	//check account vote
	if result := db.Db.Where(voteQuery).First(&contentVote); result.Error != nil {
		//user has not voted at all
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User has not voted on this content!",
		})
		return
	}

	//error if user has not already downvoted or has upvoted
	if contentVote.Upvote || !contentVote.Downvote {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User has not downvoted on this content!",
		})
		return
	}

	content.Downvote -= 1
	hive.TotalDownvotes -= 1
	contentVote.Downvote = false
	contentVote.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
	db.Db.Save(&content)
	db.Db.Save(&hive)
	db.Db.Save(&contentVote)
	c.JSON(http.StatusOK, gin.H{
		"Message": "User downvote removed sucessfully!",
	})
}

func DeleteContentByUuid(c *gin.Context) {
	var content Content
	var hive hive.Hive

	cookie, err := c.Request.Cookie("Token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if content.Deleted {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "content has already been deleted!"})
		return
	}

	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	hive.TotalContent -= 1
	content.Deleted = true
	content.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
	db.Db.Save(&content)
	db.Db.Save(&hive)
	c.JSON(http.StatusOK, content)
}

func UndeleteContentByUuid(c *gin.Context) {
	var content Content
	var hive hive.Hive

	cookie, err := c.Request.Cookie("Token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if !content.Deleted {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "content has not been deleted!"})
		return
	}

	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	hive.TotalContent += 1
	content.Deleted = false
	content.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
	db.Db.Save(&content)
	db.Db.Save(&hive)
	c.JSON(http.StatusOK, content)
}

func UpdateContentByUuid(c *gin.Context) {
	var content Content
	var updateContent Content

	cookie, err := c.Request.Cookie("Token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	if err := c.BindJSON(&updateContent); err != nil {
		return
	}

	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if val, ok := jsonDataHasKey(updateContent, "title"); ok {
		content.Title = val
	}
	if val, ok := jsonDataHasKey(updateContent, "message"); ok {
		content.Message = val
	}
	if val, ok := jsonDataHasKey(updateContent, "link"); ok {
		content.Link = val
	}
	if val, ok := jsonDataHasKey(updateContent, "imageLink"); ok {
		content.ImageLink = val
	}

	content.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}

	db.Db.Save(&content)
	c.JSON(http.StatusOK, content)
}

func jsonDataHasKey(data Content, key string) (string, bool) {
	switch key {
	case "title":
		return data.Title, true
	case "message":
		return data.Message, true
	case "link":
		return data.Link, true
	case "imageLink":
		return data.ImageLink, true
	default:
		return "null", false
	}
}
