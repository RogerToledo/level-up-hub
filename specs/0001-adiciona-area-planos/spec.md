# Spec: Planos

## 1. Visão Geral e Contexto
- **Funcionalidade:** Área para cadastrar planos futuro.
- **Objetivo:** Mapear ideias que futuramente podem virar Iniciativas ou Tasks
- **Escopo:** Criação de Planos

## 2. Histórias de Usuário & Cenários

### Cenário 1: Abrir tela de Planos
- **Dado que** o usuários está no Dash
- **Quando** clica em Planos
- **Então** deve abrir a tela de Planos
- **E** exibir os planos salvos

### Cenário 2: Abrir modal de cadastro
- **Dado que** o usuário esteja na tela de planos
- **Quando** clica em novo
- **Então** devesse abri o modal de cadastro do Plano

### Cenário 3: Cadastro valido
- **Dado que** o usuário esteja no modal de cadastro
- **E** tenha preenchido corretamente todos os campo
- **Quando** clica em Cadastrar
- **Então** o sistema deve cadastrar o Plano e exibir mensagem que cadastrou com sucesso

### Cenário 4: Cadastrar plano existente
- **Dado que** o usuário esteja no modal de cadastro
- **E** tenha preenchido um plano que já existe
- **Quando** clica em Cadastrar
- **Então** o sistema deve exibir mensagem informando que o plano já existe

### Cenário 5: Cadastro com campos invalidos
- **Dado que** o usuário esteja no modal de cadastro
- **E** tenha deixado campos obrigatórios em branco
- **Quando** clica em Cadastrar
- **Então** o sistema deve exibir mensagem informando os campos obrigatórios devem ser preenchidos

### Cenário 6: Abrindo Modal de update
- **Dado que** o usuário na tela de Planos
- **E** clica em no botão alterar
- **Então** o sistema deve abrir o modal com os campos preenchidos 

### Cenário 7: Atualizando Plano
- **Dado que** o usuário no modal de alteração
- **E** tenha alterado alguma informação
- **Quando** clica em Atualizar
- **Então** o sistema deve atualizar o plano e exibir mensagem dizendo que o plano foi alterado com sucesso

### Cenário 8: Deletando Plano
- **Dado que** o usuário está na tela de planos
- **Quando** clica em no icone de deletar
- **Então** o sistema deve deletar o plano e exibir mensagem dizendo que o plano deletado com sucesso

### Cenário 8: Visualizar Plano
- **Dado que** o usuário está na tela de planos
- **Quando** clica em no icone de visualizar
- **Então** o sistema abrir o modal de visualização
- **E** todos os campos devem estar bloqueados


## 3. Requisitos Funcionais
#### **BACKEND**
- [ ] **RF-01:** Deve ter um CRUD os campos obrigatórios: Título e Descrição, opcional iniciativa e nível.
- [ ] **RF-02:** A API deve retornar respostas paginadas por padrão.

#### **FRONTEND**
- [ ] **RF-03:** Deve validar o formato dos inputs antes de enviar a requisição.
- [ ] **RF-04:** Deve exibir o Planos abaixo do Relatório
- [ ] **RF-05:** Deve seguir o lauout das outras telas 
- [ ] **RF-06:** A API deve retornar respostas paginadas por padrão.

## 4. Requisitos Não-Funcionais & Regras
- **Segurança:** A rota exige autenticação via Middleware JWT.
- **Concorrência:** Tratamento seguro de goroutines para evitar race conditions no Go.
- **UX/UI:** Componente React com feedback visual de carregamento (*loading state*).

## 5. Casos de Borda Mapeados
- O que acontece se o token JWT expirar no meio da navegação?
- Como o frontend se comporta caso o backend retorne status `500 Internal Server Error`?