"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import PageHeader from "@/components/PageHeader";
import { api } from "@/services/api";

interface UpdateProfileRequest {
  username: string;
  email: string;
  password?: string;
  active: boolean;
  current_level: string;
  manager_name?: string;
  manager_email?: string;
}

export default function ProfilePage() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [success, setSuccess] = useState("");
  const [error, setError] = useState("");
  const [userId, setUserId] = useState("");
  
  const [formData, setFormData] = useState<UpdateProfileRequest>({
    username: "",
    email: "",
    password: "",
    active: true,
    current_level: "P1",
    manager_name: "",
    manager_email: "",
  });

  useEffect(() => {
    loadProfile();
  }, []);

  const loadProfile = async () => {
    try {
      setLoading(true);
      setError("");
      
      // Buscar ID do usuário do localStorage
      if (typeof window !== "undefined") {
        const userId = localStorage.getItem("user_id");
        const token = localStorage.getItem("career_token");
        
        if (!userId || !token) {
          console.warn("Usuário não está logado");
          router.push("/login");
          return;
        }
        
        setUserId(userId);
        
        // Buscar dados completos do backend
        const userData = await api.get(`/users/${userId}`);
        
        setFormData({
          username: userData.username || "",
          email: userData.email || "",
          password: "",
          active: userData.active ?? true,
          current_level: userData.current_level || "P1",
          manager_name: userData.manager_name || "",
          manager_email: userData.manager_email || "",
        });
      }
    } catch (err: any) {
      console.error("Erro ao carregar perfil:", err);
      setError(err.message || "Erro ao carregar perfil");
      // Se houver erro de autenticação, redireciona para login
      if (err.message?.includes("Sessão expirada") || err.message?.includes("401")) {
        router.push("/login");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      setSaving(true);
      setError("");
      setSuccess("");

      // Remove senha se estiver vazia (não atualiza)
      const dataToSend = { ...formData };
      if (!dataToSend.password || dataToSend.password.trim() === "") {
        delete dataToSend.password;
      }

      await api.put(`/users/${userId}`, dataToSend);
      
      setSuccess("Perfil atualizado com sucesso!");
      
      // Atualizar dados no localStorage
      if (typeof window !== "undefined") {
        const userStr = localStorage.getItem("career_user");
        if (userStr) {
          const user = JSON.parse(userStr);
          user.username = formData.username;
          user.email = formData.email;
          localStorage.setItem("career_user", JSON.stringify(user));
        }
      }

      // Limpar mensagem de sucesso após 3 segundos
      setTimeout(() => setSuccess(""), 3000);
    } catch (err: any) {
      console.error("Erro ao atualizar perfil:", err);
      setError(err.message || "Erro ao atualizar perfil");
    } finally {
      setSaving(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  if (loading) {
    return (
      <div className="flex min-h-screen bg-gray-950">
        <Sidebar />
        <main className="flex-1 ml-64 flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-4 text-gray-400">Carregando perfil...</p>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen bg-gray-950">
      <Sidebar />
      
      <main className="flex-1 ml-64 p-8">
        <div className="max-w-4xl mx-auto">
          <PageHeader 
            title="Meu Perfil" 
            subtitle="Gerencie suas informações pessoais e configurações"
          />

          {error && (
            <div className="mb-6 bg-red-900 bg-opacity-50 border border-red-700 text-red-300 px-4 py-3 rounded-lg">
              {error}
            </div>
          )}

          {success && (
            <div className="mb-6 bg-green-900 bg-opacity-50 border border-green-700 text-green-300 px-4 py-3 rounded-lg">
              {success}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Informações Pessoais */}
            <div className="bg-gray-900 rounded-lg shadow-sm border border-gray-800 p-6">
              <h3 className="text-lg font-semibold text-white mb-4">
                Informações Pessoais
              </h3>
              
              <div className="space-y-4">
                <div>
                  <label htmlFor="username" className="block text-sm font-medium text-gray-300 mb-1">
                    Nome de Usuário *
                  </label>
                  <input
                    type="text"
                    id="username"
                    name="username"
                    value={formData.username}
                    onChange={handleChange}
                    required
                    className="w-full px-4 py-2 bg-gray-800 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder-gray-500"
                    placeholder="Seu nome completo"
                  />
                </div>

                <div>
                  <label htmlFor="email" className="block text-sm font-medium text-gray-300 mb-1">
                    Email *
                  </label>
                  <input
                    type="email"
                    id="email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                    required
                    className="w-full px-4 py-2 bg-gray-800 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder-gray-500"
                    placeholder="seu@email.com"
                  />
                </div>

                <div>
                  <label htmlFor="current_level" className="block text-sm font-medium text-gray-300 mb-1">
                    Nível Atual
                  </label>
                  <select
                    id="current_level"
                    name="current_level"
                    value={formData.current_level}
                    onChange={handleChange}
                    className="w-full px-4 py-2 bg-gray-800 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value="P1">P1 - Júnior</option>
                    <option value="P2">P2 - Pleno</option>
                    <option value="P3">P3 - Sênior</option>
                    <option value="LT1">LT1 - Tech Lead</option>
                    <option value="LT2">LT2 - Staff</option>
                    <option value="LT3">LT3 - Principal</option>
                    <option value="LT4">LT4 - Distinguished</option>
                  </select>
                </div>

                <div>
                  <label htmlFor="password" className="block text-sm font-medium text-gray-300 mb-1">
                    Nova Senha (deixe em branco para não alterar)
                  </label>
                  <input
                    type="password"
                    id="password"
                    name="password"
                    value={formData.password}
                    onChange={handleChange}
                    className="w-full px-4 py-2 bg-gray-800 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder-gray-500"
                    placeholder="••••••••"
                    minLength={6}
                  />
                  <p className="mt-1 text-xs text-gray-400">
                    Mínimo de 6 caracteres
                  </p>
                </div>
              </div>
            </div>

            {/* Informações do Gerente */}
            <div className="bg-gray-900 rounded-lg shadow-sm border border-gray-800 p-6">
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h3 className="text-lg font-semibold text-white">
                    Gerente de Engenharia
                  </h3>
                  <p className="text-sm text-gray-400 mt-1">
                    Cadastre seu gerente para poder enviar relatórios diretamente
                  </p>
                </div>
                <div className="flex items-center gap-2">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                    📧 Envio de Relatórios
                  </span>
                </div>
              </div>
              
              <div className="space-y-4">
                <div>
                  <label htmlFor="manager_name" className="block text-sm font-medium text-gray-300 mb-1">
                    Nome do Gerente
                  </label>
                  <input
                    type="text"
                    id="manager_name"
                    name="manager_name"
                    value={formData.manager_name}
                    onChange={handleChange}
                    className="w-full px-4 py-2 bg-gray-800 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder-gray-500"
                    placeholder="Ex: Maria Santos"
                  />
                </div>

                <div>
                  <label htmlFor="manager_email" className="block text-sm font-medium text-gray-300 mb-1">
                    Email do Gerente
                  </label>
                  <input
                    type="email"
                    id="manager_email"
                    name="manager_email"
                    value={formData.manager_email}
                    onChange={handleChange}
                    className="w-full px-4 py-2 bg-gray-800 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder-gray-500"
                    placeholder="Ex: maria.santos@empresa.com"
                  />
                  {formData.manager_email && (
                    <p className="mt-1 text-xs text-green-600 flex items-center gap-1">
                      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                      </svg>
                      Pronto para enviar relatórios
                    </p>
                  )}
                </div>
              </div>
            </div>

            {/* Botões de Ação */}
            <div className="flex gap-3">
              <button
                type="submit"
                disabled={saving}
                className={`flex-1 py-3 px-6 rounded-lg font-medium text-white ${
                  saving 
                    ? "bg-gray-400 cursor-not-allowed" 
                    : "bg-blue-600 hover:bg-blue-700"
                } transition-colors`}
              >
                {saving ? "Salvando..." : "Salvar Alterações"}
              </button>
              
              <button
                type="button"
                onClick={() => router.back()}
                className="px-6 py-3 border border-gray-700 rounded-lg font-medium text-gray-300 hover:bg-gray-800 transition-colors"
              >
                Cancelar
              </button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}
