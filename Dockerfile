FROM golang:1.24-alpine AS base

FROM base AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go tool templ generate
RUN go build ./cmd/chat

FROM base AS run
COPY --from=build /app/chat /app/chat
EXPOSE 8080
CMD ["/app/chat"]