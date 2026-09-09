module github.com/S7venKing/ticket-box-go/services/identity-service

go 1.26.0

require (
	github.com/S7venKing/ticket-box-go/gen v0.0.0
	github.com/go-sql-driver/mysql v1.10.1
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.12.1
	golang.org/x/crypto v0.56.0
	google.golang.org/grpc v1.83.2
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

// gen is an unpublished module inside this repo. The replace makes
// `go mod tidy` and GOWORK=off builds (Docker) resolve it without network;
// go.work at the root points to the same directory for workspace mode.
replace github.com/S7venKing/ticket-box-go/gen => ../../gen
