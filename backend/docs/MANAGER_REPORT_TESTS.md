# Testes Manuais - Funcionalidade de Envio de Relatório

## 1. Teste de Cadastro de Gerente

### Endpoint: `PUT /v1/users/:id`

**Request:**
```bash
curl -X PUT http://localhost:8081/v1/users/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer <seu-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "João Silva",
    "email": "joao@empresa.com",
    "active": true,
    "current_level": "P2",
    "manager_name": "Maria Santos",
    "manager_email": "maria.santos@empresa.com"
  }'
```

**Resposta Esperada:**
```json
{
  "message": "Usuário atualizado com sucesso"
}
```

## 2. Teste de Envio de Relatório (Sucesso)

### Endpoint: `POST /v1/report/send-to-manager`

**Pré-requisitos:**
- Usuário deve ter gerente cadastrado
- SMTP deve estar configurado corretamente

**Request:**
```bash
curl -X POST http://localhost:8081/v1/report/send-to-manager \
  -H "Authorization: Bearer <seu-token>" \
  -H "Content-Type: application/json"
```

**Resposta Esperada:**
```json
{
  "message": "Relatório enviado com sucesso para maria.santos@empresa.com",
  "status": "success"
}
```

## 3. Teste de Envio sem Gerente Cadastrado

**Request:**
```bash
curl -X POST http://localhost:8081/v1/report/send-to-manager \
  -H "Authorization: Bearer <seu-token>" \
  -H "Content-Type: application/json"
```

**Resposta Esperada (400 Bad Request):**
```json
{
  "error": "Gerente de engenharia não cadastrado",
  "details": "Por favor, cadastre o email do seu gerente nas configurações do perfil antes de enviar o relatório."
}
```

## 4. Teste de Email Recebido

Após enviar com sucesso, o gerente deve receber um email com:

**Assunto:**
```
Dossiê de Evolução de Carreira - João Silva
```

**Corpo (HTML):**
- Saudação personalizada
- Nome do colaborador e email
- Lista de conteúdos do relatório
- Data de geração

**Anexo:**
- `dossie_Joao_Silva.pdf` (ou similar)

## 5. Teste de Configuração SMTP

### Testar Conectividade SMTP

```bash
# Verificar se consegue conectar ao servidor SMTP
telnet smtp.gmail.com 587
```

### Logs Esperados

Ao enviar email, você deve ver nos logs do servidor:

```
INFO email sent successfully to=maria.santos@empresa.com subject="Dossiê de Evolução de Carreira - João Silva"
```

## 6. Casos de Erro

### Erro de SMTP não configurado

**Logs:**
```
WARN SMTP credentials not configured, email not sent to=maria.santos@empresa.com subject="..."
```

**Resposta:**
```json
{
  "error": "Erro ao enviar email",
  "details": "SMTP not configured"
}
```

### Erro de Autenticação SMTP

**Logs:**
```
ERROR SMTP authentication failed error="..." user="your-email@gmail.com"
```

**Resposta:**
```json
{
  "error": "Erro ao enviar email",
  "details": "SMTP authentication failed: ..."
}
```

## Checklist de Testes

- [ ] Cadastrar gerente com sucesso
- [ ] Atualizar gerente existente
- [ ] Remover gerente (deixar campos vazios)
- [ ] Tentar enviar sem gerente cadastrado
- [ ] Enviar relatório com sucesso
- [ ] Verificar email recebido
- [ ] Verificar anexo PDF
- [ ] Testar com SMTP inválido
- [ ] Testar com email de gerente inválido
- [ ] Verificar logs de sucesso e erro

## Dicas para Testes

1. **Use um email de teste real** para receber os PDFs e verificar a formatação
2. **Configure Gmail SMTP** seguindo: https://support.google.com/mail/answer/7126229
3. **Verifique a pasta de spam** caso não receba o email
4. **Use ferramentas como Mailtrap** para testes sem enviar emails reais

## Variáveis de Ambiente para Testes

```env
# Teste Local com Gmail
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=seu-email-teste@gmail.com
SMTP_PASSWORD=sua-senha-de-app
SMTP_FROM=noreply@leveluphub.com
SMTP_FROM_NAME=Level Up Hub [TEST]
```

```env
# Teste com Mailtrap (recomendado para desenvolvimento)
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=seu-username-mailtrap
SMTP_PASSWORD=sua-senha-mailtrap
SMTP_FROM=test@leveluphub.com
SMTP_FROM_NAME=Level Up Hub [DEV]
```
