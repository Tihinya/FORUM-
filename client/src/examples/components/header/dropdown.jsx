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
				<button>Button 1</button>
				<button>Button 2</button>
				<button
					onClick={() => {
						localStorage.removeItem("id")
						navigate("/login")
					}}
				>
					Button 3
				</button>
			</div>
		</div>
	)
}

export default DropdownMenu
