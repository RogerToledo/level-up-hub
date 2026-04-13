// Types baseados nos DTOs do backend

export interface User {
  id: string;
  username: string;
  email: string;
  active: boolean;
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
  COMMUNICATION = "COMMUNICATION",
  TEAMWORK = "TEAMWORK",
  TECHNICAL = "TECHNICAL",
  LEADERSHIP = "LEADERSHIP",
}

export interface CreateActivityRequest {
  user_id: string;
  ladder_id: string;
  pillars: Pillar[];
  title: string;
  description?: string;
  progress_percentage: number;
  impact_summary?: string;
  is_pdi_target: boolean;
}

export interface Activity {
  id: string;
  user_id: string;
  ladder_id: string;
  title: string;
  description?: string;
  progress_percentage: number;
  impact_summary?: string;
  is_pdi_target: boolean;
  created_at: string;
  updated_at: string;
}

export interface PillarStats {
  achieved: number;
  planned: number;
  percentage: number;
}

export interface DashboardResponse {
  current_level: string;
  pdi_progress: Record<string, PillarStats>;
  max_pdi_xp: number;
  total_achieved: number;
  overdelivery: Record<string, number>;
}

export interface GapAnalysisResponse {
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
