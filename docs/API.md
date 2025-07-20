# Search Engine Service API Documentation

## Overview

The Search Engine Service API provides a comprehensive search functionality for video and text content from multiple providers. The API follows RESTful principles and includes advanced features like content scoring, filtering, and analytics.

**Base URL**: `http://localhost:8080`  
**API Version**: `v1`  
**Content-Type**: `application/json`

## Authentication

Currently, the API does not require authentication. However, rate limiting is applied:
- **Rate Limit**: 100 requests per minute per IP
- **Headers**: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`

## Endpoints

### Health Check

#### GET /health
Basic health check endpoint.

**Response**:
```json
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z",
  "service": "search-engine-service",
  "version": "2.0.0",
  "uptime": "1h 23m 45s"
}
```

#### GET /ready
Comprehensive readiness check including database and provider health.

**Response**:
```json
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z",
  "service": "search-engine-service",
  "version": "2.0.0",
  "uptime": "1h 23m 45s",
  "checks": {
    "database": {
      "status": "ok",
      "message": "Database is healthy",
      "details": {
        "open_connections": 5,
        "in_use": 2,
        "idle": 3
      }
    },
    "providers": {
      "status": "ok",
      "message": "All providers are healthy",
      "details": {
        "total_providers": 2,
        "healthy_providers": 2
      }
    }
  }
}
```

### Search API

#### GET /api/v1/search
Search for content with query parameters.

**Query Parameters**:
- `q` (string, optional): Search query
- `type` (string, optional): Content type filter (`video`, `text`, `all`)
- `page` (integer, optional): Page number (default: 1, min: 1)
- `limit` (integer, optional): Results per page (default: 10, max: 100)

**Example Request**:
```
GET /api/v1/search?q=golang&type=video&page=1&limit=10
```

**Response**:
```json
{
  "success": true,
  "data": {
    "contents": [
      {
        "id": 1,
        "provider": "json_provider",
        "provider_id": "video_001",
        "title": "Go Programlama Dili Temelleri",
        "description": "Bu videoda Go programlama dilinin temel kavramlarını öğreneceksiniz...",
        "url": "https://example.com/videos/go-basics",
        "type": "video",
        "views": 15420,
        "likes": 892,
        "duration": 1800,
        "tags": "go,golang,programlama,backend",
        "language": "tr",
        "published_at": "2024-01-15T10:30:00Z",
        "final_score": 85.2340,
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1,
    "has_next": false,
    "has_previous": false
  }
}
```

#### POST /api/v1/search/filters
Advanced search with filters.

**Request Body**:
```json
{
  "query": "programming",
  "content_type": "video",
  "min_score": 50.0,
  "max_score": 100.0,
  "min_views": 1000,
  "min_likes": 100,
  "tags": ["golang", "backend"],
  "language": "tr",
  "published_after": "2024-01-01T00:00:00Z",
  "published_before": "2024-12-31T23:59:59Z",
  "page": 1,
  "limit": 10,
  "sort_by": "final_score",
  "sort_order": "desc"
}
```

**Response**: Same as GET /api/v1/search

#### GET /api/v1/search/suggestions
Get search suggestions based on query.

**Query Parameters**:
- `q` (string, required): Partial search query

**Example Request**:
```
GET /api/v1/search/suggestions?q=gol
```

**Response**:
```json
{
  "success": true,
  "data": {
    "suggestions": [
      "golang",
      "go programming",
      "go tutorial",
      "go basics"
    ]
  }
}
```

### Content API

#### GET /api/v1/content/{id}
Get specific content by ID.

**Path Parameters**:
- `id` (integer, required): Content ID

**Example Request**:
```
GET /api/v1/content/1
```

**Response**:
```json
{
  "success": true,
  "data": {
    "content": {
      "id": 1,
      "provider": "json_provider",
      "provider_id": "video_001",
      "title": "Go Programlama Dili Temelleri",
      "description": "Bu videoda Go programlama dilinin temel kavramlarını öğreneceksiniz...",
      "url": "https://example.com/videos/go-basics",
      "type": "video",
      "views": 15420,
      "likes": 892,
      "duration": 1800,
      "tags": "go,golang,programlama,backend",
      "language": "tr",
      "published_at": "2024-01-15T10:30:00Z",
      "final_score": 85.2340,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    "score_breakdown": {
      "views_score": 25.5,
      "likes_score": 20.1,
      "duration_score": 15.2,
      "freshness_score": 12.3,
      "engagement_score": 12.1
    }
  }
}
```

#### GET /api/v1/content/popular
Get popular content based on final score.

**Query Parameters**:
- `limit` (integer, optional): Number of results (default: 10, max: 100)
- `type` (string, optional): Content type filter

**Example Request**:
```
GET /api/v1/content/popular?limit=5&type=video
```

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Go Programlama Dili Temelleri",
      "type": "video",
      "final_score": 85.2340,
      "views": 15420,
      "likes": 892
    }
  ]
}
```

#### GET /api/v1/content/trending
Get trending content based on recent activity.

**Query Parameters**:
- `limit` (integer, optional): Number of results (default: 10, max: 100)
- `period` (string, optional): Time period (`day`, `week`, `month`)

**Example Request**:
```
GET /api/v1/content/trending?limit=10&period=week
```

### Provider API

#### GET /api/v1/providers
Get list of available providers.

**Response**:
```json
{
  "success": true,
  "data": {
    "providers": [
      {
        "name": "json_provider",
        "url": "http://localhost:3001/api/videos",
        "status": "healthy",
        "last_fetch": "2024-01-15T10:30:00Z",
        "content_count": 5
      },
      {
        "name": "xml_provider",
        "url": "http://localhost:3001/api/articles",
        "status": "healthy",
        "last_fetch": "2024-01-15T10:30:00Z",
        "content_count": 5
      }
    ]
  }
}
```

#### POST /api/v1/providers/refresh
Manually refresh content from all providers.

**Response**:
```json
{
  "success": true,
  "data": {
    "message": "Content refresh completed",
    "total_fetched": 10,
    "providers_updated": 2,
    "duration": "2.5s"
  }
}
```

#### GET /api/v1/providers/stats
Get provider statistics.

**Response**:
```json
{
  "success": true,
  "data": {
    "total_content": 10,
    "providers": {
      "json_provider": {
        "content_count": 5,
        "avg_score": 82.1,
        "last_updated": "2024-01-15T10:30:00Z"
      },
      "xml_provider": {
        "content_count": 5,
        "avg_score": 84.9,
        "last_updated": "2024-01-15T10:30:00Z"
      }
    }
  }
}
```

#### GET /api/v1/providers/health
Check health status of all providers.

**Response**:
```json
{
  "success": true,
  "data": {
    "overall_status": "healthy",
    "providers": [
      {
        "name": "json_provider",
        "status": "healthy",
        "response_time": "150ms",
        "last_check": "2024-01-15T10:30:00Z"
      }
    ]
  }
}
```

### Analytics API

#### GET /api/v1/analytics/stats
Get overall service statistics.

**Response**:
```json
{
  "success": true,
  "data": {
    "total_content": 10,
    "content_by_type": {
      "video": 5,
      "text": 5
    },
    "avg_score": 83.5,
    "top_tags": [
      {"tag": "golang", "count": 3},
      {"tag": "programming", "count": 2}
    ],
    "recent_activity": {
      "last_24h": 2,
      "last_week": 8,
      "last_month": 10
    }
  }
}
```

#### GET /api/v1/analytics/trends
Get trending analytics data.

**Query Parameters**:
- `period` (string, optional): Time period (`day`, `week`, `month`)
- `metric` (string, optional): Metric to analyze (`views`, `likes`, `score`)

**Response**:
```json
{
  "success": true,
  "data": {
    "period": "week",
    "metric": "views",
    "trends": [
      {
        "date": "2024-01-15",
        "value": 15420,
        "change": 12.5
      }
    ]
  }
}
```

## Error Responses

All endpoints return consistent error responses:

### 400 Bad Request
```json
{
  "error": "Invalid content type. Must be 'video', 'text', or 'all'",
  "code": "INVALID_PARAMETER"
}
```

### 404 Not Found
```json
{
  "error": "Content not found",
  "code": "NOT_FOUND"
}
```

### 429 Too Many Requests
```json
{
  "error": "Rate limit exceeded",
  "retry_after": 60
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "code": "INTERNAL_ERROR"
}
```

## Content Scoring Algorithm

The service uses a sophisticated scoring algorithm that considers:

1. **Views Score** (25%): Based on view count
2. **Likes Score** (20%): Based on like count
3. **Duration Score** (15%): Based on content length
4. **Freshness Score** (15%): Based on publication date
5. **Engagement Score** (15%): Based on engagement metrics
6. **Quality Score** (10%): Based on content quality indicators

Final Score = (Views Score + Likes Score + Duration Score + Freshness Score + Engagement Score + Quality Score)

## Rate Limiting

- **Limit**: 100 requests per minute per IP
- **Headers**:
  - `X-RateLimit-Limit`: Request limit
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Reset time (Unix timestamp)

## Security

The API includes several security features:

- **CORS**: Configured for specific origins
- **Security Headers**: XSS protection, content type options, etc.
- **Input Sanitization**: Protection against XSS and SQL injection
- **Request ID**: Unique request tracking
- **Rate Limiting**: Protection against abuse

## Examples

### Search for Go Programming Videos
```bash
curl "http://localhost:8080/api/v1/search?q=golang&type=video&page=1&limit=5"
```

### Get Popular Content
```bash
curl "http://localhost:8080/api/v1/content/popular?limit=10"
```

### Advanced Search with Filters
```bash
curl -X POST "http://localhost:8080/api/v1/search/filters" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "programming",
    "content_type": "video",
    "min_score": 70.0,
    "tags": ["golang", "backend"],
    "page": 1,
    "limit": 10
  }'
```

### Health Check
```bash
curl "http://localhost:8080/health"
```

## SDK Examples

### JavaScript/Node.js
```javascript
const axios = require('axios');

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1'
});

// Search for content
const search = async (query, type = 'all') => {
  const response = await api.get('/search', {
    params: { q: query, type, page: 1, limit: 10 }
  });
  return response.data;
};

// Get popular content
const getPopular = async (limit = 10) => {
  const response = await api.get('/content/popular', {
    params: { limit }
  });
  return response.data;
};
```

### Python
```python
import requests

BASE_URL = "http://localhost:8080/api/v1"

def search_content(query, content_type="all", page=1, limit=10):
    response = requests.get(f"{BASE_URL}/search", params={
        "q": query,
        "type": content_type,
        "page": page,
        "limit": limit
    })
    return response.json()

def get_popular_content(limit=10):
    response = requests.get(f"{BASE_URL}/content/popular", params={
        "limit": limit
    })
    return response.json()
```

## Support

For API support and questions:
- **Documentation**: [GitHub Repository](https://github.com/your-username/search-engine-service)
- **Issues**: [GitHub Issues](https://github.com/your-username/search-engine-service/issues)
- **Email**: api-support@yourcompany.com 