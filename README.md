# PartiuFit 💪

Uma API de acompanhamento fitness baseada em Go que permite aos usuários gerenciar seus treinos, rastrear exercícios e manter sua jornada fitness. Construída com arquitetura limpa usando Go, PostgreSQL e roteador Chi.

## 🚀 Funcionalidades

- **Gerenciamento de Usuários**: Registro de usuários, autenticação e gerenciamento de perfil
- **Rastreamento de Treinos**: Criar, ler, atualizar e deletar treinos
- **Gerenciamento de Exercícios**: Rastrear exercícios individuais dentro dos treinos
- **Autenticação**: Sistema de autenticação baseado em token
- **Migrações de Banco**: Gerenciamento automatizado de esquema do banco de dados
- **Monitoramento de Saúde**: Endpoints de verificação de integridade integrados
- **Hot Reload**: Ambiente de desenvolvimento com recarregamento automático

## 🏗️ Arquitetura do Projeto

```
partiuFit/
├── internal/
│   ├── app/                    # Inicialização e configuração da aplicação
│   ├── database/               # Conexão com banco de dados e utilitários
│   ├── handlers/               # Manipuladores de requisições HTTP
│   ├── middlewares/            # Middlewares HTTP (auth, tratamento de erros)
│   ├── requests/               # Estruturas de validação de requisições
│   ├── routes/                 # Definições de rotas da API
│   ├── store/                  # Camada de acesso aos dados
│   ├── tokens/                 # Gerenciamento de tokens
│   ├── utils/                  # Funções utilitárias
│   └── valueObjects/           # Objetos de valor do domínio
├── migrations/                 # Arquivos de migração do banco
├── config/                     # Arquivos de configuração
├── bin/                        # Binários compilados
└── tmp/                        # Arquivos temporários (desenvolvimento)
```

## 🔧 Pré-requisitos

Antes de executar esta aplicação, certifique-se de ter os seguintes itens instalados:

### Dependências Obrigatórias
- **Go 1.24+** - [Instalar Go](https://golang.org/doc/install)
- **PostgreSQL 14+** - [Instalar PostgreSQL](https://www.postgresql.org/download/)
- **Docker & Docker Compose** - [Instalar Docker](https://docs.docker.com/get-docker/)

### Ferramentas de Desenvolvimento (Recomendadas)
- **Air** - Hot reloading para aplicações Go
  ```bash
  go install github.com/air-verse/air@latest
  ```
- **golangci-lint** - Linter para Go
  ```bash
  # Linux/macOS
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
  
  # Ou via Homebrew (macOS)
  brew install golangci-lint
  ```
- **goimports** - Formatação de imports
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  ```

## ⚙️ Configuração do Ambiente

1. **Clone o repositório**
   ```bash
   git clone <repository-url>
   cd partiuFit
   ```

2. **Configuração do Ambiente**
   
   Crie um arquivo `.env` no diretório raiz:
   ```bash
   cp .env.testing .env
   ```
   
   Atualize o arquivo `.env` com sua configuração:
   ```env
   DATABASE_URL="postgres://postgres:postgres@localhost:5432/partiufit?sslmode=disable"
   PORT=8080
   APP_ENV=development
   ```

3. **Configuração do Banco de Dados**
   
   **Opção A: Usando Docker (Recomendado)**
   ```bash
   # Inicie o PostgreSQL com Docker Compose
   docker-compose up db -d
   ```
   
   **Opção B: PostgreSQL Local**
   ```bash
   # Crie o banco de dados manualmente
   createdb partiufit
   ```

## 🚀 Início Rápido

### Usando Docker (Recomendado)
```bash
# Inicie todos os serviços (banco + aplicação)
docker-compose up

# Ou execute em segundo plano
docker-compose up -d
```

### Desenvolvimento Local
```bash
# Instale as dependências
go mod download

# Execute as migrações do banco (automaticamente pela aplicação)
# Inicie a aplicação com hot reload
make run

# Ou compile e execute manualmente
make build
./bin/partiuFit
```

A API estará disponível em `http://localhost:8080`

## 📋 Comandos Make Disponíveis

```bash
make help      # Mostrar comandos disponíveis
make format    # Formatar código Go usando gofmt e goimports
make lint      # Executar golangci-lint
make run       # Executar com hot reload (usa Air)
make build     # Compilar o binário da aplicação
make test      # Executar todos os testes
make clean     # Limpar artefatos de build
```

## 🔌 Endpoints da API

### Verificação de Integridade
- `GET /health` - Verificar status de integridade da aplicação

### Autenticação
- `POST /tokens` - Gerar token de autenticação (login)

### Gerenciamento de Usuários
- `POST /users` - Registrar novo usuário
- `PUT /users` - Atualizar perfil do usuário (requer autenticação)

### Gerenciamento de Treinos (Autenticação Obrigatória)
- `GET /workouts` - Obter todos os treinos do usuário
- `POST /workouts` - Criar novo treino
- `GET /workouts/{id}` - Obter treino específico por ID
- `PUT /workouts/{id}` - Atualizar treino específico
- `DELETE /workouts/{id}` - Deletar treino específico

## 🗄️ Esquema do Banco de Dados

A aplicação usa PostgreSQL com as seguintes entidades principais:

- **Users**: Contas e perfis de usuários
- **Workouts**: Sessões de treino
- **Workout_Entries**: Exercícios individuais dentro dos treinos
- **Tokens**: Tokens de autenticação

As migrações são aplicadas automaticamente na inicialização da aplicação.

## 🧪 Testes

Execute a suíte de testes:
```bash
# Execute todos os testes
make test

# Execute testes com saída detalhada
go test -v ./...

# Execute testes de pacote específico
go test -v ./internal/store
```

## 🔧 Fluxo de Desenvolvimento

1. **Inicie o ambiente de desenvolvimento**
   ```bash
   docker-compose up db -d  # Inicie o banco de dados
   make run                 # Inicie a aplicação com hot reload
   ```

2. **Formatação e linting do código**
   ```bash
   make format  # Formate o código
   make lint    # Execute o linter
   ```

3. **Executando testes**
   ```bash
   make test
   ```

## 📊 Dependências Principais

- **Framework Web**: [Chi](https://github.com/go-chi/chi) - Roteador HTTP leve
- **Banco de Dados**: [pgx](https://github.com/jackc/pgx) - Driver PostgreSQL
- **Migrações**: [Goose](https://github.com/pressly/goose) - Ferramenta de migração de banco
- **Validação**: [validator](https://github.com/go-playground/validator) - Validação de structs
- **Logging**: [Zap](https://github.com/uber-go/zap) - Logging estruturado
- **Ambiente**: [godotenv](https://github.com/joho/godotenv) - Carregamento de variáveis de ambiente
- **Senhas**: [bcrypt](https://golang.org/x/crypto/bcrypt) - Hash de senhas
- **Testes**: [Testify](https://github.com/stretchr/testify) - Kit de ferramentas para testes

## 🌐 Variáveis de Ambiente

| Variável | Descrição | Padrão | Obrigatório |
|----------|-------------|---------|----------|
| `DATABASE_URL` | String de conexão PostgreSQL | - | ✅ |
| `PORT` | Porta do servidor | `8080` | ✅ |
| `APP_ENV` | Ambiente da aplicação | `development` | ❌ |

## 🚀 Deploy de Produção

1. **Compile a aplicação**
   ```bash
   make build
   ```

2. **Defina as variáveis de ambiente**
   ```bash
   export APP_ENV=production
   export DATABASE_URL="sua-url-do-banco-producao"
   export PORT=8080
   ```

3. **Execute o binário**
   ```bash
   ./bin/partiuFit
   ```

## 🤝 Contribuindo

1. Faça um fork do repositório
2. Crie uma branch para sua feature (`git checkout -b feature/funcionalidade-incrivel`)
3. Faça suas alterações
4. Execute os testes (`make test`)
5. Formate o código (`make format`)
6. Execute o linter (`make lint`)
7. Commit suas alterações (`git commit -m 'Add funcionalidade incrível'`)
8. Envie para a branch (`git push origin feature/funcionalidade-incrivel`)
9. Abra um Pull Request

## 📝 Status do Projeto

Este é um projeto ativo de API de acompanhamento fitness construído com práticas modernas de Go. A aplicação segue princípios de arquitetura limpa com clara separação de responsabilidades.

## 🔒 Segurança

- Senhas são hash usando bcrypt
- Autenticação baseada em token
- Validação de entrada em todos os endpoints
- Prevenção de injeção SQL através de consultas parametrizadas

## 📞 Suporte

Se você encontrar algum problema ou tiver dúvidas, por favor abra uma issue no repositório.

---

**Feito com ❤️ usando Go**