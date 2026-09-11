interface Props extends React.ButtonHTMLAttributes<HTMLButtonElement> { loading?: boolean; }
export function Button({ loading, children, ...props }: Props) { return <button className="button" disabled={loading || props.disabled} {...props}>{loading ? "Loading..." : children}</button>; }
