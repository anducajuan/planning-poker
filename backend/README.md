# Planning Poker Backend

Backend da aplicação de Planning Poker desenvolvido em Go.

## Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` baseado no exemplo abaixo:

```env
# Configurações do Banco de Dados
DATABASE_URL=postgres://admin:admin@localhost:5432/database?sslmode=disable

# Porta do servidor
PORT=8080

# Modo de log (dev ou prod)
LOG_MODE=dev

# Origens permitidas para CORS (* para todas)
CORS_ORIGINS=*

# Configurações do ambiente
APP_ENV=dev
```

## Executando o Projeto

### Pré-requisitos

- Go 1.25.1+
- PostgreSQL

### Instalação

1. Instale as dependências:
```bash
go mod download
```

2. Configure o banco de dados executando as migrações:
```bash
cd migrations
./run.sh
```

3. Execute a aplicação:

Usando o Golang:
```bash
go run .
```

Usando air:

Para instalar o air, na raíz do projeto backend, execute:
```bash
go install github.com/air-verse/air@latest
```
Se estiver no linux, adicione esta linha no final do seu ~/.bash ou ~/.zshrc:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```
E com isso, após recarregar o seu shell, inicie a aplicação com o comando `air`


## API Endpoints

### Sessions

- `GET /sessions` - Lista todas as sessões
- `POST /sessions` - Cria uma nova sessão
- `DELETE /sessions/{id}` - Deleta uma sessão

### WebSocket

- `GET /ws` - Conexão WebSocket para atualizações em tempo real

## Recursos

- **CORS**: Configurado para permitir requisições do frontend
- **Logging**: Middleware de log com formato configurável (dev/prod)
- **WebSocket**: Sistema de broadcast em tempo real
- **Tratamento de Erros**: Respostas padronizadas com logs estruturados
- **Configuração**: Gerenciamento via variáveis de ambiente

## Tecnologias Utilizadas

- **Go**: Linguagem principal
- **Gorilla Mux**: Roteamento HTTP
- **Gorilla WebSocket**: Implementação WebSocket
- **pgx**: Driver PostgreSQL
- **PostgreSQL**: Banco de dados
