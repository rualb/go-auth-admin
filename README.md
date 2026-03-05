# go-auth-admin

`go-auth-admin` is the administrative backend for the authentication system, built with Go and the [Echo](https://echo.labstack.com/) framework.

## Features

- **User Management:** Administrative endpoints for managing users and roles.
- **Multi-Factor Authentication (MFA):** Integrates One-Time Passwords (OTP) using `github.com/pquerna/otp`.
- **Authentication:** JWT-based protection for administrative operations.
- **Database:** Interacts with the authentication database using GORM and PostgreSQL.
- **Metrics:** Exposes Prometheus metrics.

## Prerequisites

- Go 1.26+
- Python 3.x

## Build and Run

```sh
# Run tests
python Makefile.py test

# Run linter
python Makefile.py lint

# Build binary for Linux
python Makefile.py linux
```
