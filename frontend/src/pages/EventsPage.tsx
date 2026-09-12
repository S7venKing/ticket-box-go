import { useEffect } from "react";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { approveEvent, fetchEvents, submitEvent } from "../modules/events/eventSlice";
import { Card } from "../components/ui/Card";
import { Button } from "../components/ui/Button";
import { getEventStatusText, getTranslations } from "../config/language";

export function EventsPage() {
  const dispatch = useAppDispatch();
  const items = useAppSelector((s) => (Array.isArray(s.events.items) ? s.events.items : []));
  const loading = useAppSelector((s) => s.events.loading);
  const role = useAppSelector((s) => s.auth.user?.role);
  const language = useAppSelector((s) => s.language.language);
  const t = getTranslations(language);

  useEffect(() => {
    void dispatch(fetchEvents());
  }, [dispatch]);

  return (
    <>
      <div className="page-heading">
        <div>
          <p className="eyebrow">{t.events.eyebrow}</p>
          <h1>{t.events.title}</h1>
        </div>
      </div>

      {loading && <p>{t.events.loading}</p>}
      {!loading && items.length === 0 && <p>{t.events.empty}</p>}

      <div className="grid">
        {items.map((event) => (
          <Card key={event.id}>
            <span className={`status ${(event.status ?? "EVENT_STATUS_DRAFT").toLowerCase()}`}>
              {getEventStatusText(language, event.status ?? "EVENT_STATUS_DRAFT")}
            </span>
            <h2>{event.title}</h2>
            <p>{event.description ?? t.common.noDescription}</p>
            <p>
              {event.venue} · {new Date(event.start_at).toLocaleString(language === "vi" ? "vi-VN" : "en-US")}
            </p>
            <div className="actions">
              {role === "organizer" && event.status === "EVENT_STATUS_DRAFT" && (
                <Button onClick={() => void dispatch(submitEvent(event.id))}>{t.events.submitReview}</Button>
              )}
              {role === "admin" && event.status === "EVENT_STATUS_PENDING_APPROVAL" && (
                <Button onClick={() => void dispatch(approveEvent(event.id))}>{t.events.approve}</Button>
              )}
            </div>
          </Card>
        ))}
      </div>
    </>
  );
}
