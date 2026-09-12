import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import { defaultLanguage, supportedLanguages, type Language } from "../../config/language";

export const SUPPORTED_LANGUAGES = supportedLanguages;
export const DEFAULT_LANGUAGE = defaultLanguage;
export const LANGUAGE_STORAGE_KEY = "app-language";
export type { Language };

interface LanguageState {
  language: Language;
}

const getInitialLanguage = (): Language => {
  const saved = localStorage.getItem(LANGUAGE_STORAGE_KEY);
  return saved === "en" || saved === "vi" ? saved : DEFAULT_LANGUAGE;
};

const initialState: LanguageState = {
  language: getInitialLanguage(),
};

const slice = createSlice({
  name: "language",
  initialState,
  reducers: {
    setLanguage(state, action: PayloadAction<Language>) {
      state.language = action.payload;
      localStorage.setItem(LANGUAGE_STORAGE_KEY, action.payload);
      document.documentElement.lang = action.payload === "en" ? "en" : "vi";
    },
  },
});

export const { setLanguage } = slice.actions;
export default slice.reducer;
