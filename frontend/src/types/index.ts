// Types baseados nos DTOs do backend

export interface User {
  id: string;
  username: string;
  email: string;
  active: boolean;
  role: string;
  current_level?: string;
  manager_name?: string;
  manager_email?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  active?: boolean;
}

export enum Pillar {
  TECHNICAL = "TECHNICAL",
  RESULTS = "RESULTS",
  INFLUENCE = "INFLUENCE",
}

export interface CreateInitiativeRequest {
  user_id: string;
  ladder_id: string;
  pillars: Pillar[];
  title: string;
  description?: string;
  progress_percentage: number;
  impact_summary?: string;
  is_pdi_target: boolean;
}

export interface Initiative {
  id: string;
  user_id: string;
  title: string;
  description?: string;
  progress_percentage: number;
  is_pdi_target: boolean;
  completed_at?: string;
  created_at: string;
  task_count: number;
  has_extra: boolean;
}

export interface Task {
  id: string;
  initiative_id: string;
  ladder_id: string;
  title: string;
  execution?: string;
  impact_summary?: string;
  progress_percentage: number;
  is_extra: boolean;
  completed_at?: string;
  created_at: string;
  updated_at: string;
  evidence_count: number;
}

export interface TaskEvidence {
  id: string;
  task_id: string;
  evidence_url: string;
  description?: string;
  created_at: string;
}

export interface PillarStats {
  achieved: number;
  planned: number;
  percentage: number;
}

export interface DashboardResponse {
  official_level: string;
  target_level: string;
  current_level?: string; // Mantido para retrocompatibilidade
  pdi_progress: Record<string, PillarStats>;
  max_pdi_xp: number;
  total_achieved: number;
  overdelivery: Record<string, number>;
}

export interface GapAnalysisResponse {
  level: string;
  pillar: string;
  target: number;
  achieved: number;
  gap: number;
  status: string;
  percentage: number;
}

export interface ReadinessCheck {
  is_consistent: boolean;
  message: string;
  target_level: string;
  target_count: number;
  others_count: number;
}

export interface LevelComposition {
  level_name: string;
  activity_count: number;
  total_xp: number;
  volume_percent: number;
  xp_percent: number;
}

export interface CareerRadar {
  total_activities: number;
  total_xp: number;
  breakdown: LevelComposition[];
}

export interface LevelComparison {
  level_name: string;
  current_xp: number;
  prev_xp: number;
  diff: number;
}

export interface ComparisonReport {
  current_cycle: string;
  previous_cycle: string;
  growth_xp: number;
  percent_change: number;
  level_evolution: LevelComparison[];
}

export interface UpdateProgressRequest {
  progress_percentage: number;
}

export type LadderLevel = "P1" | "P2" | "P3" | "LT1" | "LT2" | "LT3" | "LT4";

export interface CareerLadder {
  id: string;
  level: LadderLevel;
  xp_reward: number;
  technical: string;
  expected_results: string;
  leadership_scope: string;
}

export interface CreateCareerLadderRequest {
  level: LadderLevel;
  xp_reward: number;
  technical: string;
  expected_results: string;
  leadership_scope: string;
}
