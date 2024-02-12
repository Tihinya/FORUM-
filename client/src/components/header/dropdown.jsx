import Gachi, {
	useNavigate,
	useContext,
} from "../../../Gachi.js/src/core/framework.ts"
import { getWebSocket } from "../socket/socket"

let baseURL = "https://ec2-51-20-1-125.eu-north-1.compute.amazonaws.com:8080"
const DropdownMenu = () => {
	const { isAuthenticated, setIsAuthenticated } =
		useContext("isAuthenticated")
	const navigate = useNavigate()
	const ws = getWebSocket()

	return (
		<div className="dropdown">
			<button className="dropdown-button"></button>
			<div className="dropdown-content">
				<button
					onClick={() => {
						ws.close()

						fetch(baseURL + "/logout", {
							credentials: "include",
							method: "POST",
						})

						localStorage.removeItem("id")
						setIsAuthenticated(false)
						navigate("/")
					}}
				>
					LogOut
				</button>
			</div>
		</div>
	)
}

export default DropdownMenu
