# Use a imagem base do Go
FROM golang:1.24.2-alpine

# Define o diretório de trabalho
WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

# Copia os arquivos do projeto
COPY . .

# Compila o aplicativo
RUN go build -o motor-simulacao ./cmd/main.go

# Expõe a porta que o aplicativo irá escutar
EXPOSE 8080

# Comando para executar o aplicativo
CMD ["./motor-simulacao"]
