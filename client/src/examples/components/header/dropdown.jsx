import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

const DropdownMenu = () => {
	const navigate = useNavigate()

	// const [selectedOption, setSelectedOption] = useState("") // To keep track of the selected option

	// const handleOptionChange = (event) => {
	// 	setSelectedOption(event.target.value)
	// }

	return (
		<div className="dropdown">
			<button className="dropdown-button"></button>
			<div className="dropdown-content">
				{/* <button onClick={() => navigate("/")}>Profile Page</button> */}
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
