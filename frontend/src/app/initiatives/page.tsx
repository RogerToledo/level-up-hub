"use client";

import { useEffect, useState } from "react";
import Sidebar from "@/components/Sidebar";
import PageHeader from "@/components/PageHeader";
import { api } from "@/services/api";
import { useToast } from "@/components/Toast";
import { Pillar, Initiative, Task, TaskEvidence } from "@/types";

interface Ladder {
  id: string;
  level: string;
  xp_reward: number;
}

interface InitiativeForm {
  title: string;
  description: string;
  pillars: string[];
  impact_summary: string;
  is_pdi_target: boolean;
  ladder_id: string;
}

const emptyForm: InitiativeForm = {
  title: "",
  description: "",
  pillars: [],
  impact_summary: "",
  is_pdi_target: false,
  ladder_id: "",
};

export default function InitiativesPage() {
  const { toast } = useToast();
  const [initiatives, setInitiatives] = useState<Initiative[]>([]);
  const [ladders, setLadders] = useState<Ladder[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [formData, setFormData] = useState<InitiativeForm>(emptyForm);
  const [editFormData, setEditFormData] = useState<InitiativeForm>(emptyForm);
  const [editingInitiative, setEditingInitiative] = useState<Initiative | null>(null);
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [tasks, setTasks] = useState<Record<string, Task[]>>({});
  const [showTaskModal, setShowTaskModal] = useState(false);
  const [taskForm, setTaskForm] = useState({ title: "", execution: "", progress_percentage: 0 });
  const [taskInitiativeId, setTaskInitiativeId] = useState<string>("");
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [showEvidenceModal, setShowEvidenceModal] = useState(false);
  const [evidenceTaskId, setEvidenceTaskId] = useState<string>("");
  const [evidences, setEvidences] = useState<TaskEvidence[]>([]);
  const [evidenceForm, setEvidenceForm] = useState({ url: "", description: "" });

  useEffect(() => {
    fetchInitiatives();
    fetchLadders();
  }, []);

  const fetchInitiatives = async () => {
    try {
      const response = await api.get("/initiatives");
      setInitiatives(Array.isArray(response) ? response : []);
    } catch (error) {
      console.error("Erro ao carregar iniciativas:", error);
      setInitiatives([]);
    } finally {
      setLoading(false);
    }
  };

  const fetchLadders = async () => {
    try {
      const response = await api.get("/ladders");
      setLadders(Array.isArray(response) ? response : []);
    } catch (error) {
      console.error("Erro ao carregar niveis:", error);
      setLadders([]);
    }
  };

  const fetchTasks = async (initiativeId: string) => {
    try {
      const response = await api.get(`/initiatives/${initiativeId}/tasks`);
      setTasks(prev => ({ ...prev, [initiativeId]: Array.isArray(response) ? response : [] }));
    } catch (error) {
      console.error("Erro ao carregar tasks:", error);
      setTasks(prev => ({ ...prev, [initiativeId]: [] }));
    }
  };

  const toggleExpand = (id: string) => {
    if (expandedId === id) {
      setExpandedId(null);
    } else {
      setExpandedId(id);
      if (!tasks[id]) fetchTasks(id);
    }
  };

  // Initiative CRUD
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.title || formData.pillars.length === 0 || !formData.ladder_id) {
      toast("Preencha os campos obrigatorios: Titulo, Pilares e Nivel", "warning");
      return;
    }
    setSubmitting(true);
    try {
      const userId = localStorage.getItem("user_id");
      if (!userId) { toast("Faca login novamente.", "error"); return; }
      await api.post("/initiatives", {
        user_id: userId,
        ladder_id: formData.ladder_id,
        pillars: formData.pillars,
        title: formData.title,
        description: formData.description || undefined,
        progress_percentage: 0,
        impact_summary: formData.impact_summary || undefined,
        is_pdi_target: formData.is_pdi_target,
      });
      setFormData(emptyForm);
      setShowModal(false);
      fetchInitiatives();
      toast("Iniciativa criada com sucesso!");
    } catch (error) {
      console.error("Erro ao criar iniciativa:", error);
      toast("Erro ao criar iniciativa. Tente novamente.", "error");
    } finally {
      setSubmitting(false);
    }
  };

  const handleEditClick = async (init: Initiative) => {
    setEditingInitiative(init);
    setEditFormData({
      title: init.title,
      description: init.description || "",
      impact_summary: init.impact_summary || "",
      is_pdi_target: init.is_pdi_target,
      ladder_id: init.ladder_id,
      pillars: [],
    });
    setShowEditModal(true);
    try {
      const pillars = await api.get(`/initiatives/${init.id}/pillars`);
      setEditFormData(prev => ({ ...prev, pillars: Array.isArray(pillars) ? pillars : [] }));
    } catch {}
  };

  const handleUpdate = async () => {
    if (!editingInitiative) return;
    if (!editFormData.title.trim()) { toast("O titulo e obrigatorio", "warning"); return; }
    setSubmitting(true);
    try {
      await api.put(`/initiatives/${editingInitiative.id}`, {
        title: editFormData.title,
        description: editFormData.description || undefined,
        progress_percentage: editingInitiative.progress_percentage,
        impact_summary: editFormData.impact_summary || undefined,
        is_pdi_target: editFormData.is_pdi_target,
        ladder_id: editFormData.ladder_id || undefined,
        pillars: editFormData.pillars.length > 0 ? editFormData.pillars : undefined,
      });
      setShowEditModal(false);
      setEditingInitiative(null);
      fetchInitiatives();
      toast("Iniciativa atualizada com sucesso!");
    } catch (error) {
      console.error("Erro ao atualizar:", error);
      toast("Erro ao atualizar iniciativa.", "error");
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Tem certeza que deseja deletar esta iniciativa?")) return;
    try {
      await api.delete(`/initiatives/${id}`);
      setInitiatives(initiatives.filter(i => i.id !== id));
      toast("Iniciativa deletada com sucesso!");
    } catch (error) {
      console.error("Erro ao deletar:", error);
      toast("Erro ao deletar iniciativa.", "error");
    }
  };

  // Task CRUD
  const openTaskModal = (initiativeId: string, task?: Task) => {
    setTaskInitiativeId(initiativeId);
    if (task) {
      setEditingTask(task);
      setTaskForm({ title: task.title, execution: task.execution || "", progress_percentage: task.progress_percentage });
    } else {
      setEditingTask(null);
      setTaskForm({ title: "", execution: "", progress_percentage: 0 });
    }
    setShowTaskModal(true);
  };

  const handleTaskSubmit = async () => {
    if (!taskForm.title.trim()) { toast("O titulo da task e obrigatorio", "warning"); return; }
    setSubmitting(true);
    try {
      if (editingTask) {
        await api.put(`/tasks/${editingTask.id}`, taskForm);
        toast("Task atualizada com sucesso!");
      } else {
        await api.post("/tasks", { ...taskForm, initiative_id: taskInitiativeId });
        toast("Task criada com sucesso!");
      }
      setShowTaskModal(false);
      fetchTasks(taskInitiativeId);
      fetchInitiatives(); // refresh progress
    } catch (error) {
      console.error("Erro na task:", error);
      toast("Erro ao salvar task.", "error");
    } finally {
      setSubmitting(false);
    }
  };

  const handleTaskDelete = async (taskId: string, initiativeId: string) => {
    if (!confirm("Deletar esta task?")) return;
    try {
      await api.delete(`/tasks/${taskId}`);
      fetchTasks(initiativeId);
      fetchInitiatives();
      toast("Task deletada!");
    } catch (error) {
      console.error("Erro ao deletar task:", error);
      toast("Erro ao deletar task.", "error");
    }
  };

  // Evidence
  const openEvidenceModal = async (taskId: string) => {
    setEvidenceTaskId(taskId);
    setEvidenceForm({ url: "", description: "" });
    setShowEvidenceModal(true);
    try {
      const response = await api.get(`/tasks/${taskId}/evidences`);
      setEvidences(Array.isArray(response) ? response : []);
    } catch {
      setEvidences([]);
    }
  };

  const handleAddEvidence = async () => {
    if (!evidenceForm.url.trim()) { toast("A URL e obrigatoria", "warning"); return; }
    try { new URL(evidenceForm.url); } catch { toast("URL invalida", "warning"); return; }
    setSubmitting(true);
    try {
      await api.post(`/tasks/${evidenceTaskId}/evidence`, {
        url: evidenceForm.url,
        description: evidenceForm.description || undefined,
      });
      setEvidenceForm({ url: "", description: "" });
      const response = await api.get(`/tasks/${evidenceTaskId}/evidences`);
      setEvidences(Array.isArray(response) ? response : []);
      // refresh task evidence count
      if (expandedId) fetchTasks(expandedId);
      toast("Evidencia adicionada!");
    } catch (error) {
      console.error("Erro ao adicionar evidencia:", error);
      toast("Erro ao adicionar evidencia.", "error");
    } finally {
      setSubmitting(false);
    }
  };

  const handlePillarToggle = (pillar: string, isEdit = false) => {
    if (isEdit) {
      setEditFormData(prev => ({
        ...prev,
        pillars: prev.pillars.includes(pillar) ? prev.pillars.filter(p => p !== pillar) : [...prev.pillars, pillar],
      }));
    } else {
      setFormData(prev => ({
        ...prev,
        pillars: prev.pillars.includes(pillar) ? prev.pillars.filter(p => p !== pillar) : [...prev.pillars, pillar],
      }));
    }
  };

  const formatDate = (d: string) => d ? new Date(d).toLocaleDateString("pt-BR", { day: "2-digit", month: "short", year: "numeric" }) : "";

  return (
    <div className="flex min-h-screen bg-gray-950">
      <Sidebar />
      <main className="flex-1 ml-64 p-8">
        <PageHeader
          title="Iniciativas"
          subtitle="Gerencie suas iniciativas de desenvolvimento"
          action={
            <button onClick={() => setShowModal(true)} className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all flex items-center gap-2">
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" /></svg>
              Nova Iniciativa
            </button>
          }
        />

        {loading ? (
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
            <p className="text-gray-400">Carregando iniciativas...</p>
          </div>
        ) : initiatives.length === 0 ? (
          <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
            <h3 className="text-xl font-semibold text-white mb-2">Nenhuma iniciativa cadastrada</h3>
            <p className="text-gray-400 mb-6">Comece criando sua primeira iniciativa</p>
            <button onClick={() => setShowModal(true)} className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all">Criar Iniciativa</button>
          </div>
        ) : (
          <div className="grid gap-4">
            {initiatives.map((init) => (
              <div key={init.id} className="bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
                {/* Initiative header */}
                <div className="p-6">
                  <div className="flex items-start justify-between mb-3">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <h3 className="text-lg font-semibold text-white">{init.title}</h3>
                        {init.is_pdi_target && <span className="px-2 py-0.5 bg-blue-900 text-blue-300 text-xs rounded">PDI</span>}
                        {init.progress_percentage === 100 && <span className="px-2 py-0.5 bg-green-900 text-green-300 text-xs rounded">Completo</span>}
                      </div>
                      {init.description && <p className="text-gray-400 text-sm">{init.description}</p>}
                    </div>
                    <div className="flex gap-1">
                      <button onClick={() => handleEditClick(init)} className="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded" title="Editar">
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                      </button>
                      <button onClick={() => handleDelete(init.id)} className="p-2 text-gray-400 hover:text-red-400 hover:bg-gray-700 rounded" title="Deletar">
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                      </button>
                    </div>
                  </div>
                  {/* Progress bar */}
                  <div className="flex items-center gap-3 mb-3">
                    <div className="flex-1 bg-gray-700 rounded-full h-2">
                      <div className={`h-2 rounded-full ${init.progress_percentage === 100 ? "bg-green-600" : init.progress_percentage >= 50 ? "bg-blue-600" : "bg-yellow-600"}`} style={{ width: `${init.progress_percentage}%` }}></div>
                    </div>
                    <span className="text-sm text-gray-300 min-w-12">{init.progress_percentage}%</span>
                  </div>
                  {/* Expand tasks button */}
                  <button onClick={() => toggleExpand(init.id)} className="text-sm text-gray-400 hover:text-white flex items-center gap-1">
                    <svg className={`w-4 h-4 transition-transform ${expandedId === init.id ? "rotate-180" : ""}`} fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
                    {init.task_count} task{init.task_count !== 1 ? "s" : ""} - {expandedId === init.id ? "Ocultar" : "Ver tasks"}
                  </button>
                </div>

                {/* Expanded tasks section */}
                {expandedId === init.id && (
                  <div className="border-t border-gray-700 bg-gray-850 px-6 py-4">
                    <div className="flex items-center justify-between mb-3">
                      <h4 className="text-sm font-medium text-gray-300">Tasks</h4>
                      <button onClick={() => openTaskModal(init.id)} className="text-xs px-3 py-1.5 bg-blue-600 text-white rounded hover:bg-blue-700 flex items-center gap-1">
                        <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" /></svg>
                        Nova Task
                      </button>
                    </div>
                    {(tasks[init.id] || []).length === 0 ? (
                      <p className="text-gray-500 text-sm py-4 text-center">Nenhuma task cadastrada</p>
                    ) : (
                      <div className="space-y-2">
                        {(tasks[init.id] || []).map((task) => (
                          <div key={task.id} className="bg-gray-700 rounded-lg p-3">
                            <div className="flex items-center justify-between mb-1">
                              <span className="text-sm font-medium text-white">{task.title}</span>
                              <div className="flex items-center gap-1">
                                {task.evidence_count > 0 && (
                                  <span className="text-xs text-purple-300 bg-purple-900 px-1.5 py-0.5 rounded">{task.evidence_count} ev</span>
                                )}
                                <button onClick={() => openEvidenceModal(task.id)} className="p-1 text-gray-400 hover:text-blue-400 rounded" title="Evidencias">
                                  <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" /></svg>
                                </button>
                                <button onClick={() => openTaskModal(init.id, task)} className="p-1 text-gray-400 hover:text-white rounded" title="Editar">
                                  <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                                </button>
                                <button onClick={() => handleTaskDelete(task.id, init.id)} className="p-1 text-gray-400 hover:text-red-400 rounded" title="Deletar">
                                  <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                                </button>
                              </div>
                            </div>
                            {task.execution && <p className="text-xs text-gray-400 mb-1">{task.execution}</p>}
                            <div className="flex items-center gap-2">
                              <div className="flex-1 bg-gray-600 rounded-full h-1.5">
                                <div className={`h-1.5 rounded-full ${task.progress_percentage === 100 ? "bg-green-500" : "bg-blue-500"}`} style={{ width: `${task.progress_percentage}%` }}></div>
                              </div>
                              <span className="text-xs text-gray-400">{task.progress_percentage}%</span>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </main>

      {/* Initiative Create/Edit Modal */}
      {(showModal || showEditModal) && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-lg w-full max-h-[90vh] overflow-y-auto border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-xl font-bold text-white">{showEditModal ? "Editar Iniciativa" : "Nova Iniciativa"}</h2>
              <button onClick={() => { setShowModal(false); setShowEditModal(false); setEditingInitiative(null); }} className="text-gray-400 hover:text-white">
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <form onSubmit={showEditModal ? (e) => { e.preventDefault(); handleUpdate(); } : handleSubmit} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Titulo *</label>
                <input type="text" value={showEditModal ? editFormData.title : formData.title} onChange={(e) => showEditModal ? setEditFormData({ ...editFormData, title: e.target.value }) : setFormData({ ...formData, title: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500" placeholder="Ex: Migracao para microservicos" required />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Descricao</label>
                <textarea value={showEditModal ? editFormData.description : formData.description} onChange={(e) => showEditModal ? setEditFormData({ ...editFormData, description: e.target.value }) : setFormData({ ...formData, description: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500" rows={2} />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Nivel *</label>
                <select value={showEditModal ? editFormData.ladder_id : formData.ladder_id} onChange={(e) => showEditModal ? setEditFormData({ ...editFormData, ladder_id: e.target.value }) : setFormData({ ...formData, ladder_id: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500" required>
                  <option value="">Selecione</option>
                  {ladders.map(l => <option key={l.id} value={l.id}>{l.level} - {l.xp_reward} XP</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Pilares *</label>
                <div className="grid grid-cols-3 gap-2">
                  {Object.values(Pillar).map(p => (
                    <button key={p} type="button" onClick={() => handlePillarToggle(p, showEditModal)} className={`px-3 py-2 rounded-lg border text-sm transition-all ${(showEditModal ? editFormData : formData).pillars.includes(p) ? "bg-blue-600 border-blue-500 text-white" : "bg-gray-700 border-gray-600 text-gray-300"}`}>{p}</button>
                  ))}
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Resumo do Impacto</label>
                <textarea value={showEditModal ? editFormData.impact_summary : formData.impact_summary} onChange={(e) => showEditModal ? setEditFormData({ ...editFormData, impact_summary: e.target.value }) : setFormData({ ...formData, impact_summary: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500" rows={2} />
              </div>
              <div className="flex items-center gap-3">
                <input type="checkbox" id="pdi" checked={showEditModal ? editFormData.is_pdi_target : formData.is_pdi_target} onChange={(e) => showEditModal ? setEditFormData({ ...editFormData, is_pdi_target: e.target.checked }) : setFormData({ ...formData, is_pdi_target: e.target.checked })} className="w-4 h-4 text-blue-600 bg-gray-700 border-gray-600 rounded" />
                <label htmlFor="pdi" className="text-sm text-gray-300">Faz parte do meu PDI</label>
              </div>
              <div className="flex gap-3 pt-2">
                <button type="button" onClick={() => { setShowModal(false); setShowEditModal(false); setEditingInitiative(null); }} className="flex-1 px-4 py-3 bg-gray-700 text-gray-300 rounded-lg hover:bg-gray-600">Cancelar</button>
                <button type="submit" disabled={submitting} className="flex-1 px-4 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 disabled:opacity-50">{submitting ? "Salvando..." : "Salvar"}</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Task Modal */}
      {showTaskModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-xl font-bold text-white">{editingTask ? "Editar Task" : "Nova Task"}</h2>
              <button onClick={() => setShowTaskModal(false)} className="text-gray-400 hover:text-white">
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Titulo *</label>
                <input type="text" value={taskForm.title} onChange={(e) => setTaskForm({ ...taskForm, title: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500" placeholder="O que foi feito?" />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Execucao (o que foi feito)</label>
                <textarea value={taskForm.execution} onChange={(e) => setTaskForm({ ...taskForm, execution: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500" rows={3} placeholder="Descreva em detalhes o que foi executado..." />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Progresso: {taskForm.progress_percentage}%</label>
                <input type="range" min="0" max="100" value={taskForm.progress_percentage} onChange={(e) => setTaskForm({ ...taskForm, progress_percentage: Number(e.target.value) })} className="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer" />
              </div>
              <div className="flex gap-3 pt-2">
                <button type="button" onClick={() => setShowTaskModal(false)} className="flex-1 px-4 py-3 bg-gray-700 text-gray-300 rounded-lg hover:bg-gray-600">Cancelar</button>
                <button onClick={handleTaskSubmit} disabled={submitting} className="flex-1 px-4 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 disabled:opacity-50">{submitting ? "Salvando..." : "Salvar"}</button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Evidence Modal */}
      {showEvidenceModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-lg w-full max-h-[80vh] overflow-y-auto border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-xl font-bold text-white">Evidencias</h2>
              <button onClick={() => setShowEvidenceModal(false)} className="text-gray-400 hover:text-white">
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div className="bg-gray-700 p-4 rounded-lg space-y-3">
                <input type="url" value={evidenceForm.url} onChange={(e) => setEvidenceForm({ ...evidenceForm, url: e.target.value })} className="w-full px-4 py-2 bg-gray-600 border border-gray-500 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm" placeholder="https://..." />
                <input type="text" value={evidenceForm.description} onChange={(e) => setEvidenceForm({ ...evidenceForm, description: e.target.value })} className="w-full px-4 py-2 bg-gray-600 border border-gray-500 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm" placeholder="Descricao (opcional)" />
                <button onClick={handleAddEvidence} disabled={submitting} className="w-full px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-semibold hover:bg-blue-700 disabled:opacity-50">{submitting ? "Adicionando..." : "Adicionar Evidencia"}</button>
              </div>
              {evidences.length > 0 ? (
                <div className="space-y-2">
                  {evidences.map((ev, i) => (
                    <div key={ev.id} className="bg-gray-700 p-3 rounded-lg flex items-start gap-2">
                      <span className="text-xs text-gray-400 mt-0.5">{i + 1}.</span>
                      <div className="flex-1 min-w-0">
                        <a href={ev.evidence_url} target="_blank" rel="noopener noreferrer" className="text-blue-400 hover:underline text-sm break-all">{ev.evidence_url}</a>
                        {ev.description && <p className="text-xs text-gray-400 mt-1">{ev.description}</p>}
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-gray-500 text-sm text-center py-4">Nenhuma evidencia cadastrada</p>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
