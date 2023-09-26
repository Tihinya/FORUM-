import Gachi, { useContext, useState, useEffect } from "../../../core/framework"

import { fetchData } from "../../additional-funcitons/api.js"

export function NavBar() {
	const { activeSubj, setActiveSubj } = useContext("currentCategory")
	const detailsVisible = useState(false)
	const [categories, setCategories] = useState([])

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
			<div
				className={`detailed-thread ${detailsVisible ? "show" : ""}`}
			></div>
		</div>
	)
}
