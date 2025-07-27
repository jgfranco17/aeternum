# JWT Authentication & Supabase Integration

The Aeternum API now includes JWT (JSON Web Token) authentication for securing endpoints, powered by Supabase Auth, and Supabase database integration for storing and retrieving test results.

## Environment Setup

The JWT secret and Supabase configuration are set via environment variables in your `.envrc` file:

```bash
export AETERNUM_JWT_SECRET="your-secret-key-here"
export AETERNUM_DB_URL="https://your-project.supabase.co"
export AETERNUM_DB_KEY="your-anon-key"
```

## Authentication Flow

### 1. Registration

To create a new user account, make a POST request to the `/register` endpoint:

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "your-password"
  }'
```

**Response:**

```json
{
  "id": "user-uuid",
  "email": "user@example.com"
}
```

### 2. Login

To obtain a JWT token, make a POST request to the `/login` endpoint:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "your-password"
  }'
```

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "type": "Bearer"
}
```

### 3. Using Protected Endpoints

All `/v0/*` endpoints require authentication. Include the JWT token in the Authorization header:

```bash
curl -X GET http://localhost:8080/v0/tests/results \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Protected Endpoints

The following endpoints require JWT authentication:

- `POST /v0/tests/run` - Run tests and store results
- `GET /v0/tests/results?id=<request_id>` - Get specific test result
- `GET /v0/tests/history?limit=10` - Get user's test history

## Test Results Storage

### Running Tests with Storage

When you run tests via `POST /v0/tests/run`, the results are automatically stored in Supabase:

```bash
curl -X POST http://localhost:8080/v0/tests/run \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "base_url": "https://api.example.com",
    "endpoints": [
      {
        "path": "/health",
        "expected_status": 200
      }
    ]
  }'
```

**Response includes RequestID for later retrieval:**

```json
{
  "request_id": "aeternum-v0-uuid-here",
  "base_url": "https://api.example.com",
  "status": "PASS",
  "results": [...]
}
```

### Retrieving Test Results

Get a specific test result by its RequestID:

```bash
curl -X GET "http://localhost:8080/v0/tests/results?id=aeternum-v0-uuid-here" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**

```json
{
  "id": "aeternum-v0-uuid-here",
  "user_id": "user-uuid",
  "request_id": "aeternum-v0-uuid-here",
  "base_url": "https://api.example.com",
  "status": "PASS",
  "results": [...],
  "created_at": "2024-01-01T12:00:00Z",
  "metadata": {
    "endpoint_count": 1,
    "passed_count": 1,
    "failed_count": 0
  }
}
```

### Getting Test History

Retrieve all test results for the authenticated user:

```bash
curl -X GET "http://localhost:8080/v0/tests/history?limit=5" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**

```json
{
  "results": [
    {
      "id": "aeternum-v0-uuid-1",
      "user_id": "user-uuid",
      "request_id": "aeternum-v0-uuid-1",
      "base_url": "https://api.example.com",
      "status": "PASS",
      "results": [...],
      "created_at": "2024-01-01T12:00:00Z",
      "metadata": {...}
    }
  ],
  "count": 1
}
```

## Error Responses

### Unauthorized (401)

- Missing Authorization header
- Invalid token format
- Expired or invalid token

```json
{
  "message": "Authorization header required"
}
```

### Bad Request (400)

- Invalid login credentials
- Malformed request body
- Registration errors
- Invalid query parameters

```json
{
  "message": "Invalid credentials"
}
```

### Not Found (404)

- Test result not found for the given ID

```json
{
  "message": "No result found for ID aeternum-v0-uuid-here"
}
```

## Token Details

- **Algorithm**: HS256
- **Expiration**: 24 hours from issuance
- **Claims**: User ID and email address
- **Provider**: Supabase Auth

## User Management

User accounts are managed through Supabase Auth, which provides:

- Secure password hashing
- Email verification (configurable)
- Password reset functionality
- User profile management
- Admin user management through Supabase dashboard

## Database Integration

Test results are stored in Supabase with the following features:

- **User Isolation**: Each user can only access their own test results
- **Metadata Storage**: Additional metrics like endpoint count, pass/fail counts
- **Timestamp Tracking**: Automatic creation timestamps for all results
- **Scalable Storage**: Leverages Supabase's PostgreSQL backend
- **ORM Interface**: Uses the official [Supabase Go SDK](https://github.com/supabase-community/supabase-go) for clean database operations

### Implementation Status

**✅ Fully Implemented:**

- Authentication system with Supabase Auth
- Real database operations using Supabase Go SDK
- API endpoints and routing
- Database interface and data persistence
- Error handling and logging
- User context and security

### Database Schema

The following table must be created in your Supabase project:

```sql
CREATE TABLE test_results (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    request_id TEXT NOT NULL,
    base_url TEXT NOT NULL,
    status TEXT NOT NULL,
    results JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB
);

-- Enable Row Level Security
ALTER TABLE test_results ENABLE ROW LEVEL SECURITY;

-- Create policy to allow users to only see their own results
CREATE POLICY "Users can only access their own test results" ON test_results
    FOR ALL USING (auth.uid()::text = user_id);
```

### Database Operations

The system uses the [Supabase Go SDK](https://github.com/supabase-community/supabase-go) for clean, ORM-like database operations:

1. **INSERT**: `client.From("test_results").Insert(data).Execute()`
2. **SELECT**: `client.From("test_results").Select("*").Eq("id", value).Execute()`
3. **SELECT with filtering**: `client.From("test_results").Select("*").Eq("user_id", userID).Limit(limit).Execute()`
4. **Error handling**: Proper error responses for database failures

## Security Features

- **Supabase Auth**: Industry-standard authentication with secure password hashing
- **JWT Tokens**: Stateless authentication tokens
- **Environment-based Configuration**: No hardcoded secrets
- **Error Handling**: Consistent error responses with proper HTTP status codes
- **Input Validation**: Request body validation using Gin's binding
- **User Data Isolation**: Database queries are scoped to authenticated user
- **Row Level Security**: Database-level security policies

## Development Notes

The authentication and database system uses the [Supabase Go SDK](https://github.com/supabase-community/supabase-go) and follows the v0 route pattern with function factories and error handling middleware. The database integration is fully implemented using the official SDK's ORM-like interface.

### Technical Implementation

- **Supabase Go SDK**: Uses the official Go client library for clean database operations
- **ORM-like Interface**: Fluent API for building database queries
- **JSON Handling**: Automatic marshaling/unmarshaling of data structures
- **Error Handling**: Comprehensive error handling for database operations
- **Logging**: Detailed logging for debugging and monitoring

### Performance Considerations

- **Connection Management**: SDK handles connection pooling automatically
- **Query Optimization**: Efficient queries with proper filtering and pagination
- **Response Parsing**: Automatic JSON parsing of database responses
- **Error Recovery**: Graceful handling of database connection issues

### SDK Features Used

- **Postgrest Integration**: Access database using REST API generated from schema
- **Authentication**: User authentication with email/password
- **Real-time**: Support for real-time database changes (future enhancement)
- **Storage**: File storage capabilities (future enhancement)
