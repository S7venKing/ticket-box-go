import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { get } from "../../lib/api";
import type { User } from "../../types";

interface State {
  items: User[];
  loading: boolean;
  error: string | null;
}

const initialState: State = { items: [], loading: false, error: null };

export const fetchUsers = createAsyncThunk(
  "users/list",
  async () => (await get<{ users: User[] }>("/admin/users")).users,
);

const slice = createSlice({
  name: "users",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchUsers.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchUsers.fulfilled, (state, action) => {
        state.loading = false;
        state.items = Array.isArray(action.payload) ? action.payload : [];
      })
      .addCase(fetchUsers.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message ?? "Không thể tải danh sách người dùng.";
      });
  },
});

export default slice.reducer;
