"use client";

import { useEffect, useState } from "react";
import Sidebar from "@/components/Sidebar";
import PageHeader from "@/components/PageHeader";
import { api } from "@/services/api";
import { useToast } from "@/components/Toast";

interface User {
  id: string;
  username: string;
  email: string;
  active: boolean;
}

interface UserForm {
  username: string;
  email: string;
  password: string;
  active: boolean;
  current_level: string;
}

export default function UsersPage() {
  const { toast } = useToast();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showViewModal, setShowViewModal] = useState(false);
  const [viewingUser, setViewingUser] = useState<User | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [formData, setFormData] = useState<UserForm>({
    username: "",
    email: "",
    password: "",
    active: true,
    current_level: "P1",
  });
  const [editFormData, setEditFormData] = useState({
    username: "",
    email: "",
    password: "",
    active: true,
    current_level: "P1",
  });

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const response = await api.get("/users");
      
      // A resposta pode vir paginada: response.items ou diretamente como array
      if (Array.isArray(response)) {
        setUsers(response);
      } else if (response.items && Array.isArray(response.items)) {
        setUsers(response.items);
      } else if (response.data && Array.isArray(response.data)) {
        setUsers(response.data);
      } else {
        console.error("Formato de resposta inesperado:", response);
        setUsers([]);
      }
    } catch (error) {
      console.error("Erro ao carregar usuários:", error);
      setUsers([]);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Tem certeza que deseja deletar este usuário?")) return;
    
    try {
      await api.delete(`/users/${id}`);
      setUsers(users.filter(u => u.id !== id));
      toast("Usuario deletado com sucesso!");
    } catch (error) {
      console.error("Erro ao deletar:", error);
      toast("Erro ao deletar usuario", "error");
    }
  };

  const handleEditClick = async (user: User) => {
    setEditingUser(user);
    setEditFormData({
      username: user.username,
      email: user.email,
      password: "",
      active: user.active,
      current_level: "P1",
    });
    setShowEditModal(true);
  };

  const handleUpdateUser = async () => {
    if (!editingUser) return;

    if (!editFormData.username.trim() || !editFormData.email.trim()) {
      toast("Nome de usuario e email sao obrigatorios", "warning");
      return;
    }

    setSubmitting(true);
    try {
      await api.put(`/users/${editingUser.id}`, {
        username: editFormData.username,
        email: editFormData.email,
        password: editFormData.password || undefined,
        active: editFormData.active,
        current_level: editFormData.current_level,
      });

      await fetchUsers();
      setShowEditModal(false);
      setEditingUser(null);
      toast("Usuario atualizado com sucesso!");
    } catch (error) {
      console.error("Erro ao atualizar usuário:", error);
      toast("Erro ao atualizar usuario. Tente novamente.", "error");
    } finally {
      setSubmitting(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.username.trim() || !formData.email.trim() || !formData.password.trim()) {
      toast("Preencha todos os campos obrigatorios", "warning");
      return;
    }

    setSubmitting(true);
    try {
      await api.post("/users", {
        username: formData.username,
        email: formData.email,
        password: formData.password,
        active: formData.active,
        current_level: formData.current_level,
      });
      
      setFormData({
        username: "",
        email: "",
        password: "",
        active: true,
        current_level: "P1",
      });
      
      setShowModal(false);
      fetchUsers();
      toast("Usuario criado com sucesso!");
    } catch (error) {
      console.error("Erro ao criar usuário:", error);
      toast("Erro ao criar usuario. Tente novamente.", "error");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="flex min-h-screen bg-gray-950">
      <Sidebar />
      
      <main className="flex-1 ml-64 p-8">
        <PageHeader 
          title="Usuários" 
          subtitle="Gerencie os usuários do sistema"
          action={
            <button 
              onClick={() => setShowModal(true)}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all flex items-center gap-2"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              Novo Usuário
            </button>
          }
        />

        {loading ? (
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
            <p className="text-gray-400">Carregando usuários...</p>
          </div>
        ) : users.length === 0 ? (
          <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
            <svg className="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <h3 className="text-xl font-semibold text-white mb-2">Nenhum usuário cadastrado</h3>
            <p className="text-gray-400 mb-6">Comece criando o primeiro usuário</p>
            <button 
              onClick={() => setShowModal(true)}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all"
            >
              Criar Usuário
            </button>
          </div>
        ) : Array.isArray(users) ? (
          <div className="bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
            <table className="w-full">
              <thead className="bg-gray-700">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">
                    Nome
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">
                    Email
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-300 uppercase tracking-wider">
                    Ações
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-700">
                {users.map((user) => (
                  <tr key={user.id} className="hover:bg-gray-750 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <div className="w-10 h-10 bg-blue-600 rounded-full flex items-center justify-center text-white font-semibold">
                          {user.username.charAt(0).toUpperCase()}
                        </div>
                        <div className="ml-4">
                          <div className="text-sm font-medium text-white">{user.username}</div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-300">{user.email}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        user.active 
                          ? 'bg-green-900 text-green-300' 
                          : 'bg-red-900 text-red-300'
                      }`}>
                        {user.active ? 'Ativo' : 'Inativo'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <div className="flex justify-end gap-2">
                        <button
                          onClick={() => { setViewingUser(user); setShowViewModal(true); }}
                          className="p-2 text-gray-400 hover:text-green-400 hover:bg-gray-700 rounded transition-all"
                          title="Visualizar Usuário"
                        >
                          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                          </svg>
                        </button>
                        <button
                          onClick={() => handleEditClick(user)}
                          className="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded transition-all"
                          title="Editar Usuário"
                        >
                          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                          </svg>
                        </button>
                        <button
                          onClick={() => handleDelete(user.id)}
                          className="p-2 text-gray-400 hover:text-red-400 hover:bg-gray-700 rounded transition-all"
                          title="Deletar Usuário"
                        >
                          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                          </svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
            <div className="w-16 h-16 bg-red-900 text-red-400 rounded-full flex items-center justify-center mx-auto mb-4 text-2xl">
              ⚠️
            </div>
            <h3 className="text-xl font-semibold text-white mb-2">Erro ao carregar usuários</h3>
            <p className="text-gray-400 mb-6">Formato de resposta inesperado</p>
            <button 
              onClick={fetchUsers}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all"
            >
              Tentar Novamente
            </button>
          </div>
        )}
      </main>

      {/* Modal de Criação de Usuário */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-2xl font-bold text-white">Novo Usuário</h2>
              <button
                onClick={() => setShowModal(false)}
                className="text-gray-400 hover:text-white transition-colors"
              >
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Nome de Usuário <span className="text-red-400">*</span>
                </label>
                <input
                  type="text"
                  value={formData.username}
                  onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="johndoe"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Email <span className="text-red-400">*</span>
                </label>
                <input
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="john@exemplo.com"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Senha <span className="text-red-400">*</span>
                </label>
                <input
                  type="password"
                  value={formData.password}
                  onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="••••••••"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Nível Atual
                </label>
                <select
                  value={formData.current_level}
                  onChange={(e) => setFormData({ ...formData, current_level: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="P1">P1 - Júnior</option>
                  <option value="P2">P2 - Pleno</option>
                  <option value="P3">P3 - Sênior</option>
                  <option value="LT1">LT1 - Tech Lead</option>
                  <option value="LT2">LT2 - Engineering Manager</option>
                  <option value="LT3">LT3 - Director</option>
                  <option value="LT4">LT4 - VP</option>
                </select>
              </div>

              <div className="flex items-center gap-3">
                <input
                  type="checkbox"
                  id="active"
                  checked={formData.active}
                  onChange={(e) => setFormData({ ...formData, active: e.target.checked })}
                  className="w-5 h-5 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-2 focus:ring-blue-500"
                />
                <label htmlFor="active" className="text-sm font-medium text-gray-300">
                  Usuário ativo
                </label>
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="flex-1 px-6 py-3 bg-gray-700 text-white rounded-lg font-semibold hover:bg-gray-600 transition-all"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  disabled={submitting}
                  className="flex-1 px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {submitting ? "Criando..." : "Criar Usuário"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Modal de Edição de Usuário */}
      {showEditModal && editingUser && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-2xl font-bold text-white">Editar Usuário</h2>
              <button
                onClick={() => {
                  setShowEditModal(false);
                  setEditingUser(null);
                }}
                className="text-gray-400 hover:text-white transition-colors"
              >
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Nome de Usuário <span className="text-red-400">*</span>
                </label>
                <input
                  type="text"
                  value={editFormData.username}
                  onChange={(e) => setEditFormData({ ...editFormData, username: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Email <span className="text-red-400">*</span>
                </label>
                <input
                  type="email"
                  value={editFormData.email}
                  onChange={(e) => setEditFormData({ ...editFormData, email: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Nova Senha (deixe em branco para manter a atual)
                </label>
                <input
                  type="password"
                  value={editFormData.password}
                  onChange={(e) => setEditFormData({ ...editFormData, password: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="••••••••"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Nível Atual
                </label>
                <select
                  value={editFormData.current_level}
                  onChange={(e) => setEditFormData({ ...editFormData, current_level: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="P1">P1 - Júnior</option>
                  <option value="P2">P2 - Pleno</option>
                  <option value="P3">P3 - Sênior</option>
                  <option value="LT1">LT1 - Tech Lead</option>
                  <option value="LT2">LT2 - Engineering Manager</option>
                  <option value="LT3">LT3 - Director</option>
                  <option value="LT4">LT4 - VP</option>
                </select>
              </div>

              <div className="flex items-center gap-3">
                <input
                  type="checkbox"
                  id="edit_active"
                  checked={editFormData.active}
                  onChange={(e) => setEditFormData({ ...editFormData, active: e.target.checked })}
                  className="w-5 h-5 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-2 focus:ring-blue-500"
                />
                <label htmlFor="edit_active" className="text-sm font-medium text-gray-300">
                  Usuário ativo
                </label>
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowEditModal(false);
                    setEditingUser(null);
                  }}
                  className="flex-1 px-6 py-3 bg-gray-700 text-white rounded-lg font-semibold hover:bg-gray-600 transition-all"
                >
                  Cancelar
                </button>
                <button
                  onClick={handleUpdateUser}
                  disabled={submitting}
                  className="flex-1 px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {submitting ? "Salvando..." : "Salvar Alterações"}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
      {/* Modal de Visualizacao de Usuario */}
      {showViewModal && viewingUser && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-xl font-bold text-white">Detalhes do Usuario</h2>
              <button onClick={() => { setShowViewModal(false); setViewingUser(null); }} className="text-gray-400 hover:text-white">
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div className="flex items-center gap-4 mb-2">
                <div className="w-14 h-14 bg-blue-600 rounded-full flex items-center justify-center text-white text-xl font-bold">
                  {viewingUser.username.charAt(0).toUpperCase()}
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-white">{viewingUser.username}</h3>
                  <p className="text-sm text-gray-400">{viewingUser.email}</p>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div className="bg-gray-700 rounded-lg p-3">
                  <p className="text-xs text-gray-400 mb-1">Status</p>
                  <span className={`px-2 py-1 text-xs font-semibold rounded-full ${viewingUser.active ? "bg-green-900 text-green-300" : "bg-red-900 text-red-300"}`}>
                    {viewingUser.active ? "Ativo" : "Inativo"}
                  </span>
                </div>
                <div className="bg-gray-700 rounded-lg p-3">
                  <p className="text-xs text-gray-400 mb-1">ID</p>
                  <span className="text-xs text-gray-300 font-mono">{viewingUser.id.slice(0, 8)}...</span>
                </div>
              </div>
            </div>
            <div className="border-t border-gray-700 px-6 py-4">
              <button onClick={() => { setShowViewModal(false); setViewingUser(null); }} className="w-full px-4 py-3 bg-gray-700 text-gray-300 rounded-lg font-medium hover:bg-gray-600 transition-all">Fechar</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
