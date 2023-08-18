import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../core/framework"

export function NavBar() {
	const [activeSubj, setActiveSubj] = useState(0)
	const detailsVisible = useState(false)

	const subjects = ["UX/UI", "JavaScript", "Golang", "Wisdom"]

	const handleSubjectClick = (index) => {
		if (activeSubj !== index) {
			setActiveSubj(index)
		}
	}

	return (
		<div className="nav__menu">
			<div className="menu__logo">Menu</div>
			<div className="nav__options">
				{subjects.map((subject, index) => (
					<p
						key={index}
						className={`nav__options_hover ${
							activeSubj === index ? "nav__options_active" : ""
						}`}
						onClick={() => handleSubjectClick(index)}
					>
						{subject}
					</p>
				))}
			</div>
			<div className={`detailed-thread ${detailsVisible ? "show" : ""}`}>
				{/* Detailed thread content */}
			</div>
		</div>
	)
}
