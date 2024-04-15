FROM golang:1.18-alpine

WORKDIR /app

# 依存関係のコピーとインストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
RUN go build -o backend .

EXPOSE 8080

CMD ["./backend"]
