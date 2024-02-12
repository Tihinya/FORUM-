import Gachi from "../../../Gachi.js/src/core/framework.ts"

export default function ErrorWindow({ errorMessage, onClose }) {
	return (
		<div className="error-window">
			<button onClick={onClose}>[X]</button>
			{errorMessage}
		</div>
	)
}
