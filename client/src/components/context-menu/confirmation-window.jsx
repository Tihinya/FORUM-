import Gachi from "../../../Gachi.js/src/core/framework.ts"

export default function ConfirmationWindow({ message, onYes, onNo }) {
	return (
		<div className="confirmation-window">
			{message}
			<button onClick={onYes}>Yes</button>
			&nbsp;
			<button onClick={onNo}>No</button>
		</div>
	)
}
