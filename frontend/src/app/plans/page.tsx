"use client";

import { useEffect, useState } from "react";
import Sidebar from "@/components/Sidebar";
import PageHeader from "@/components/PageHeader";
import { api } from "@/services/api";
import { useToast } from "@/components/Toast";

interface Plan {
  id: string;
  title: string;
  description?: string;
  initiative_id?: string;
  initiative_title?: string;
  level_target?: string;
  status: string;
  position: number;
  created_at: string;
  updated_at: string;
}

interface PlanForm {
  title: string;
  description: string;
  initiative_id: string;
  level_target: string;
}

interface Initiative {
  id: string;
  title: string;
}

interface Ladder {
  id: string;
  level: string;
  xp_reward: number;
}

const emptyPlanForm: PlanForm = { title: "", description: "", initiative_id: "", level_target: "" };

export default function PlansPage() {
  const { toast } = useToast();
  const [plans, setPlans] = useState<Plan[]>([]);
  const [initiatives, setInitiatives] = useState<Initiative[]>([]);
  const [ladders, setLadders] = useState<Ladder[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showViewModal, setShowViewModal] = useState(false);
  const [viewingPlan, setViewingPlan] = useState<Plan | null>(null);
  const [editingPlan, setEditingPlan] = useState<Plan | null>(null);
  const [createForm, setCreateForm] = useState<PlanForm>(emptyPlanForm);
  const [editForm, setEditForm] = useState<PlanForm>(emptyPlanForm);

  useEffect(() => {
    fetchPlans();
    fetchInitiatives();
    fetchLadders();
  }, []);

  const fetchPlans = async () => {
    try {
      const r = await api.get("/plans");
      setPlans(Array.isArray(r) ? r : []);
    } catch {
      setPlans([]);
    } finally {
      setLoading(false);
    }
  };

  const fetchInitiatives = async () => {
    try {
      const r = await api.get("/initiatives");
      setInitiatives(Array.isArray(r) ? r : []);
    } catch {
      setInitiatives([]);
    }
  };

  const fetchLadders = async () => {
    try {
      const r = await api.get("/ladders");
      setLadders(Array.isArray(r) ? r : []);
    } catch {
      setLadders([]);
    }
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!createForm.title.trim()) {
      toast("Titulo e obrigatorio", "warning");
      return;
    }
    setSubmitting(true);
    try {
      await api.post("/plans", {
        title: createForm.title,
        description: createForm.description || undefined,
        initiative_id: createForm.initiative_id || undefined,
        level_target: createForm.level_target || undefined,
      });
      setCreateForm(emptyPlanForm);
      setShowCreateModal(false);
      fetchPlans();
      toast("Plano criado com sucesso!");
    } catch (err) {
      const error = err as Error;
      toast(error.message || "Erro ao criar plano.", "error");
    } finally {
      setSubmitting(false);
    }
  };

  const handleEditClick = (plan: Plan) => {
    setEditingPlan(plan);
    setEditForm({
      title: plan.title,
      description: plan.description || "",
      initiative_id: plan.initiative_id || "",
      level_target: plan.level_target || "",
    });
    setShowEditModal(true);
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingPlan || !editForm.title.trim()) {
      toast("Titulo e obrigatorio", "warning");
      return;
    }
    setSubmitting(true);
    try {
      await api.put(`/plans/${editingPlan.id}`, {
        title: editForm.title,
        description: editForm.description || undefined,
        initiative_id: editForm.initiative_id || undefined,
        level_target: editForm.level_target || undefined,
      });
      setShowEditModal(false);
      setEditingPlan(null);
      fetchPlans();
      toast("Plano atualizado!");
    } catch (err) {
      const error = err as Error;
      toast(error.message || "Erro ao atualizar.", "error");
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Deletar este plano?")) return;
    try {
      await api.delete(`/plans/${id}`);
      setPlans((prev) => prev.filter((p) => p.id !== id));
      toast("Plano deletado!");
    } catch {
      toast("Erro ao deletar.", "error");
    }
  };

  const handleView = (plan: Plan) => {
    setViewingPlan(plan);
    setShowViewModal(true);
  };

  const handleMoveUp = async (plan: Plan) => {
    if (plan.position === 0) return;
    try {
      await api.put(`/plans/${plan.id}/move-up`, {});
      fetchPlans();
    } catch {
      toast("Erro ao mover plano.", "error");
    }
  };

  const handleMoveDown = async (plan: Plan) => {
    if (plan.position >= plans.length - 1) return;
    try {
      await api.put(`/plans/${plan.id}/move-down`, {});
      fetchPlans();
    } catch {
      toast("Erro ao mover plano.", "error");
    }
  };

  return (
    <div className="flex min-h-screen bg-gray-950">
      <Sidebar />
      <main className="flex-1 ml-64 p-8">
        <PageHeader
          title="Planos"
          subtitle="Gerencie seus planos futuros de desenvolvimento"
          action={
            <button
              onClick={() => setShowCreateModal(true)}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all flex items-center gap-2"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              Novo Plano
            </button>
          }
        />

        {loading ? (
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
            <p className="text-gray-400">Carregando...</p>
          </div>
        ) : plans.length === 0 ? (
          <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
            <h3 className="text-xl font-semibold text-white mb-2">Nenhum plano cadastrado</h3>
            <p className="text-gray-400 mb-6">Comece criando seu primeiro plano</p>
            <button
              onClick={() => setShowCreateModal(true)}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700"
            >
              Criar Plano
            </button>
          </div>
        ) : (
          <div className="grid gap-4">
            {plans.map((plan) => (
              <div key={plan.id} className="bg-gray-800 border border-gray-700 rounded-lg p-6">
                <div className="flex items-start justify-between mb-3">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <h3 className="text-lg font-semibold text-white">{plan.title}</h3>
                      {plan.initiative_title && (
                        <span className="px-2 py-0.5 bg-purple-900 text-purple-300 text-xs rounded">
                          {plan.initiative_title}
                        </span>
                      )}
                      {plan.level_target && (
                        <span className="px-2 py-0.5 bg-yellow-900 text-yellow-300 text-xs rounded">
                          {plan.level_target}
                        </span>
                      )}
                    </div>
                    {plan.description && <p className="text-gray-400 text-sm">{plan.description}</p>}
                  </div>
                  <div className="flex gap-1">
                    <button
                      onClick={() => handleMoveUp(plan)}
                      disabled={plan.position === 0}
                      className="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded disabled:opacity-30 disabled:cursor-not-allowed"
                      title="Mover para cima"
                    >
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button
                      onClick={() => handleMoveDown(plan)}
                      disabled={plan.position >= plans.length - 1}
                      className="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded disabled:opacity-30 disabled:cursor-not-allowed"
                      title="Mover para baixo"
                    >
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                    <button
                      onClick={() => handleView(plan)}
                      className="p-2 text-gray-400 hover:text-green-400 hover:bg-gray-700 rounded"
                      title="Visualizar"
                    >
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                        />
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                        />
                      </svg>
                    </button>
                    <button
                      onClick={() => handleEditClick(plan)}
                      className="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded"
                      title="Editar"
                    >
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                        />
                      </svg>
                    </button>
                    <button
                      onClick={() => handleDelete(plan.id)}
                      className="p-2 text-gray-400 hover:text-red-400 hover:bg-gray-700 rounded"
                      title="Deletar"
                    >
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>

      {/* Create Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-xl font-bold text-white">Novo Plano</h2>
              <button
                onClick={() => setShowCreateModal(false)}
                className="text-gray-400 hover:text-white"
              >
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <form onSubmit={handleCreate} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Titulo *</label>
                <input
                  type="text"
                  value={createForm.title}
                  onChange={(e) => setCreateForm({ ...createForm, title: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Ex: Estudar arquitetura de microservicos"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Descricao</label>
                <textarea
                  value={createForm.description}
                  onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  rows={3}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Iniciativa (opcional)</label>
                <select
                  value={createForm.initiative_id}
                  onChange={(e) => setCreateForm({ ...createForm, initiative_id: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Selecione</option>
                  {initiatives.map((i) => (
                    <option key={i.id} value={i.id}>
                      {i.title}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Nivel Alvo (opcional)</label>
                <select
                  value={createForm.level_target}
                  onChange={(e) => setCreateForm({ ...createForm, level_target: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Selecione</option>
                  {ladders.map((l) => (
                    <option key={l.id} value={l.level}>
                      {l.level} - {l.xp_reward} XP
                    </option>
                  ))}
                </select>
              </div>
              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowCreateModal(false)}
                  className="flex-1 px-4 py-3 bg-gray-700 text-gray-300 rounded-lg hover:bg-gray-600"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  disabled={submitting}
                  className="flex-1 px-4 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 disabled:opacity-50"
                >
                  {submitting ? "Salvando..." : "Cadastrar"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Edit Modal */}
      {showEditModal && editingPlan && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-xl font-bold text-white">Editar Plano</h2>
              <button
                onClick={() => {
                  setShowEditModal(false);
                  setEditingPlan(null);
                }}
                className="text-gray-400 hover:text-white"
              >
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <form onSubmit={handleUpdate} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Titulo *</label>
                <input
                  type="text"
                  value={editForm.title}
                  onChange={(e) => setEditForm({ ...editForm, title: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Descricao</label>
                <textarea
                  value={editForm.description}
                  onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  rows={3}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Iniciativa (opcional)</label>
                <select
                  value={editForm.initiative_id}
                  onChange={(e) => setEditForm({ ...editForm, initiative_id: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Selecione</option>
                  {initiatives.map((i) => (
                    <option key={i.id} value={i.id}>
                      {i.title}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Nivel Alvo (opcional)</label>
                <select
                  value={editForm.level_target}
                  onChange={(e) => setEditForm({ ...editForm, level_target: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Selecione</option>
                  {ladders.map((l) => (
                    <option key={l.id} value={l.level}>
                      {l.level} - {l.xp_reward} XP
                    </option>
                  ))}
                </select>
              </div>
              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => {
                    setShowEditModal(false);
                    setEditingPlan(null);
                  }}
                  className="flex-1 px-4 py-3 bg-gray-700 text-gray-300 rounded-lg hover:bg-gray-600"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  disabled={submitting}
                  className="flex-1 px-4 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 disabled:opacity-50"
                >
                  {submitting ? "Salvando..." : "Atualizar"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* View Modal */}
      {showViewModal && viewingPlan && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-xl font-bold text-white">Detalhes do Plano</h2>
              <button
                onClick={() => {
                  setShowViewModal(false);
                  setViewingPlan(null);
                }}
                className="text-gray-400 hover:text-white"
              >
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <h3 className="text-lg font-semibold text-white mb-1">{viewingPlan.title}</h3>
                <div className="flex gap-2 mb-3">
                  {viewingPlan.initiative_title && (
                    <span className="px-2 py-0.5 bg-purple-900 text-purple-300 text-xs rounded">
                      {viewingPlan.initiative_title}
                    </span>
                  )}
                  {viewingPlan.level_target && (
                    <span className="px-2 py-0.5 bg-yellow-900 text-yellow-300 text-xs rounded">
                      {viewingPlan.level_target}
                    </span>
                  )}
                </div>
              </div>
              {viewingPlan.description && (
                <div className="bg-gray-700 rounded-lg p-3">
                  <p className="text-xs text-gray-400 mb-1">Descricao</p>
                  <p className="text-sm text-gray-200">{viewingPlan.description}</p>
                </div>
              )}
              <div className="grid grid-cols-2 gap-3">
                <div className="bg-gray-700 rounded-lg p-3">
                  <p className="text-xs text-gray-400 mb-1">Status</p>
                  <p className="text-sm font-semibold text-white capitalize">{viewingPlan.status}</p>
                </div>
                <div className="bg-gray-700 rounded-lg p-3">
                  <p className="text-xs text-gray-400 mb-1">Criado em</p>
                  <p className="text-sm text-white">
                    {new Date(viewingPlan.created_at).toLocaleDateString("pt-BR")}
                  </p>
                </div>
              </div>
            </div>
            <div className="border-t border-gray-700 px-6 py-4">
              <button
                onClick={() => {
                  setShowViewModal(false);
                  setViewingPlan(null);
                }}
                className="w-full px-4 py-3 bg-gray-700 text-gray-300 rounded-lg font-medium hover:bg-gray-600 transition-all"
              >
                Fechar
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}