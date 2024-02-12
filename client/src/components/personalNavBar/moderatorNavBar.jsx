import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../Gachi.js/src/core/framework.ts"

export default function ModeratorNavBar() {
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
				<p>Moderator Navigator</p>
			</div>
			<div className="personal-navigator__categories">
				<div
					className={`nav__options_hover ${
						props === "postreport/get" ? "nav__options_active" : ""
					}`}
					onClick={() => {
						toggleFilter("postreport/answer/get")
					}}
				>
					Reports
				</div>
			</div>
		</div>
	)
}
