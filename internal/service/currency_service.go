package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Tarifsiz/go-currency-api/internal/model"
	"github.com/Tarifsiz/go-currency-api/internal/repository"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// CurrencyServiceInterface defines the business logic for currency operations
type CurrencyServiceInterface interface {
	// Basic CRUD operations
	CreateCurrency(ctx context.Context, currency *model.Currency) error
	GetCurrencyByID(ctx context.Context, id uuid.UUID) (*model.Currency, error)
	GetCurrencyByCode(ctx context.Context, code string) (*model.Currency, error)
	GetAllCurrencies(ctx context.Context, limit, offset int) ([]*model.Currency, error)
	UpdateCurrency(ctx context.Context, currency *model.Currency) error
	DeleteCurrency(ctx context.Context, id uuid.UUID) error
	
	// Business logic operations
	SearchCurrencies(ctx context.Context, query string) ([]*model.Currency, error)
	GetCurrenciesByFactor(ctx context.Context, factor int) ([]*model.Currency, error)
	GetCurrencyCount(ctx context.Context) (int64, error)
}

// CurrencyService implements the CurrencyServiceInterface
type CurrencyService struct {
	currencyRepo repository.CurrencyRepositoryInterface
	redisClient  *redis.Client
	cacheTimeout time.Duration
}

// NewCurrencyService creates a new currency service instance
func NewCurrencyService(currencyRepo repository.CurrencyRepositoryInterface, redisClient *redis.Client) CurrencyServiceInterface {
	return &CurrencyService{
		currencyRepo: currencyRepo,
		redisClient:  redisClient,
		cacheTimeout: 15 * time.Minute, // Cache currencies for 15 minutes
	}
}

// CreateCurrency creates a new currency
func (s *CurrencyService) CreateCurrency(ctx context.Context, currency *model.Currency) error {
	// Validate required fields
	if currency.Code == "" {
		return fmt.Errorf("currency code is required")
	}
	if currency.Description == "" {
		return fmt.Errorf("currency description is required")
	}
	
	// Set default values
	if currency.Factor == 0 {
		currency.Factor = 100 // Default to 2 decimal places
	}
	if currency.AmountDisplayFormat == "" {
		currency.AmountDisplayFormat = "###,###.##"
	}
	if currency.CreatedBy == uuid.Nil {
		// Set a default created_by UUID (in real app, this would come from auth context)
		currency.CreatedBy = uuid.MustParse("1609b0e1-30c4-402c-a76e-8f5b4d6cfc24")
	}
	
	// Create currency
	if err := s.currencyRepo.Create(ctx, currency); err != nil {
		return fmt.Errorf("failed to create currency: %w", err)
	}
	
	// Invalidate cache
	s.invalidateCache(ctx, currency.Code)
	
	return nil
}

// GetCurrencyByID retrieves a currency by ID
func (s *CurrencyService) GetCurrencyByID(ctx context.Context, id uuid.UUID) (*model.Currency, error) {
	return s.currencyRepo.GetByID(ctx, id)
}

// GetCurrencyByCode retrieves a currency by code with caching
func (s *CurrencyService) GetCurrencyByCode(ctx context.Context, code string) (*model.Currency, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("currency:code:%s", code)
	cachedCurrency, err := s.redisClient.Get(ctx, cacheKey).Result()
	
	if err == nil {
		// Cache hit - unmarshal and return
		var currency model.Currency
		if err := json.Unmarshal([]byte(cachedCurrency), &currency); err == nil {
			return &currency, nil
		}
	}
	
	// Cache miss - get from database
	currency, err := s.currencyRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	s.cacheCurrency(ctx, cacheKey, currency)
	
	return currency, nil
}

// GetAllCurrencies retrieves all currencies with pagination and caching
func (s *CurrencyService) GetAllCurrencies(ctx context.Context, limit, offset int) ([]*model.Currency, error) {
	// For simplicity, only cache the first page (offset = 0) with default limit
	if offset == 0 && limit <= 100 {
		cacheKey := fmt.Sprintf("currencies:all:%d:%d", limit, offset)
		cachedCurrencies, err := s.redisClient.Get(ctx, cacheKey).Result()
		
		if err == nil {
			// Cache hit
			var currencies []*model.Currency
			if err := json.Unmarshal([]byte(cachedCurrencies), &currencies); err == nil {
				return currencies, nil
			}
		}
		
		// Cache miss - get from database
		currencies, err := s.currencyRepo.GetAll(ctx, limit, offset)
		if err != nil {
			return nil, err
		}
		
		// Cache the result
		currenciesJSON, _ := json.Marshal(currencies)
		s.redisClient.Set(ctx, cacheKey, currenciesJSON, s.cacheTimeout)
		
		return currencies, nil
	}
	
	// For other pages, don't cache
	return s.currencyRepo.GetAll(ctx, limit, offset)
}

// UpdateCurrency updates an existing currency
func (s *CurrencyService) UpdateCurrency(ctx context.Context, currency *model.Currency) error {
	// Validate required fields
	if currency.Code == "" {
		return fmt.Errorf("currency code is required")
	}
	if currency.Description == "" {
		return fmt.Errorf("currency description is required")
	}
	
	// Update currency
	if err := s.currencyRepo.Update(ctx, currency); err != nil {
		return fmt.Errorf("failed to update currency: %w", err)
	}
	
	// Invalidate cache
	s.invalidateCache(ctx, currency.Code)
	
	return nil
}

// DeleteCurrency deletes a currency
func (s *CurrencyService) DeleteCurrency(ctx context.Context, id uuid.UUID) error {
	// Get currency first to get the code for cache invalidation
	currency, err := s.currencyRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get currency before deletion: %w", err)
	}
	
	// Delete currency
	if err := s.currencyRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete currency: %w", err)
	}
	
	// Invalidate cache
	s.invalidateCache(ctx, currency.Code)
	
	return nil
}

// SearchCurrencies searches currencies by name/description
func (s *CurrencyService) SearchCurrencies(ctx context.Context, query string) ([]*model.Currency, error) {
	if query == "" {
		return []*model.Currency{}, nil
	}
	
	return s.currencyRepo.SearchByName(ctx, query)
}

// GetCurrenciesByFactor retrieves currencies by decimal factor
func (s *CurrencyService) GetCurrenciesByFactor(ctx context.Context, factor int) ([]*model.Currency, error) {
	return s.currencyRepo.GetCurrenciesByFactor(ctx, factor)
}

// GetCurrencyCount returns total count of currencies
func (s *CurrencyService) GetCurrencyCount(ctx context.Context) (int64, error) {
	return s.currencyRepo.GetCount(ctx)
}

// Helper methods for caching

func (s *CurrencyService) cacheCurrency(ctx context.Context, cacheKey string, currency *model.Currency) {
	currencyJSON, err := json.Marshal(currency)
	if err == nil {
		s.redisClient.Set(ctx, cacheKey, currencyJSON, s.cacheTimeout)
	}
}

func (s *CurrencyService) invalidateCache(ctx context.Context, currencyCode string) {
	// Invalidate specific currency cache
	cacheKey := fmt.Sprintf("currency:code:%s", currencyCode)
	s.redisClient.Del(ctx, cacheKey)
	
	// Invalidate list cache (simple approach - delete all list caches)
	pattern := "currencies:all:*"
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err == nil && len(keys) > 0 {
		s.redisClient.Del(ctx, keys...)
	}
}