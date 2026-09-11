import { Navigate, Outlet } from "react-router-dom";
import { useAppSelector } from "../../hooks/redux";
import type { Role } from "../../types";
export function ProtectedRoute({ roles }: { roles?: Role[] }) {
  const { user, token } = useAppSelector((s) => s.auth);
  if (!token) return <Navigate to="/login" replace />;
  if (roles && user && !roles.includes(user.role)) return <Navigate to="/" replace />;
  return <Outlet />;
}
