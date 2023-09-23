import Gachi, { useEffect, useState } from "../../../core/framework.ts"
import { convertTime } from "../../additional-funcitons/post.js"

var navigatePostId = 0

export function sendPostId(postId) {
	navigatePostId = postId
}

export function Comment() {
	const [posts, setPosts] = useState([])
	const [comments, setComments] = useState([])

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
			.catch((error) => console.error("Error fetching posts:", error))
	}

	useEffect(() => {
		fetchPost()
		fetchComments()
	}, [])

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
									<img src="../img/thumbs-up.svg" />
									<p>{post.likes}</p>
									<img src="../img/thumbs-down.svg" />
									<p>{post.dislikes}</p>
								</div>
							</div>
						</div>
					</div>
					<p>Comments:</p>
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
							</div>
						))}
				</div>
			))}
		</div>
	)
}
