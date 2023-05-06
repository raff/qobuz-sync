FROM golang:1.20-alpine AS builder
# --- GO BASE ---
RUN apk update && \
	apk add --no-cache \
	"git" \
    "ca-certificates"
WORKDIR /app
COPY  . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o qobuz-sync ./cmd && chmod o+x qobuz-sync

# --- RUNTIME ---
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/qobuz-sync /

ENTRYPOINT [ "/qobuz-sync"]
