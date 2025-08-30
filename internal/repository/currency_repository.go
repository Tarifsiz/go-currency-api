package repository

import (
	"context"
	"fmt"

	"github.com/Tarifsiz/go-currency-api/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CurrencyRepositoryInterface defines the contract for currency data operations
type CurrencyRepositoryInterface interface {
	// Basic CRUD operations
	Create(ctx context.Context, currency *model.Currency) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Currency, error)
	GetByCode(ctx context.Context, code string) (*model.Currency, error)
	GetAll(ctx context.Context, limit, offset int) ([]*model.Currency, error)
	Update(ctx context.Context, currency *model.Currency) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Business logic operations
	GetCurrenciesByFactor(ctx context.Context, factor int) ([]*model.Currency, error)
	SearchByName(ctx context.Context, name string) ([]*model.Currency, error)
	GetByCodes(ctx context.Context, codes []string) ([]*model.Currency, error)
	CreateBatch(ctx context.Context, currencies []*model.Currency) error
	GetCount(ctx context.Context) (int64, error)
}

// CurrencyRepository implements the CurrencyRepositoryInterface
type CurrencyRepository struct {
	db *gorm.DB
}

// NewCurrencyRepository creates a new currency repository instance
func NewCurrencyRepository(db *gorm.DB) CurrencyRepositoryInterface {
	return &CurrencyRepository{
		db: db,
	}
}

// Create creates a new currency record
func (r *CurrencyRepository) Create(ctx context.Context, currency *model.Currency) error {
	if err := r.db.WithContext(ctx).Create(currency).Error; err != nil {
		return fmt.Errorf("failed to create currency: %w", err)
	}
	return nil
}

// GetByID retrieves a currency by its UUID
func (r *CurrencyRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Currency, error) {
	var currency model.Currency
	err := r.db.WithContext(ctx).First(&currency, "id = ?", id).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("currency not found with id %s", id.String())
		}
		return nil, fmt.Errorf("failed to get currency by id: %w", err)
	}
	
	return &currency, nil
}

// GetByCode retrieves a currency by its code (e.g., "USD", "EUR")
func (r *CurrencyRepository) GetByCode(ctx context.Context, code string) (*model.Currency, error) {
	var currency model.Currency
	err := r.db.WithContext(ctx).First(&currency, "code = ?", code).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("currency not found with code %s", code)
		}
		return nil, fmt.Errorf("failed to get currency by code: %w", err)
	}
	
	return &currency, nil
}

// GetAll retrieves all currencies with pagination
func (r *CurrencyRepository) GetAll(ctx context.Context, limit, offset int) ([]*model.Currency, error) {
	var currencies []*model.Currency
	
	query := r.db.WithContext(ctx).Order("code ASC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Find(&currencies).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all currencies: %w", err)
	}
	
	return currencies, nil
}

// Update updates an existing currency record
func (r *CurrencyRepository) Update(ctx context.Context, currency *model.Currency) error {
	err := r.db.WithContext(ctx).
		Model(currency).
		Where("id = ?", currency.ID).
		Updates(currency).Error
	
	if err != nil {
		return fmt.Errorf("failed to update currency: %w", err)
	}
	
	return nil
}

// Delete deletes a currency record
func (r *CurrencyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Currency{}, "id = ?", id)
	
	if result.Error != nil {
		return fmt.Errorf("failed to delete currency: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("currency not found with id %s", id.String())
	}
	
	return nil
}

// GetCurrenciesByFactor retrieves currencies with a specific decimal factor
func (r *CurrencyRepository) GetCurrenciesByFactor(ctx context.Context, factor int) ([]*model.Currency, error) {
	var currencies []*model.Currency
	err := r.db.WithContext(ctx).
		Where("factor = ?", factor).
		Order("code ASC").
		Find(&currencies).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get currencies by factor: %w", err)
	}
	
	return currencies, nil
}

// SearchByName searches currencies by description/name
func (r *CurrencyRepository) SearchByName(ctx context.Context, name string) ([]*model.Currency, error) {
	var currencies []*model.Currency
	err := r.db.WithContext(ctx).
		Where("description ILIKE ?", "%"+name+"%").
		Order("code ASC").
		Find(&currencies).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to search currencies by name: %w", err)
	}
	
	return currencies, nil
}

// GetByCodes retrieves multiple currencies by their codes
func (r *CurrencyRepository) GetByCodes(ctx context.Context, codes []string) ([]*model.Currency, error) {
	if len(codes) == 0 {
		return []*model.Currency{}, nil
	}
	
	var currencies []*model.Currency
	err := r.db.WithContext(ctx).
		Where("code IN ?", codes).
		Order("code ASC").
		Find(&currencies).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get currencies by codes: %w", err)
	}
	
	return currencies, nil
}

// CreateBatch creates multiple currency records in a single transaction
func (r *CurrencyRepository) CreateBatch(ctx context.Context, currencies []*model.Currency) error {
	if len(currencies) == 0 {
		return nil
	}
	
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, currency := range currencies {
			if err := tx.Create(currency).Error; err != nil {
				return fmt.Errorf("failed to create currency %s: %w", currency.Code, err)
			}
		}
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to create currencies in batch: %w", err)
	}
	
	return nil
}

// GetCount returns the total count of currencies
func (r *CurrencyRepository) GetCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Currency{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get currency count: %w", err)
	}
	return count, nil
}