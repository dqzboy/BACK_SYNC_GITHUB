FROM node:20-alpine AS frontend
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.26-alpine AS backend
WORKDIR /app
COPY server/go.mod server/go.sum ./server/
WORKDIR /app/server
RUN go mod download
COPY server/ ./
COPY --from=frontend /app/web/dist ./webroot/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/git-backup-server .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates git openssh-client tzdata \
    && adduser -D -u 1000 appuser
WORKDIR /app
COPY --from=backend /out/git-backup-server /app/git-backup-server
# 数据目录（SQLite 会在启动时自动建库），建议挂载卷持久化
RUN mkdir -p /app/data && chown -R appuser:appuser /app/data
USER appuser
EXPOSE 8080
VOLUME ["/app/data"]
ENTRYPOINT ["/app/git-backup-server"]