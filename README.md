# SayMore Server

🚀 **SayMore Server** is the backend service for the [SayMore](https://github.com/your-org/saymore) English learning platform.  
It provides RESTful APIs to support students, teachers, and administrators in an interactive, modular language training ecosystem.

---

## 🔧 Features

- 🗣️ Manage student profiles, bookings, and lesson records
- 🧑‍🏫 Teacher feedback and lesson scheduling
- 🧾 Admin control panel: course rules, notifications, and reports
- 🔐 JWT-based authentication with token refresh
- 💬 Aliyun SMS for user verification and notifications
- ⏱️ Cron jobs for auto-reminders and class monitoring
- ☁️ Aliyun OSS for media upload and playback
- ⚡ High-performance setup using Go and Redis

---

## 🧱 Tech Stack

- **Language**: Golang 1.23
- **Database**: MySQL 8.x
- **Cache**: Redis
- **Cloud**: Aliyun OSS & SMS
- **Auth**: JWT
- **Containerization**: Docker / Docker Compose
- **Config**: TOML (`config.toml`)

---

## Project Introduction

SayMore Server is a backend service for an online education platform, providing course management, user management, payment, and other functionalities. The project is developed in Go language, based on the Gin framework, using MySQL as the primary database and Redis for caching and session storage.

## Project Structure

```
.
├── cmd/                # Command line entry
│   └── server/        # Server startup entry
├── config/            # Configuration files
├── internal/          # Internal packages
│   ├── app/          # Application layer
│   │   ├── controllers/  # Controllers
│   │   ├── models/      # Data models
│   │   └── services/    # Service layer
│   └── pkg/          # Public packages
│       ├── initialization/  # Initialization
│       ├── token/          # Token management
│       └── wechat/         # WeChat related
├── routes/            # Route configuration
├── utils/             # Utility functions
├── config.toml.example # Configuration file example
├── go.mod            # Go module file
└── README.md         # Project documentation
```

## Requirements

- Go 1.21 or higher
- MySQL 5.7 or higher
- Redis 6.0 or higher
- Docker (optional, for containerized deployment)

## Configuration

1. Copy the configuration file:
   ```bash
   cp config/config.toml.example config/config.toml
   ```

2. Modify the configuration items in `config/config.toml`:

   ```toml
   [app]
   port = "8080"
   mode = "debug"  # debug or release
   log_level = "debug"

   [mysql]
   host = "localhost"
   port = "3306"
   user = "root"
   password = ""  # Database password
   dbname = "saymore"
   max_open_conns = 1000
   max_idle_conns = 100
   max_life_time = 1

   [redis]
   address = "127.0.0.1:6379"
   password = ""  # Redis password
   db = 0
   prefix = "say-more"

   [access_token]
   token_expire = 168  # JWT token expiration time (hours)
   token_refresh = 24  # JWT token refresh interval (hours)

   [ali_oss]
   accesskeyid = ""     # Aliyun OSS AccessKey
   accesskeysecret = "" # Aliyun OSS Secret
   bucketname = ""      # OSS Bucket name
   endpoint = ""        # OSS access domain
   bucketurl = ""       # OSS Bucket URL

   [ali_textmsg]
   enable = true
   access_key_id = ""     # Aliyun SMS AccessKey
   access_key_secret = "" # Aliyun SMS Secret
   alarm_template_code = ""
   endpoint = ""
   identity_template_code = ""
   sign_name = ""
   ttl = 1

   [wechat]
   app_id = ""     # WeChat Mini Program AppID
   app_secret = "" # WeChat Mini Program Secret
   ```

## Environment Variables

The project supports overriding configuration values through environment variables:

- `MYSQL_HOST`: MySQL host address
- `MYSQL_PORT`: MySQL port
- `MYSQL_USER`: MySQL username
- `MYSQL_PASSWORD`: MySQL password
- `MYSQL_DBNAME`: MySQL database name
- `REDIS_ADDRESS`: Redis address
- `REDIS_PASSWORD`: Redis password
- `ALI_OSS_ACCESSKEYID`: Aliyun OSS AccessKey
- `ALI_OSS_ACCESSKEYSECRET`: Aliyun OSS Secret
- `WECHAT_APPID`: WeChat Mini Program AppID
- `WECHAT_APPSECRET`: WeChat Mini Program Secret

## Running Steps

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Create database:
   ```sql
   CREATE DATABASE saymore CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

3. Run the server:
   ```bash
   # Using default configuration file
   go run cmd/server/main.go

   # Specify configuration file
   go run cmd/server/main.go -config=/path/to/config.toml
   ```

## Docker Deployment

1. Build the image:
   ```bash
   docker build -t saymore-server .
   ```

2. Run the container:
   ```bash
   docker run -d \
     -p 8080:8080 \
     -v /path/to/config.toml:/app/config/config.toml \
     saymore-server
   ```

## API Documentation

API documentation is generated using Swagger, accessible at: `http://localhost:8080/swagger/index.html`

Main functional modules:

- User Authentication
- Course Management
- Order Management
- Payment Interface
- File Upload
- SMS Service

## Development Guide

1. Code Standards
   - Follow Go official code standards
   - Use gofmt for code formatting
   - Run go vet and golangci-lint before committing

2. Branch Management
   - main: Production branch
   - develop: Development branch
   - feature/*: Feature branches
   - hotfix/*: Hotfix branches

3. Commit Standards
   - feat: New feature
   - fix: Bug fix
   - docs: Documentation updates
   - style: Code style changes
   - refactor: Code refactoring
   - test: Test related
   - chore: Build process or auxiliary tool changes

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
