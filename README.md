go mod init github.com/Eursukkul/spd-fiber-booking-system

go get github.com/gofiber/fiber/v2
go get github.com/gofiber/swagger
go get github.com/swaggo/swag/cmd/swag@latest
go get github.com/stretchr/testify

mkdir -p cmd
mkdir -p dto
mkdir -p handler
mkdir -p middleware
mkdir -p router
mkdir -p repository
mkdir -p service
mkdir -p models
mkdir -p utils
mkdir -p config
mkdir -p tests
touch go.mod
touch go.sum
touch README.md

spd-fiber-booking-system/
├── cmd/
│ └── main.go
├── dto/
│ └── booking.go
├── handler/
│ └── booking_handler.go
├── middleware/
│ ├── auth.go
│ └── logging.go
├── router/
│ └── router.go
├── repository/
│ └── booking_repository.go
├── service/
│ └── booking_service.go
├── models/
│ └── booking.go
├── utils/
│ └── cache.go
│ └── hash.go
├── config/
│ └── config.go
├── tests/
│ └── booking_service_test.go
├── go.mod
├── go.sum
└── README.md
