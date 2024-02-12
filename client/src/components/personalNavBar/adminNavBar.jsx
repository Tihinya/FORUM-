import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../Gachi.js/src/core/framework.ts"

export default function AdminNavBar() {
	const { props, setProps } = useContext("currentProps")

	const toggleFilter = (filterType) => {
		if (props !== filterType) {
			setProps("")
			setProps(filterType)
		} else {
			setProps("")
		}
	}
	return (
		<div>
			<div className="personal-navigator__menu">
				<p>Admin Navigator</p>
			</div>
			<div className="personal-navigator__categories">
				<div
					className={`nav__options_hover ${
						props === "promotions" ? "nav__options_active" : ""
					}`}
					onClick={() => {
						toggleFilter("promotions")
					}}
				>
					Promotions
				</div>

				<div
					className={`nav__options_hover ${
						props === "demotions" ? "nav__options_active" : ""
					}`}
					onClick={() => {
						toggleFilter("demotions")
					}}
				>
					Demotions
				</div>

				<div
					className={`nav__options_hover ${
						props === "categories" ? "nav__options_active" : ""
					}`}
					onClick={() => {
						toggleFilter("categories")
					}}
				>
					Create Category
				</div>

				<div
					className={`nav__options_hover ${
						props === "postreport/get" ? "nav__options_active" : ""
					}`}
					onClick={() => {
						toggleFilter("postreport/get")
					}}
				>
					Reports
				</div>
			</div>
		</div>
	)
}
