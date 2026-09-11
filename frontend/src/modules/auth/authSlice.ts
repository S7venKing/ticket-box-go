import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { api, post } from "../../lib/api";
import type { Role, User } from "../../types";

interface AuthState {
  user: User | null;
  token: string | null;
  loading: boolean;
  error: string | null;
}
const initialState: AuthState = {
  user: null,
  token: localStorage.getItem("access_token"),
  loading: false,
  error: null,
};

export const login = createAsyncThunk(
  "auth/login",
  async (input: { email: string; password: string }) => {
    const result = await post<{ user: User; access_token: string }>("/login", input);
    localStorage.setItem("access_token", result.access_token);
    return result;
  },
);
export const register = createAsyncThunk(
  "auth/register",
  async (input: { email: string; password: string; full_name: string; phone?: string }) =>
    post<{ user: User }>("/register", input),
);
export const loadMe = createAsyncThunk("auth/me", async () => api<{ user: User }>("/me"));

const slice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    logout(state) {
      state.user = null;
      state.token = null;
      localStorage.removeItem("access_token");
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(login.pending, (s) => {
        s.loading = true;
        s.error = null;
      })
      .addCase(login.fulfilled, (s, a) => {
        s.loading = false;
        s.user = a.payload.user;
        s.token = a.payload.access_token;
      })
      .addCase(login.rejected, (s, a) => {
        s.loading = false;
        s.error = a.error.message ?? "Login failed";
      })
      .addCase(loadMe.fulfilled, (s, a) => {
        s.user = a.payload.user;
      })
      .addCase(loadMe.rejected, (s) => {
        s.user = null;
        s.token = null;
        localStorage.removeItem("access_token");
      });
  },
});
export const { logout } = slice.actions;
export const selectRole = (state: { auth: AuthState }): Role | null =>
  state.auth.user?.role ?? null;
export default slice.reducer;
