import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "../components/ui/Button";
import { Card } from "../components/ui/Card";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { login } from "../modules/auth/authSlice";
export function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { loading, error } = useAppSelector((s) => s.auth);
  const submit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await dispatch(login({ email, password })).unwrap();
      navigate("/events");
    } catch {
      // The rejected thunk already exposes the message through Redux state.
    }
  };
  return (
    <Card title="Sign in">
      <form onSubmit={submit} className="form">
        <input placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
        <input
          placeholder="Password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <Button loading={loading}>Sign in</Button>
        {error && <p className="error">{error}</p>}
      </form>
    </Card>
  );
}
