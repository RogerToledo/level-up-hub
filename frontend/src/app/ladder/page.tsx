"use client";

import { useEffect, useState } from "react";
import Sidebar from "@/components/Sidebar";
import PageHeader from "@/components/PageHeader";
import { api } from "@/services/api";
import { useToast } from "@/components/Toast";
import { CareerLadder, LadderLevel, CreateCareerLadderRequest } from "@/types";

const LADDER_LEVELS: LadderLevel[] = ["P1", "P2", "P3", "LT1", "LT2", "LT3", "LT4"];

const emptyForm: CreateCareerLadderRequest = {
  level: "P1",
  xp_reward: 0,
  technical: "",
  expected_results: "",
  leadership_scope: "",
};

function LadderCard({
  ladder,
  isAdmin,
  onEdit,
  onDelete,
}: {
  ladder: CareerLadder;
  isAdmin: boolean;
  onEdit: (l: CareerLadder) => void;
  onDelete: (id: string) => void;
}) {
  const [expanded, setExpanded] = useState(false);

  return (
    <div className="bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
      <div className="px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <span className="px-3 py-1 bg-blue-900 text-blue-300 text-sm font-semibold rounded-full">
            {ladder.level}
          </span>
          <span className="text-sm text-white font-medium">{ladder.xp_reward} XP</span>
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={() => setExpanded(!expanded)}
            className="px-3 py-1.5 text-sm text-gray-300 hover:text-white hover:bg-gray-700 rounded-lg transition-all flex items-center gap-1"
          >
            <svg
              className={`w-4 h-4 transition-transform ${expanded ? "rotate-180" : ""}`}
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
            </svg>
            {expanded ? "Ocultar" : "Ver detalhes"}
          </button>

          {isAdmin && (
            <>
              <button
                onClick={() => onEdit(ladder)}
                className="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded transition-all"
                title="Editar"
              >
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </button>
              <button
                onClick={() => onDelete(ladder.id)}
                className="p-2 text-gray-400 hover:text-red-400 hover:bg-gray-700 rounded transition-all"
                title="Deletar"
              >
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </>
          )}
        </div>
      </div>

      {expanded && (
        <div className="border-t border-gray-700 px-6 py-4 grid gap-4 md:grid-cols-3">
          <div>
            <h4 className="text-xs font-medium text-gray-400 uppercase mb-1">Competencias Tecnicas</h4>
            <p className="text-sm text-gray-200 whitespace-pre-wrap">{ladder.technical}</p>
          </div>
          <div>
            <h4 className="text-xs font-medium text-gray-400 uppercase mb-1">Resultados Esperados</h4>
            <p className="text-sm text-gray-200 whitespace-pre-wrap">{ladder.expected_results}</p>
          </div>
          <div>
            <h4 className="text-xs font-medium text-gray-400 uppercase mb-1">Escopo de Lideranca</h4>
            <p className="text-sm text-gray-200 whitespace-pre-wrap">{ladder.leadership_scope}</p>
          </div>
        </div>
      )}
    </div>
  );
}

export default function LadderPage() {
  const { toast } = useToast();
  const [ladders, setLadders] = useState<CareerLadder[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [editingLadder, setEditingLadder] = useState<CareerLadder | null>(null);
  const [formData, setFormData] = useState<CreateCareerLadderRequest>(emptyForm);
  const [editFormData, setEditFormData] = useState<CreateCareerLadderRequest>(emptyForm);
  const [isAdmin, setIsAdmin] = useState(false);

  useEffect(() => {
    fetchLadders();
    setIsAdmin(localStorage.getItem("user_role") === "admin");
  }, []);

  const fetchLadders = async () => {
    try {
      const response = await api.get("/ladders");
      console.log("[Ladder] Response type:", typeof response, "isArray:", Array.isArray(response), "value:", response);
      if (Array.isArray(response)) {
        setLadders(response);
      } else {
        console.warn("[Ladder] Response is not an array, setting empty");
        setLadders([]);
      }
    } catch (error) {
      console.error("[Ladder] Erro ao carregar níveis:", error);
      setLadders([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.technical.trim() || !formData.expected_results.trim() || !formData.leadership_scope.trim()) {
      toast("Preencha todos os campos obrigatorios", "warning");
      return;
    }

    setSubmitting(true);
    try {
      await api.post("/ladder", formData);
      setFormData(emptyForm);
      setShowModal(false);
      fetchLadders();
      toast("Nivel criado com sucesso!");
    } catch (error) {
      console.error("Erro ao criar nível:", error);
      toast("Erro ao criar nivel. Tente novamente.", "error");
    } finally {
      setSubmitting(false);
    }
  };

  const handleEditClick = (ladder: CareerLadder) => {
    setEditingLadder(ladder);
    setEditFormData({
      level: ladder.level,
      xp_reward: ladder.xp_reward,
      technical: ladder.technical,
      expected_results: ladder.expected_results,
      leadership_scope: ladder.leadership_scope,
    });
    setShowEditModal(true);
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingLadder) return;

    if (!editFormData.technical.trim() || !editFormData.expected_results.trim() || !editFormData.leadership_scope.trim()) {
      toast("Preencha todos os campos obrigatorios", "warning");
      return;
    }

    setSubmitting(true);
    try {
      await api.put(`/ladder/${editingLadder.id}`, editFormData);
      setShowEditModal(false);
      setEditingLadder(null);
      fetchLadders();
      toast("Nivel atualizado com sucesso!");
    } catch (error) {
      console.error("Erro ao atualizar nível:", error);
      toast("Erro ao atualizar nivel. Tente novamente.", "error");
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Tem certeza que deseja deletar este nível?")) return;

    try {
      await api.delete(`/ladder/${id}`);
      setLadders(ladders.filter((l) => l.id !== id));
      toast("Nivel deletado com sucesso!");
    } catch (error) {
      console.error("Erro ao deletar:", error);
      toast("Erro ao deletar nivel. Verifique se nao ha dependencias.", "error");
    }
  };

  const renderForm = (
    data: CreateCareerLadderRequest,
    setData: (d: CreateCareerLadderRequest) => void,
    onSubmit: (e: React.FormEvent) => void,
    title: string,
    onClose: () => void
  ) => (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-lg w-full border border-gray-700 max-h-[90vh] overflow-y-auto">
        <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between">
          <h2 className="text-2xl font-bold text-white">{title}</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white transition-colors"
          >
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form onSubmit={onSubmit} className="p-6 space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Nivel <span className="text-red-400">*</span>
              </label>
              <select
                value={data.level}
                onChange={(e) => setData({ ...data, level: e.target.value as LadderLevel })}
                className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {LADDER_LEVELS.map((level) => (
                  <option key={level} value={level}>{level}</option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                XP Reward <span className="text-red-400">*</span>
              </label>
              <input
                type="number"
                min={0}
                value={data.xp_reward}
                onChange={(e) => setData({ ...data, xp_reward: Number(e.target.value) })}
                className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Competencias Tecnicas <span className="text-red-400">*</span>
            </label>
            <textarea
              value={data.technical}
              onChange={(e) => setData({ ...data, technical: e.target.value })}
              rows={3}
              className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
              placeholder="Descreva as competencias tecnicas esperadas..."
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Resultados Esperados <span className="text-red-400">*</span>
            </label>
            <textarea
              value={data.expected_results}
              onChange={(e) => setData({ ...data, expected_results: e.target.value })}
              rows={3}
              className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
              placeholder="Descreva os resultados esperados..."
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Escopo de Lideranca <span className="text-red-400">*</span>
            </label>
            <textarea
              value={data.leadership_scope}
              onChange={(e) => setData({ ...data, leadership_scope: e.target.value })}
              rows={3}
              className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
              placeholder="Descreva o escopo de lideranca..."
              required
            />
          </div>

          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-3 bg-gray-700 text-gray-300 rounded-lg font-medium hover:bg-gray-600 transition-all"
            >
              Cancelar
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="flex-1 px-4 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all disabled:opacity-50"
            >
              {submitting ? "Salvando..." : "Salvar"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );

  return (
    <div className="flex min-h-screen bg-gray-950">
      <Sidebar />

      <main className="flex-1 ml-64 p-8">
        <PageHeader
          title="Niveis de Carreira"
          subtitle="Visualize os niveis da career ladder"
          action={
            isAdmin ? (
              <button
                onClick={() => setShowModal(true)}
                className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all flex items-center gap-2"
              >
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                </svg>
                Novo Nivel
              </button>
            ) : undefined
          }
        />

        {loading ? (
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
            <p className="text-gray-400">Carregando niveis...</p>
          </div>
        ) : ladders.length === 0 ? (
          <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
            <svg className="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
            </svg>
            <h3 className="text-xl font-semibold text-white mb-2">Nenhum nivel cadastrado</h3>
            <p className="text-gray-400 mb-6">Comece criando o primeiro nivel da career ladder</p>
            {isAdmin && (
              <button
                onClick={() => setShowModal(true)}
                className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all"
              >
                Criar Nivel
              </button>
            )}
          </div>
        ) : (
          <div className="grid gap-4">
            {ladders.map((ladder) => (
              <LadderCard
                key={ladder.id}
                ladder={ladder}
                isAdmin={isAdmin}
                onEdit={handleEditClick}
                onDelete={handleDelete}
              />
            ))}
          </div>
        )}
      </main>

      {showModal && renderForm(formData, setFormData, handleSubmit, "Novo Nivel", () => setShowModal(false))}
      {showEditModal && renderForm(editFormData, setEditFormData, handleUpdate, "Editar Nivel", () => setShowEditModal(false))}
    </div>
  );
}
