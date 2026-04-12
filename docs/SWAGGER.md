# API Documentation - Swagger/OpenAPI

## 📚 Documentação Interativa

A API Level Up Hub possui documentação completa no padrão OpenAPI 3.0 com interface interativa Swagger UI.

### 🌐 Acessar Documentação

```
http://localhost:8081/swagger/index.html
```

## 🚀 Como Usar

### 1. **Acessar a UI**

Inicie o servidor e acesse no navegador:
```bash
./api
# Abra: http://localhost:8081/swagger/index.html
```

### 2. **Testar Endpoints**

#### **Sem Autenticação** (Public)
- `POST /v1/login` - Fazer login
- `POST /v1/register` - Criar conta

#### **Com Autenticação** (Protected)
1. Faça login em `/v1/login`
2. Copie o token retornado
3. Clique em "Authorize" (cadeado no topo)
4. Cole `Bearer SEU_TOKEN`
5. Teste os endpoints protegidos

### 3. **Exemplo de Fluxo**

```bash
# 1. Registrar usuário
POST /v1/register
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "senha123",
  "active": true
}

# 2. Fazer login
POST /v1/login
{
  "email": "john@example.com",
  "password": "senha123"
}

# 3. Copiar token da resposta
# 4. Clicar em "Authorize" e colar: Bearer eyJhbGc...
# 5. Agora você pode testar endpoints protegidos!
```

---

## 🛠️ Regenerar Documentação

### Quando Regenerar

Sempre que você:
- Adicionar novos endpoints
- Modificar parâmetros ou respostas
- Atualizar descrições

### Como Regenerar

```bash
# Executar o comando swag
swag init -g cmd/api/main.go --output docs

# Ou use o Makefile (se configurado)
make swagger
```

---

## 📝 Como Documentar Endpoints

### Estrutura Básica

```go
// HandlerName godoc
// @Summary      Resumo curto do endpoint
// @Description  Descrição detalhada
// @Tags         categoria
// @Accept       json
// @Produce      json
// @Param        nome  local  tipo  required  "descrição"
// @Success      200   {object}  TipoResposta  "descrição"
// @Failure      400   {object}  TipoErro  "descrição"
// @Router       /rota [método]
func (h *Handler) HandlerName(c *gin.Context) {
    // código
}
```

### Exemplo Completo

```go
// CreateActivity godoc
// @Summary      Criar nova atividade
// @Description  Cria uma atividade associada ao career ladder
// @Tags         activities
// @Accept       json
// @Produce      json
// @Param        activity  body      CreateActivityDTO  true  "Dados da atividade"
// @Security     BearerAuth
// @Success      201       {object}  map[string]interface{}  "Atividade criada"
// @Failure      400       {object}  map[string]interface{}  "Dados inválidos"
// @Failure      401       {object}  map[string]interface{}  "Não autorizado"
// @Failure      500       {object}  map[string]interface{}  "Erro interno"
// @Router       /activities [post]
func (h *ActivityHandler) Create(c *gin.Context) {
    // implementação
}
```

### Tipos de Parâmetros

| Local | Descrição | Exemplo |
|-------|-----------|---------|
| `path` | Na URL | `/users/{id}` |
| `query` | Query string | `/users?page=1` |
| `body` | No corpo | JSON payload |
| `header` | No header | `Authorization` |

### Tags Comuns

```go
// @Summary      - Título curto (obrigatório)
// @Description  - Descrição detalhada
// @Tags         - Categoria do endpoint
// @Accept       - Content-Type aceito (json, xml, etc)
// @Produce      - Content-Type da resposta
// @Param        - Parâmetro do endpoint
// @Success      - Resposta de sucesso
// @Failure      - Resposta de erro
// @Security     - Esquema de segurança
// @Router       - Rota e método HTTP
```

---

## 🔐 Autenticação

A API usa **Bearer Token (JWT)**:

```yaml
securityDefinitions:
  BearerAuth:
    type: apiKey
    in: header
    name: Authorization
    description: "Digite 'Bearer' seguido do token JWT"
```

**Uso no Swagger UI:**
1. Clique no botão **"Authorize"** (cadeado)
2. Digite: `Bearer SEU_TOKEN_AQUI`
3. Clique em "Authorize"
4. Feche o modal

---

## 📦 Arquivos Gerados

```
docs/
├── docs.go        # Código Go gerado
├── swagger.json   # Especificação OpenAPI em JSON
└── swagger.yaml   # Especificação OpenAPI em YAML
```

**Importante:**
- `docs.go` é importado automaticamente
- Os arquivos são regerados a cada `swag init`
- Adicione ao `.gitignore` se preferir gerar na build

---

## 🎨 Personalização

### Informações Gerais (main.go)

```go
// @title           Seu Título
// @version         1.0
// @description     Descrição da API

// @contact.name   Nome do Suporte
// @contact.email  suporte@exemplo.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8081
// @BasePath  /v1
```

### Diferentes Ambientes

```go
// Desenvolvimento
// @host      localhost:8081

// Staging
// @host      staging-api.exemplo.com

// Produção
// @host      api.exemplo.com
```

---

## 🔍 Validação

### Validar Especificação OpenAPI

```bash
# Usar validador online
# Acesse: https://editor.swagger.io/
# Copie o conteúdo de docs/swagger.yaml

# Ou use CLI
npm install -g @apidevtools/swagger-cli
swagger-cli validate docs/swagger.yaml
```

---

## 📚 Referências

- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [swaggo/swag](https://github.com/swaggo/swag)
- [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)
- [Declarative Comments Format](https://github.com/swaggo/swag#declarative-comments-format)

---

## 🚧 Próximos Passos

### Melhorias Futuras

1. **Adicionar exemplos de request/response**
   ```go
   // @Success 200 {object} User "exemplo: {\"id\": \"123\", \"name\": \"John\"}"
   ```

2. **Documentar todos os endpoints**
   - Adicionar comentários em todos os handlers
   - Regenerar documentação

3. **Versionamento da API**
   - Manter documentação de versões anteriores
   - Endpoint `/swagger/v1/` e `/swagger/v2/`

4. **CI/CD**
   - Gerar docs automaticamente na build
   - Publicar em GitHub Pages ou similar

5. **Mock Server**
   - Usar Prism para mock baseado no OpenAPI
   ```bash
   npm install -g @stoplight/prism-cli
   prism mock docs/swagger.yaml
   ```
