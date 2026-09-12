import { Link } from "react-router-dom";
import { Card } from "../components/ui/Card";
import { useAppSelector } from "../hooks/redux";
import { getTranslations } from "../config/language";

export function HomePage() {
  const language = useAppSelector((s) => s.language.language);
  const t = getTranslations(language).home;

  return (
    <div className="hero">
      <p className="eyebrow">{t.eyebrow}</p>
      <h1>{t.title}</h1>
      <p>{t.description}</p>
      <Link className="button" to="/events">{t.explore}</Link>
      <Card title={t.workflow}>
        <p>{t.workflowDescription}</p>
      </Card>
    </div>
  );
}
