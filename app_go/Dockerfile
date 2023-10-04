# Latest golang image on apline linux
FROM golang:1.17-alpine

# Work directory
WORKDIR /app_go

# Installing dependencies
COPY ./app_go/go.mod ./app_go/go.sum ./
RUN go mod download

# Copying all the files
COPY . .

# Starting our application
CMD ["go", "run", "./app_go/main.go"]

# Exposing server port
EXPOSE 3000