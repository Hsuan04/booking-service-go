# 第一階段：編譯階段 (Builder)
FROM golang:1.22-alpine AS builder

# 安裝憑證與基本工具（確保之後能正常連線到 Supabase 等外部服務）
RUN apk add --no-cache ca-certificates && update-ca-certificates

# 設定工作目錄
WORKDIR /app

# 複製 go.mod 與 go.sum，利用 Docker 快取機制加速未來建置
COPY go.mod go.sum ./
RUN go mod download

# 複製所有程式碼
COPY . .

# 編譯 Go 程式
# CGO_ENABLED=0 為了確保編譯出來的二進位檔能獨立運行
# -ldflags="-w -s" 用於移除除錯資訊，進一步縮小檔案體積
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api

# 第二階段：運行階段 (Runtime)
FROM alpine:latest

# 確保運行環境也有安全憑證
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 從 builder 階段把編譯好的二進位檔複製過來
COPY --from=builder /app/main .

# 宣告預計使用的 Port (Cloud Run 預設監聽 8080)
EXPOSE 8080

# 執行 API 服務
CMD ["./main"]