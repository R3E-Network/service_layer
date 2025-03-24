# Price Feed Service

The Price Feed Service is responsible for fetching, caching, and providing cryptocurrency price data to other services within the application.

## Features

- Fetches price data from multiple sources (CoinMarketCap, CoinGecko, etc.)
- Maintains an in-memory cache of the latest price data
- Provides a REST API for querying price data
- Updates prices at configurable intervals
- Handles fallback to alternative data sources when primary sources fail

## Configuration

The service is configured through the `PriceFeedConfig` section in the application's configuration:

```json
{
  "pricefeed": {
    "updateIntervalSec": 300,
    "dataSources": ["coinmarketcap", "coingecko"],
    "supportedTokens": ["NEO", "GAS", "ETH", "BTC"]
  }
}
```

## Implementation Details

The service implements the `common.Service` interface:

- `Name()`: Returns "pricefeed"
- `Start(ctx context.Context)`: Initializes the price feeds and starts the update goroutine
- `Stop()`: Stops the update goroutine and cleans up resources
- `Health()`: Returns the health status based on successful price data retrieval

## API Usage

The service exposes methods for other internal components to use:

```go
// Get the latest price for a token in USD
price, err := priceFeedService.GetPrice("NEO")

// Get the latest prices for all supported tokens
prices, err := priceFeedService.GetAllPrices()
```

## Integration with Blockchain

The service also provides methods for blockchain contract interaction to enable on-chain price oracle functionality.
