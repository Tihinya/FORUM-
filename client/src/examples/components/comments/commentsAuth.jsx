import Gachi, {
	useContext,
	useEffect,
	useNavigate,
	useState,
} from "../../../core/framework.ts"
import isLogin from "../../additional-funcitons/isLogin.js"
import { convertTime } from "../../additional-funcitons/post.js"

export function CommentAuth({ postId: navigatePostId }) {
	const isLoggin = isLogin()
	const [posts, setPosts] = useState([])
	const [likedPosts, setLikedPosts] = useState([])
	const [dislikedPosts, setDislikedPosts] = useState([])
	const [likedComments, setLikedComments] = useState([])
	const [dislikedComments, setDislikedComments] = useState([])
	const [comments, setComments] = useState([])
	const [commentValue, setCommentValue] = useState("")

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

	const fetchLikedComments = () => {
		fetch("https://localhost:8080/user/likedComments", {
			credentials: "include",
		})
			.then((response) => response.json())
			.then((data) => setLikedComments(data))
			.catch((error) =>
				console.error("Error fetching liked comments:", error)
			)
	}

	const fetchDislikedComments = () => {
		fetch("https://localhost:8080/user/dislikedComments", {
			credentials: "include",
		})
			.then((response) => response.json())
			.then((data) => setDislikedComments(data))
			.catch((error) =>
				console.error("Error fetching liked comments:", error)
			)
	}

	const fetchPost = () => {
		fetch(`https://localhost:8080/post/${navigatePostId}`)
			.then((response) => response.json())
			.then((data) => setPosts(data))
			.catch((error) => console.error("Error fetching posts:", error))
	}

	const fetchComments = () => {
		fetch(`https://localhost:8080/comments/${navigatePostId}`)
			.then((response) => response.json())
			.then((data) => setComments(data))
			.catch((error) => console.error("Error fetching comments:", error))
	}

	useEffect(() => {
		fetchPost()
		fetchComments()
		fetchDislikedPosts()
		fetchLikedPosts()
		fetchDislikedComments()
		fetchLikedComments()
	}, [])

	function handleSubmit(e) {
		// Prevent the browser from reloading the page
		e.preventDefault()
		// Read the form data
		const form = e.target
		const formData = new FormData(form)
		const formJson = Object.fromEntries(formData.entries())
		createComment(formJson.content)
	}

	const createComment = async (content) => {
		const response = await fetch(
			`https://localhost:8080/comment/${navigatePostId}`,
			{
				method: "POST",
				credentials: "include",
				headers: {
					"Content-Type": "application/json",
					Accept: "application/json",
				},
				body: JSON.stringify({ content: content }),
			}
		)

		if (response.ok) {
			setCommentValue("")
			fetchComments()
		} else {
			const errorData = await response.json()
			console.error(
				response.status,
				response.statusText,
				"-",
				errorData.message
			)
		}
	}

	const handleLikeComment = async (type, commentId) => {
		try {
			if (
				!likedComments.includes(commentId) &&
				!dislikedComments.includes(commentId)
			) {
				const response = await fetch(
					`https://localhost:8080/comment/${commentId}/${type}`,
					{
						method: "POST",
						credentials: "include",
					}
				)

				const errorData = await response.json()

				if (response.ok) {
					setComments((prevComments) => {
						return prevComments.map((comment) => {
							if (comment.id === commentId) {
								if (type === "like") {
									return {
										...comment,
										likes: comment.likes + 1,
									}
								} else {
									return {
										...comment,
										dislikes: comment.dislikes + 1,
									}
								}
							}
							return comment
						})
					})
					fetchLikedComments()
					fetchDislikedComments()
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
					`https://localhost:8080/comment/${commentId}/un${type}`,
					{
						method: "POST",
						credentials: "include",
					}
				)

				const errorData = await response.json()

				if (response.ok) {
					setComments((prevComments) => {
						return prevComments.map((comment) => {
							if (comment.id === commentId) {
								if (type === "like") {
									return {
										...comment,
										likes: comment.likes - 1,
									}
								} else {
									return {
										...comment,
										dislikes: comment.dislikes - 1,
									}
								}
							}
							return comment
						})
					})
					fetchLikedComments()
					fetchDislikedComments()
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

	const handleLikePost = async (type, postId) => {
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
		} catch (error) {
			console.error(error, "You are most definitely not logged in")
		}
	}

	return (
		<div className="post__container">
			{posts.map((post) => (
				<div>
					<div className="post-section">
						<div className="post__box">
							<div className="post__header">
								<div className="user__info">
									<div className="user__info_picture">
										<a href="../html/profile-page.html">
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
								<p className="post-text">{post.content}</p>
							</div>
							<div className="post__info">
								<div className="post__tags">
									{post.categories.map((categories) => (
										<p className="tag">{categories}</p>
									))}
								</div>
								<div className="post__likes">
									<img
										onClick={() =>
											handleLikePost("like", post.id)
										}
										src="../img/thumbs-up.svg"
									/>
									<p
										onClick={() =>
											handleLikePost("like", post.id)
										}
									>
										{post.likes}
									</p>
									<img
										onClick={() =>
											handleLikePost("dislike", post.id)
										}
										src="../img/thumbs-down.svg"
									/>
									<p
										onClick={() =>
											handleLikePost("dislike", post.id)
										}
									>
										{post.dislikes}
									</p>
								</div>
							</div>
						</div>
						<form
							onSubmit={handleSubmit}
							className={
								isLoggin ? "post__box" : "post__box_closed "
							}
						>
							<p className="post__box_comment-message">
								Leave your comment
							</p>
							<div className="input-fields">
								<textarea
									value={commentValue}
									onChange={(e) =>
										setCommentValue(e.target.value)
									}
									className="text-area"
									name="content"
									rows="5"
									cols="200"
									placeholder="Type here"
								></textarea>
							</div>
							<div className="promotion-message__buttons">
								<button
									type="submit"
									className="sign__button-orange"
								>
									{" "}
									Leave comment
								</button>
							</div>
						</form>
					</div>
					{comments
						.sort(
							(a, b) =>
								new Date(b.creation_date) -
								new Date(a.creation_date)
						)
						.map((comment) => (
							<div className="post__box">
								<div className="post__header">
									<div className="user__info">
										<div className="user__info_picture">
											<a href="../html/profile-page.html">
												<img src="../img/avatarka.jpeg" />
											</a>
										</div>
										<div className="user__info_name">
											<p className="name">
												{comment.user_info.username}
											</p>
											<p className="date">
												{convertTime(
													comment.creation_date
												)}
											</p>
										</div>
									</div>
								</div>
								<div className="post__content">
									<p className="post-text">
										{comment.content}
									</p>
								</div>
								<div className="post__likes">
									<img
										onClick={() =>
											handleLikeComment(
												"like",
												comment.id
											)
										}
										src="../img/thumbs-up.svg"
									/>
									<p
										onClick={() =>
											handleLikeComment(
												"like",
												comment.id
											)
										}
									>
										{comment.likes}
									</p>
									<img
										onClick={() =>
											handleLikeComment(
												"dislike",
												comment.id
											)
										}
										src="../img/thumbs-down.svg"
									/>
									<p
										onClick={() =>
											handleLikeComment(
												"dislike",
												comment.id
											)
										}
									>
										{comment.dislikes}
									</p>
								</div>
							</div>
						))}
				</div>
			))}
		</div>
	)
}
