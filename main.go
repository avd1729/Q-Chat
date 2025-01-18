package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

// Global variables
var (
	mutex          sync.Mutex
	activeProtocol *BB84Protocol
	secureChannel  *SecureChannel
)

// Request and Response Models
type InitRequest struct {
	Bits int `json:"bits"`
}

type EncryptRequest struct {
	Plaintext string `json:"plaintext"`
	Sender    string `json:"sender"`
}

type EncryptResponse struct {
	Ciphertext string `json:"ciphertext"`
}

type DecryptRequest struct {
	Ciphertext string `json:"ciphertext"`
}

type DecryptResponse struct {
	Plaintext string `json:"plaintext"`
}

// InitializeProtocol initializes the BB84 protocol
func InitializeProtocol(bits int) error {
	mutex.Lock()
	defer mutex.Unlock()

	activeProtocol = NewBB84Protocol(bits)
	if err := activeProtocol.RunProtocol(); err != nil {
		return err
	}

	secureChannel = NewSecureChannel(activeProtocol.SharedKey)
	return nil
}

// API Endpoints

// Initialize the protocol
func initializeProtocolHandler(c *gin.Context) {
	var req InitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := InitializeProtocol(req.Bits); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize protocol"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Protocol initialized successfully", "sharedKey": activeProtocol.SharedKey})
}

// Encrypt a message
func encryptHandler(c *gin.Context) {
	var req EncryptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if secureChannel == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Secure channel not initialized"})
		return
	}

	msg, err := secureChannel.EncryptMessage(req.Plaintext, req.Sender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt message"})
		return
	}

	c.JSON(http.StatusOK, EncryptResponse{Ciphertext: msg.Ciphertext})
}

// Decrypt a message
func decryptHandler(c *gin.Context) {
	var req DecryptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if secureChannel == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Secure channel not initialized"})
		return
	}

	msg := &Message{
		Ciphertext: req.Ciphertext,
	}

	plaintext, err := secureChannel.DecryptMessage(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt message"})
		return
	}

	c.JSON(http.StatusOK, DecryptResponse{Plaintext: plaintext})
}

func getMessagesHandler(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	if secureChannel == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Secure channel not initialized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": secureChannel.Messages,
	})
}

// Main Function
func main() {
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // In production, replace with specific origin
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}

	r.Use(cors.New(config))

	// Your existing routes
	r.POST("/initialize", initializeProtocolHandler)
	r.POST("/encrypt", encryptHandler)
	r.POST("/decrypt", decryptHandler)
	r.GET("/messages", getMessagesHandler)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
