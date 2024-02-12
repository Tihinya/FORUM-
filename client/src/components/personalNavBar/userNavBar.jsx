import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../Gachi.js/src/core/framework.ts"

export default function UserNavBar() {
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
				<p>Promotion Navigator</p>
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
			</div>
		</div>
	)
}
