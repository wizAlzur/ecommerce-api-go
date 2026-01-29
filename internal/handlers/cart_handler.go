package handlers

import (
	"ecommerce-api/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service services.CartService
}

func NewCartHandler(cartService services.CartService) *CartHandler {
	return &CartHandler{service: cartService}
}

func (ch *CartHandler) AddToCart(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		ProductID int64 `json:"product_id" binding:"required"`
		Quantity  int   `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := ch.service.AddItem(ctx, userID, req.ProductID, req.Quantity); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "item added to cart"})
}

func (ch *CartHandler) UpdateCartItem(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid product_id"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,gte=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := ch.service.UpdateItem(ctx, userID, productID, req.Quantity); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "cart item updated"})
}

func (ch *CartHandler) RemoveFromCart(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid product_id"})
		return
	}

	ctx := c.Request.Context()
	if err := ch.service.RemoveItem(ctx, userID, productID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "item removed from cart"})
}

func (ch *CartHandler) GetCart(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	ctx := c.Request.Context()
	response, err := ch.service.GetCartResponse(ctx, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

func (ch *CartHandler) ClearCart(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	ctx := c.Request.Context()
	if err := ch.service.ClearCart(ctx, userID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "cart cleared"})
}

func getUserID(c *gin.Context) int64 {
	if val, ok := c.Get("user_id"); ok {
		return val.(int64)
	}
	return 0
}
