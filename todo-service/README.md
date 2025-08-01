# Kratos Project Template

## Install Kratos
```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```
## Create a service
```
# Create a template project
kratos new server

cd server
# Add a proto template
kratos proto add api/server/server.proto
# Generate the proto code
kratos proto client api/server/server.proto
# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```
## Generate other auxiliary files by Makefile
```
# Download and update dependencies
make init
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
### Generate all files
```
make all
```

## Database Initialization
```bash
# Make sure MySQL is running
# Execute the initialization script
cd scripts
./init_db.sh
```

> Note: The default database configuration is root:123456@tcp(127.0.0.1:3306). You can modify the database connection string in configs/config.yaml if needed.

