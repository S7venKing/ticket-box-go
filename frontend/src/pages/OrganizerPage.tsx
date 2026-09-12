import { FormEvent, useEffect, useState } from "react";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { createEvent, fetchEvents } from "../modules/events/eventSlice";
import { Card } from "../components/ui/Card";
import { Button } from "../components/ui/Button";
import { getEventStatusText, getTranslations } from "../config/language";

export function OrganizerPage() {
  const dispatch = useAppDispatch();
  const user = useAppSelector((s) => s.auth.user);
  const events = useAppSelector((s) => (Array.isArray(s.events.items) ? s.events.items : []));
  const loading = useAppSelector((s) => s.events.loading);
  const language = useAppSelector((s) => s.language.language);
  const t = getTranslations(language);
  const [message, setMessage] = useState("");
  const [form, setForm] = useState({
    organizer_id: user?.organizer_id ?? "",
    title: "",
    description: "",
    venue: "",
    start_at: "",
    end_at: "",
    capacity: 100,
  });
  useEffect(() => {
    void dispatch(fetchEvents());
  }, [dispatch]);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await dispatch(createEvent(form)).unwrap();
      await dispatch(fetchEvents()).unwrap();
      setMessage(t.organizerPage.createSuccess);
      setForm({ ...form, title: "", description: "", venue: "", start_at: "", end_at: "" });
    } catch (error) {
      setMessage(error instanceof Error ? error.message : t.organizerPage.createError);
    }
  };
  if (!user?.organizer_id) {
    return <Card title={t.organizerPage.title}><p>{t.organizerPage.notAssigned}</p></Card>;
  }
  return (
    <div className="stack">
      <Card title={`${t.organizerPage.createFor}${user ? " " + user.email : ""}`}>
        <form className="form" onSubmit={submit}>
          {([
            ["title", t.organizerPage.titleField],
            ["description", t.organizerPage.descriptionField],
            ["venue", t.organizerPage.venueField],
            ["start_at", t.organizerPage.startField],
            ["end_at", t.organizerPage.endField],
          ] as const).map(([key, placeholder]) => (
            <input
              key={key}
              placeholder={placeholder}
              value={form[key]}
              onChange={(e) => setForm({ ...form, [key]: e.target.value })}
            />
          ))}
        <input
          type="number"
          value={form.capacity}
          onChange={(e) => setForm({ ...form, capacity: Number(e.target.value) })}
        />
          <Button>{t.organizerPage.createDraft}</Button>
          {message && <p>{message}</p>}
        </form>
      </Card>
      <Card title={t.organizerPage.eventList}>
        {loading && <p>{t.organizerPage.loading}</p>}
        {!loading && events.length === 0 && <p>{t.organizerPage.empty}</p>}
        {!loading && events.length > 0 && (
          <div className="table-wrapper">
            <table className="data-table">
              <thead><tr><th>{t.organizerPage.event}</th><th>{t.organizerPage.venue}</th><th>{t.organizerPage.start}</th><th>{t.organizerPage.capacity}</th><th>{t.organizerPage.status}</th></tr></thead>
              <tbody>{events.map((event) => (
                <tr key={event.id}>
                  <td className="table-primary">{event.title}</td>
                  <td>{event.venue}</td>
                  <td>{new Date(event.start_at).toLocaleString(language === "vi" ? "vi-VN" : "en-US")}</td>
                  <td>{event.capacity}</td>
                  <td><span className="status">{getEventStatusText(language, event.status)}</span></td>
                </tr>
              ))}</tbody>
            </table>
          </div>
        )}
      </Card>
    </div>
  );
}
