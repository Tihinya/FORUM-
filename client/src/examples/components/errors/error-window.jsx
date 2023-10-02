import Gachi from "../../../core/framework"

export default function ErrorWindow({ errorMessage, onClose }) {
	return (
		<div className="error-window">
			<button onClick={onClose}>[X]</button>
			{errorMessage}
		</div>
	)
}
