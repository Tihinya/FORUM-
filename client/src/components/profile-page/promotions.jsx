import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"
import { fetchData } from "../../additional-funcitons/api"

export default function Promotions() {
	const { promotions, setPromotions } = useContext("promotions")
	const { props } = useContext("currentProps")
	const { userRole } = useContext("currentUserRole")
	const { userId } = useContext("currentUserId")
	const [response, setResponse] = useState("")

	function handlePromotionClick(requestType, userId) {
		const endpoint = `promotion/${userId}/${requestType}`

		fetchData(null, endpoint, "PATCH").then((resultInJson) => {
			fetchData(null, `promotions/get`, "GET").then((resultInJson) => {
				setPromotions(resultInJson)
			})
		})
	}

	function handleRequestClick() {
		const endpoint = `promotion/${userId}`
		const bodyData = {
			new_role: "moderator",
		}

		fetchData(bodyData, endpoint, "POST")
			.then((response) => setResponse(response))
			.catch((err) => console.error(err))
	}

	return (
		<>
			{props == "promotions" ? (
				<div className="post__container">
					{userRole !== "admin" ? (
						<div className="post__box promotion">
							{response == "" ? (
								<h3>
									Do you want to request for a promotion to
									moderator status?
								</h3>
							) : (
								<h3>You have a request pending</h3>
							)}

							<div className="category__buttons">
								<button
									className="sign__button category"
									onClick={() => handleRequestClick()}
								>
									Oh yes sir!
								</button>
							</div>
						</div>
					) : (
						<>
							{promotions.map((promotion) => {
								return (
									<div className="post__box promotion">
										<div className="promotion__text">
											<h3>
												User {promotion.UserName}{" "}
												requested a promotion to{" "}
												{promotion.RoleName}
											</h3>
										</div>
										<div className="promotion__buttons">
											<div
												className="promotion-button accept"
												onClick={() =>
													handlePromotionClick(
														"accept",
														promotion.UserID
													)
												}
											>
												Accept
											</div>
											<div
												className="promotion-button decline"
												onClick={() =>
													handlePromotionClick(
														"decline",
														promotion.UserID
													)
												}
											>
												Decline
											</div>
										</div>
									</div>
								)
							})}
						</>
					)}
				</div>
			) : null}
		</>
	)
}
