import React from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { store } from "./store";
import { AppLayout } from "./components/layout/AppLayout";
import { ProtectedRoute } from "./components/auth/ProtectedRoute";
import { HomePage } from "./pages/HomePage";
import { LoginPage } from "./pages/LoginPage";
import { EventsPage } from "./pages/EventsPage";
import { OrganizerPage } from "./pages/OrganizerPage";
import { AdminDashboardPage } from "./features/admin/pages/AdminDashboardPage";
import { AdminUsersPage } from "./features/admin/pages/AdminUsersPage";
import { AdminOrganizersPage } from "./features/admin/pages/AdminOrganizersPage";
import { AdminEventsPage } from "./features/admin/pages/AdminEventsPage";
import { AdminLayout } from "./features/admin/components/AdminLayout";
import "./styles.css";

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Provider store={store}>
      <BrowserRouter>
        <Routes>
          <Route element={<AppLayout />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/events" element={<EventsPage />} />

            <Route element={<ProtectedRoute roles={["admin"]} />}>
              <Route path="/admin" element={<AdminLayout />}>
                <Route index element={<AdminDashboardPage />} />
                <Route path="users" element={<AdminUsersPage />} />
                <Route path="organizers" element={<AdminOrganizersPage />} />
                <Route path="events" element={<AdminEventsPage />} />
              </Route>
            </Route>

            <Route element={<ProtectedRoute />}>
              <Route path="/organizer" element={<OrganizerPage />} />
            </Route>
          </Route>
        </Routes>
      </BrowserRouter>
    </Provider>
  </React.StrictMode>,
);
