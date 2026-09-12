import { Card } from "../../../components/ui/Card";
import { useAppSelector } from "../../../hooks/redux";
import { getTranslations } from "../../../config/language";

export function AdminDashboardPage() {
    const language = useAppSelector((s) => s.language.language);
    const t = getTranslations(language).dashboard;

    return (
        <div className="admin-dashboard">
            <div className="admin-page-heading">
                <span className="eyebrow">{t.eyebrow}</span>
                <h1>{t.title}</h1>
                <p>{t.welcome}</p>
            </div>
            <div className="admin-stat-grid">
                <Card title={t.users}>
                    <p className="admin-stat-description">{t.usersDescription}</p>
                </Card>
                <Card title={t.organizers}>
                    <p className="admin-stat-description">{t.organizersDescription}</p>
                </Card>
                <Card title={t.events}>
                    <p className="admin-stat-description">{t.eventsDescription}</p>
                </Card>
            </div>
        </div>
    );
}
