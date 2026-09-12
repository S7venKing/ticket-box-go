import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "../components/ui/Button";
import { Card } from "../components/ui/Card";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { login } from "../modules/auth/authSlice";
import { getTranslations } from "../config/language";

export function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { loading, error } = useAppSelector((s) => s.auth);
  const language = useAppSelector((s) => s.language.language);
  const t = getTranslations(language);
  const submit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      const result = await dispatch(login({ email, password })).unwrap();
      navigate(result.user.role === "admin" ? "/admin" : "/events");
    } catch {
      // The rejected thunk already exposes the message through Redux state.
    }
  };
  return (
    <Card title={t.login.title}>
      <form onSubmit={submit} className="form">
        <input placeholder={t.common.email} value={email} onChange={(e) => setEmail(e.target.value)} />
        <input
          placeholder={t.common.password}
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <Button loading={loading}>{t.login.submit}</Button>
        {error && <p className="error">{error}</p>}
      </form>
    </Card>
  );
}
