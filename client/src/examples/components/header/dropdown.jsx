import Gachi, { useNavigate } from "../../../core/framework"

const DropdownMenu = () => {
	const navigate = useNavigate()

	return (
		<div className="dropdown">
			<button className="dropdown-button"></button>
			<div className="dropdown-content">
				<button
					onClick={() => {
						localStorage.removeItem("id")
						navigate("/login")
						fetch("https://localhost:8080/logout", {
							credentials: "include",
						})
					}}
				>
					LogOut
				</button>
			</div>
		</div>
	)
}

export default DropdownMenu
