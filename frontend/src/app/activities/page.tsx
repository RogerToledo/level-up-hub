"use client";

import { useEffect, useState } from "react";
import Sidebar from "@/components/Sidebar";
import PageHeader from "@/components/PageHeader";
import { api } from "@/services/api";
import { Pillar } from "@/types";

interface Activity {
  id: string;
  title: string;
  description?: string;
  progress_percentage: number;
  is_pdi_target: boolean;
  created_at: string;
  evidence_count?: number;
}

interface Evidence {
  id: string;
  evidence_url: string;
  description: string;
}

interface Ladder {
  id: string;
  level: string;
  xp_reward: number;
}

interface ActivityForm {
  title: string;
  description: string;
  pillars: string[];
  progress_percentage: number;
  impact_summary: string;
  is_pdi_target: boolean;
  ladder_id: string;
}

export default function ActivitiesPage() {
  const [activities, setActivities] = useState<Activity[]>([]);
  const [ladders, setLadders] = useState<Ladder[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showEvidenceModal, setShowEvidenceModal] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [editingActivity, setEditingActivity] = useState<Activity | null>(null);
  const [selectedActivityForEvidence, setSelectedActivityForEvidence] = useState<Activity | null>(null);
  const [evidenceForm, setEvidenceForm] = useState({
    url: "",
    description: "",
  });
  const [activityEvidences, setActivityEvidences] = useState<Record<string, Evidence[]>>({});
  const [editFormData, setEditFormData] = useState({
    title: "",
    description: "",
    progress_percentage: 0,
    impact_summary: "",
    is_pdi_target: false,
  });
  const [formData, setFormData] = useState<ActivityForm>({
    title: "",
    description: "",
    pillars: [],
    progress_percentage: 0,
    impact_summary: "",
    is_pdi_target: false,
    ladder_id: "",
  });

  useEffect(() => {
    fetchActivities();
    fetchLadders();
  }, []);

  const fetchActivities = async () => {
    try {
      const response = await api.get("/activities");
      console.log("Activities response:", response);
      
      // A resposta pode vir paginada: response.items ou diretamente como array
      if (Array.isArray(response)) {
        setActivities(response);
      } else if (response.items && Array.isArray(response.items)) {
        setActivities(response.items);
      } else if (response.data && Array.isArray(response.data)) {
        setActivities(response.data);
      } else {
        console.error("Formato de resposta inesperado:", response);
        setActivities([]);
      }
    } catch (error) {
      console.error("Erro ao carregar atividades:", error);
      setActivities([]);
    } finally {
      setLoading(false);
    }
  };

  const fetchLadders = async () => {
    try {
      const data = await api.get("/ladders");
      setLadders(data);
    } catch (error) {
      console.error("Erro ao carregar níveis:", error);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Tem certeza que deseja deletar esta atividade?")) return;
    
    try {
      await api.delete(`/activities/${id}`);
      setActivities(activities.filter(a => a.id !== id));
    } catch (error) {
      console.error("Erro ao deletar:", error);
      alert("Erro ao deletar atividade");
    }
  };

  const handleEditClick = (activity: Activity) => {
    setEditingActivity(activity);
    setEditFormData({
      title: activity.title,
      description: activity.description || "",
      progress_percentage: activity.progress_percentage,
      impact_summary: "", // Não temos impact_summary na interface Activity
      is_pdi_target: activity.is_pdi_target,
    });
    setShowEditModal(true);
  };

  const handleUpdateActivity = async () => {
    if (!editingActivity) return;

    if (!editFormData.title.trim()) {
      alert("O título é obrigatório");
      return;
    }

    setSubmitting(true);
    try {
      await api.put(`/activities/${editingActivity.id}`, {
        title: editFormData.title,
        description: editFormData.description || undefined,
        progress_percentage: editFormData.progress_percentage,
        impact_summary: editFormData.impact_summary || undefined,
        is_pdi_target: editFormData.is_pdi_target,
      });

      // Atualiza a lista de atividades
      await fetchActivities();

      setShowEditModal(false);
      setEditingActivity(null);
      alert("Atividade atualizada com sucesso!");
    } catch (error) {
      console.error("Erro ao atualizar atividade:", error);
      alert("Erro ao atualizar atividade. Tente novamente.");
    } finally {
      setSubmitting(false);
    }
  };

  const handleEvidenceClick = async (activity: Activity) => {
    setSelectedActivityForEvidence(activity);
    setEvidenceForm({ url: "", description: "" });
    setShowEvidenceModal(true);
    
    // Busca evidências existentes
    try {
      const evidences = await api.get(`/activities/${activity.id}/evidences`);
      setActivityEvidences(prev => ({ ...prev, [activity.id]: evidences }));
    } catch (error) {
      console.error("Erro ao carregar evidências:", error);
    }
  };

  const handleAddEvidence = async () => {
    if (!selectedActivityForEvidence) return;

    if (!evidenceForm.url.trim()) {
      alert("A URL é obrigatória");
      return;
    }

    // Validação básica de URL
    try {
      new URL(evidenceForm.url);
    } catch {
      alert("Por favor, insira uma URL válida");
      return;
    }

    setSubmitting(true);
    try {
      await api.post(`/activities/${selectedActivityForEvidence.id}/evidence`, {
        url: evidenceForm.url,
        description: evidenceForm.description || undefined,
      });

      // Atualiza lista de evidências
      const evidences = await api.get(`/activities/${selectedActivityForEvidence.id}/evidences`);
      setActivityEvidences(prev => ({ ...prev, [selectedActivityForEvidence.id]: evidences }));

      // Atualiza a lista de atividades para refletir a nova contagem de evidências
      await fetchActivities();

      setEvidenceForm({ url: "", description: "" });
      alert("Evidência adicionada com sucesso!");
    } catch (error) {
      console.error("Erro ao adicionar evidência:", error);
      alert("Erro ao adicionar evidência. Tente novamente.");
    } finally {
      setSubmitting(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.title || formData.pillars.length === 0 || !formData.ladder_id) {
      alert("Preencha os campos obrigatórios: Título, Pilares e Nível");
      return;
    }

    setSubmitting(true);
    
    try {
      const userId = localStorage.getItem("user_id");
      if (!userId) {
        alert("Usuário não identificado. Faça login novamente.");
        return;
      }

      const payload = {
        user_id: userId,
        ladder_id: formData.ladder_id,
        pillars: formData.pillars,
        title: formData.title,
        description: formData.description || undefined,
        progress_percentage: formData.progress_percentage,
        impact_summary: formData.impact_summary || undefined,
        is_pdi_target: formData.is_pdi_target,
      };

      await api.post("/activities", payload);
      
      // Limpa o formulário
      setFormData({
        title: "",
        description: "",
        pillars: [],
        progress_percentage: 0,
        impact_summary: "",
        is_pdi_target: false,
        ladder_id: "",
      });
      
      setShowModal(false);
      fetchActivities();
      alert("Atividade criada com sucesso!");
    } catch (error) {
      console.error("Erro ao criar atividade:", error);
      alert("Erro ao criar atividade. Tente novamente.");
    } finally {
      setSubmitting(false);
    }
  };

  const handlePillarToggle = (pillar: string) => {
    setFormData(prev => ({
      ...prev,
      pillars: prev.pillars.includes(pillar)
        ? prev.pillars.filter(p => p !== pillar)
        : [...prev.pillars, pillar]
    }));
  };

  const formatDate = (dateString: string) => {
    if (!dateString) return "";
    const date = new Date(dateString);
    return date.toLocaleDateString("pt-BR", {
      day: "2-digit",
      month: "short",
      year: "numeric"
    });
  };

  return (
    <div className="flex min-h-screen bg-gray-950">
      <Sidebar />
      
      <main className="flex-1 ml-64 p-8">
        <PageHeader 
          title="Atividades" 
          subtitle="Gerencie suas atividades de desenvolvimento"
          action={
            <button 
              onClick={() => setShowModal(true)}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all flex items-center gap-2"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              Nova Atividade
            </button>
          }
        />

        {loading ? (
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
            <p className="text-gray-400">Carregando atividades...</p>
          </div>
        ) : activities.length === 0 ? (
          <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
            <svg className="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <h3 className="text-xl font-semibold text-white mb-2">Nenhuma atividade cadastrada</h3>
            <p className="text-gray-400 mb-6">Comece criando sua primeira atividade</p>
            <button 
              onClick={() => setShowModal(true)}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all"
            >
              Criar Atividade
            </button>
          </div>
        ) : Array.isArray(activities) ? (
          <div className="grid gap-4">
            {activities.map((activity) => (
              <div key={activity.id} className="bg-gray-800 border border-gray-700 rounded-lg p-6 hover:border-gray-600 transition-all">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <h3 className="text-lg font-semibold text-white">{activity.title}</h3>
                      {activity.is_pdi_target && (
                        <span className="px-2 py-1 bg-blue-900 text-blue-300 text-xs font-medium rounded">PDI</span>
                      )}
                      {activity.progress_percentage === 100 && (
                        <span className="px-2 py-1 bg-green-900 text-green-300 text-xs font-medium rounded">✓ Completo</span>
                      )}
                      {(activity.evidence_count ?? 0) > 0 && (
                        <span className="px-2 py-1 bg-purple-900 text-purple-300 text-xs font-medium rounded flex items-center gap-1">
                          <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                          </svg>
                          {activity.evidence_count}
                        </span>
                      )}
                    </div>
                    {activity.description && (
                      <p className="text-gray-400 text-sm mb-3">{activity.description}</p>
                    )}
                    <p className="text-gray-500 text-xs">
                      Criada em {formatDate(activity.created_at)}
                    </p>
                  </div>
                  <div className="flex gap-2">
                    <button 
                      onClick={() => handleEvidenceClick(activity)}
                      className="p-2 text-gray-400 hover:text-blue-400 hover:bg-gray-700 rounded transition-all"
                      title="Gerenciar Evidências"
                    >
                      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                      </svg>
                    </button>
                    <button 
                      onClick={() => handleEditClick(activity)}
                      className="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded transition-all"
                      title="Editar Atividade"
                    >
                      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                    </button>
                    <button 
                      onClick={() => handleDelete(activity.id)}
                      className="p-2 text-gray-400 hover:text-red-400 hover:bg-gray-700 rounded transition-all"
                      title="Deletar Atividade"
                    >
                      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                  </div>
                </div>
                
                <div className="flex items-center gap-3">
                  <div className="flex-1 bg-gray-700 rounded-full h-2.5">
                    <div
                      className={`h-2.5 rounded-full transition-all ${
                        activity.progress_percentage === 100 
                          ? 'bg-green-600' 
                          : activity.progress_percentage >= 50 
                          ? 'bg-blue-600' 
                          : 'bg-yellow-600'
                      }`}
                      style={{ width: `${activity.progress_percentage}%` }}
                    ></div>
                  </div>
                  <span className="text-sm font-medium text-gray-300 min-w-12">
                    {activity.progress_percentage}%
                  </span>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
            <div className="w-16 h-16 bg-red-900 text-red-400 rounded-full flex items-center justify-center mx-auto mb-4 text-2xl">
              ⚠️
            </div>
            <h3 className="text-xl font-semibold text-white mb-2">Erro ao carregar atividades</h3>
            <p className="text-gray-400 mb-6">Formato de resposta inesperado</p>
            <button 
              onClick={fetchActivities}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all"
            >
              Tentar Novamente
            </button>
          </div>
        )}
      </main>

      {/* Modal de Criação de Atividade */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto border border-gray-700">
            <div className="sticky top-0 bg-gray-800 border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-2xl font-bold text-white">Nova Atividade</h2>
              <button
                onClick={() => setShowModal(false)}
                className="text-gray-400 hover:text-white transition-colors"
              >
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-6">
              {/* Título */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Título <span className="text-red-400">*</span>
                </label>
                <input
                  type="text"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Ex: Refatoração do módulo de autenticação"
                  required
                />
              </div>

              {/* Descrição */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Descrição
                </label>
                <textarea
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Descreva a atividade..."
                  rows={3}
                />
              </div>

              {/* Nível/Ladder */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Nível <span className="text-red-400">*</span>
                </label>
                <select
                  value={formData.ladder_id}
                  onChange={(e) => setFormData({ ...formData, ladder_id: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                >
                  <option value="">Selecione um nível</option>
                  {ladders.map((ladder) => (
                    <option key={ladder.id} value={ladder.id}>
                      {ladder.level} - {ladder.xp_reward} XP
                    </option>
                  ))}
                </select>
              </div>

              {/* Pilares */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Pilares <span className="text-red-400">*</span>
                </label>
                <div className="grid grid-cols-2 gap-3">
                  {Object.values(Pillar).map((pillar) => (
                    <button
                      key={pillar}
                      type="button"
                      onClick={() => handlePillarToggle(pillar)}
                      className={`px-4 py-3 rounded-lg border-2 transition-all ${
                        formData.pillars.includes(pillar)
                          ? "bg-blue-600 border-blue-500 text-white"
                          : "bg-gray-700 border-gray-600 text-gray-300 hover:border-gray-500"
                      }`}
                    >
                      {pillar}
                    </button>
                  ))}
                </div>
              </div>

              {/* Progresso */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Progresso: {formData.progress_percentage}%
                </label>
                <input
                  type="range"
                  min="0"
                  max="100"
                  value={formData.progress_percentage}
                  onChange={(e) => setFormData({ ...formData, progress_percentage: Number(e.target.value) })}
                  className="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer"
                />
              </div>

              {/* Resumo do Impacto */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Resumo do Impacto
                </label>
                <textarea
                  value={formData.impact_summary}
                  onChange={(e) => setFormData({ ...formData, impact_summary: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Qual foi o impacto desta atividade?"
                  rows={2}
                />
              </div>

              {/* PDI Target */}
              <div className="flex items-center gap-3">
                <input
                  type="checkbox"
                  id="pdi_target"
                  checked={formData.is_pdi_target}
                  onChange={(e) => setFormData({ ...formData, is_pdi_target: e.target.checked })}
                  className="w-5 h-5 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-2 focus:ring-blue-500"
                />
                <label htmlFor="pdi_target" className="text-sm font-medium text-gray-300">
                  Esta atividade faz parte do meu PDI
                </label>
              </div>

              {/* Botões */}
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
                  {submitting ? "Salvando..." : "Criar Atividade"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Modal de Edição de Atividade */}
      {showEditModal && editingActivity && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto border border-gray-700">
            <div className="sticky top-0 bg-gray-800 border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-2xl font-bold text-white">Editar Atividade</h2>
              <button
                onClick={() => {
                  setShowEditModal(false);
                  setEditingActivity(null);
                }}
                className="text-gray-400 hover:text-white transition-colors"
              >
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div className="p-6 space-y-6">
              {/* Título */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Título <span className="text-red-400">*</span>
                </label>
                <input
                  type="text"
                  value={editFormData.title}
                  onChange={(e) => setEditFormData({ ...editFormData, title: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Título da atividade"
                />
              </div>

              {/* Descrição */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Descrição
                </label>
                <textarea
                  value={editFormData.description}
                  onChange={(e) => setEditFormData({ ...editFormData, description: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Descreva a atividade..."
                  rows={3}
                />
              </div>

              {/* Progresso */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Progresso: {editFormData.progress_percentage}%
                </label>
                <input
                  type="range"
                  min="0"
                  max="100"
                  value={editFormData.progress_percentage}
                  onChange={(e) => setEditFormData({ ...editFormData, progress_percentage: Number(e.target.value) })}
                  className="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer"
                />
                <div className="flex justify-between text-xs text-gray-500 mt-2">
                  <span>0%</span>
                  <span>50%</span>
                  <span>100%</span>
                </div>
              </div>

              {/* Resumo do Impacto */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Resumo do Impacto
                </label>
                <textarea
                  value={editFormData.impact_summary}
                  onChange={(e) => setEditFormData({ ...editFormData, impact_summary: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Qual foi o impacto desta atividade?"
                  rows={2}
                />
              </div>

              {/* PDI Target */}
              <div className="flex items-center gap-3">
                <input
                  type="checkbox"
                  id="edit_pdi_target"
                  checked={editFormData.is_pdi_target}
                  onChange={(e) => setEditFormData({ ...editFormData, is_pdi_target: e.target.checked })}
                  className="w-5 h-5 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-2 focus:ring-blue-500"
                />
                <label htmlFor="edit_pdi_target" className="text-sm font-medium text-gray-300">
                  Esta atividade faz parte do meu PDI
                </label>
              </div>

              {/* Preview da Barra de Progresso */}
              <div>
                <p className="text-sm text-gray-400 mb-2">Preview do Progresso</p>
                <div className="w-full bg-gray-700 rounded-full h-3">
                  <div
                    className={`h-3 rounded-full transition-all ${
                      editFormData.progress_percentage === 100 
                        ? 'bg-green-600' 
                        : editFormData.progress_percentage >= 50 
                        ? 'bg-blue-600' 
                        : 'bg-yellow-600'
                    }`}
                    style={{ width: `${editFormData.progress_percentage}%` }}
                  ></div>
                </div>
              </div>

              {/* Botões */}
              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowEditModal(false);
                    setEditingActivity(null);
                  }}
                  className="flex-1 px-6 py-3 bg-gray-700 text-white rounded-lg font-semibold hover:bg-gray-600 transition-all"
                >
                  Cancelar
                </button>
                <button
                  onClick={handleUpdateActivity}
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

      {/* Modal de Evidências */}
      {showEvidenceModal && selectedActivityForEvidence && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-3xl w-full max-h-[90vh] overflow-y-auto border border-gray-700">
            <div className="sticky top-0 bg-gray-800 border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-bold text-white">Evidências</h2>
                <p className="text-sm text-gray-400 mt-1">{selectedActivityForEvidence.title}</p>
              </div>
              <button
                onClick={() => {
                  setShowEvidenceModal(false);
                  setSelectedActivityForEvidence(null);
                  // Recarrega a lista para atualizar contagem
                  fetchActivities();
                }}
                className="text-gray-400 hover:text-white transition-colors"
              >
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div className="p-6 space-y-6">
              {/* Formulário para Adicionar Evidência */}
              <div className="bg-gray-700 p-4 rounded-lg space-y-4">
                <h3 className="text-lg font-semibold text-white">Adicionar Nova Evidência</h3>
                
                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-2">
                    URL <span className="text-red-400">*</span>
                  </label>
                  <input
                    type="url"
                    value={evidenceForm.url}
                    onChange={(e) => setEvidenceForm({ ...evidenceForm, url: e.target.value })}
                    className="w-full px-4 py-3 bg-gray-600 border border-gray-500 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="https://exemplo.com/documento"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-2">
                    Descrição
                  </label>
                  <textarea
                    value={evidenceForm.description}
                    onChange={(e) => setEvidenceForm({ ...evidenceForm, description: e.target.value })}
                    className="w-full px-4 py-3 bg-gray-600 border border-gray-500 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Descreva esta evidência..."
                    rows={2}
                  />
                </div>

                <button
                  onClick={handleAddEvidence}
                  disabled={submitting}
                  className="w-full px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                >
                  {submitting ? (
                    "Adicionando..."
                  ) : (
                    <>
                      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                      </svg>
                      Adicionar Evidência
                    </>
                  )}
                </button>
              </div>

              {/* Lista de Evidências */}
              <div>
                <h3 className="text-lg font-semibold text-white mb-4">Evidências Cadastradas</h3>
                {activityEvidences[selectedActivityForEvidence.id]?.length > 0 ? (
                  <div className="space-y-3">
                    {activityEvidences[selectedActivityForEvidence.id].map((evidence, index) => (
                      <div key={evidence.id} className="bg-gray-700 p-4 rounded-lg border border-gray-600">
                        <div className="flex items-start gap-3">
                          <div className="shrink-0 w-10 h-10 bg-blue-900 text-blue-300 rounded-lg flex items-center justify-center">
                            {index + 1}
                          </div>
                          <div className="flex-1 min-w-0">
                            <a
                              href={evidence.evidence_url}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="text-blue-400 hover:text-blue-300 hover:underline text-sm break-all flex items-center gap-2"
                            >
                              {evidence.evidence_url}
                              <svg className="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                              </svg>
                            </a>
                            {evidence.description && (
                              <p className="text-gray-400 text-sm mt-2">{evidence.description}</p>
                            )}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8 bg-gray-700 rounded-lg">
                    <svg className="w-12 h-12 text-gray-500 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                    </svg>
                    <p className="text-gray-400">Nenhuma evidência cadastrada ainda</p>
                  </div>
                )}
              </div>
            </div>

            <div className="sticky bottom-0 bg-gray-800 border-t border-gray-700 px-6 py-4">
              <button
                onClick={() => {
                  setShowEvidenceModal(false);
                  setSelectedActivityForEvidence(null);
                  // Recarrega a lista para atualizar contagem
                  fetchActivities();
                }}
                className="w-full px-6 py-3 bg-gray-700 text-white rounded-lg font-semibold hover:bg-gray-600 transition-all"
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
