import { FormEvent, useEffect, useState } from "react";
import { Button } from "../../../components/ui/Button";
import { Card } from "../../../components/ui/Card";
import { Modal } from "../../../components/ui/Modal";
import { useAppDispatch, useAppSelector } from "../../../hooks/redux";
import { approveEvent, createEvent, fetchEvents } from "../../../modules/events/eventSlice";
import { getEventStatusText, getTranslations } from "../../../config/language";

export function AdminEventsPage() {
    const dispatch = useAppDispatch();
    const items = useAppSelector((s) => (Array.isArray(s.events.items) ? s.events.items : []));
    const loading = useAppSelector((s) => s.events.loading);
    const loadError = useAppSelector((s) => s.events.error);
    const language = useAppSelector((s) => s.language.language);
    const t = getTranslations(language);
    const [open, setOpen] = useState(false);
    const [message, setMessage] = useState("");
    const [form, setForm] = useState({
        organizer_id: "",
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
            setForm({
                organizer_id: "",
                title: "",
                description: "",
                venue: "",
                start_at: "",
                end_at: "",
                capacity: 100,
            });
            setOpen(false);
            setMessage(t.adminEvents.createSuccess);
        } catch (err) {
            setMessage(err instanceof Error ? err.message : t.adminEvents.createError);
        }
    };

    return (
        <div className="stack">
            <Card title={t.adminEvents.title}>
                <div className="stack">
                    <Button type="button" onClick={() => setOpen(true)}>{t.adminEvents.add}</Button>
                    {message && <p>{message}</p>}
                </div>
            </Card>

            <Card title={t.adminEvents.list}>
                {loading && <p>{t.adminEvents.loading}</p>}
                {loadError && <p className="error">{loadError}</p>}
                {!loading && !loadError && items.length === 0 && <p>{t.adminEvents.empty}</p>}
                {!loading && !loadError && items.length > 0 && (
                    <div className="table-wrapper">
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>{t.adminEvents.event}</th>
                                    <th>{t.adminEvents.venue}</th>
                                    <th>{t.adminEvents.start}</th>
                                    <th>{t.adminEvents.capacity}</th>
                                    <th>{t.adminEvents.status}</th>
                                    <th>{t.adminEvents.actions}</th>
                                </tr>
                            </thead>
                            <tbody>
                                {items.map((event) => {
                                    const status = event.status ?? "EVENT_STATUS_DRAFT";
                                    return (
                                        <tr key={event.id}>
                                            <td className="table-primary">
                                                <strong>{event.title}</strong>
                                                <small className="table-secondary">{event.description ?? t.common.noDescription}</small>
                                            </td>
                                            <td>{event.venue}</td>
                                            <td>{new Date(event.start_at).toLocaleString(language === "vi" ? "vi-VN" : "en-US")}</td>
                                            <td>{event.capacity}</td>
                                            <td>
                                                <span className={`status ${status.toLowerCase()}`}>
                                                    {getEventStatusText(language, status)}
                                                </span>
                                            </td>
                                            <td>
                                                {status === "EVENT_STATUS_PENDING_APPROVAL" ? (
                                                    <Button onClick={() => void dispatch(approveEvent(event.id))}>{t.adminEvents.approve}</Button>
                                                ) : (
                                                    <span className="table-muted">{t.adminEvents.none}</span>
                                                )}
                                            </td>
                                        </tr>
                                    );
                                })}
                            </tbody>
                        </table>
                    </div>
                )}
            </Card>

            <Modal open={open} title={t.adminEvents.addNew} onClose={() => setOpen(false)}>
                <form className="form" onSubmit={submit}>
                    <input placeholder={t.adminEvents.organizerId} value={form.organizer_id} onChange={(e) => setForm({ ...form, organizer_id: e.target.value })} />
                    <input placeholder={t.adminEvents.eventTitle} value={form.title} onChange={(e) => setForm({ ...form, title: e.target.value })} />
                    <input placeholder={t.adminEvents.description} value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
                    <input placeholder={t.adminEvents.venue} value={form.venue} onChange={(e) => setForm({ ...form, venue: e.target.value })} />
                    <input type="datetime-local" value={form.start_at} onChange={(e) => setForm({ ...form, start_at: e.target.value })} />
                    <input type="datetime-local" value={form.end_at} onChange={(e) => setForm({ ...form, end_at: e.target.value })} />
                    <input type="number" min={1} value={form.capacity} onChange={(e) => setForm({ ...form, capacity: Number(e.target.value) })} />
                    <div className="modal-actions">
                        <Button type="button" onClick={() => setOpen(false)}>{t.common.cancel}</Button>
                        <Button type="submit">{t.common.save}</Button>
                    </div>
                </form>
            </Modal>
        </div>
    );
}
