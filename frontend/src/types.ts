export type Role = "user" | "organizer" | "admin";
export type EventStatus =
  | "EVENT_STATUS_DRAFT"
  | "EVENT_STATUS_PENDING_APPROVAL"
  | "EVENT_STATUS_PUBLISHED"
  | "EVENT_STATUS_CANCELLED"
  | "EVENT_STATUS_CLOSED";

export interface User {
  id: string;
  email: string;
  full_name: string;
  phone?: string;
  role: Role;
  is_active: boolean;
}
export interface Organizer {
  id: string;
  name: string;
  email: string;
  phone?: string;
  slug: string;
  is_active: boolean;
}
export interface Event {
  id: string;
  organizer_id: string;
  title: string;
  description?: string;
  venue: string;
  start_at: string;
  end_at: string;
  capacity: number;
  status: EventStatus;
}
