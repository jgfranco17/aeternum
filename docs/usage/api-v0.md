# Using the API

## 📝 Submit a Test Request

```http
GET /results/{request_id}
```

Creates a test execution request. The request runs asynchronously, and you can
fetch results later using the `request_id`.

```bash
curl -X POST https://aeternum-api.onrender.com/tests/run -d '{
  "base_url": "https://target-api.com",
  "endpoints": [
    { "path": "/status", "expected_code": 200 },
    { "path": "/data", "expected_code": 404 }
  ]
}' -H "Content-Type: application/json"
```

Alternatively, you can define your request body in a separate JSON file.

```json title="request.json"
{
    "base_url": "https://example.com/api",
    "endpoints": [
        {
            "path": "/health",
            "expected_status": 200
        },
        {
            "path": "/users",
            "expected_status": 200
        }
    ]
}
```

You can then call the `curl` command with the following:

```bash
curl -X POST https://aeternum-api.onrender.com/tests/run \
    -d "@request.json" \
    -H "Content-Type: application/json"
```

## Limitations

Currently we only support `GET` requests as the primary action for targets, we are
working on support for other request methods as well.
