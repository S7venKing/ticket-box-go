import { useAppSelector } from "../../hooks/redux";
import { getTranslations } from "../../config/language";

interface Props extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  loading?: boolean;
}

export function Button({ loading, children, ...props }: Props) {
  const language = useAppSelector((s) => s.language.language);

  return (
    <button className="button" disabled={loading || props.disabled} {...props}>
      {loading ? getTranslations(language).common.loading : children}
    </button>
  );
}
