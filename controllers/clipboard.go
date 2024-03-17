package controllers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"

	"ClipCon-Server/database"
	"ClipCon-Server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateClipboardItem handles POST requests to store clipboard content.
func CreateClipboardItem(c *gin.Context) {
	var clipboardItem models.ClipboardItem
	if err := c.BindJSON(&clipboardItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clipboardItem.CreatedAt = time.Now()
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	result, err := database.ClipboardCollection.InsertOne(ctx, clipboardItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	clipboardItem.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusOK, clipboardItem)
}

// GetClipboardItems handles GET requests to fetch clipboard content.
func GetClipboardItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	page, err := getPageAndLimitFromQuery(c.Query("page"), c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.D{}
	findOptions := options.Find().SetSort(bson.D{{"createdAt", -1}}).SetSkip(int64(page.Skip)).SetLimit(int64(page.Limit))

	cursor, err := database.ClipboardCollection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var clipboardItems = make([]models.ClipboardItem, 0)
	if err = cursor.All(ctx, &clipboardItems); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clipboardItems)
}

// GetTotalClipboardItems handles GET requests to fetch the total count of clipboard items.
func GetTotalClipboardItems(c *gin.Context) {
	// Create a context with a timeout to ensure the request does not hang indefinitely
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Using an empty bson.D{} as filter condition to count all documents in the collection
	filter := bson.D{}

	// Retrieve the total number of documents that match the filter criteria
	total, err := database.ClipboardCollection.CountDocuments(ctx, filter)
	if err != nil {
		// If an error occurs, return an internal server error status and the error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total count: " + err.Error()})
		return
	}

	// Return the total count of documents as a JSON response
	c.JSON(http.StatusOK, gin.H{"total": total})
}

func getPageAndLimitFromQuery(pageQuery, limitQuery string) (pageValue, error) {
	page := 1
	limit := 10

	if pageQuery != "" {
		parsedPage, err := strconv.Atoi(pageQuery)
		if err != nil {
			return pageValue{}, err
		}
		page = parsedPage
	}

	if limitQuery != "" {
		parsedLimit, err := strconv.Atoi(limitQuery)
		if err != nil {
			return pageValue{}, err
		}
		limit = parsedLimit
	}

	return pageValue{Skip: (page - 1) * limit, Limit: limit}, nil
}

type pageValue struct {
	Skip  int
	Limit int
}
