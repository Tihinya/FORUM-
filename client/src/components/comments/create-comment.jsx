import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"

import ErrorWindow from "../errors/error-window"
import { fetchData } from "../../additional-funcitons/api.js"

export default function CreateComment({ endPointUrl, userId }) {
	const navigate = useNavigate()
	const { setComments } = useContext("currentComments")
	const [errorMessage, setErrorMessage] = useState("")
	const isLoggin = useContext("isAuthenticated").isAuthenticated

	const point = `comments/${userId}`

	const endpoint = `${endPointUrl}/${userId}`

	const handleSubmitClick = (e) => {
		e.preventDefault()

		const form = e.target
		const formData = new FormData(form)
		const formJson = Object.fromEntries(formData.entries())

		fetchData(formJson, endpoint, "POST")
			.then((resultInJson) => {
				if (resultInJson.status === "success") {
					fetchData(null, point, "GET").then((resultInJson) => {
						setComments(resultInJson)
					})

					form.reset()
				} else if (resultInJson.status === "error") {
					setErrorMessage(resultInJson.message)
				}
			})
			.catch((error) => {
				navigate("serverded")
				console.error("Error :", error)
			})
	}
	const handleErrorMessageClose = () => {
		setErrorMessage("")
	}

	return (
		<div className="post__container">
			{errorMessage != "" ? (
				<ErrorWindow
					errorMessage={errorMessage}
					onClose={handleErrorMessageClose}
				/>
			) : (
				""
			)}
			<form
				onSubmit={handleSubmitClick}
				className={isLoggin ? "post__box" : "post__box_closed "}
			>
				<p className="post__box_comment-message">Leave your comment</p>
				<div className="input-fields">
					<textarea
						className="text-area"
						name="content"
						rows="5"
						cols="200"
						placeholder="Type here"
					></textarea>
				</div>
				<div className="promotion-message__buttons">
					<button type="submit" className="sign__button-orange">
						Leave comment
					</button>
				</div>
			</form>
		</div>
	)
}
