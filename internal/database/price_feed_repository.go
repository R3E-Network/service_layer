package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/R3E-Network/service_layer/internal/models"
	"github.com/jmoiron/sqlx"
)

// PriceFeedRepository implements the models.PriceFeedRepository interface
type PriceFeedRepository struct {
	db *sqlx.DB
}

// NewPriceFeedRepository creates a new price feed repository
func NewPriceFeedRepository(db *sqlx.DB) *PriceFeedRepository {
	return &PriceFeedRepository{
		db: db,
	}
}

// CreatePriceFeed creates a new price feed in the database
func (r *PriceFeedRepository) CreatePriceFeed(ctx interface{}, feed *models.PriceFeed) (*models.PriceFeed, error) {
	query := `
		INSERT INTO price_feeds 
		(symbol, update_interval, deviation_threshold, update_interval, contract_address, active, created_at, updated_at) 
		VALUES 
		($1, $2, $3, $4, $5, $6, $7, $7) 
		RETURNING id
	`

	now := time.Now().UTC()
	feed.CreatedAt = now
	feed.UpdatedAt = now

	var id string
	err := r.db.QueryRowx(query,
		feed.Symbol,
		feed.UpdateInterval,
		feed.DeviationThreshold,
		feed.UpdateInterval,
		feed.ContractAddress,
		feed.Active,
		now,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to create price feed: %w", err)
	}

	feed.ID = id
	return feed, nil
}

// GetPriceFeed gets a price feed by ID
func (r *PriceFeedRepository) GetPriceFeed(ctx interface{}, id string) (*models.PriceFeed, error) {
	query := `SELECT * FROM price_feeds WHERE id = $1`

	var feed models.PriceFeed
	err := r.db.Get(&feed, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrPriceFeedNotFound
		}
		return nil, fmt.Errorf("failed to get price feed: %w", err)
	}

	return &feed, nil
}

// GetPriceFeedBySymbol gets a price feed by symbol
func (r *PriceFeedRepository) GetPriceFeedBySymbol(ctx interface{}, symbol string) (*models.PriceFeed, error) {
	query := `SELECT * FROM price_feeds WHERE symbol = $1`

	var feed models.PriceFeed
	err := r.db.Get(&feed, query, symbol)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get price feed by symbol: %w", err)
	}

	return &feed, nil
}

// ListPriceFeeds lists all price feeds
func (r *PriceFeedRepository) ListPriceFeeds(ctx interface{}) ([]*models.PriceFeed, error) {
	query := `SELECT * FROM price_feeds ORDER BY id`

	var feeds []*models.PriceFeed
	err := r.db.Select(&feeds, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list price feeds: %w", err)
	}

	return feeds, nil
}

// UpdatePriceFeed updates a price feed
func (r *PriceFeedRepository) UpdatePriceFeed(ctx interface{}, feed *models.PriceFeed) (*models.PriceFeed, error) {
	query := `
		UPDATE price_feeds 
		SET symbol = $1, update_interval = $2, deviation_threshold = $3, update_interval = $4, 
		    contract_address = $5, active = $6, updated_at = $7 
		WHERE id = $8
	`

	now := time.Now().UTC()
	feed.UpdatedAt = now

	_, err := r.db.Exec(query,
		feed.Symbol,
		feed.UpdateInterval,
		feed.DeviationThreshold,
		feed.UpdateInterval,
		feed.ContractAddress,
		feed.Active,
		now,
		feed.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update price feed: %w", err)
	}

	return feed, nil
}

// DeletePriceFeed deletes a price feed
func (r *PriceFeedRepository) DeletePriceFeed(ctx interface{}, id string) error {
	query := `DELETE FROM price_feeds WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete price feed: %w", err)
	}

	return nil
}

// CreatePriceData creates a new price data entry
func (r *PriceFeedRepository) CreatePriceData(ctx interface{}, data *models.PriceData) (*models.PriceData, error) {
	query := `
		INSERT INTO price_data 
		(price_feed_id, price, timestamp, round_id, tx_hash, source, created_at) 
		VALUES 
		($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id
	`

	now := time.Now().UTC()
	data.CreatedAt = now

	var id string
	err := r.db.QueryRowx(query,
		data.PriceFeedID,
		data.Price,
		data.Timestamp,
		data.RoundID,
		data.TxHash,
		data.Source,
		now,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to create price data: %w", err)
	}

	data.ID = id
	return data, nil
}

// GetLatestPriceData gets the latest price data for a price feed
func (r *PriceFeedRepository) GetLatestPriceData(ctx interface{}, priceFeedID string) (*models.PriceData, error) {
	query := `
		SELECT * FROM price_data 
		WHERE price_feed_id = $1 
		ORDER BY timestamp DESC, id DESC 
		LIMIT 1
	`

	var data models.PriceData
	err := r.db.Get(&data, query, priceFeedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest price data: %w", err)
	}

	return &data, nil
}

// GetPriceDataHistory gets the price data history for a price feed
func (r *PriceFeedRepository) GetPriceDataHistory(ctx interface{}, priceFeedID string, limit int, offset int) ([]*models.PriceData, error) {
	query := `
		SELECT * FROM price_data 
		WHERE price_feed_id = $1 
		ORDER BY timestamp DESC, id DESC 
		LIMIT $2 OFFSET $3
	`

	var dataHistory []*models.PriceData
	err := r.db.Select(&dataHistory, query, priceFeedID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get price data history: %w", err)
	}

	return dataHistory, nil
}

// CreatePriceSource creates a new price source
func (r *PriceFeedRepository) CreatePriceSource(ctx interface{}, source *models.PriceSource) (*models.PriceSource, error) {
	query := `
		INSERT INTO price_sources 
		(name, url, weight, timeout, active, created_at, updated_at) 
		VALUES 
		($1, $2, $3, $4, $5, $6, $6) 
		RETURNING id
	`

	now := time.Now().UTC()
	source.CreatedAt = now
	source.UpdatedAt = now

	var id int
	err := r.db.QueryRowx(query,
		source.Name,
		source.URL,
		source.Weight,
		source.Timeout,
		source.Active,
		now,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to create price source: %w", err)
	}

	source.ID = fmt.Sprintf("%d", id)  // Convert int to string
	return source, nil
}

// GetPriceSourceByID gets a price source by ID
func (r *PriceFeedRepository) GetPriceSourceByID(ctx interface{}, id int) (*models.PriceSource, error) {
	query := `SELECT * FROM price_sources WHERE id = $1`

	var source models.PriceSource
	err := r.db.Get(&source, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get price source: %w", err)
	}

	return &source, nil
}

// GetPriceSourceByName gets a price source by name
func (r *PriceFeedRepository) GetPriceSourceByName(ctx interface{}, name string) (*models.PriceSource, error) {
	query := `SELECT * FROM price_sources WHERE name = $1`

	var source models.PriceSource
	err := r.db.Get(&source, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get price source by name: %w", err)
	}

	return &source, nil
}

// ListPriceSources lists all price sources
func (r *PriceFeedRepository) ListPriceSources(ctx interface{}) ([]*models.PriceSource, error) {
	query := `SELECT * FROM price_sources ORDER BY id`

	var sources []*models.PriceSource
	err := r.db.Select(&sources, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list price sources: %w", err)
	}

	return sources, nil
}

// UpdatePriceSource updates a price source
func (r *PriceFeedRepository) UpdatePriceSource(ctx interface{}, source *models.PriceSource) (*models.PriceSource, error) {
	query := `
		UPDATE price_sources 
		SET name = $1, url = $2, weight = $3, timeout = $4, active = $5, updated_at = $6 
		WHERE id = $7
	`

	now := time.Now().UTC()
	source.UpdatedAt = now

	_, err := r.db.Exec(query,
		source.Name,
		source.URL,
		source.Weight,
		source.Timeout,
		source.Active,
		now,
		source.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update price source: %w", err)
	}

	return source, nil
}

// DeletePriceSource deletes a price source
func (r *PriceFeedRepository) DeletePriceSource(ctx interface{}, id int) error {
	query := `DELETE FROM price_sources WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete price source: %w", err)
	}

	return nil
}

// AddPriceSource adds a price source to a price feed
func (r *PriceFeedRepository) AddPriceSource(ctx interface{}, priceFeedID string, source *models.PriceSource) (*models.PriceSource, error) {
	// Start a transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Create the price source if it doesn't exist
	if source.ID == "" {
		source.CreatedAt = time.Now().UTC()
		source.UpdatedAt = source.CreatedAt
		
		query := `
			INSERT INTO price_sources
			(name, url, path, weight, timeout, active, created_at, updated_at)
			VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
		
		err = tx.QueryRowx(query,
			source.Name,
			source.URL,
			source.Path,
			source.Weight,
			source.Timeout,
			source.Active,
			source.CreatedAt,
			source.UpdatedAt,
		).Scan(&source.ID)
		
		if err != nil {
			return nil, fmt.Errorf("failed to create price source: %w", err)
		}
	}

	// Link the price source to the price feed
	query := `
		INSERT INTO price_feed_sources
		(price_feed_id, price_source_id, created_at)
		VALUES
		($1, $2, $3)
		ON CONFLICT (price_feed_id, price_source_id) DO NOTHING
	`
	
	_, err = tx.Exec(query, priceFeedID, source.ID, time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("failed to add price source to price feed: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return source, nil
}

// RemovePriceSource removes a price source from a price feed
func (r *PriceFeedRepository) RemovePriceSource(ctx interface{}, priceFeedID string, sourceID string) error {
	query := `DELETE FROM price_feed_sources WHERE price_feed_id = $1 AND price_source_id = $2`
	_, err := r.db.Exec(query, priceFeedID, sourceID)
	if err != nil {
		return fmt.Errorf("failed to remove price source from price feed: %w", err)
	}
	return nil
}

// UpdateLastPrice updates the last price for a price feed
func (r *PriceFeedRepository) UpdateLastPrice(ctx interface{}, id string, price float64, txHash string) error {
	query := `
		UPDATE price_feeds
		SET last_price = $2, last_update = NOW(), last_tx_hash = $3
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, price, txHash)
	if err != nil {
		return fmt.Errorf("failed to update last price: %w", err)
	}
	return nil
}

// RecordPriceUpdate records a price update
func (r *PriceFeedRepository) RecordPriceUpdate(ctx interface{}, id string, price float64, txID string) error {
	update := &models.PriceUpdate{
		PriceFeedID:   id,
		Price:         price,
		TransactionID: txID,
		Success:       true,
		Timestamp:     time.Now().UTC(),
	}
	
	query := `
		INSERT INTO price_updates
		(price_feed_id, price, transaction_id, success, timestamp)
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING id
	`
	
	err := r.db.QueryRowx(query,
		update.PriceFeedID,
		update.Price,
		update.TransactionID,
		update.Success,
		update.Timestamp,
	).Scan(&update.ID)
	
	if err != nil {
		return fmt.Errorf("failed to record price update: %w", err)
	}
	
	return nil
}

// GetPriceHistory gets the price history for a price feed
func (r *PriceFeedRepository) GetPriceHistory(ctx interface{}, id string, limit int) ([]*models.PriceUpdate, error) {
	query := `
		SELECT * FROM price_updates
		WHERE price_feed_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`
	
	var updates []*models.PriceUpdate
	err := r.db.Select(&updates, query, id, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history: %w", err)
	}
	
	return updates, nil
}
