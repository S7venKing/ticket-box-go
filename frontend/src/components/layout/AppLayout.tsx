import { Link, Outlet } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../../hooks/redux";
import { logout } from "../../modules/auth/authSlice";
export function AppLayout() {
  const user = useAppSelector(s => s.auth.user); const dispatch = useAppDispatch();
  return <><header className="topbar"><Link className="brand" to="/">Ticket Box</Link><nav><Link to="/events">Events</Link>{user?.role === "admin" && <Link to="/admin">Admin</Link>}{user?.role === "organizer" && <Link to="/organizer">Organizer</Link>}</nav>{user ? <button className="link-button" onClick={() => dispatch(logout())}>Logout</button> : <Link to="/login">Login</Link>}</header><main className="container"><Outlet /></main></>;
}
