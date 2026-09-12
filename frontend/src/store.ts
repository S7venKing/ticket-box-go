import { configureStore } from "@reduxjs/toolkit";
import auth from "./modules/auth/authSlice";
import organizers from "./modules/organizers/organizerSlice";
import events from "./modules/events/eventSlice";
import users from "./modules/users/userSlice";
import language from "./modules/language/languageSlice";

export const store = configureStore({ reducer: { auth, organizers, events, users, language } });
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
