import Gachi, {
	useContext,
	useEffect,
	useNavigate,
	useState,

	// Add hover-over date to get full creation date
} from "../../../core/framework.ts"
import { convertTime } from "../../additional-funcitons/post.jsx"
import { sendPostId } from "../comments/commentsAuth.jsx"

export function PostsAuth() {
	const [posts, setPosts] = useState([])
	const [categories, setCategories] = useState([])
	const [likedPosts, setLikedPosts] = useState([])
	const [dislikedPosts, setDislikedPosts] = useState([])
	const [threadClicked, setThreadClicked] = useState(false)
	const [selectedCategories, setSelectedCategories] = useState([])
	const [threadTitleValue, setThreadTitleValue] = useState("")
	const [threadContentValue, setThreadContentValue] = useState("")
	const navigate = useNavigate()

	// For displaying liked icon, if the post is already liked (TODO)
	const fetchLikedPosts = () => {
		fetch("https://localhost:8080/user/liked", {
			credentials: "include",
		})
			.then((response) => response.json())
			.then((data) => setLikedPosts(data))
			.catch((error) =>
				console.error("Error fetching liked posts:", error)
			)
	}

	const fetchDislikedPosts = () => {
		fetch("https://localhost:8080/user/disliked", {
			credentials: "include",
		})
			.then((response) => response.json())
			.then((data) => setDislikedPosts(data))
			.catch((error) =>
				console.error("Error fetching liked posts:", error)
			)
	}

	const fetchPosts = () => {
		fetch("https://localhost:8080/posts")
			.then((response) => response.json())
			.then((data) => setPosts(data))
			.catch((error) => console.error("Error fetching posts:", error))
	}

	const fetchCategoriesAndPostCategories = () => {
		Promise.all([
			fetch("https://localhost:8080/categories"),
			fetch("https://localhost:8080/postcategories"),
		])
			.then(([categoriesResponse, postCategoriesResponse]) => {
				return Promise.all([
					categoriesResponse.json(),
					postCategoriesResponse.json(),
				])
			})
			.then(([categoriesData, postCategoriesData]) => {
				sortCategories(categoriesData, postCategoriesData)
			})
			.catch((error) => {
				console.error("Error fetching data:", error)
			})
	}

	const handleThreadButtonClick = () => {
		if (!threadClicked) {
			setThreadClicked(true)
		} else {
			setThreadClicked(false)
		}
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

	const sortCategories = (categoryObj, postCategoryObj) => {
		const arrayedPostCategories = postCategoryObj.map(
			(postCategory) => postCategory.CategoryId
		)
		const countedPostCategories = arrayedPostCategories.reduce(function (
			obj,
			val
		) {
			obj[val] = (obj[val] || 0) + 1
			return obj
		},
		{})
		const ascendingCategories = Object.keys(countedPostCategories).sort(
			(a, b) => countedPostCategories[b] - countedPostCategories[a]
		)

		const popularCategories = ascendingCategories.map((id) => {
			const category = categoryObj.find((category) => category.id == id)
			return category.category
		})

		setCategories(popularCategories)
	}

	function handleSubmit(e) {
		// Prevent the browser from reloading the page
		e.preventDefault()
		// Read the form data
		const form = e.target
		const formData = new FormData(form)
		const formJson = Object.fromEntries(formData.entries())
		createPost(formJson.title, formJson.content, selectedCategories) // TODO CATEGORIES
	}

	// Initialize posts/likes/dislikes upon page load
	useEffect(() => {
		fetchPosts()
		fetchCategoriesAndPostCategories()
		fetchDislikedPosts()
		fetchLikedPosts()
	}, [])

	const createPost = async (title, content, categories) => {
		// img to be added
		try {
			const response = await fetch(`https://localhost:8080/post`, {
				method: "POST",
				credentials: "include",
				headers: {
					"Content-Type": "application/json",
					Accept: "application/json",
				},
				body: JSON.stringify({
					title: title,
					content: content,
					categories: categories,
				}),
			})

			if (response.ok) {
				setThreadContentValue("")
				setThreadTitleValue("")
				setSelectedCategories([])
				handleThreadButtonClick()
				fetchPosts()
			} else {
				const errorData = await response.json()
				console.error(
					response.status,
					response.statusText,
					"-",
					errorData.message
				)
			}
		} catch {
			console.error("you broke the system!")
		}
	}

	const handleLike = async (type, postId) => {
		try {
			if (
				!likedPosts.includes(postId) &&
				!dislikedPosts.includes(postId)
			) {
				const response = await fetch(
					`https://localhost:8080/post/${postId}/${type}`,
					{
						method: "POST",
						credentials: "include",
					}
				)

				const errorData = await response.json()

				if (response.ok) {
					setPosts((prevPosts) => {
						return prevPosts.map((post) => {
							if (post.id === postId) {
								if (type === "like") {
									return { ...post, likes: post.likes + 1 }
								} else {
									return {
										...post,
										dislikes: post.dislikes + 1,
									}
								}
							}
							return post
						})
					})
					fetchLikedPosts()
					fetchDislikedPosts()
				} else {
					console.error(
						response.status,
						response.statusText,
						"-",
						errorData.message
					)
				}
			} else {
				const response = await fetch(
					`https://localhost:8080/post/${postId}/un${type}`,
					{
						method: "POST",
						credentials: "include",
					}
				)

				const errorData = await response.json()

				if (response.ok) {
					setPosts((prevPosts) => {
						return prevPosts.map((post) => {
							if (post.id === postId) {
								if (type === "like") {
									return { ...post, likes: post.likes - 1 }
								} else {
									return {
										...post,
										dislikes: post.dislikes - 1,
									}
								}
							}
							return post
						})
					})
					fetchLikedPosts()
					fetchDislikedPosts()
				} else {
					console.error(
						response.status,
						response.statusText,
						"-",
						errorData.message
					)
				}
			}
		} catch {
			console.error("You are most definitely not logged in")
		}
	}

	return (
		<div className="post__container">
			<form onSubmit={handleSubmit} className="add-thread">
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
					value={threadTitleValue}
					onChange={(e) => setThreadTitleValue(e.target.value)}
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
							<img src="../img/add picture.svg" />
						</div>
						<textarea
							value={threadContentValue}
							onChange={(e) =>
								setThreadContentValue(e.target.value)
							}
							className="thread-text"
							name="content"
							placeholder="Description here"
							rows={10}
						/>
						<div className="thread-tags">
							{categories.slice(0, 5).map((category) => (
								<p
									className={`thread-subject ${
										selectedCategories.includes(category)
											? "active"
											: ""
									}`}
									onClick={() => selectCategory(category)}
								>
									{category}
								</p>
							))}
						</div>
					</div>
					<div className="create-post-button">
						<button className="sign__button" type="submit">
							Create Post
						</button>
					</div>
				</div>
			</form>
			{posts
				.sort(
					(a, b) =>
						new Date(b.creation_date) - new Date(a.creation_date)
				)
				.map((post) => (
					<div className="post__box">
						<div className="post__header">
							<div className="user__info">
								<div className="user__info_picture">
									<a
										onClick={() =>
											navigate("/profile-page")
										}
									>
										<img src="../img/avatarka.jpeg" />
									</a>
								</div>
								<div className="user__info_name">
									<p className="name">
										{post.user_info.username}
									</p>
									<p className="date">
										{convertTime(post.creation_date)}
									</p>
								</div>
							</div>
						</div>
						<div className="post__content">
							<h3>{post.title}</h3>
							<p className="post__text">{post.content}</p>
						</div>
						<div className="post__info">
							<div className="post__tags">
								{post.categories.map((categories) => (
									<p className="tag">{categories}</p>
								))}
							</div>
							<div className="post__likes">
								<a
									onClick={() => {
										sendPostId(post.id)
										navigate(`/comments-authorized`)
									}}
								>
									<img src="../img/message-square.svg" />
								</a>
								<p
									onClick={() => {
										sendPostId(post.id)
										navigate(`/comments-authorized`)
									}}
								>
									{post.comment_count}
								</p>
								<img
									onClick={() => handleLike("like", post.id)}
									src="../img/thumbs-up.svg"
								/>
								<p onClick={() => handleLike("like", post.id)}>
									{post.likes}
								</p>
								<img
									onClick={() =>
										handleLike("dislike", post.id)
									}
									src="../img/thumbs-down.svg"
								/>
								<p
									onClick={() =>
										handleLike("dislike", post.id)
									}
								>
									{post.dislikes}
								</p>
							</div>
						</div>
					</div>
				))}
		</div>
	)
}
