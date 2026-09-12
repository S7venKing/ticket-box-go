import { Link, Outlet } from "react-router-dom";
import { useEffect } from "react";
import { useAppDispatch, useAppSelector } from "../../hooks/redux";
import { loadMe, logout } from "../../modules/auth/authSlice";
import { LanguageSwitcher } from "../ui/LanguageSwitcher";
import { getText } from "../../config/language";

export function AppLayout() {
  const user = useAppSelector((s) => s.auth.user);
  const token = useAppSelector((s) => s.auth.token);
  const language = useAppSelector((s) => s.language.language);
  const dispatch = useAppDispatch();

  useEffect(() => {
    document.documentElement.lang = language === "en" ? "en" : "vi";
  }, [language]);

  useEffect(() => {
    if (token && !user) {
      void dispatch(loadMe());
    }
  }, [dispatch, token, user]);

  return (
    <>
      <header className="topbar">
        <div className="topbar-inner">
          <Link className="brand" to="/">
            <span className="brand-mark">T</span>
            Ticket Box
          </Link>

          <nav className="topbar-nav">
            <Link to="/events">{getText(language, "events")}</Link>
            {user?.role === "admin" && <Link to="/admin">{getText(language, "admin")}</Link>}
            {(user?.role === "organizer" || user?.organizer_id) && <Link to="/organizer">{getText(language, "organizer")}</Link>}
          </nav>

          <div className="topbar-actions">
            <LanguageSwitcher label={getText(language, "language")} />
            {user ? (
              <>
                <span className="user-pill">{user.email}</span>
                <button className="link-button" onClick={() => dispatch(logout())}>
                  {getText(language, "logout")}
                </button>
              </>
            ) : (
              <Link className="login-link" to="/login">
                {getText(language, "login")}
              </Link>
            )}
          </div>
        </div>
      </header>

      <main className="page-shell">
        <div className="container">
          <Outlet />
        </div>
      </main>
    </>
  );
}
