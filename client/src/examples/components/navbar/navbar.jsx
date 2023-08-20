import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

export function NavBar() {
	// const [activeSubj, setActiveSubj] = useState(0)
	const { activeSubj, setActiveSubj } = useContext("currentCategory")
	const detailsVisible = useState(false)
	const [categories, setCategories] = useState([])

	// const subjects = ["UX/UI", "JavaScript", "Golang", "Wisdom"]
	useEffect(() => {
		fetch("http://localhost:8080/categories")
			.then((response) => response.json())
			.then((data) => setCategories(data))
			.catch((error) => console.error("Error fetching posts:", error))
	}, [])

	const handleSubjectClick = (index) => {
		if (activeSubj !== index) {
			setActiveSubj(index)
		} else {
			setActiveSubj("")
		}
	}

	return (
		<div className="nav__menu">
			<div className="menu__logo">Menu</div>
			<div className="nav__options">
				{categories.map(({ category: subject }, index) => (
					<p
						key={index}
						className={`nav__options_hover ${
							activeSubj === subject ? "nav__options_active" : ""
						}`}
						onClick={() => handleSubjectClick(subject)}
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
