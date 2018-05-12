package api

import (
	"net/http"
	"os"

	"github.com/RTradeLtd/Temporal/rtfs_cluster"

	"github.com/RTradeLtd/Temporal/database"
	"github.com/RTradeLtd/Temporal/rtfs"
	"github.com/gin-contrib/rollbar"
	"github.com/gin-gonic/gin"
	"github.com/stvp/roll"
)

// Setup is used to build our routes
func Setup() *gin.Engine {
	token := os.Getenv("ROLLBAR_TOKEN")
	roll.Token = token
	roll.Environment = "development"
	r := gin.Default()
	r.Use(rollbar.Recovery(false))
	setupRoutes(r)
	return r
}

func setupRoutes(g *gin.Engine) {

	g.POST("/api/v1/ipfs/pin/:hash", pinHashLocally)
	g.POST("/api/v1/ipfs/add-file", addFileLocally)
	g.POST("/api/v1/ipfs-cluster/pin/:hash", pinHashToCluster)
	g.GET("/api/v1/ipfs-cluster/status/:hash", getLocalStatusForClusterPin)
	g.GET("/api/v1/ipfs/pins", getLocalPins)
	g.GET("/api/v1/database/uploads", getUploads)
	g.GET("/api/v1/database/uploads/:address", getUploadsForAddress)
}

func getUploads(c *gin.Context) {
	uploads := database.GetUploads()
	if uploads == nil {
		c.JSON(http.StatusNotFound, nil)
	}
	c.JSON(http.StatusFound, gin.H{"uploads": uploads})
}

func getUploadsForAddress(c *gin.Context) {
	uploads := database.GetUploadsForAddress(c.Param("address"))
	if uploads == nil {
		c.JSON(http.StatusNotFound, nil)
	}
	c.JSON(http.StatusFound, gin.H{"uploads": uploads})
}

func getLocalPins(c *gin.Context) {
	manager := rtfs.Initialize()
	pinInfo, err := manager.Shell.Pins()
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, gin.H{"pins": pinInfo})
}

func pinHashToCluster(c *gin.Context) {
	hash := c.Param("hash")
	err := database.AddHash(c)
	if err != nil {
		c.Error(err)
	}
	manager := rtfs_cluster.Initialize()
	contentIdentifier := manager.DecodeHashString(hash)
	manager.Client.Pin(contentIdentifier, -1, -1, hash)
	c.JSON(http.StatusOK, gin.H{"hash": hash})
}

func getLocalStatusForClusterPin(c *gin.Context) {
	hash := c.Param("hash")
	manager := rtfs_cluster.Initialize()
	status, err := manager.GetStatusForCidLocally(hash)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusFound, gin.H{"status": status})
}

func pinHashLocally(c *gin.Context) {
	hash := c.Param("hash")
	err := database.AddHash(c)
	if err != nil {
		return
	}
	manager := rtfs.Initialize()
	err = manager.Shell.Pin(hash)
	if err != nil {
		c.Error(err)
	}
	upload := database.GetUpload(hash, c.PostForm("uploadAddress"))
	c.JSON(http.StatusOK, gin.H{
		"hash":      upload.Hash,
		"uploader":  upload.UploadAddress,
		"retention": upload.HoldTimeInMonths})
}

func addFileLocally(c *gin.Context) {
	fileHandler, err := c.FormFile("file")
	if err != nil {
		c.Error(err)
	}
	openFile, err := fileHandler.Open()
	if err != nil {
		c.Error(err)
	}
	manager := rtfs.Initialize()
	resp, err := manager.Shell.Add(openFile)
	if err != nil {
		c.Error(err)
	}
	database.AddFileHash(c, resp)
	c.JSON(http.StatusOK, gin.H{"response": resp})
}
