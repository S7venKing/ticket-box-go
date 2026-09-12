import { Navigate, Outlet } from "react-router-dom";
import { useAppSelector } from "../../hooks/redux";
import type { Role } from "../../types";
import { getTranslations } from "../../config/language";

export function ProtectedRoute({ roles }: { roles?: Role[] }) {
  const { user, token, hydrating } = useAppSelector((s) => s.auth);
  const language = useAppSelector((s) => s.language.language);
  if (!token) return <Navigate to="/login" replace />;
  if (hydrating || !user) return <div className="route-loading">{getTranslations(language).protectedRoute.authenticating}</div>;
  if (roles && user && !roles.includes(user.role)) return <Navigate to="/" replace />;
  return <Outlet />;
}
