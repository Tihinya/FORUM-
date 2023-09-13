import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../core/framework"

export default function ErrorWindow({ errorMessage, onClose }) {
	return (
		<div className="error-window">
			<button onClick={onClose}>[X]</button>
			{errorMessage}
		</div>
	)
}
