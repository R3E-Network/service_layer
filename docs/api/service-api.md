# Service Layer API

## Authentication
All API requests require authentication using JWT tokens.

## Endpoints

### Function Management
- `POST /functions` - Deploy a new JavaScript function
- `GET /functions/{id}` - Get function details
- `PUT /functions/{id}` - Update function
- `DELETE /functions/{id}` - Delete function
- `POST /functions/{id}/execute` - Manually execute function

### Secret Management
- `POST /secrets` - Store a new secret
- `GET /secrets` - List user's secrets (returns only metadata)
- `DELETE /secrets/{id}` - Delete a secret

### Trigger Management
- `POST /triggers` - Create a new trigger
- `GET /triggers` - List all triggers
- `GET /triggers/{id}` - Get trigger details
- `PUT /triggers/{id}` - Update trigger
- `DELETE /triggers/{id}` - Delete trigger

### Price Feed
- `GET /pricefeeds` - List available price feeds
- `GET /pricefeeds/{id}` - Get specific price feed data
- `POST /pricefeeds` - Create custom price feed

### GasBank
- `GET /gasbank/balance` - Get user's gas balance
- `POST /gasbank/deposit` - Deposit gas
- `POST /gasbank/withdraw` - Withdraw gas 