import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"
import isLogin from "../../additional-funcitons/isLogin"
import { fetchData } from "../../additional-funcitons/api.js"
import ErrorWindow from "../errors/error-window"

export default function CreatePost() {
	const navigate = useNavigate()

	const isLoggin = isLogin()
	const categorieUrl = "categories"
	const createPostUrl = "post"

	const { setPosts } = useContext("currentPosts")
	const [selectedCategories, setSelectedCategories] = useState([])
	const [errorMessage, setErrorMessage] = useState("")
	const [categories, setCategories] = useState([])
	const [threadClicked, setThreadClicked] = useState(false)
	const [selectedImage, setSelectedImage] = useState("")
	const [formData, setFormData] = useState({
		title: "",
		content: "",
		image: "",
	})

	useEffect(() => {
		fetchData(null, categorieUrl, "GET")
			.then((resultInJson) => {
				setCategories(resultInJson)
			})
			.catch((error) => {
				setErrorMessage("Failed to fetch categories: " + error.message)
			})
	}, [])

	const handleThreadButtonClick = () => {
		if (!threadClicked) {
			setThreadClicked(true)
		} else {
			setThreadClicked(false)
		}
	}

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
					alert("File is too big, max size is 20Mb")
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
		fetchData(formData, createPostUrl, "POST")
			.then((resultInJson) => {
				if (resultInJson.status === "success") {
					setThreadClicked("")
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
	const imageHandler = (e) => {
		const file = e.target.files[0]
		const imageSize = 20 * 1024 * 1024 // 20MB

		if (file) {
			if (file.size > imageSize) {
				alert("File is too big, max size is 20Mb")
				return
			}

			const reader = new FileReader()

			reader.onload = (event) => {
				const imageURL = event.target.result

				setSelectedImage(imageURL)

				setFormData((prevData) => ({
					...prevData,
					image: imageURL,
				}))
			}

			reader.readAsDataURL(file)
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
					onClick={handleThreadButtonClick}
				>
					+
				</div>
				<input
					type="text"
					name="title"
					maxlength="120"
					value={formData.title}
					onChange={handleInputChange}
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
									onChange={imageHandler}
									id="image-file-upload"
								/>
							</label>
						</div>
						<textarea
							value={formData.content}
							onChange={handleInputChange}
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
