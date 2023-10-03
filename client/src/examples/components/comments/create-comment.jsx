import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

import ErrorWindow from "../errors/error-window"
import { fetchData } from "../../additional-funcitons/api.js"

export default function CreateComment({ endPointUrl, userId }) {
	const navigate = useNavigate()
	const { setComments } = useContext("currentComment")
	const [errorMessage, setErrorMessage] = useState("")
	const isLoggin = useContext("isAuthenticated").isAuthenticated

	const [formData, setFormData] = useState({
		content: "",
	})
	const point = `comments/${userId}`

	const endpoint = `${endPointUrl}/${userId}`

	const handleSubmitClick = (e) => {
		e.preventDefault()

		fetchData(formData, endpoint, "POST")
			.then((resultInJson) => {
				if (resultInJson.status === "success") {
					setFormData({
						...formData,
						content: "",
					})

					fetchData(null, point, "GET").then((resultInJson) => {
						setComments(resultInJson)
					})
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
	const handleInputChange = (e) => {
		const { name, value } = e.target
		setFormData((prevData) => ({
			...prevData,
			[name]: value,
		}))
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
						value={formData.content}
						onChange={handleInputChange}
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
