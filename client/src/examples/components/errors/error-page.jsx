import Gachi from "../../../core/framework"

export default function ErrorPage({ error }) {
	return (
		<div className="error__page">
			<div className="error__message">{error.status}</div>
			<div className="message">{error.message}</div>
		</div>
	)
}
