# Stage 1: Build Lit WebComponents UI
FROM node:22-alpine AS frontend-builder
WORKDIR /app/ui
COPY ui/package*.json ./
RUN npm ci
COPY ui/ ./
RUN npm run build

# Stage 2: Compile Montage Go Binary with Embedded UI
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
COPY --from=frontend-builder /app/ui/dist ./ui/dist
RUN CGO_ENABLED=0 GOOS=linux go build -tags embedui -ldflags="-s -w" -o /bin/montage .

# Stage 3: Minimal Production Container
FROM alpine:3.21
RUN apk add --no-cache ca-certificates git github-cli graphviz

WORKDIR /app
COPY --from=backend-builder /bin/montage /usr/local/bin/montage

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["montage", "serve"]
