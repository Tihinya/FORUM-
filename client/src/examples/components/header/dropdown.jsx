import Gachi, { useNavigate, useContext } from "../../../core/framework"
import { logoutRequest } from "../../additional-funcitons/authorization"

const DropdownMenu = () => {
	const { isAuthenticated, setIsAuthenticated } = useContext("isAuthenticated")
	const navigate = useNavigate()

	return (
		<div className="dropdown">
			<button className="dropdown-button"></button>
			<div className="dropdown-content">
				<button
					onClick={() => {
						logoutRequest()
						setIsAuthenticated(false)
						navigate("/login")
						fetch("https://localhost:8080/logout", {
							credentials: "include",
						})
						window.location.reload()
					}}
				>
					LogOut
				</button>
			</div>
		</div>
	)
}

export default DropdownMenu
