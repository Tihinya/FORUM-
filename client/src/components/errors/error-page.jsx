import Gachi from "../../../Gachi.js/src/core/framework.ts"

export default function ErrorPage({ error }) {
	return (
		<div className="error__page">
			<div className="error__message">{error.status}</div>
			<div className="message">{error.message}</div>
		</div>
	)
}
