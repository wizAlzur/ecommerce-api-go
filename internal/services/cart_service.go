package services

import (
	"context"
	"ecommerce-api/internal/models"
	"ecommerce-api/internal/repositories"
)

type CartService interface {
	AddItem(ctx context.Context, userID int64, productID int64, quantity int) error
	UpdateItem(ctx context.Context, userID int64, productID int64, quantity int) error
	RemoveItem(ctx context.Context, userID int64, productID int64) error
	GetCartResponse(ctx context.Context, userID int64) (*models.CartResponse, error)
	ClearCart(ctx context.Context, userID int64) error
}

type cartService struct {
	cartRepo    repositories.CartRepository
	productRepo repositories.ProductRepository
}

func NewCartService(cartRepo repositories.CartRepository, productRepo repositories.ProductRepository) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (cs *cartService) AddItem(ctx context.Context, userID int64, productID int64, quantity int) error {
	return cs.cartRepo.AddItem(ctx, userID, productID, quantity)
}

func (cs *cartService) UpdateItem(ctx context.Context, userID int64, productID int64, quantity int) error {
	return cs.cartRepo.UpdateItem(ctx, userID, productID, quantity)
}

func (cs *cartService) RemoveItem(ctx context.Context, userID int64, productID int64) error {
	return cs.cartRepo.RemoveItem(ctx, userID, productID)
}

func (cs *cartService) ClearCart(ctx context.Context, userID int64) error {
	return cs.cartRepo.ClearCart(ctx, userID)
}

func (cs *cartService) GetCartResponse(ctx context.Context, userID int64) (*models.CartResponse, error) {
	cartMap, err := cs.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(cartMap) == 0 {
		return &models.CartResponse{}, nil
	}

	productIDs := make([]int64, 0, len(cartMap))
	for id := range cartMap {
		productIDs = append(productIDs, id)
	}

	products, err := cs.productRepo.GetByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	productMap := make(map[int64]*models.Product)
	for _, p := range products {
		productMap[p.ID] = p
	}

	response := &models.CartResponse{
		Items: make([]models.CartResponseItem, 0, len(cartMap)),
	}

	var total float64
	var itemCount int

	for productID, quantity := range cartMap {
		product, found := productMap[productID]
		if !found {
			continue
		}

		subtotal := float64(quantity) * product.Price
		total += subtotal
		itemCount += quantity

		response.Items = append(response.Items, models.CartResponseItem{
			ProductID:   product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    quantity,
			Subtotal:    subtotal,
		})
	}

	response.Total = total
	response.ItemCount = itemCount

	return response, nil
}
