import { useEffect } from "react";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { approveEvent, fetchEvents, submitEvent } from "../modules/events/eventSlice";
import { Card } from "../components/ui/Card";
import { Button } from "../components/ui/Button";
export function EventsPage() {
  const dispatch = useAppDispatch();
  const { items, loading } = useAppSelector((s) => s.events);
  const role = useAppSelector((s) => s.auth.user?.role);
  useEffect(() => {
    dispatch(fetchEvents());
  }, [dispatch]);
  return (
    <>
      <div className="page-heading">
        <div>
          <p className="eyebrow">DISCOVER</p>
          <h1>Events</h1>
        </div>
      </div>
      <div className="grid">
        {items.map((event) => (
          <Card key={event.id}>
            <span className={`status ${event.status.toLowerCase()}`}>
              {event.status.replace("EVENT_STATUS_", "")}
            </span>
            <h2>{event.title}</h2>
            <p>{event.description}</p>
            <p>
              {event.venue} · {new Date(event.start_at).toLocaleString()}
            </p>
            <div className="actions">
              {role === "organizer" && event.status === "EVENT_STATUS_DRAFT" && (
                <Button onClick={() => dispatch(submitEvent(event.id))}>Submit for review</Button>
              )}
              {role === "admin" && event.status === "EVENT_STATUS_PENDING_APPROVAL" && (
                <Button onClick={() => dispatch(approveEvent(event.id))}>Approve</Button>
              )}
            </div>
          </Card>
        ))}
      </div>
      {loading && <p>Loading events...</p>}
    </>
  );
}
