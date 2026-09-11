import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { get, post } from "../../lib/api";
import type { Organizer } from "../../types";
interface State {
  items: Organizer[];
  loading: boolean;
  error: string | null;
}
const initialState: State = { items: [], loading: false, error: null };
export const fetchOrganizers = createAsyncThunk(
  "organizers/list",
  async () => (await get<{ organizers: Organizer[] }>("/organizers")).organizers,
);
export const createOrganizer = createAsyncThunk(
  "organizers/create",
  async (input: { name: string; email: string; phone?: string; slug: string }) =>
    (await post<{ organizer: Organizer }>("/organizers", input)).organizer,
);
const slice = createSlice({
  name: "organizers",
  initialState,
  reducers: {},
  extraReducers: (b) => {
    b.addCase(fetchOrganizers.pending, (s) => {
      s.loading = true;
    })
      .addCase(fetchOrganizers.fulfilled, (s, a) => {
        s.loading = false;
        s.items = a.payload;
      })
      .addCase(fetchOrganizers.rejected, (s, a) => {
        s.loading = false;
        s.error = a.error.message ?? "Failed to load organizers";
      })
      .addCase(createOrganizer.fulfilled, (s, a) => {
        s.items.unshift(a.payload);
      });
  },
});
export default slice.reducer;
