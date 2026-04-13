"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/services/api";
import Sidebar from "@/components/Sidebar";
import StatCard from "@/components/StatCard";
import PageHeader from "@/components/PageHeader";
import type { DashboardResponse } from "@/types";

export default function DashboardPage() {
  const router = useRouter();
  const [data, setData] = useState<DashboardResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchDashboard = async () => {
      try {
        const response = await api.get("/dashboard");
        setData(response);
      } catch (err: unknown) {
        if (err instanceof Error) setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchDashboard();
  }, []);

  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-950">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Carregando...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-950 p-4">
        <div className="bg-gray-800 p-8 rounded-2xl shadow-lg border border-red-900 text-center max-w-md">
          <div className="w-16 h-16 bg-red-900 text-red-400 rounded-full flex items-center justify-center mx-auto mb-4 text-2xl">
            ⚠️
          </div>
          <h2 className="text-xl font-bold text-white mb-2">Erro ao Carregar</h2>
          <p className="text-gray-400 mb-6">{error}</p>
          <button
            onClick={() => router.push("/login")}
            className="w-full py-3 bg-blue-600 text-white rounded-xl font-semibold hover:bg-blue-700 transition-colors"
          >
            Voltar para o Login
          </button>
        </div>
      </div>
    );
  }

  if (!data) {
    return null;
  }

  // Calcula total de XP dos pilares PDI com verificação de segurança
  const pdiProgress = data.pdi_progress || {};
  const overdelivery = data.overdelivery || {};
  
  const totalPdiXp = Object.values(pdiProgress).reduce((sum, pillar) => sum + pillar.achieved, 0);
  const avgPdiPercentage = Object.keys(pdiProgress).length > 0
    ? Object.values(pdiProgress).reduce((sum, pillar) => sum + pillar.percentage, 0) / Object.keys(pdiProgress).length
    : 0;

  return (
    <div className="flex min-h-screen bg-gray-950">
      <Sidebar />
      
      <main className="flex-1 ml-64 p-8">
        <PageHeader 
          title="Dashboard" 
          subtitle={`Nível Atual: ${data.current_level}`}
        />

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <StatCard
            title="Nível Atual"
            value={data.current_level}
            icon={
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            }
          />
          
          <StatCard
            title="XP Total Conquistado"
            value={data.total_achieved || 0}
            subtitle={`de ${data.max_pdi_xp || 0} max`}
            icon={
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
              </svg>
            }
          />

          <StatCard
            title="Progresso PDI"
            value={`${avgPdiPercentage.toFixed(0)}%`}
            subtitle={`${totalPdiXp} XP nos pilares`}
            icon={
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            }
          />

          <StatCard
            title="XP Máximo PDI"
            value={data.max_pdi_xp || 0}
            subtitle="Meta do ciclo"
            icon={
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            }
          />
        </div>

        {/* Progresso por Pilar */}
        <div className="bg-gray-800 border border-gray-700 rounded-lg p-6 mb-8">
          <h2 className="text-xl font-bold text-white mb-6">Progresso por Pilar PDI</h2>
          {Object.keys(pdiProgress).length === 0 ? (
            <div className="text-center py-8">
              <p className="text-gray-400">Nenhum dado de progresso disponível</p>
            </div>
          ) : (
            <div className="space-y-4">
              {Object.entries(pdiProgress).map(([pillar, stats]) => (
                <div key={pillar}>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-gray-300 capitalize">
                      {pillar.toLowerCase()}
                    </span>
                    <span className="text-sm text-gray-400">
                      {stats.achieved} / {stats.planned} XP ({stats.percentage.toFixed(0)}%)
                    </span>
                  </div>
                  <div className="w-full bg-gray-700 rounded-full h-2.5">
                    <div
                      className="bg-blue-600 h-2.5 rounded-full transition-all"
                      style={{ width: `${Math.min(stats.percentage, 100)}%` }}
                    ></div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Overdelivery */}
        {Object.keys(overdelivery).length > 0 && (
          <div className="bg-gray-800 border border-gray-700 rounded-lg p-6">
            <h2 className="text-xl font-bold text-white mb-6">Overdelivery 🚀</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              {Object.entries(overdelivery).map(([level, xp]) => (
                <div key={level} className="bg-gray-700 rounded-lg p-4 border border-gray-600">
                  <p className="text-sm text-gray-400">Nível {level}</p>
                  <p className="text-2xl font-bold text-green-400">+{xp} XP</p>
                </div>
              ))}
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
