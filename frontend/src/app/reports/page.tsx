"use client";

import { useEffect, useState } from "react";
import Sidebar from "@/components/Sidebar";
import PageHeader from "@/components/PageHeader";
import { api } from "@/services/api";
import { useToast } from "@/components/Toast";
import type { GapAnalysisResponse, CareerRadar, ComparisonReport } from "@/types";

export default function ReportsPage() {
  const { toast } = useToast();
  const [activeTab, setActiveTab] = useState<"detailed" | "gap" | "radar" | "comparison">("detailed");
  const [gapData, setGapData] = useState<GapAnalysisResponse[]>([]);
  const [radarData, setRadarData] = useState<CareerRadar | null>(null);
  const [comparisonData, setComparisonData] = useState<ComparisonReport | null>(null);
  const [detailedData, setDetailedData] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [sendingToManager, setSendingToManager] = useState(false);
  const [hasManager, setHasManager] = useState(false);
  const [showConfirmSend, setShowConfirmSend] = useState(false);

  useEffect(() => { loadData(); checkManagerInfo(); }, [activeTab]);

  useEffect(() => {
    const handleFocus = () => { checkManagerInfo(); };
    window.addEventListener("focus", handleFocus);
    return () => window.removeEventListener("focus", handleFocus);
  }, []);

  const checkManagerInfo = async () => {
    try {
      if (typeof window !== "undefined") {
        const userId = localStorage.getItem("user_id");
        if (userId) {
          const userData = await api.get(`/users/${userId}`);
          setHasManager(!!(userData.manager_email && userData.manager_email.trim() !== ""));
        }
      }
    } catch { setHasManager(false); }
  };

  const loadData = async () => {
    setLoading(true);
    setError("");
    try {
      if (activeTab === "detailed") {
        const r = await api.get("/report");
        setDetailedData(Array.isArray(r) ? r : []);
      } else if (activeTab === "gap") {
        const r = await api.get(`/gap-analysis?year=${new Date().getFullYear()}`);
        setGapData(Array.isArray(r) ? r : []);
      } else if (activeTab === "radar") {
        setRadarData(await api.get("/career-radar"));
      } else if (activeTab === "comparison") {
        setComparisonData(await api.get("/cycle-comparison"));
      }
    } catch {
      setError("Erro ao carregar relatorio. Tente novamente.");
    } finally {
      setLoading(false);
    }
  };

  const handleDownloadPDF = async () => {
    try {
      const token = localStorage.getItem("career_token");
      const baseUrl = window.location.hostname !== "localhost" && window.location.hostname !== "127.0.0.1" ? "/api" : "http://localhost:8081/v1";
      const response = await fetch(`${baseUrl}/report/pdf`, { headers: { Authorization: `Bearer ${token}` } });
      if (!response.ok) throw new Error("Erro ao baixar PDF");
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `relatorio_${new Date().toISOString().split("T")[0]}.pdf`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch {
      toast("Erro ao baixar relatorio em PDF", "error");
    }
  };

  const handleSendToManager = async () => {
    if (!hasManager) { toast("Cadastre seu gerente no perfil antes de enviar.", "warning"); window.location.href = "/profile"; return; }
    setShowConfirmSend(true);
  };

  const confirmSendToManager = async () => {
    setShowConfirmSend(false);
    try {
      setSendingToManager(true);
      const response = await api.post("/report/send-to-manager", {});
      toast(response.message || "Relatorio enviado com sucesso!");
    } catch (error: any) {
      toast(error.message || "Erro ao enviar relatorio para o gerente", "error");
    } finally {
      setSendingToManager(false);
    }
  };

  const tabs = [
    { id: "detailed" as const, label: "Relatorio", icon: "📄" },
    { id: "gap" as const, label: "Gap Analysis", icon: "📊" },
    { id: "radar" as const, label: "Career Radar", icon: "🎯" },
    { id: "comparison" as const, label: "Ciclos", icon: "📈" },
  ];

  // Stats for detailed report
  const totalTasks = detailedData.length;
  const completedTasks = detailedData.filter((d) => d.progress_percentage === 100).length;
  const totalXP = detailedData.filter((d) => d.progress_percentage === 100).reduce((sum: number, d: any) => sum + (d.xp_reward || 0), 0);
  const avgProgress = totalTasks > 0 ? Math.round(detailedData.reduce((sum: number, d: any) => sum + d.progress_percentage, 0) / totalTasks) : 0;

  const EmptyState = ({ title, subtitle }: { title: string; subtitle: string }) => (
    <div className="bg-gray-800 border border-gray-700 rounded-xl p-16 text-center">
      <div className="w-16 h-16 bg-gray-700 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg className="w-8 h-8 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
      </div>
      <h3 className="text-lg font-semibold text-white mb-2">{title}</h3>
      <p className="text-gray-400 text-sm">{subtitle}</p>
    </div>
  );

  return (
    <div className="flex min-h-screen bg-gray-950">
      <Sidebar />
      <main className="flex-1 ml-64 p-8">
        <PageHeader
          title="Relatorios"
          subtitle="Analises e insights sobre seu desenvolvimento"
          action={
            activeTab === "detailed" && detailedData.length > 0 ? (
              <div className="flex gap-3">
                <button onClick={handleSendToManager} disabled={sendingToManager} className={`px-5 py-2.5 ${hasManager ? "bg-green-600 hover:bg-green-700" : "bg-gray-600 hover:bg-gray-700"} text-white rounded-lg font-medium transition-all flex items-center gap-2 text-sm ${sendingToManager ? "opacity-50" : ""}`}>
                  {sendingToManager ? <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div> : <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>}
                  {hasManager ? "Enviar" : "Cadastrar Gerente"}
                </button>
                <button onClick={handleDownloadPDF} className="px-5 py-2.5 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-all flex items-center gap-2 text-sm">
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
                  Baixar PDF
                </button>
              </div>
            ) : undefined
          }
        />

        {/* Tabs */}
        <div className="flex gap-1 mb-8 bg-gray-800 rounded-lg p-1 border border-gray-700 w-fit">
          {tabs.map((tab) => (
            <button key={tab.id} onClick={() => setActiveTab(tab.id)} className={`px-5 py-2.5 rounded-md text-sm font-medium transition-all ${activeTab === tab.id ? "bg-blue-600 text-white shadow" : "text-gray-400 hover:text-white hover:bg-gray-700"}`}>
              <span className="mr-1.5">{tab.icon}</span>{tab.label}
            </button>
          ))}
        </div>

        {loading ? (
          <div className="text-center py-16"><div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500 mx-auto mb-4"></div><p className="text-gray-400 text-sm">Carregando...</p></div>
        ) : error ? (
          <div className="bg-gray-800 border border-red-900 rounded-xl p-8 text-center">
            <p className="text-red-400 mb-4">{error}</p>
            <button onClick={loadData} className="px-5 py-2.5 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700">Tentar Novamente</button>
          </div>
        ) : (
          <>
            {/* ===== DETAILED REPORT ===== */}
            {activeTab === "detailed" && (
              detailedData.length === 0 ? (
                <EmptyState title="Sem dados" subtitle="Cadastre tasks nas suas iniciativas para gerar o relatorio" />
              ) : (
                <div className="space-y-6">
                  {/* Summary cards */}
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div className="bg-gray-800 border border-gray-700 rounded-xl p-5">
                      <p className="text-xs text-gray-400 uppercase tracking-wide mb-1">Total Tasks</p>
                      <p className="text-3xl font-bold text-white">{totalTasks}</p>
                    </div>
                    <div className="bg-gray-800 border border-gray-700 rounded-xl p-5">
                      <p className="text-xs text-gray-400 uppercase tracking-wide mb-1">Concluidas</p>
                      <p className="text-3xl font-bold text-green-400">{completedTasks}</p>
                    </div>
                    <div className="bg-gray-800 border border-gray-700 rounded-xl p-5">
                      <p className="text-xs text-gray-400 uppercase tracking-wide mb-1">XP Conquistado</p>
                      <p className="text-3xl font-bold text-blue-400">{totalXP}</p>
                    </div>
                    <div className="bg-gray-800 border border-gray-700 rounded-xl p-5">
                      <p className="text-xs text-gray-400 uppercase tracking-wide mb-1">Progresso Medio</p>
                      <p className="text-3xl font-bold text-yellow-400">{avgProgress}%</p>
                    </div>
                  </div>

                  {/* Task list */}
                  <div className="space-y-3">
                    {detailedData.map((item: any, index: number) => (
                      <div key={index} className="bg-gray-800 border border-gray-700 rounded-xl p-5 hover:border-gray-600 transition-all">
                        <div className="flex items-start justify-between gap-4">
                          <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-2 mb-2 flex-wrap">
                              <h3 className="text-base font-semibold text-white">{item.title}</h3>
                              <span className="px-2 py-0.5 bg-gray-700 text-gray-300 text-xs rounded font-medium">{item.level}</span>
                              {item.is_pdi_target && <span className="px-2 py-0.5 bg-blue-900 text-blue-300 text-xs rounded">PDI</span>}
                              {item.pillars && Array.isArray(item.pillars) && item.pillars.map((p: string) => (
                                <span key={p} className="px-2 py-0.5 bg-purple-900 text-purple-300 text-xs rounded">{p}</span>
                              ))}
                            </div>
                            {/* Progress bar */}
                            <div className="flex items-center gap-3">
                              <div className="flex-1 bg-gray-700 rounded-full h-1.5">
                                <div className={`h-1.5 rounded-full transition-all ${item.progress_percentage === 100 ? "bg-green-500" : item.progress_percentage >= 50 ? "bg-blue-500" : "bg-yellow-500"}`} style={{ width: `${item.progress_percentage}%` }}></div>
                              </div>
                              <span className="text-xs text-gray-400 min-w-8">{item.progress_percentage}%</span>
                            </div>
                          </div>
                          <div className="text-right shrink-0">
                            <p className="text-lg font-bold text-blue-400">{item.xp_reward}</p>
                            <p className="text-xs text-gray-500">XP</p>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )
            )}

            {/* ===== GAP ANALYSIS ===== */}
            {activeTab === "gap" && (
              gapData.length === 0 ? (
                <EmptyState title="Sem dados de gap" subtitle="Marque iniciativas como PDI e cadastre tasks para ver a analise" />
              ) : (
                <div className="space-y-4">
                  {gapData.map((gap, i) => (
                    <div key={i} className="bg-gray-800 border border-gray-700 rounded-xl p-5">
                      <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-3">
                          <span className="px-3 py-1 bg-blue-900 text-blue-300 text-sm font-semibold rounded-full">{gap.level}</span>
                          <span className="text-sm text-gray-300 capitalize">{gap.pillar.toLowerCase()}</span>
                        </div>
                        <span className={`px-3 py-1 rounded-full text-xs font-medium ${gap.status === "DONE" ? "bg-green-900 text-green-300" : gap.status === "CRITICAL" ? "bg-red-900 text-red-300" : "bg-yellow-900 text-yellow-300"}`}>
                          {gap.status === "DONE" ? "Atingido" : gap.status === "CRITICAL" ? "Critico" : "Em progresso"}
                        </span>
                      </div>
                      <div className="grid grid-cols-3 gap-4 mb-4 text-center">
                        <div>
                          <p className="text-2xl font-bold text-white">{gap.target}</p>
                          <p className="text-xs text-gray-500">Meta XP</p>
                        </div>
                        <div>
                          <p className="text-2xl font-bold text-green-400">{gap.achieved}</p>
                          <p className="text-xs text-gray-500">Conquistado</p>
                        </div>
                        <div>
                          <p className={`text-2xl font-bold ${gap.gap <= 0 ? "text-green-400" : "text-orange-400"}`}>{gap.gap}</p>
                          <p className="text-xs text-gray-500">Gap</p>
                        </div>
                      </div>
                      <div className="relative">
                        <div className="w-full bg-gray-700 rounded-full h-2.5">
                          <div className={`h-2.5 rounded-full transition-all ${gap.status === "DONE" ? "bg-green-500" : gap.status === "CRITICAL" ? "bg-red-500" : "bg-yellow-500"}`} style={{ width: `${Math.min(gap.percentage, 100)}%` }}></div>
                        </div>
                        <p className="text-right text-xs text-gray-400 mt-1">{gap.percentage}%</p>
                      </div>
                    </div>
                  ))}
                </div>
              )
            )}

            {/* ===== CAREER RADAR ===== */}
            {activeTab === "radar" && (
              !radarData || !radarData.breakdown?.length ? (
                <EmptyState title="Sem dados do radar" subtitle="Conclua tasks para ver a distribuicao por nivel" />
              ) : (
                <div className="space-y-6">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="bg-gray-800 border border-gray-700 rounded-xl p-6 text-center">
                      <p className="text-xs text-gray-400 uppercase tracking-wide mb-2">Tasks Concluidas</p>
                      <p className="text-4xl font-bold text-white">{radarData.total_activities}</p>
                    </div>
                    <div className="bg-gray-800 border border-gray-700 rounded-xl p-6 text-center">
                      <p className="text-xs text-gray-400 uppercase tracking-wide mb-2">XP Total</p>
                      <p className="text-4xl font-bold text-blue-400">{radarData.total_xp}</p>
                    </div>
                  </div>

                  <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
                    <h3 className="text-base font-semibold text-white mb-5">Distribuicao por Nivel</h3>
                    <div className="space-y-4">
                      {radarData.breakdown.map((level, i) => (
                        <div key={i}>
                          <div className="flex items-center justify-between mb-1.5">
                            <div className="flex items-center gap-2">
                              <span className="px-2 py-0.5 bg-blue-900 text-blue-300 text-xs font-semibold rounded">{level.level_name}</span>
                              <span className="text-xs text-gray-400">{level.activity_count} tasks</span>
                            </div>
                            <span className="text-sm font-medium text-white">{level.total_xp} XP</span>
                          </div>
                          <div className="w-full bg-gray-700 rounded-full h-2">
                            <div className="bg-blue-500 h-2 rounded-full transition-all" style={{ width: `${level.xp_percent}%` }}></div>
                          </div>
                          <p className="text-right text-xs text-gray-500 mt-0.5">{level.xp_percent.toFixed(1)}% do XP total</p>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              )
            )}

            {/* ===== CYCLE COMPARISON ===== */}
            {activeTab === "comparison" && (
              !comparisonData || !comparisonData.level_evolution?.length ? (
                <EmptyState title="Sem dados de comparacao" subtitle="Necessario ter ciclos de avaliacao e tasks concluidas em periodos diferentes" />
              ) : (
                <div className="space-y-6">
                  <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
                    <div className="grid grid-cols-3 gap-6 text-center">
                      <div>
                        <p className="text-xs text-gray-400 uppercase tracking-wide mb-1">Ciclo Atual</p>
                        <p className="text-lg font-bold text-white">{comparisonData.current_cycle || "—"}</p>
                      </div>
                      <div>
                        <p className="text-xs text-gray-400 uppercase tracking-wide mb-1">Anterior</p>
                        <p className="text-lg font-bold text-gray-400">{comparisonData.previous_cycle || "—"}</p>
                      </div>
                      <div>
                        <p className="text-xs text-gray-400 uppercase tracking-wide mb-1">Crescimento</p>
                        <p className={`text-xl font-bold ${comparisonData.growth_xp >= 0 ? "text-green-400" : "text-red-400"}`}>
                          {comparisonData.growth_xp >= 0 ? "+" : ""}{comparisonData.growth_xp} XP
                        </p>
                        <p className="text-xs text-gray-500">{comparisonData.percent_change.toFixed(1)}%</p>
                      </div>
                    </div>
                  </div>

                  <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
                    <h3 className="text-base font-semibold text-white mb-4">Evolucao por Nivel</h3>
                    <div className="space-y-3">
                      {comparisonData.level_evolution.map((level, i) => (
                        <div key={i} className="flex items-center justify-between p-4 bg-gray-700 rounded-lg">
                          <span className="px-2.5 py-1 bg-blue-900 text-blue-300 text-xs font-semibold rounded">{level.level_name}</span>
                          <div className="flex items-center gap-4">
                            <span className="text-sm text-gray-400 w-16 text-right">{level.prev_xp} XP</span>
                            <svg className="w-4 h-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" /></svg>
                            <span className="text-sm font-medium text-white w-16">{level.current_xp} XP</span>
                            <span className={`px-2.5 py-1 rounded text-xs font-semibold min-w-14 text-center ${level.diff > 0 ? "bg-green-900 text-green-300" : level.diff < 0 ? "bg-red-900 text-red-300" : "bg-gray-600 text-gray-300"}`}>
                              {level.diff > 0 ? "+" : ""}{level.diff}
                            </span>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              )
            )}
          </>
        )}
      </main>

      {/* Confirm send modal */}
      {showConfirmSend && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-sm w-full border border-gray-700">
            <div className="p-6 text-center">
              <div className="w-12 h-12 bg-blue-900 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg className="w-6 h-6 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
              </div>
              <h3 className="text-base font-semibold text-white mb-1">Enviar relatorio?</h3>
              <p className="text-sm text-gray-400 mb-5">O PDF sera enviado por email para seu gerente.</p>
              <div className="flex gap-3">
                <button onClick={() => setShowConfirmSend(false)} className="flex-1 px-4 py-2.5 bg-gray-700 text-gray-300 rounded-lg text-sm font-medium hover:bg-gray-600">Cancelar</button>
                <button onClick={confirmSendToManager} className="flex-1 px-4 py-2.5 bg-blue-600 text-white rounded-lg text-sm font-semibold hover:bg-blue-700">Enviar</button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
