# V0 API

## Submit a Test Request

```http
POST /v0/tests/run
```

Run a test againbst a target URL. Aeternum will validate the endpoint and respond
with a success message or an error.

=== "Bash"
    ```bash
    curl -X POST https://aeternum-api.onrender.com/v0/tests/run -d '{
      "base_url": "https://target-api.com",
      "endpoints": [
        { "path": "/status", "expected_code": 200 },
        { "path": "/data", "expected_code": 404 }
      ]
    }' -H "Content-Type: application/json"
    ```

=== "Python 3"
    ```python
    import requests
    url = "https://aeternum-api.onrender.com/v0/tests/run"
    headers = {"Content-Type": "application/json"}
    data = {
        "base_url": "https://target-api.com",
        "endpoints": [
            {"path": "/status", "expected_code": 200},
            {"path": "/data", "expected_code": 404}
        ]
    }
    response = requests.post(url, json=data, headers=headers)
    print(response.status_code)
    print(response.json())
    ```

=== "Javascript"
    ```javascript
    const fetch = require("node-fetch");
    const url = "https://aeternum-api.onrender.com/v0/tests/run";
    const data = {
        base_url: "https://target-api.com",
        endpoints: [
            { path: "/status", expected_code: 200 },
            { path: "/data", expected_code: 404 }
        ]
    };
    fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => console.log(data))
    .catch(error => console.error("Error:", error));
    ```

## Limitations

Currently we only support `GET` requests as the primary action for targets, we are
working on support for other request methods as well.
