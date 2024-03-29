import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"
import { fetchData } from "../../additional-funcitons/api.js"
import ErrorWindow from "../errors/error-window"

export default function CreatePost() {
	const navigate = useNavigate()

	const isLoggin = useContext("isAuthenticated").isAuthenticated
	const categorieUrl = "categories"
	const createPostUrl = "post"

	const { setPosts } = useContext("currentPosts")
	const [selectedCategories, setSelectedCategories] = useState([])
	const [categories, setCategories] = useState([])
	const [errorMessage, setErrorMessage] = useState("")
	const [threadClicked, setThreadClicked] = useState(false)
	const [selectedImage, setSelectedImage] = useState("")
	const [formData, setFormData] = useState({
		title: "",
		content: "",
		image: "",
		categories: [],
	})
	useEffect(() => {
		fetchData(null, categorieUrl, "GET")
			.then((resultInJson) => {
				setCategories(resultInJson)
			})
			.catch((error) => {
				setErrorMessage("Failed to fetch categories: " + error.message)
			})

		setFormData((prevData) => ({
			...prevData,
			categories: selectedCategories,
		}))
	}, [selectedCategories])

	const handleInputChange = (e) => {
		const { name, value, type } = e.target

		// Check if the input field is of type "file" (for images)
		if (type === "file") {
			const file = e.target.files[0]
			const imageSize = 20 * 1024 * 1024 // 20MB

			if (file) {
				// Handle image file upload
				const reader = new FileReader()

				// Check if the image size is within limits
				if (file.size > imageSize) {
					setErrorMessage("File is too big, max size is 20Mb")
					return
				}

				reader.onload = (event) => {
					const imageURL = event.target.result

					setFormData((prevData) => ({
						...prevData,
						[name]: imageURL, // Store the image data (URL) in the form data
					}))
					setSelectedImage(imageURL) // Optionally, update the selectedImage state
				}
				reader.readAsDataURL(file)
			}
		} else {
			// For non-image fields, update the form data as usual
			setFormData((prevData) => ({
				...prevData,
				[name]: value,
			}))
		}
	}

	const handleSubmitClick = (e) => {
		e.preventDefault()

		// This is for input data disappearing on every notification fetch
		// Going to switch back to old method when implementing React
		const form = e.target
		const formDatafied = new FormData(form)
		const formJson = Object.fromEntries(formDatafied.entries())

		formJson.image = formData.image
		formJson.categories = formData.categories

		fetchData(formJson, createPostUrl, "POST")
			.then((resultInJson) => {
				if (resultInJson.status === "success") {
					setThreadClicked("")
					setFormData({
						title: "",
						content: "",
						image: "",
						categories: [],
					})
					setSelectedCategories([])
					form.reset()

					fetchData(null, "posts", "GET").then((resultInJson) => {
						setPosts(resultInJson)
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

	const selectCategory = (category) => {
		if (selectedCategories.includes(category)) {
			setSelectedCategories((prevCategories) =>
				prevCategories.filter((cat) => cat !== category)
			)
		} else {
			setSelectedCategories((prevCategories) => [
				...prevCategories,
				category,
			])
		}
	}

	const imageHandlerDelete = () => {
		setSelectedImage("")
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
				className={
					isLoggin ? "add-thread" : "add-thread-closed add-thread"
				}
				id="thread-window"
			>
				<div
					className="thread-button"
					id="add-a-thread"
					onClick={() => setThreadClicked(!threadClicked)}
				>
					+
				</div>
				<input
					type="text"
					name="title"
					maxlength="120"
					id="titleValue"
					placeholder="Add a thread"
				/>
				<div
					className={
						threadClicked ? "thread-window-open" : "thread-window"
					}
					id="detailed-thread"
				>
					<div className="thread-options">
						<div className="upload-image">
							<label
								for="image-file-upload"
								className="custom-upload-file-button"
							>
								<input
									type="file"
									name="image"
									onChange={handleInputChange}
									id="image-file-upload"
								/>
							</label>
						</div>
						<textarea
							id="contentValue"
							className="thread-text"
							name="content"
							placeholder="Description here"
							rows={10}
						/>
						<div className="thread-tags">
							{categories.map((category) => (
								<p
									className={`thread-subject ${
										selectedCategories.includes(
											category.category
										)
											? "active"
											: ""
									}`}
									onClick={() =>
										selectCategory(category.category)
									}
								>
									{category.category}
								</p>
							))}
						</div>
					</div>

					<div className="create-post-button">
						<div className="create-post-image-added-container">
							{selectedImage && (
								<div className="create-post-image-added">
									<img
										className="create-post-image"
										src={selectedImage}
										alt="Select Image"
									/>
									<button
										className="sign__button"
										onClick={imageHandlerDelete}
									>
										Remove Image
									</button>
								</div>
							)}
						</div>

						<button className="sign__button" type="submit">
							Create Post
						</button>
					</div>
				</div>
			</form>
		</div>
	)
}
