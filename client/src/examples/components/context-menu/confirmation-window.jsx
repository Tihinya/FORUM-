import Gachi from "../../../core/framework"

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
