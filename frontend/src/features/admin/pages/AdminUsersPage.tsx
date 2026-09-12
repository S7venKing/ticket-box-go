import { FormEvent, useEffect, useState } from "react";
import { Button } from "../../../components/ui/Button";
import { Card } from "../../../components/ui/Card";
import { Modal } from "../../../components/ui/Modal";
import { useAppDispatch, useAppSelector } from "../../../hooks/redux";
import { post } from "../../../lib/api";
import { register } from "../../../modules/auth/authSlice";
import { fetchUsers } from "../../../modules/users/userSlice";
import { getTranslations } from "../../../config/language";

export function AdminUsersPage() {
    const dispatch = useAppDispatch();
    const users = useAppSelector((s) => (Array.isArray(s.users.items) ? s.users.items : []));
    const loading = useAppSelector((s) => s.users.loading);
    const loadError = useAppSelector((s) => s.users.error);
    const [open, setOpen] = useState(false);
    const [form, setForm] = useState({ email: "", password: "", full_name: "", phone: "" });
    const [message, setMessage] = useState("");
    const [organizerIDs, setOrganizerIDs] = useState<Record<string, string>>({});
    const language = useAppSelector((s) => s.language.language);

    const t = getTranslations(language).adminUsers;

    useEffect(() => {
        void dispatch(fetchUsers());
    }, [dispatch]);

    const submit = async (e: FormEvent) => {
        e.preventDefault();
        try {
            await dispatch(register(form)).unwrap();
            await dispatch(fetchUsers()).unwrap();
            setForm({ email: "", password: "", full_name: "", phone: "" });
            setOpen(false);
            setMessage(t.userCreated);
        } catch (err) {
            setMessage(err instanceof Error ? err.message : t.createUserFailed);
        }
    };

    return (
        <div className="stack">
            <Card title={t.title}>
                <div className="stack">
                    <Button type="button" onClick={() => setOpen(true)}>{t.addUser}</Button>
                    {message && <p>{message}</p>}
                </div>
            </Card>

            <Card title={t.userList}>
                {loading && <p>{t.loading}</p>}
                {loadError && <p className="error">{loadError}</p>}
                {!loading && !loadError && users.length === 0 && <p>{t.empty}</p>}
                {!loading && !loadError && users.length > 0 && (
                    <div className="table-wrapper">
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>{t.fullName}</th>
                                    <th>{t.email}</th>
                                    <th>{t.phone}</th>
                                    <th>{t.role}</th>
                                    <th>{t.status}</th>
                                    <th>{t.assignOrganizer}</th>
                                </tr>
                            </thead>
                            <tbody>
                                {users.map((user) => (
                                    <tr key={user.id}>
                                        <td className="table-primary">{user.full_name || t.noName}</td>
                                        <td>{user.email}</td>
                                        <td>{user.phone || t.noPhone}</td>
                                        <td><span className="status">{user.role}</span></td>
                                        <td>
                                            <span className={`status ${user.is_active ? "status-active" : "status-inactive"}`}>
                                                {user.is_active ? t.active : t.inactive}
                                            </span>
                                        </td>
                                        <td>
                                            <input placeholder={t.organizerId} value={organizerIDs[user.id] ?? ""}
                                                onChange={(e) => setOrganizerIDs({ ...organizerIDs, [user.id]: e.target.value })} />
                                            <Button type="button" onClick={async () => {
                                                try {
                                                    await post(`/admin/organizers/${organizerIDs[user.id]}/members`, { user_id: user.id });
                                                    setMessage(t.assignedUser);
                                                } catch (err) { setMessage(err instanceof Error ? err.message : t.assignFailed); }
                                            }}>{t.assign}</Button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </Card>

            <Modal open={open} title={t.addNewUser} onClose={() => setOpen(false)}>
                <form className="form" onSubmit={submit}>
                    <input placeholder={t.email} value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
                    <input type="password" placeholder={t.password} value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} />
                    <input placeholder={t.fullName} value={form.full_name} onChange={(e) => setForm({ ...form, full_name: e.target.value })} />
                    <input placeholder={t.phone} value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
                    <div className="modal-actions">
                        <Button type="button" onClick={() => setOpen(false)}>{t.cancel}</Button>
                        <Button type="submit">{t.save}</Button>
                    </div>
                </form>
            </Modal>
        </div>
    );
}
