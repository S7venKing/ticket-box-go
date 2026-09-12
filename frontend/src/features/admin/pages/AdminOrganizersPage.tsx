import { FormEvent, useEffect, useState } from "react";
import { Button } from "../../../components/ui/Button";
import { Card } from "../../../components/ui/Card";
import { Modal } from "../../../components/ui/Modal";
import { useAppDispatch, useAppSelector } from "../../../hooks/redux";
import { getTranslations } from "../../../config/language";
import { createOrganizer, fetchOrganizers, updateOrganizer } from "../../../modules/organizers/organizerSlice";

export function AdminOrganizersPage() {
    const dispatch = useAppDispatch();
    const items = useAppSelector((s) => (Array.isArray(s.organizers.items) ? s.organizers.items : []));
    const loading = useAppSelector((s) => s.organizers.loading);
    const loadError = useAppSelector((s) => s.organizers.error);
    const [open, setOpen] = useState(false);
    const [editingId, setEditingId] = useState<string | null>(null);
    const [form, setForm] = useState({ name: "", email: "", phone: "", slug: "" });
    const [message, setMessage] = useState("");
    const language = useAppSelector((s) => s.language.language);
    const t = getTranslations(language).adminOrganizers;

    useEffect(() => {
        void dispatch(fetchOrganizers());
    }, [dispatch]);

    const openCreate = () => {
        setEditingId(null);
        setForm({ name: "", email: "", phone: "", slug: "" });
        setOpen(true);
    };

    const openEdit = (organizer: (typeof items)[number]) => {
        setEditingId(organizer.id);
        setForm({
            name: organizer.name,
            email: organizer.email,
            phone: organizer.phone ?? "",
            slug: organizer.slug,
        });
        setOpen(true);
    };

    const submit = async (e: FormEvent) => {
        e.preventDefault();
        try {
            if (editingId) {
                await dispatch(updateOrganizer({ id: editingId, ...form })).unwrap();
                setMessage(t.updateSuccess);
            } else {
                await dispatch(createOrganizer(form)).unwrap();
                setMessage(t.createSuccess);
            }
            setForm({ name: "", email: "", phone: "", slug: "" });
            setEditingId(null);
            setOpen(false);
        } catch (err) {
            setMessage(err instanceof Error ? err.message : t.saveError);
        }
    };

    return (
        <div className="stack">
            <Card title={t.title}>
                <div className="stack">
                    <Button type="button" onClick={openCreate}>{t.addOrganizer}</Button>
                    {message && <p>{message}</p>}
                </div>
            </Card>

            <Card title={t.organizerList}>
                {loading && <p>{t.loading}</p>}
                {loadError && <p className="error">{loadError}</p>}
                {!loading && !loadError && items.length === 0 && <p>{t.empty}</p>}
                {!loading && !loadError && items.length > 0 && (
                    <div className="table-wrapper">
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>{t.name}</th>
                                    <th>{t.email}</th>
                                    <th>{t.phone}</th>
                                    <th>{t.slug}</th>
                                    <th>{t.status}</th>
                                    <th>{t.actions}</th>
                                </tr>
                            </thead>
                            <tbody>
                                {items.map((organizer) => (
                                    <tr key={organizer.id}>
                                        <td className="table-primary">{organizer.name}</td>
                                        <td>{organizer.email}</td>
                                        <td>{organizer.phone ?? t.noPhone}</td>
                                        <td>{organizer.slug}</td>
                                        <td>
                                            <span className={`status ${organizer.is_active ? "status-active" : "status-inactive"}`}>
                                                {organizer.is_active ? t.active : t.inactive}
                                            </span>
                                        </td>
                                        <td>
                                            <Button type="button" onClick={() => openEdit(organizer)}>{t.edit}</Button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </Card>

            <Modal open={open} title={editingId ? t.editOrganizer : t.addNewOrganizer} onClose={() => setOpen(false)}>
                <form className="form" onSubmit={submit}>
                    <input placeholder={t.name} value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
                    <input placeholder={t.email} value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
                    <input placeholder={t.phone} value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
                    <input placeholder={t.slug} value={form.slug} onChange={(e) => setForm({ ...form, slug: e.target.value })} />
                    <div className="modal-actions">
                        <Button type="button" onClick={() => setOpen(false)}>{t.cancel}</Button>
                        <Button type="submit">{t.save}</Button>
                    </div>
                </form>
            </Modal>
        </div>
    );
}
