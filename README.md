# Aeternum API

_Ensuring your APIs stand the test of time._

![STATUS](https://img.shields.io/badge/status-active-brightgreen?style=for-the-badge)
![LICENSE](https://img.shields.io/badge/license-BSD3-blue?style=for-the-badge)

---

## Introduction

### About

Aeternum API is a lightweight SaaS tool designed to continuously test API availability
and correctness. With a simple configuration, users can define their API endpoints,
expected HTTP status codes, and response times, ensuring their services remain operational
without manual monitoring.

Built using Golang, Aeternum API runs without requiring a database or external services;
this makes it lightweight, fast, and easy to deploy.

### API

Check out the [documentation page](https://jgfranco17.github.io/aeternum-api/) for more
information about using the API.

### CLI

To run Aeternum locally, we also provide a CLI tool. This allows you to run your API tests
llocally (from your local machine) or remotely (via request to the API).

To download the CLI, an install script has been provided.

```bash
wget -O - https://raw.githubusercontent.com/jgfranco17/aeternum/main/install.sh | bash
```

## License

This project is licensed under the BSD-3 License. See the LICENSE file for more details.
