// frontend/src/services/api.ts

const API_URL = process.env.NEXT_PUBLIC_API_URL;

// Função auxiliar para pegar o token do navegador com segurança
const getHeaders = () => {
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  // Garante que o código só tente acessar o localStorage no navegador
  if (typeof window !== "undefined") {
    const token = localStorage.getItem("career_token");
    if (token) {
      headers["Authorization"] = `Bearer ${token}`; 
    }
  }
  return headers;
};

const handleResponse = async (response: Response) => {
  const data = await response.json();
  
  // Se o token for inválido, desloga o usuário
  if (response.status === 401) {
    if (typeof window !== "undefined") {
      localStorage.removeItem("career_token");
      window.location.href = "/login";
    }
    throw new Error("Sessão expirada");
  }

  if (!response.ok) throw new Error(data.error || "Erro na requisição");
  return data;
};

export const api = {
  // === MÉTODO POST ===
  async post<T>(endpoint: string, body: T) {
    const response = await fetch(`${API_URL}${endpoint}`, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(body),
    });
    return handleResponse(response);
  },

  // === MÉTODO GET ===
  async get(endpoint: string) {
    const response = await fetch(`${API_URL}${endpoint}`, {
      method: "GET",
      headers: getHeaders(),
    });
    return handleResponse(response);
  },

  // === MÉTODO PATCH ===
  async patch<T>(endpoint: string, body: T) {
    const response = await fetch(`${API_URL}${endpoint}`, {
      method: "PATCH",
      headers: getHeaders(),
      body: JSON.stringify(body),
    });
    return handleResponse(response);
  },

  // === MÉTODO DELETE ===
  async delete(endpoint: string) {
    const response = await fetch(`${API_URL}${endpoint}`, {
      method: "DELETE",
      headers: getHeaders(),
    });
    return handleResponse(response);
  },

  // === MÉTODO PUT ===
  async put<T>(endpoint: string, body: T) {
    const response = await fetch(`${API_URL}${endpoint}`, {
      method: "PUT",
      headers: getHeaders(),
      body: JSON.stringify(body),
    });
    return handleResponse(response);
  },
};