import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

const DropdownMenu = () => {
	const navigate = useNavigate()

	const [selectedOption, setSelectedOption] = useState("") // To keep track of the selected option

	const handleOptionChange = (event) => {
		setSelectedOption(event.target.value)
	}

	return (
		<div className="profile-popup">
			<select value={selectedOption} onChange={handleOptionChange}>
				<option value="option1">Logined as {}</option>
				<option
					value="option2"
					onClick={() => navigate("/profile-page")}
				>
					Personal Cabinet
				</option>
				<option value="option3" onClick={() => navigate("/")}>
					Logout
				</option>
			</select>
		</div>
	)
}

export default DropdownMenu
