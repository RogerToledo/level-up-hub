"use client";

import { useEffect, useState } from "react";
import Sidebar from "@/components/Sidebar";
import PageHeader from "@/components/PageHeader";
import { api } from "@/services/api";
import type { GapAnalysisResponse, CareerRadar, ComparisonReport } from "@/types";

export default function ReportsPage() {
  const [activeTab, setActiveTab] = useState<"gap" | "radar" | "comparison">("gap");
  const [gapData, setGapData] = useState<GapAnalysisResponse[]>([]);
  const [radarData, setRadarData] = useState<CareerRadar | null>(null);
  const [comparisonData, setComparisonData] = useState<ComparisonReport | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, [activeTab]);

  const loadData = async () => {
    setLoading(true);
    try {
      if (activeTab === "gap") {
        const data = await api.get("/gap-analysis");
        setGapData(data);
      } else if (activeTab === "radar") {
        const data = await api.get("/career-radar");
        setRadarData(data);
      } else if (activeTab === "comparison") {
        const data = await api.get("/cycle-comparison");
        setComparisonData(data);
      }
    } catch (error) {
      console.error("Erro ao carregar relatório:", error);
    } finally {
      setLoading(false);
    }
  };

  const tabs = [
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
        ) : (
          <>
            {/* Gap Analysis */}
            {activeTab === "gap" && (
              <div className="space-y-4">
                {gapData.map((gap, index) => (
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
            )}

            {/* Career Radar */}
            {activeTab === "radar" && radarData && (
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
                    {radarData.breakdown.map((level, index) => (
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
                    ))}
                  </div>
                </div>
              </div>
            )}

            {/* Cycle Comparison */}
            {activeTab === "comparison" && comparisonData && (
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
                    {comparisonData.level_evolution.map((level, index) => (
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
                    ))}
                  </div>
                </div>
              </div>
            )}
          </>
        )}
      </main>
    </div>
  );
}
