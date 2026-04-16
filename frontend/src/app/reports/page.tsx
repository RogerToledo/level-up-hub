"use client";

import { useEffect, useState } from "react";
import Sidebar from "@/components/Sidebar";
import PageHeader from "@/components/PageHeader";
import { api } from "@/services/api";
import type { GapAnalysisResponse, CareerRadar, ComparisonReport } from "@/types";

export default function ReportsPage() {
  const [activeTab, setActiveTab] = useState<"detailed" | "gap" | "radar" | "comparison">("detailed");
  const [gapData, setGapData] = useState<GapAnalysisResponse[]>([]);
  const [radarData, setRadarData] = useState<CareerRadar | null>(null);
  const [comparisonData, setComparisonData] = useState<ComparisonReport | null>(null);
  const [detailedData, setDetailedData] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [sendingToManager, setSendingToManager] = useState(false);
  const [hasManager, setHasManager] = useState(false);

  useEffect(() => {
    loadData();
    checkManagerInfo();
  }, [activeTab]);

  // Recarrega info do gerente quando a janela recebe foco
  useEffect(() => {
    const handleFocus = () => {
      checkManagerInfo();
    };
    
    window.addEventListener('focus', handleFocus);
    return () => window.removeEventListener('focus', handleFocus);
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
    } catch (err) {
      console.error("Erro ao verificar gerente:", err);
      setHasManager(false);
    }
  };

  const loadData = async () => {
    setLoading(true);
    setError("");
    try {
      if (activeTab === "detailed") {
        const response = await api.get("/report");
        console.log("Detailed report response:", response);
        
        // Trata diferentes formatos de resposta
        if (Array.isArray(response)) {
          setDetailedData(response);
        } else if (response.items && Array.isArray(response.items)) {
          setDetailedData(response.items);
        } else if (response.data && Array.isArray(response.data)) {
          setDetailedData(response.data);
        } else {
          console.error("Formato de resposta inesperado:", response);
          setDetailedData([]);
        }
      } else if (activeTab === "gap") {
        const currentYear = new Date().getFullYear();
        const response = await api.get(`/gap-analysis?year=${currentYear}`);
        console.log("Gap analysis response:", response);
        
        if (Array.isArray(response)) {
          setGapData(response);
        } else if (response.items && Array.isArray(response.items)) {
          setGapData(response.items);
        } else if (response.data && Array.isArray(response.data)) {
          setGapData(response.data);
        } else {
          console.error("Formato de resposta inesperado:", response);
          setGapData([]);
        }
      } else if (activeTab === "radar") {
        const data = await api.get("/career-radar");
        setRadarData(data);
      } else if (activeTab === "comparison") {
        const data = await api.get("/cycle-comparison");
        setComparisonData(data);
      }
    } catch (err) {
      console.error("Erro ao carregar relatório:", err);
      setError("Erro ao carregar relatório. Tente novamente.");
    } finally {
      setLoading(false);
    }
  };

  const handleDownloadPDF = async () => {
    try {
      const token = localStorage.getItem("career_token");
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/report/pdf`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (!response.ok) {
        throw new Error('Erro ao baixar PDF');
      }
      
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `relatorio_${new Date().toISOString().split('T')[0]}.pdf`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (error) {
      console.error('Erro ao baixar PDF:', error);
      alert('Erro ao baixar relatório em PDF');
    }
  };

  const handleSendToManager = async () => {
    if (!hasManager) {
      if (confirm('Você não tem um gerente cadastrado. Deseja cadastrar agora?')) {
        window.location.href = '/profile';
      }
      return;
    }

    if (!confirm('Deseja enviar o relatório completo para seu gerente por email?')) {
      return;
    }

    try {
      setSendingToManager(true);
      const response = await api.post('/report/send-to-manager', {});
      alert(response.message || 'Relatório enviado com sucesso!');
    } catch (error: any) {
      console.error('Erro ao enviar relatório:', error);
      alert(error.message || 'Erro ao enviar relatório para o gerente');
    } finally {
      setSendingToManager(false);
    }
  };

  const tabs = [
    { id: "detailed" as const, label: "Relatório Detalhado", icon: "📄" },
    { id: "gap" as const, label: "Gap Analysis", icon: "📊" },
    { id: "radar" as const, label: "Career Radar", icon: "🎯" },
    { id: "comparison" as const, label: "Comparação de Ciclos", icon: "📈" },
  ];

  return (
    <div className="flex min-h-screen bg-gray-950">
      <Sidebar />
      
      <main className="flex-1 ml-64 p-8">
        <PageHeader 
          title="Relatórios" 
          subtitle="Análises e insights sobre seu desenvolvimento"
          action={
            activeTab === "detailed" && detailedData.length > 0 ? (
              <div className="flex gap-3">
                <button
                  onClick={handleSendToManager}
                  disabled={sendingToManager}
                  className={`px-6 py-3 ${
                    hasManager 
                      ? 'bg-green-600 hover:bg-green-700' 
                      : 'bg-gray-600 hover:bg-gray-700'
                  } text-white rounded-lg font-semibold transition-all flex items-center gap-2 ${
                    sendingToManager ? 'opacity-50 cursor-not-allowed' : ''
                  }`}
                  title={hasManager ? 'Enviar relatório por email' : 'Cadastre um gerente no Perfil primeiro'}
                >
                  {sendingToManager ? (
                    <>
                      <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                      Enviando...
                    </>
                  ) : (
                    <>
                      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                      </svg>
                      {hasManager ? 'Enviar para Gerente' : 'Cadastrar Gerente'}
                    </>
                  )}
                </button>
                <button
                  onClick={handleDownloadPDF}
                  className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all flex items-center gap-2"
                >
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  Baixar PDF
                </button>
              </div>
            ) : undefined
          }
        />

        {/* Tabs */}
        <div className="flex gap-2 mb-8 border-b border-gray-700">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-6 py-3 font-medium transition-all relative ${
                activeTab === tab.id
                  ? "text-blue-400"
                  : "text-gray-400 hover:text-gray-300"
              }`}
            >
              <span className="mr-2">{tab.icon}</span>
              {tab.label}
              {activeTab === tab.id && (
                <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-500"></div>
              )}
            </button>
          ))}
        </div>

        {loading ? (
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
            <p className="text-gray-400">Carregando relatório...</p>
          </div>
        ) : error ? (
          <div className="bg-gray-800 border border-red-700 rounded-lg p-8 text-center">
            <div className="w-16 h-16 bg-red-900 text-red-400 rounded-full flex items-center justify-center mx-auto mb-4 text-2xl">
              ⚠️
            </div>
            <p className="text-red-400 mb-4">{error}</p>
            <button
              onClick={loadData}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all"
            >
              Tentar Novamente
            </button>
          </div>
        ) : (
          <>
            {/* Detailed Report */}
            {activeTab === "detailed" && (
              detailedData.length === 0 ? (
                <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
                  <svg className="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  <h3 className="text-xl font-semibold text-white mb-2">Sem atividades cadastradas</h3>
                  <p className="text-gray-400">Cadastre atividades para visualizar o relatório detalhado</p>
                </div>
              ) : Array.isArray(detailedData) ? (
                <div className="space-y-4">
                  {detailedData.map((activity, index) => (
                    <div key={index} className="bg-gray-800 border border-gray-700 rounded-lg p-6">
                      <div className="flex items-start justify-between mb-4">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-2">
                            <h3 className="text-lg font-bold text-white">{activity.title}</h3>
                            {activity.is_pdi_target && (
                              <span className="px-2 py-1 bg-blue-900 text-blue-300 text-xs font-medium rounded">PDI</span>
                            )}
                          </div>
                          <p className="text-sm text-gray-400 mb-2">Nível: {activity.level_name}</p>
                          {activity.description && (
                            <p className="text-gray-400 text-sm">{activity.description}</p>
                          )}
                        </div>
                        <div className="text-right">
                          <p className="text-2xl font-bold text-blue-400">{activity.xp_reward} XP</p>
                          <p className="text-sm text-gray-500">
                            {activity.progress_percentage}% completo
                          </p>
                        </div>
                      </div>
                      
                      {activity.impact_summary && (
                        <div className="mt-4 p-4 bg-gray-700 rounded-lg">
                          <p className="text-sm font-medium text-gray-300 mb-1">Impacto:</p>
                          <p className="text-sm text-gray-400">{activity.impact_summary}</p>
                        </div>
                      )}
                      
                      <div className="mt-4 flex items-center gap-3">
                        <div className="flex-1 bg-gray-700 rounded-full h-2">
                          <div
                            className={`h-2 rounded-full transition-all ${
                              activity.progress_percentage === 100 
                                ? 'bg-green-600' 
                                : activity.progress_percentage >= 50 
                                ? 'bg-blue-600' 
                                : 'bg-yellow-600'
                            }`}
                            style={{ width: `${activity.progress_percentage}%` }}
                          ></div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
                  <div className="w-16 h-16 bg-red-900 text-red-400 rounded-full flex items-center justify-center mx-auto mb-4 text-2xl">
                    ⚠️
                  </div>
                  <h3 className="text-xl font-semibold text-white mb-2">Erro ao carregar dados</h3>
                  <p className="text-gray-400 mb-6">Formato de resposta inesperado</p>
                  <button 
                    onClick={loadData}
                    className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all"
                  >
                    Tentar Novamente
                  </button>
                </div>
              )
            )}

            {/* Gap Analysis */}
            {activeTab === "gap" && (
              gapData.length === 0 ? (
                <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
                  <svg className="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                  <h3 className="text-xl font-semibold text-white mb-2">Sem dados de análise de gap</h3>
                  <p className="text-gray-400">Cadastre atividades para visualizar este relatório</p>
                </div>
              ) : Array.isArray(gapData) ? (
              <div className="space-y-4">{gapData.map((gap, index) => (
                  <div key={index} className="bg-gray-800 border border-gray-700 rounded-lg p-6">
                    <div className="flex items-center justify-between mb-4">
                      <h3 className="text-xl font-bold text-white capitalize">{gap.pillar.toLowerCase()}</h3>
                      <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                        gap.status === "on_track" ? "bg-green-900 text-green-300" :
                        gap.status === "at_risk" ? "bg-yellow-900 text-yellow-300" :
                        "bg-red-900 text-red-300"
                      }`}>
                        {gap.status === "on_track" ? "No Caminho" : gap.status === "at_risk" ? "Em Risco" : "Crítico"}
                      </span>
                    </div>
                    <div className="grid grid-cols-3 gap-4 mb-4">
                      <div>
                        <p className="text-sm text-gray-400">Meta</p>
                        <p className="text-2xl font-bold text-white">{gap.target}</p>
                      </div>
                      <div>
                        <p className="text-sm text-gray-400">Conquistado</p>
                        <p className="text-2xl font-bold text-blue-400">{gap.achieved}</p>
                      </div>
                      <div>
                        <p className="text-sm text-gray-400">Gap</p>
                        <p className="text-2xl font-bold text-orange-400">{gap.gap}</p>
                      </div>
                    </div>
                    <div className="w-full bg-gray-700 rounded-full h-3">
                      <div
                        className="bg-blue-600 h-3 rounded-full transition-all"
                        style={{ width: `${gap.percentage}%` }}
                      ></div>
                    </div>
                  </div>
                ))}
              </div>
              ) : (
                <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
                  <div className="w-16 h-16 bg-red-900 text-red-400 rounded-full flex items-center justify-center mx-auto mb-4 text-2xl">
                    ⚠️
                  </div>
                  <h3 className="text-xl font-semibold text-white mb-2">Erro ao carregar dados</h3>
                  <p className="text-gray-400 mb-6">Formato de resposta inesperado</p>
                  <button 
                    onClick={loadData}
                    className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all"
                  >
                    Tentar Novamente
                  </button>
                </div>
              )
            )}

            {/* Career Radar */}
            {activeTab === "radar" && (
              !radarData ? (
                <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
                  <svg className="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                  <h3 className="text-xl font-semibold text-white mb-2">Sem dados do radar de carreira</h3>
                  <p className="text-gray-400">Cadastre atividades para visualizar este relatório</p>
                </div>
              ) : (
              <div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
                  <div className="bg-gray-800 border border-gray-700 rounded-lg p-6">
                    <p className="text-sm text-gray-400 mb-1">Total de Atividades</p>
                    <p className="text-4xl font-bold text-white">{radarData.total_activities}</p>
                  </div>
                  <div className="bg-gray-800 border border-gray-700 rounded-lg p-6">
                    <p className="text-sm text-gray-400 mb-1">XP Total</p>
                    <p className="text-4xl font-bold text-blue-400">{radarData.total_xp}</p>
                  </div>
                </div>

                <div className="bg-gray-800 border border-gray-700 rounded-lg p-6">
                  <h3 className="text-xl font-bold text-white mb-6">Distribuição por Nível</h3>
                  <div className="space-y-4">
                    {radarData.breakdown && Array.isArray(radarData.breakdown) ? radarData.breakdown.map((level, index) => (
                      <div key={index} className="border-b border-gray-700 pb-4 last:border-0">
                        <div className="flex items-center justify-between mb-2">
                          <span className="font-medium text-white">{level.level_name}</span>
                          <span className="text-sm text-gray-400">{level.activity_count} atividades</span>
                        </div>
                        <div className="grid grid-cols-2 gap-4 text-sm">
                          <div>
                            <span className="text-gray-400">XP: </span>
                            <span className="text-blue-400 font-medium">{level.total_xp}</span>
                          </div>
                          <div>
                            <span className="text-gray-400">% Volume: </span>
                            <span className="text-white font-medium">{level.volume_percent.toFixed(1)}%</span>
                          </div>
                        </div>
                      </div>
                    )) : (
                      <p className="text-gray-400 text-center py-4">Nenhum dado de distribuição disponível</p>
                    )}
                  </div>
                </div>
              </div>
              )
            )}

            {/* Cycle Comparison */}
            {activeTab === "comparison" && (
              !comparisonData ? (
                <div className="bg-gray-800 border border-gray-700 rounded-lg p-12 text-center">
                  <svg className="w-16 h-16 text-gray-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                  <h3 className="text-xl font-semibold text-white mb-2">Sem dados de comparação</h3>
                  <p className="text-gray-400">Cadastre atividades em diferentes ciclos para visualizar este relatório</p>
                </div>
              ) : (
              <div>
                <div className="bg-gray-800 border border-gray-700 rounded-lg p-6 mb-6">
                  <div className="grid grid-cols-3 gap-6">
                    <div>
                      <p className="text-sm text-gray-400 mb-1">Ciclo Atual</p>
                      <p className="text-xl font-bold text-white">{comparisonData.current_cycle}</p>
                    </div>
                    <div>
                      <p className="text-sm text-gray-400 mb-1">Ciclo Anterior</p>
                      <p className="text-xl font-bold text-gray-300">{comparisonData.previous_cycle}</p>
                    </div>
                    <div>
                      <p className="text-sm text-gray-400 mb-1">Crescimento</p>
                      <p className={`text-2xl font-bold ${comparisonData.growth_xp >= 0 ? "text-green-400" : "text-red-400"}`}>
                        {comparisonData.growth_xp >= 0 ? "+" : ""}{comparisonData.growth_xp} XP
                        <span className="text-sm ml-2">({comparisonData.percent_change.toFixed(1)}%)</span>
                      </p>
                    </div>
                  </div>
                </div>

                <div className="bg-gray-800 border border-gray-700 rounded-lg p-6">
                  <h3 className="text-xl font-bold text-white mb-6">Evolução por Nível</h3>
                  <div className="space-y-4">
                    {comparisonData.level_evolution && Array.isArray(comparisonData.level_evolution) ? comparisonData.level_evolution.map((level, index) => (
                      <div key={index} className="flex items-center justify-between p-4 bg-gray-700 rounded-lg">
                        <span className="font-medium text-white">{level.level_name}</span>
                        <div className="flex items-center gap-4">
                          <span className="text-sm text-gray-400">{level.prev_xp} XP</span>
                          <svg className="w-5 h-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                          </svg>
                          <span className="text-sm font-medium text-white">{level.current_xp} XP</span>
                          <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                            level.diff > 0 ? "bg-green-900 text-green-300" : 
                            level.diff < 0 ? "bg-red-900 text-red-300" : 
                            "bg-gray-600 text-gray-300"
                          }`}>
                            {level.diff > 0 ? "+" : ""}{level.diff}
                          </span>
                        </div>
                      </div>
                    )) : (
                      <p className="text-gray-400 text-center py-4">Nenhum dado de evolução disponível</p>
                    )}
                  </div>
                </div>
              </div>
              )
            )}
          </>
        )}
      </main>
    </div>
  );
}
