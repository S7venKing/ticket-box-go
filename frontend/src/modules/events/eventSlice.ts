import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { get, post } from "../../lib/api";
import type { Event } from "../../types";
interface State { items: Event[]; loading: boolean; error: string | null; }
const initialState: State = { items: [], loading: false, error: null };
export const fetchEvents = createAsyncThunk("events/list", async () => (await get<{ events: Event[] }>("/events")).events);
export const createEvent = createAsyncThunk("events/create", async (input: Omit<Event, "id" | "status">) => (await post<{ event: Event }>("/events", input)).event);
export const submitEvent = createAsyncThunk("events/submit", async (id: string) => (await post<{ event: Event }>(`/events/${id}/submit`)).event);
export const approveEvent = createAsyncThunk("events/approve", async (id: string) => (await post<{ event: Event }>(`/events/${id}/approve`)).event);
const slice = createSlice({ name: "events", initialState, reducers: {}, extraReducers: b => {
  b.addCase(fetchEvents.pending, s => { s.loading = true; }).addCase(fetchEvents.fulfilled, (s, a) => { s.loading = false; s.items = a.payload; }).addCase(fetchEvents.rejected, (s, a) => { s.loading = false; s.error = a.error.message ?? "Failed to load events"; })
    .addCase(createEvent.fulfilled, (s, a) => { s.items.unshift(a.payload); }).addCase(submitEvent.fulfilled, (s, a) => { const i = s.items.findIndex(e => e.id === a.payload.id); if (i >= 0) s.items[i] = a.payload; }).addCase(approveEvent.fulfilled, (s, a) => { const i = s.items.findIndex(e => e.id === a.payload.id); if (i >= 0) s.items[i] = a.payload; });
} });
export default slice.reducer;
