import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../core/framework"

export default function ProfilePage() {
	const navigate = useNavigate()

	const bad = localStorage.getItem("id")

	if (!bad) {
		navigate("/bad")
	}
	return <div>Profile Page</div>
}
