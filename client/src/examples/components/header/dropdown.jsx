import Gachi, { useNavigate, useContext } from "../../../core/framework"

const DropdownMenu = () => {
	const { isAuthenticated, setIsAuthenticated } = useContext("isAuthenticated")
	const navigate = useNavigate()

	return (
		<div className="dropdown">
			<button className="dropdown-button"></button>
			<div className="dropdown-content">
				<button
					onClick={() => {
						fetch("https://localhost:8080/logout", {
							credentials: "include",
							method: "POST"
						})
						setIsAuthenticated(false)
						navigate("/")
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
