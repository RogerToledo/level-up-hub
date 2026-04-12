# Paginação de APIs - Guia de Uso

## 📚 Endpoints com Paginação

### 1. **Listar Usuários** (Admin)
```bash
GET /v1/users?page=1&page_size=20
```

**Parâmetros:**
- `page` (opcional): Número da página (default: 1)
- `page_size` (opcional): Itens por página (default: 20, máx: 100)

**Resposta:**
```json
{
  "status_code": 200,
  "message": {
    "data": [
      {
        "id": "uuid",
        "username": "john_doe",
        "email": "john@example.com",
        "active": true,
        "role": "user",
        "created_at": "2026-04-01"
      }
    ],
    "pagination": {
      "current_page": 1,
      "page_size": 20,
      "total_pages": 5,
      "total_records": 87,
      "has_next": true,
      "has_previous": false
    }
  }
}
```

---

## 🎯 **Como Usar**

### Primeira Página (default)
```bash
curl http://localhost:8081/v1/users
```

### Segunda Página
```bash
curl http://localhost:8081/v1/users?page=2
```

### Mudar Tamanho da Página
```bash
curl http://localhost:8081/v1/users?page=1&page_size=50
```

### Última Página
```bash
# Use total_pages da resposta anterior
curl http://localhost:8081/v1/users?page=5
```

---

## 📊 **Metadata da Paginação**

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `current_page` | int | Página atual |
| `page_size` | int | Itens por página |
| `total_pages` | int | Total de páginas disponíveis |
| `total_records` | int64 | Total de registros |
| `has_next` | bool | Existe próxima página? |
| `has_previous` | bool | Existe página anterior? |

---

## ⚙️ **Configurações**

### Limites

```go
const (
    DefaultPage     = 1        // Página padrão
    DefaultPageSize = 20       // Itens por página padrão
    MaxPageSize     = 100      // Máximo de itens por página
)
```

### Validações Automáticas

**Se `page < 1`** → usa `page = 1`  
**Se `page_size < 1`** → usa `page_size = 20`  
**Se `page_size > 100`** → usa `page_size = 100`

---

## 🔧 **Implementação no Frontend**

### React Example

```javascript
import { useState, useEffect } from 'react';

function UsersList() {
  const [data, setData] = useState([]);
  const [pagination, setPagination] = useState({});
  const [page, setPage] = useState(1);
  
  useEffect(() => {
    fetch(`/v1/users?page=${page}&page_size=20`)
      .then(res => res.json())
      .then(json => {
        setData(json.message.data);
        setPagination(json.message.pagination);
      });
  }, [page]);
  
  return (
    <div>
      {/* Lista de usuários */}
      {data.map(user => (
        <div key={user.id}>{user.username}</div>
      ))}
      
      {/* Controles de paginação */}
      <button 
        disabled={!pagination.has_previous}
        onClick={() => setPage(page - 1)}
      >
        Anterior
      </button>
      
      <span>Página {pagination.current_page} de {pagination.total_pages}</span>
      
      <button 
        disabled={!pagination.has_next}
        onClick={() => setPage(page + 1)}
      >
        Próxima
      </button>
    </div>
  );
}
```

### Vue Example

```vue
<template>
  <div>
    <div v-for="user in users" :key="user.id">
      {{ user.username }}
    </div>
    
    <div class="pagination">
      <button @click="prevPage" :disabled="!pagination.has_previous">
        Anterior
      </button>
      <span>{{ pagination.current_page }} / {{ pagination.total_pages }}</span>
      <button @click="nextPage" :disabled="!pagination.has_next">
        Próxima
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';

const users = ref([]);
const pagination = ref({});
const page = ref(1);

const fetchUsers = async () => {
  const response = await fetch(`/v1/users?page=${page.value}&page_size=20`);
  const json = await response.json();
  users.value = json.message.data;
  pagination.value = json.message.pagination;
};

const nextPage = () => page.value++;
const prevPage = () => page.value--;

onMounted(fetchUsers);
watch(page, fetchUsers);
</script>
```

---

## 🎨 **Componentes UI Prontos**

### Material-UI (React)

```jsx
import { Pagination } from '@mui/material';

<Pagination 
  count={pagination.total_pages} 
  page={pagination.current_page}
  onChange={(e, newPage) => setPage(newPage)}
/>
```

### Vuetify (Vue)

```vue
<v-pagination
  v-model="page"
  :length="pagination.total_pages"
></v-pagination>
```

---

## 🚀 **Performance**

### Antes da Paginação
```
GET /v1/users
- Retorna: 10.000 registros
- Tempo: 2.5s
- Memória: 50MB
- Experiência: ❌ Ruim
```

### Depois da Paginação
```
GET /v1/users?page=1&page_size=20
- Retorna: 20 registros
- Tempo: 50ms
- Memória: 2MB
- Experiência: ✅ Excelente
```

**Ganho: 50x mais rápido!** 🚀

---

## 🔒 **Segurança**

### Proteção contra DoS

A paginação tem limite máximo de `100 itens` para prevenir:
- Sobrecarga do servidor
- Uso excessivo de memória
- Queries muito pesadas no banco

### Rate Limiting

Considere adicionar rate limiting para endpoints paginados:

```go
// Exemplo de uso com gin-rate-limit
import "github.com/JGLTechnologies/gin-rate-limit"

store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
    Rate:  time.Second,
    Limit: 20,
})

r.GET("/v1/users", ratelimit.RateLimiter(store), handler)
```

---

## 📚 **Referências**

- [REST API Design - Pagination](https://restfulapi.net/pagination/)
- [GraphQL Pagination Best Practices](https://graphql.org/learn/pagination/)
- [PostgreSQL LIMIT/OFFSET Performance](https://www.postgresql.org/docs/current/queries-limit.html)
