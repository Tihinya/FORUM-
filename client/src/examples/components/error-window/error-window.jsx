import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../core/framework"

export default function ErrorWindow({ errorMessage }) {
	return <div className="error-window">{errorMessage}</div>
}
