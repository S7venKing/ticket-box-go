import type { CSSProperties } from "react";
import { Button } from "./Button";
import { useAppDispatch, useAppSelector } from "../../hooks/redux";
import { setLanguage, type Language } from "../../modules/language/languageSlice";

interface LanguageSwitcherProps {
  label?: string;
  style?: CSSProperties;
}

export function LanguageSwitcher({ label = "Language", style }: LanguageSwitcherProps) {
  const dispatch = useAppDispatch();
  const language = useAppSelector((s) => s.language.language);

  const handleChange = (nextLanguage: Language) => {
    dispatch(setLanguage(nextLanguage));
  };

  return (
    <div style={{ display: "flex", gap: "0.5rem", alignItems: "center", justifyContent: "space-between", flexWrap: "wrap", ...style }}>
      <span style={{ fontWeight: 600 }}>{label}:</span>
      <div style={{ display: "flex", gap: "0.5rem" }}>
        <Button type="button" onClick={() => handleChange("en")} style={{ opacity: language === "en" ? 1 : 0.7 }}>
          English
        </Button>
        <Button type="button" onClick={() => handleChange("vi")} style={{ opacity: language === "vi" ? 1 : 0.7 }}>
          Tiếng Việt
        </Button>
      </div>
    </div>
  );
}
