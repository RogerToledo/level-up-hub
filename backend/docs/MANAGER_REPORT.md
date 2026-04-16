# Funcionalidade de Envio de Relatório ao Gerente de Engenharia

## Resumo

Esta funcionalidade permite que os colaboradores cadastrem as informações do seu Gerente de Engenharia (GE) e enviem o dossiê de evolução de carreira diretamente por email.

## Componentes Implementados

### 1. Migration de Banco de Dados
- **Arquivo**: `db/migrations/008_add_manager_fields.sql`
- Adiciona campos `manager_name` e `manager_email` na tabela `users`
- Inclui índice para otimizar consultas por email do gerente

### 2. Queries SQLC Atualizadas
- **Arquivo**: `db/queries/user.sql`
- `FindUserByID`: Agora retorna os campos de gerente
- `UpdateUser`: Permite atualizar nome e email do gerente

### 3. Serviço de Email
- **Arquivo**: `internal/email/service.go`
- Serviço completo de envio de emails via SMTP
- Suporte a anexos (PDF)
- Template HTML profissional para o email
- Função específica `SendReportToManager` para enviar dossiês

### 4. Configuração
- **Arquivo**: `config/config.go`
- Novas variáveis de ambiente para SMTP:
  - `SMTP_HOST`: Servidor SMTP (padrão: smtp.gmail.com)
  - `SMTP_PORT`: Porta SMTP (padrão: 587)
  - `SMTP_USER`: Usuário para autenticação
  - `SMTP_PASSWORD`: Senha para autenticação
  - `SMTP_FROM`: Email de origem
  - `SMTP_FROM_NAME`: Nome do remetente (padrão: Level Up Hub)

### 5. API Endpoint
- **Rota**: `POST /v1/report/send-to-manager`
- **Autenticação**: Requerida (Bearer Token)
- **Funcionalidade**:
  1. Verifica se o usuário tem gerente cadastrado
  2. Gera o PDF do dossiê completo
  3. Envia por email para o gerente
  4. Retorna confirmação de envio

### 6. DTOs Atualizados
- **Arquivo**: `internal/account/dto.go`
- `UpdateUserRequest`: Campos opcionais `manager_name` e `manager_email`
- `UserResponse`: Inclui informações do gerente nas respostas

## Como Usar

### 1. Configurar SMTP

Adicione as variáveis de ambiente no arquivo `.env`:

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=seu-email@gmail.com
SMTP_PASSWORD=sua-senha-de-app
SMTP_FROM=noreply@leveluphub.com
SMTP_FROM_NAME=Level Up Hub
```

**Nota para Gmail**: É necessário usar uma "Senha de App" em vez da senha normal. Acesse: https://myaccount.google.com/apppasswords

### 2. Executar a Migration

```bash
# Execute a migration 008 para adicionar os campos de gerente
# (O método depende da ferramenta de migration usada no projeto)
```

### 3. Cadastrar Gerente

Atualize o perfil do usuário incluindo as informações do gerente:

```bash
PUT /v1/users/:id
{
  "username": "João Silva",
  "email": "joao@empresa.com",
  "active": true,
  "current_level": "P2",
  "manager_name": "Maria Santos",
  "manager_email": "maria.santos@empresa.com"
}
```

### 4. Enviar Relatório

```bash
POST /v1/report/send-to-manager
Authorization: Bearer <token>
```

Resposta de sucesso:
```json
{
  "message": "Relatório enviado com sucesso para maria.santos@empresa.com",
  "status": "success"
}
```

## Validações

- O sistema verifica se o gerente está cadastrado antes de tentar enviar
- Validação de formato de email
- Mensagens de erro claras para o usuário

## Email Enviado

O gerente receberá um email profissional contendo:
- Saudação personalizada com o nome do gerente
- Informações sobre o colaborador
- Descrição do conteúdo do relatório
- PDF anexado com o dossiê completo
- Data de geração

## Segurança

- Autenticação SMTP com TLS
- Validação de campos antes do envio
- Logs estruturados de todas as operações de email
- Tratamento de erros robusto

## Frontend (Sugestão de Implementação)

No frontend, você pode adicionar:

1. **Página de Configurações do Perfil**:
   - Campos para cadastrar nome e email do gerente
   - Botão "Salvar Alterações"

2. **Página de Relatórios**:
   - Botão "Baixar PDF" (já existente)
   - Botão "Enviar para Gerente" (novo)
   - Tooltip informando se o gerente está cadastrado
   - Modal de confirmação antes de enviar

Exemplo de código React:
```typescript
const sendToManager = async () => {
  try {
    const response = await api.post('/report/send-to-manager');
    toast.success(response.data.message);
  } catch (error) {
    if (error.response?.status === 400) {
      toast.error(error.response.data.message);
      // Redirecionar para página de configurações
    } else {
      toast.error('Erro ao enviar relatório');
    }
  }
};
```

## Troubleshooting

### Email não enviado
- Verifique as credenciais SMTP no `.env`
- Para Gmail, certifique-se de usar Senha de App
- Verifique os logs do servidor para detalhes do erro

### Gerente não cadastrado
- O usuário precisa primeiro atualizar seu perfil
- Campos `manager_name` e `manager_email` são opcionais, mas necessários para envio

## Próximas Melhorias Possíveis

1. Histórico de envios (quando foi enviado, para quem)
2. Agendamento de envios periódicos
3. Template de email customizável
4. Envio para múltiplos destinatários
5. Confirmação de leitura do email

