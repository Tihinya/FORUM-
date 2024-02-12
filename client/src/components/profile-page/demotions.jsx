import Gachi, {
	useContext,
	useState,
} from "../../../Gachi.js/src/core/framework.ts"
import { fetchData } from "../../additional-funcitons/api"

export default function Demotions() {
	const { props } = useContext("currentProps")
	const { setDisplayNavbar } = useContext("displayNavbar")
	const { setModerators } = useContext("currentModerators")
	const { selectedModerator, setSelectedModerator } =
		useContext("selectedModerator")
	const [response, setResponse] = useState("")

	setDisplayNavbar(false)

	function handleDemoteClick(moderator) {
		const demoteEndpoint = `promotion/${moderator.id}/demote`
		const fetchEndpoint = `users/get`
		const bodyData = {
			new_role: "user",
		}

		fetchData(bodyData, demoteEndpoint, "POST")
			.then((response) => {
				if (response !== "success") {
					setResponse(response)
				}
				fetchData(null, fetchEndpoint, "GET").then((usersInJson) => {
					const moderators = usersInJson.filter(
						(user) => user.role_id == 2
					) // Filter users who have moderator role
					setModerators(moderators)
				})
			})
			.catch((err) => console.error(err))
	}

	return (
		<div className="post__container">
			{props == "demotions" ? (
				<div className="post__box category">
					{response == "" ? (
						<h3>
							Select a moderator on the left-side navigation bar
							to demote
						</h3>
					) : (
						<h3>You demoted the moderator, savage!</h3>
					)}

					<div className="category__buttons">
						<button
							className="sign__button category delete"
							onClick={() => handleDemoteClick(selectedModerator)}
						>
							Demote User
						</button>
					</div>
				</div>
			) : null}
		</div>
	)
}
