import { Link } from "react-router-dom";
import { Card } from "../components/ui/Card";
export function HomePage() {
  return (
    <div className="hero">
      <p className="eyebrow">EVENT PLATFORM</p>
      <h1>Create. Review. Publish.</h1>
      <p>Ticket Box centralizes organizer management and admin-approved events.</p>
      <Link className="button" to="/events">
        Explore events
      </Link>
      <Card title="Workflow">
        <p>Organizer creates a draft, submits it for review, and an admin publishes it.</p>
      </Card>
    </div>
  );
}
