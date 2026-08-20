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
  is_pdi_target: boolean;
}

interface TaskForm {
  title: string;
  execution: string;
  impact_summary: string;
  ladder_id: string;
  pillars: string[];
  progress_percentage: number;
  is_extra: boolean;
}

const emptyInitForm: InitiativeForm = { title: "", description: "", is_pdi_target: false };
const emptyTaskForm: TaskForm = { title: "", execution: "", impact_summary: "", ladder_id: "", pillars: [], progress_percentage: 0, is_extra: false };

export default function InitiativesPage() {
  const { toast } = useToast();
  const [initiatives, setInitiatives] = useState<Initiative[]>([]);
  const [ladders, setLadders] = useState<Ladder[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  // Initiative modals
  const [showInitModal, setShowInitModal] = useState(false);
  const [showEditInitModal, setShowEditInitModal] = useState(false);
  const [initForm, setInitForm] = useState<InitiativeForm>(emptyInitForm);
  const [editInitForm, setEditInitForm] = useState<InitiativeForm>(emptyInitForm);
  const [editingInit, setEditingInit] = useState<Initiative | null>(null);
  // Task state
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [tasks, setTasks] = useState<Record<string, Task[]>>({});
  const [showTaskModal, setShowTaskModal] = useState(false);
  const [taskForm, setTaskForm] = useState<TaskForm>(emptyTaskForm);
  const [taskInitiativeId, setTaskInitiativeId] = useState("");
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  // Evidence state
  const [showEvidenceModal, setShowEvidenceModal] = useState(false);
  const [evidenceTaskId, setEvidenceTaskId] = useState("");
  const [evidences, setEvidences] = useState<TaskEvidence[]>([]);
  const [evidenceForm, setEvidenceForm] = useState({ url: "", description: "" });
  const [editingEvidence, setEditingEvidence] = useState<TaskEvidence | null>(null);
  // View state
  const [showViewInitModal, setShowViewInitModal] = useState(false);
  const [viewingInit, setViewingInit] = useState<Initiative | null>(null);
  const [showViewTaskModal, setShowViewTaskModal] = useState(false);
  const [viewingTask, setViewingTask] = useState<Task | null>(null);

  // AI Classification state
  const [showClassifyModal, setShowClassifyModal] = useState(false);
  const [classifying, setClassifying] = useState(false);
  const [classifyForm, setClassifyForm] = useState({ title: "", execution: "" });
  const [classifyResult, setClassifyResult] = useState<{ level: string; pillars: string[] } | null>(null);

  useEffect(() => { fetchInitiatives(); fetchLadders(); }, []);

  const fetchInitiatives = async () => {
    try {
      const r = await api.get("/initiatives");
      setInitiatives(Array.isArray(r) ? r : []);
    } catch { setInitiatives([]); } finally { setLoading(false); }
  };
  const fetchLadders = async () => {
    try {
      const r = await api.get("/ladders");
      setLadders(Array.isArray(r) ? r : []);
    } catch { setLadders([]); }
  };
  const fetchTasks = async (id: string) => {
    try {
      const r = await api.get(`/initiatives/${id}/tasks`);
      setTasks(prev => ({ ...prev, [id]: Array.isArray(r) ? r : [] }));
    } catch { setTasks(prev => ({ ...prev, [id]: [] })); }
  };
  const toggleExpand = (id: string) => {
    if (expandedId === id) { setExpandedId(null); } else { setExpandedId(id); if (!tasks[id]) fetchTasks(id); }
  };

  // Initiative CRUD
  const handleCreateInit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!initForm.title.trim()) { toast("Titulo e obrigatorio", "warning"); return; }
    setSubmitting(true);
    try {
      const userId = localStorage.getItem("user_id");
      if (!userId) { toast("Faca login novamente.", "error"); return; }
      await api.post("/initiatives", { user_id: userId, title: initForm.title, description: initForm.description || undefined, is_pdi_target: initForm.is_pdi_target });
      setInitForm(emptyInitForm);
      setShowInitModal(false);
      fetchInitiatives();
      toast("Iniciativa criada com sucesso!");
    } catch { toast("Erro ao criar iniciativa.", "error"); } finally { setSubmitting(false); }
  };
  const handleEditInitClick = (init: Initiative) => {
    setEditingInit(init);
    setEditInitForm({ title: init.title, description: init.description || "", is_pdi_target: init.is_pdi_target });
    setShowEditInitModal(true);
  };
  const handleUpdateInit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingInit || !editInitForm.title.trim()) { toast("Titulo e obrigatorio", "warning"); return; }
    setSubmitting(true);
    try {
      await api.put(`/initiatives/${editingInit.id}`, { title: editInitForm.title, description: editInitForm.description || undefined, is_pdi_target: editInitForm.is_pdi_target });
      setShowEditInitModal(false); setEditingInit(null);
      fetchInitiatives();
      toast("Iniciativa atualizada!");
    } catch { toast("Erro ao atualizar.", "error"); } finally { setSubmitting(false); }
  };
  const handleDeleteInit = async (id: string) => {
    if (!confirm("Deletar esta iniciativa e todas as suas tasks?")) return;
    try { await api.delete(`/initiatives/${id}`); setInitiatives(prev => prev.filter(i => i.id !== id)); toast("Iniciativa deletada!"); } catch { toast("Erro ao deletar.", "error"); }
  };

  // Task CRUD
  const openTaskModal = async (initiativeId: string, task?: Task) => {
    setTaskInitiativeId(initiativeId);
    if (task) {
      setEditingTask(task);
      setTaskForm({ title: task.title, execution: task.execution || "", impact_summary: task.impact_summary || "", ladder_id: task.ladder_id || "", pillars: [], progress_percentage: task.progress_percentage, is_extra: task.is_extra || false });
      try { const p = await api.get(`/tasks/${task.id}/pillars`); setTaskForm(prev => ({ ...prev, pillars: Array.isArray(p) ? p : [] })); } catch {}
    } else {
      setEditingTask(null);
      setTaskForm(emptyTaskForm);
    }
    setShowTaskModal(true);
  };
  const handleTaskSubmit = async () => {
    if (!taskForm.title.trim() || !taskForm.ladder_id || taskForm.pillars.length === 0) { toast("Preencha titulo, nivel e pilares", "warning"); return; }
    setSubmitting(true);
    try {
      if (editingTask) {
        await api.put(`/tasks/${editingTask.id}`, { title: taskForm.title, ladder_id: taskForm.ladder_id, pillars: taskForm.pillars, execution: taskForm.execution || undefined, impact_summary: taskForm.impact_summary || undefined, progress_percentage: taskForm.progress_percentage, is_extra: taskForm.is_extra });
        toast("Task atualizada!");
      } else {
        await api.post("/tasks", { initiative_id: taskInitiativeId, ladder_id: taskForm.ladder_id, pillars: taskForm.pillars, title: taskForm.title, execution: taskForm.execution || undefined, impact_summary: taskForm.impact_summary || undefined, progress_percentage: taskForm.progress_percentage, is_extra: taskForm.is_extra });
        toast("Task criada!");
      }
      setShowTaskModal(false);
      fetchTasks(taskInitiativeId);
      fetchInitiatives();
    } catch { toast("Erro ao salvar task.", "error"); } finally { setSubmitting(false); }
  };
  const handleTaskDelete = async (taskId: string, initiativeId: string) => {
    if (!confirm("Deletar esta task?")) return;
    try { await api.delete(`/tasks/${taskId}`); fetchTasks(initiativeId); fetchInitiatives(); toast("Task deletada!"); } catch { toast("Erro ao deletar task.", "error"); }
  };

  // Evidence
  const openEvidenceModal = async (taskId: string) => {
    setEvidenceTaskId(taskId); setEvidenceForm({ url: "", description: "" }); setEditingEvidence(null); setShowEvidenceModal(true);
    try { const r = await api.get(`/tasks/${taskId}/evidences`); setEvidences(Array.isArray(r) ? r : []); } catch { setEvidences([]); }
  };
  const handleAddEvidence = async () => {
    if (!evidenceForm.url.trim()) { toast("URL e obrigatoria", "warning"); return; }
    try { new URL(evidenceForm.url); } catch { toast("URL invalida", "warning"); return; }
    setSubmitting(true);
    try {
      if (editingEvidence) {
        await api.put(`/evidences/${editingEvidence.id}`, { url: evidenceForm.url, description: evidenceForm.description || undefined });
        toast("Evidencia atualizada!");
      } else {
        await api.post(`/tasks/${evidenceTaskId}/evidence`, { url: evidenceForm.url, description: evidenceForm.description || undefined });
        toast("Evidencia adicionada!");
      }
      setEvidenceForm({ url: "", description: "" });
      setEditingEvidence(null);
      const r = await api.get(`/tasks/${evidenceTaskId}/evidences`); setEvidences(Array.isArray(r) ? r : []);
      if (expandedId) fetchTasks(expandedId);
    } catch { toast("Erro ao salvar evidencia.", "error"); } finally { setSubmitting(false); }
  };
  const handleEditEvidence = (ev: TaskEvidence) => {
    setEditingEvidence(ev);
    setEvidenceForm({ url: ev.evidence_url, description: ev.description || "" });
  };
  const handleDeleteEvidence = async (evidenceId: string) => {
    if (!confirm("Deletar esta evidencia?")) return;
    try {
      await api.delete(`/evidences/${evidenceId}`);
      setEvidences(prev => prev.filter(e => e.id !== evidenceId));
      if (expandedId) fetchTasks(expandedId);
      toast("Evidencia deletada!");
    } catch { toast("Erro ao deletar evidencia.", "error"); }
  };

  const handlePillarToggle = (pillar: string) => {
    setTaskForm(prev => ({ ...prev, pillars: prev.pillars.includes(pillar) ? prev.pillars.filter(p => p !== pillar) : [...prev.pillars, pillar] }));
  };

  const openClassifyModal = () => {
    setClassifyForm({ title: taskForm.title, execution: taskForm.execution });
    setClassifyResult(null);
    setShowClassifyModal(true);
  };

  const handleClassify = async () => {
    if (!classifyForm.title.trim()) { toast("Preencha o titulo para classificar", "warning"); return; }
    setClassifying(true);
    try {
      const result = await api.post("/tasks/classify", { title: classifyForm.title, execution: classifyForm.execution || undefined });
      setClassifyResult(result as { level: string; pillars: string[] });
    } catch {
      toast("Erro ao classificar. Tente novamente.", "error");
    } finally {
      setClassifying(false);
    }
  };

  const confirmClassification = () => {
    if (!classifyResult) return;
    // Find ladder_id by level name
    const ladder = ladders.find(l => l.level === classifyResult.level);
    setTaskForm(prev => ({
      ...prev,
      ladder_id: ladder?.id || prev.ladder_id,
      pillars: classifyResult.pillars || prev.pillars,
    }));
    setShowClassifyModal(false);
    toast("Nivel e pilares preenchidos!");
  };

  return (
    <div className="flex min-h-screen bg-gray-950">
      <Sidebar />
      <main className="flex-1 ml-64 p-8">
        <PageHeader title="Iniciativas" subtitle="Gerencie suas iniciativas de desenvolvimento" action={
          <button onClick={() => setShowInitModal(true)} className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all flex items-center gap-2">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" /></svg>
            Nova Iniciativa
          </button>
        } />

        {loading ? (
          <div className="text-center py-12"><div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div><p className="text-gray-400">Carregando...</p></div>
        ) : initiatives.length === 0 ? (
          <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
            <h3 className="text-xl font-semibold text-white mb-2">Nenhuma iniciativa cadastrada</h3>
            <p className="text-gray-400 mb-6">Comece criando sua primeira iniciativa</p>
            <button onClick={() => setShowInitModal(true)} className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700">Criar Iniciativa</button>
          </div>
        ) : (
          <div className="grid gap-4">
            {initiatives.map((init) => (
              <div key={init.id} className="bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
                <div className="p-6">
                  <div className="flex items-start justify-between mb-3">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <h3 className="text-lg font-semibold text-white">{init.title}</h3>
                        {init.is_pdi_target && <span className="px-2 py-0.5 bg-blue-900 text-blue-300 text-xs rounded">PDI</span>}
                        {init.has_extra && <span className="px-2 py-0.5 bg-orange-900 text-orange-300 text-xs rounded">Extra</span>}
                        {init.progress_percentage === 100 && <span className="px-2 py-0.5 bg-green-900 text-green-300 text-xs rounded">Completo</span>}
                      </div>
                      {init.description && <p className="text-gray-400 text-sm">{init.description}</p>}
                    </div>
                    <div className="flex gap-1">
                      <button onClick={() => { setViewingInit(init); setShowViewInitModal(true); }} className="p-2 text-gray-400 hover:text-green-400 hover:bg-gray-700 rounded" title="Visualizar">
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" /></svg>
                      </button>
                      <button onClick={() => handleEditInitClick(init)} className="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded" title="Editar">
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                      </button>
                      <button onClick={() => handleDeleteInit(init.id)} className="p-2 text-gray-400 hover:text-red-400 hover:bg-gray-700 rounded" title="Deletar">
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                      </button>
                    </div>
                  </div>
                  <div className="flex items-center gap-3 mb-3">
                    <div className="flex-1 bg-gray-700 rounded-full h-2">
                      <div className={`h-2 rounded-full ${init.progress_percentage === 100 ? "bg-green-600" : init.progress_percentage >= 50 ? "bg-blue-600" : "bg-yellow-600"}`} style={{ width: `${init.progress_percentage}%` }}></div>
                    </div>
                    <span className="text-sm text-gray-300 min-w-12">{init.progress_percentage}%</span>
                  </div>
                  <button onClick={() => toggleExpand(init.id)} className="text-sm text-gray-400 hover:text-white flex items-center gap-1">
                    <svg className={`w-4 h-4 transition-transform ${expandedId === init.id ? "rotate-180" : ""}`} fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
                    {init.task_count} task{init.task_count !== 1 ? "s" : ""} - {expandedId === init.id ? "Ocultar" : "Ver tasks"}
                  </button>
                </div>

                {/* Expanded tasks */}
                {expandedId === init.id && (
                  <div className="border-t border-gray-700 px-6 py-4">
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
                              <div className="flex items-center gap-2">
                                <span className="text-sm font-medium text-white">{task.title}</span>
                                {task.is_extra && <span className="text-xs text-orange-300 bg-orange-900 px-1.5 py-0.5 rounded">Extra</span>}
                                {task.evidence_count > 0 && <span className="text-xs text-purple-300 bg-purple-900 px-1.5 py-0.5 rounded">{task.evidence_count} ev</span>}
                              </div>
                              <div className="flex items-center gap-1">
                                <button onClick={() => openEvidenceModal(task.id)} className="p-1 text-gray-400 hover:text-blue-400 rounded" title="Evidencias">
                                  <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" /></svg>
                                </button>
                                <button onClick={() => { setViewingTask(task); setShowViewTaskModal(true); }} className="p-1 text-gray-400 hover:text-green-400 rounded" title="Visualizar">
                                  <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" /></svg>
                                </button>
                                <button onClick={() => openTaskModal(init.id, task)} className="p-1 text-gray-400 hover:text-white rounded" title="Editar">
                                  <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                                </button>
                                <button onClick={() => handleTaskDelete(task.id, init.id)} className="p-1 text-gray-400 hover:text-red-400 rounded" title="Deletar">
                                  <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                                </button>
                              </div>
                            </div>
                            {task.execution && <p className="text-xs text-gray-400 mb-1"><span className="text-gray-500">Execucao:</span> {task.execution}</p>}
                            {task.impact_summary && <p className="text-xs text-gray-400 mb-1"><span className="text-gray-500">Impacto:</span> {task.impact_summary}</p>}
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
      {(showInitModal || showEditInitModal) && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-xl font-bold text-white">{showEditInitModal ? "Editar Iniciativa" : "Nova Iniciativa"}</h2>
              <button onClick={() => { setShowInitModal(false); setShowEditInitModal(false); setEditingInit(null); }} className="text-gray-400 hover:text-white">
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <form onSubmit={showEditInitModal ? handleUpdateInit : handleCreateInit} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Titulo *</label>
                <input type="text" value={showEditInitModal ? editInitForm.title : initForm.title} onChange={(e) => showEditInitModal ? setEditInitForm({ ...editInitForm, title: e.target.value }) : setInitForm({ ...initForm, title: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500" placeholder="Ex: Migracao para microservicos" required />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Descricao</label>
                <textarea value={showEditInitModal ? editInitForm.description : initForm.description} onChange={(e) => showEditInitModal ? setEditInitForm({ ...editInitForm, description: e.target.value }) : setInitForm({ ...initForm, description: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500" rows={3} />
              </div>
              <div className="flex items-center gap-3">
                <input type="checkbox" id="pdi" checked={showEditInitModal ? editInitForm.is_pdi_target : initForm.is_pdi_target} onChange={(e) => showEditInitModal ? setEditInitForm({ ...editInitForm, is_pdi_target: e.target.checked }) : setInitForm({ ...initForm, is_pdi_target: e.target.checked })} className="w-4 h-4 text-blue-600 bg-gray-700 border-gray-600 rounded" />
                <label htmlFor="pdi" className="text-sm text-gray-300">Faz parte do meu PDI</label>
              </div>
              <div className="flex gap-3 pt-2">
                <button type="button" onClick={() => { setShowInitModal(false); setShowEditInitModal(false); setEditingInit(null); }} className="flex-1 px-4 py-3 bg-gray-700 text-gray-300 rounded-lg hover:bg-gray-600">Cancelar</button>
                <button type="submit" disabled={submitting} className="flex-1 px-4 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 disabled:opacity-50">{submitting ? "Salvando..." : "Salvar"}</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Task Modal */}
      {showTaskModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-lg w-full max-h-[90vh] overflow-y-auto border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-xl font-bold text-white">{editingTask ? "Editar Task" : "Nova Task"}</h2>
              <button onClick={() => setShowTaskModal(false)} className="text-gray-400 hover:text-white">
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Titulo *</label>
                <input type="text" value={taskForm.title} onChange={(e) => setTaskForm({ ...taskForm, title: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500" placeholder="Titulo da atividade" />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Nivel *</label>
                <select value={taskForm.ladder_id} onChange={(e) => setTaskForm({ ...taskForm, ladder_id: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500">
                  <option value="">Selecione</option>
                  {ladders.map(l => <option key={l.id} value={l.id}>{l.level} - {l.xp_reward} XP</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Pilares *</label>
                <div className="grid grid-cols-3 gap-2">
                  {Object.values(Pillar).map(p => (
                    <button key={p} type="button" onClick={() => handlePillarToggle(p)} className={`px-3 py-2 rounded-lg border text-sm transition-all ${taskForm.pillars.includes(p) ? "bg-blue-600 border-blue-500 text-white" : "bg-gray-700 border-gray-600 text-gray-300"}`}>{p}</button>
                  ))}
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Execução</label>
                <textarea value={taskForm.execution} onChange={(e) => setTaskForm({ ...taskForm, execution: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500" rows={3} placeholder="O que foi feito?" />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Resumo do Impacto</label>
                <textarea value={taskForm.impact_summary} onChange={(e) => setTaskForm({ ...taskForm, impact_summary: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500" rows={2} placeholder="Qual foi o impacto?" />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Progresso: {taskForm.progress_percentage}%</label>
                <input type="range" min="0" max="100" value={taskForm.progress_percentage} onChange={(e) => setTaskForm({ ...taskForm, progress_percentage: Number(e.target.value) })} className="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer" />
              </div>
              <div className="flex items-center gap-3">
                <input type="checkbox" id="is_extra" checked={taskForm.is_extra} onChange={(e) => setTaskForm({ ...taskForm, is_extra: e.target.checked })} className="w-4 h-4 text-orange-600 bg-gray-700 border-gray-600 rounded" />
                <label htmlFor="is_extra" className="text-sm text-gray-300">Task extra (overdelivery)</label>
              </div>
              <div className="flex gap-3 pt-2">
                <button type="button" disabled className="px-4 py-3 bg-gray-600 text-gray-400 rounded-lg text-sm font-medium cursor-not-allowed flex items-center gap-1 opacity-50" title="IA indisponivel no momento">
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
                  IA
                </button>
                <button type="button" onClick={() => setShowTaskModal(false)} className="flex-1 px-4 py-3 bg-gray-700 text-gray-300 rounded-lg hover:bg-gray-600">Cancelar</button>
                <button onClick={handleTaskSubmit} disabled={submitting} className="flex-1 px-4 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 disabled:opacity-50">{submitting ? "Salvando..." : "Salvar"}</button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* View Initiative Modal */}
      {showViewInitModal && viewingInit && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-xl font-bold text-white">Detalhes da Iniciativa</h2>
              <button onClick={() => { setShowViewInitModal(false); setViewingInit(null); }} className="text-gray-400 hover:text-white">
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <h3 className="text-lg font-semibold text-white mb-1">{viewingInit.title}</h3>
                <div className="flex gap-2 mb-3">
                  {viewingInit.is_pdi_target && <span className="px-2 py-0.5 bg-blue-900 text-blue-300 text-xs rounded">PDI</span>}
                  {viewingInit.has_extra && <span className="px-2 py-0.5 bg-orange-900 text-orange-300 text-xs rounded">Extra</span>}
                  {viewingInit.progress_percentage === 100 && <span className="px-2 py-0.5 bg-green-900 text-green-300 text-xs rounded">Completo</span>}
                </div>
              </div>
              {viewingInit.description && (
                <div className="bg-gray-700 rounded-lg p-3">
                  <p className="text-xs text-gray-400 mb-1">Descricao</p>
                  <p className="text-sm text-gray-200">{viewingInit.description}</p>
                </div>
              )}
              <div className="grid grid-cols-2 gap-3">
                <div className="bg-gray-700 rounded-lg p-3">
                  <p className="text-xs text-gray-400 mb-1">Progresso</p>
                  <p className="text-lg font-bold text-white">{viewingInit.progress_percentage}%</p>
                </div>
                <div className="bg-gray-700 rounded-lg p-3">
                  <p className="text-xs text-gray-400 mb-1">Tasks</p>
                  <p className="text-lg font-bold text-white">{viewingInit.task_count}</p>
                </div>
              </div>
              <div className="bg-gray-700 rounded-lg p-3">
                <p className="text-xs text-gray-400 mb-2">Progresso</p>
                <div className="w-full bg-gray-600 rounded-full h-2.5">
                  <div className={`h-2.5 rounded-full ${viewingInit.progress_percentage === 100 ? "bg-green-600" : viewingInit.progress_percentage >= 50 ? "bg-blue-600" : "bg-yellow-600"}`} style={{ width: `${viewingInit.progress_percentage}%` }}></div>
                </div>
              </div>
            </div>
            <div className="border-t border-gray-700 px-6 py-4">
              <button onClick={() => { setShowViewInitModal(false); setViewingInit(null); }} className="w-full px-4 py-3 bg-gray-700 text-gray-300 rounded-lg font-medium hover:bg-gray-600 transition-all">Fechar</button>
            </div>
          </div>
        </div>
      )}

      {/* View Task Modal */}
      {showViewTaskModal && viewingTask && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-lg w-full border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <h2 className="text-xl font-bold text-white">Detalhes da Task</h2>
              <button onClick={() => { setShowViewTaskModal(false); setViewingTask(null); }} className="text-gray-400 hover:text-white">
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <h3 className="text-lg font-semibold text-white mb-1">{viewingTask.title}</h3>
                <div className="flex gap-2">
                  {viewingTask.is_extra && <span className="px-2 py-0.5 bg-orange-900 text-orange-300 text-xs rounded">Extra</span>}
                  {viewingTask.progress_percentage === 100 && <span className="px-2 py-0.5 bg-green-900 text-green-300 text-xs rounded">Completo</span>}
                </div>
              </div>
              {viewingTask.execution && (
                <div className="bg-gray-700 rounded-lg p-3">
                  <p className="text-xs text-gray-400 mb-1">Execucao</p>
                  <p className="text-sm text-gray-200 whitespace-pre-wrap">{viewingTask.execution}</p>
                </div>
              )}
              {viewingTask.impact_summary && (
                <div className="bg-gray-700 rounded-lg p-3">
                  <p className="text-xs text-gray-400 mb-1">Resumo do Impacto</p>
                  <p className="text-sm text-gray-200 whitespace-pre-wrap">{viewingTask.impact_summary}</p>
                </div>
              )}
              <div className="grid grid-cols-2 gap-3">
                <div className="bg-gray-700 rounded-lg p-3">
                  <p className="text-xs text-gray-400 mb-1">Progresso</p>
                  <p className="text-lg font-bold text-white">{viewingTask.progress_percentage}%</p>
                </div>
                <div className="bg-gray-700 rounded-lg p-3">
                  <p className="text-xs text-gray-400 mb-1">Evidencias</p>
                  <p className="text-lg font-bold text-white">{viewingTask.evidence_count}</p>
                </div>
              </div>
              <div className="bg-gray-700 rounded-lg p-3">
                <p className="text-xs text-gray-400 mb-2">Progresso</p>
                <div className="w-full bg-gray-600 rounded-full h-2.5">
                  <div className={`h-2.5 rounded-full ${viewingTask.progress_percentage === 100 ? "bg-green-600" : "bg-blue-600"}`} style={{ width: `${viewingTask.progress_percentage}%` }}></div>
                </div>
              </div>
            </div>
            <div className="border-t border-gray-700 px-6 py-4">
              <button onClick={() => { setShowViewTaskModal(false); setViewingTask(null); }} className="w-full px-4 py-3 bg-gray-700 text-gray-300 rounded-lg font-medium hover:bg-gray-600 transition-all">Fechar</button>
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
                <div className="flex gap-2">
                  {editingEvidence && (
                    <button onClick={() => { setEditingEvidence(null); setEvidenceForm({ url: "", description: "" }); }} className="px-4 py-2 bg-gray-600 text-gray-300 rounded-lg text-sm hover:bg-gray-500">Cancelar</button>
                  )}
                  <button onClick={handleAddEvidence} disabled={submitting} className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-semibold hover:bg-blue-700 disabled:opacity-50">{submitting ? "Salvando..." : editingEvidence ? "Atualizar" : "Adicionar"}</button>
                </div>
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
                      <div className="flex gap-1 shrink-0">
                        <button onClick={() => handleEditEvidence(ev)} className="p-1 text-gray-400 hover:text-white rounded" title="Editar">
                          <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                        </button>
                        <button onClick={() => handleDeleteEvidence(ev.id)} className="p-1 text-gray-400 hover:text-red-400 rounded" title="Deletar">
                          <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                        </button>
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
      {/* Classify AI Modal */}
      {showClassifyModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[60] p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full border border-gray-700">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
              <div className="flex items-center gap-2">
                <svg className="w-5 h-5 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
                <h2 className="text-lg font-bold text-white">Classificar com IA</h2>
              </div>
              <button onClick={() => setShowClassifyModal(false)} className="text-gray-400 hover:text-white">
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Titulo da atividade *</label>
                <input type="text" value={classifyForm.title} onChange={(e) => setClassifyForm({ ...classifyForm, title: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500" placeholder="O que foi feito?" />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Resultados / Execucao</label>
                <textarea value={classifyForm.execution} onChange={(e) => setClassifyForm({ ...classifyForm, execution: e.target.value })} className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500" rows={3} placeholder="Descreva os resultados obtidos..." />
              </div>

              {!classifyResult && (
                <button onClick={handleClassify} disabled={classifying} className="w-full px-4 py-3 bg-purple-600 text-white rounded-lg font-semibold hover:bg-purple-700 disabled:opacity-50 flex items-center justify-center gap-2">
                  {classifying ? (
                    <><div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div> Classificando...</>
                  ) : (
                    <><svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" /></svg> Classificar</>
                  )}
                </button>
              )}

              {classifyResult && (
                <div className="space-y-3">
                  <div className="bg-gray-700 rounded-lg p-4 border border-purple-900">
                    <p className="text-xs text-gray-400 mb-2">Resultado da classificacao:</p>
                    <div className="flex items-center gap-3 mb-3">
                      <span className="text-sm text-gray-300">Nivel:</span>
                      <span className="px-3 py-1 bg-blue-900 text-blue-300 text-sm font-semibold rounded-full">{classifyResult.level}</span>
                    </div>
                    <div className="flex items-center gap-2 flex-wrap">
                      <span className="text-sm text-gray-300">Pilares:</span>
                      {classifyResult.pillars.map(p => (
                        <span key={p} className="px-2 py-1 bg-purple-900 text-purple-300 text-xs rounded">{p}</span>
                      ))}
                    </div>
                  </div>
                  <p className="text-xs text-gray-500 text-center">Voce pode alterar os valores no formulario apos confirmar</p>
                  <div className="flex gap-3">
                    <button onClick={() => setShowClassifyModal(false)} className="flex-1 px-4 py-3 bg-gray-700 text-gray-300 rounded-lg font-medium hover:bg-gray-600">Cancelar</button>
                    <button onClick={confirmClassification} className="flex-1 px-4 py-3 bg-purple-600 text-white rounded-lg font-semibold hover:bg-purple-700">Confirmar</button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
