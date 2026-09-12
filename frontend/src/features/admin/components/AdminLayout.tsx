import { NavLink, Outlet } from "react-router-dom";
import { useAppSelector } from "../../../hooks/redux";
import { getTranslations } from "../../../config/language";

export function AdminLayout() {
  const language = useAppSelector((s) => s.language.language);
  const t = getTranslations(language).adminLayout;

  const adminLinks = [
    { to: "/admin", label: t.overview, end: true },
    { to: "/admin/users", label: t.users },
    { to: "/admin/organizers", label: t.organizers },
    { to: "/admin/events", label: t.events },
  ];

  return (
    <div className="admin-layout">
      <aside className="admin-sidebar">
        <div className="admin-sidebar-heading">
          <span className="admin-sidebar-kicker">ADMIN</span>
          <h2>{t.area}</h2>
          <p>{t.description}</p>
        </div>

        <nav className="admin-nav" aria-label={t.navigation}>
          {adminLinks.map((link) => (
            <NavLink
              key={link.to}
              to={link.to}
              end={link.end}
              className={({ isActive }) =>
                `admin-nav-link${isActive ? " active" : ""}`
              }
            >
              {link.label}
            </NavLink>
          ))}
        </nav>
      </aside>

      <section className="admin-content">
        <Outlet />
      </section>
    </div>
  );
}
