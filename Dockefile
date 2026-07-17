FROM golang:1.26-alpine

WORKDIR /app

# Copia as dependências da raiz do monorepo
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o resto do projeto (incluindo a pasta backend)
COPY . .

# Compila o Go apontando exatamente para onde o seu main está
RUN go build -o out backend/cmd/main.go

EXPOSE 8080

CMD ["./out"]