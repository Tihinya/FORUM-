import Gachi, {
	useContext,
	useState,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"

import { fetchData } from "../../additional-funcitons/api.js"

export function NavBar() {
	const { activeSubj, setActiveSubj } = useContext("currentCategory")
	const { selectedModerator, setSelectedModerator } =
		useContext("selectedModerator")
	const detailsVisible = useState(false)
	const { categories, setCategories } = useContext("categories")
	const { displayNavbar } = useContext("displayNavbar")
	const { moderators } = useContext("currentModerators")

	useEffect(() => {
		fetchData(null, "categories", "GET").then((resultInJson) => {
			setCategories(resultInJson)
		})
	}, [])

	const handleSubjectClick = (index) => {
		if (activeSubj !== index) {
			setActiveSubj(index)
		} else {
			setActiveSubj("")
		}
	}

	const handleModeratorClick = (moderator) => {
		if (selectedModerator !== moderator) {
			setSelectedModerator(moderator)
		} else {
			setSelectedModerator([])
		}
	}

	if (displayNavbar) {
		// Returns regular categories navbar
		return (
			<div className="nav__menu">
				<div className="menu__logo">Menu</div>
				<div className="nav__options">
					{categories.map(({ category: subject }, index) => (
						<p
							key={index}
							className={`nav__options_hover ${
								activeSubj === subject
									? "nav__options_active"
									: ""
							}`}
							onClick={() => handleSubjectClick(subject)}
						>
							{subject}
						</p>
					))}
				</div>
				<div
					className={`detailed-thread ${
						detailsVisible ? "show" : ""
					}`}
				></div>
			</div>
		)
	} else {
		// Returns moderators list navbar
		return (
			<div className="nav__menu">
				<div className="menu__logo">Moderators</div>
				<div className="nav__options">
					{moderators.map((moderator, index) => (
						<p
							key={index}
							className={`nav__options_hover ${
								selectedModerator === moderator
									? "nav__options_active"
									: ""
							}`}
							onClick={() => handleModeratorClick(moderator)}
						>
							{moderator.username}
						</p>
					))}
				</div>
				<div
					className={`detailed-thread ${
						detailsVisible ? "show" : ""
					}`}
				></div>
			</div>
		)
	}
}
