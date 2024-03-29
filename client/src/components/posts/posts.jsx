import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"

import { convertTime } from "../../additional-funcitons/post.js"
import { fetchData } from "../../additional-funcitons/api.js"
import LikesAndDislikes from "../post-likes/post-likes"
import CommentsIcon from "../comments/comment-icon.jsx"
import Categories from "./categories.jsx"
import ContextMenu from "../context-menu/context-menu"

import profilePicture from "../../img/avatarka.jpeg"

export default function Posts({ endPointUrl, userId }) {
	if (endPointUrl === "") {
		return <h1 style={"text-align: center"}>Posts not found</h1>
	}

	const { isAuthenticated: isLoggin } = useContext("isAuthenticated")
	const { posts, setPosts } = useContext("currentPosts")
	const { activeSubj } = useContext("currentCategory")
	const { comments, setComments } = useContext("currentComments")
	const navigate = useNavigate()
	const postOrComment = endPointUrl !== "comments"
	const endpoint =
		userId !== ""
			? `${endPointUrl}/${userId}`
			: `${activeSubj}` !== ""
			? "posts?categories=" + activeSubj
			: `${endPointUrl}`
	const badEndpoints = ["promotions", "categories", "demotions"]

	if (badEndpoints.includes(endPointUrl)) {
		return
	}

	useEffect(() => {
		fetchData(null, endpoint, "GET").then((resultInJson) => {
			if (
				endPointUrl === "posts" ||
				endPointUrl === "user/posts" ||
				endPointUrl === "post"
			) {
				setPosts(resultInJson)
			} else {
				setComments(resultInJson)
			}
		})
	}, [activeSubj])

	const data =
		endPointUrl === "posts" ||
		endPointUrl === "user/posts" ||
		endPointUrl === "post"
			? posts
			: comments

	if (Array.isArray(data)) {
		data.sort((a, b) => {
			const dateA = new Date(a.creation_date)
			const dateB = new Date(b.creation_date)
			return dateB - dateA
		})
	}

	if (!data.length) {
		return <h1 style={"text-align: center"}>Posts not found</h1>
	}

	return (
		<div>
			<div className="post__container">
				{data.map((post) => (
					<div className="post__box">
						<div className="post__header">
							<div className="user__info">
								<div className="user__info_picture">
									<a
										onClick={() => {
											if (isLoggin) {
												navigate("/profile-page")
											}
										}}
									>
										<img src={profilePicture} />
									</a>
								</div>
								<div className="user__info_name">
									<div className="name">
										{post.user_info?.username}
									</div>
									<div className="date">
										{convertTime(post.creation_date)}
									</div>
								</div>
							</div>
							<ContextMenu obj={post} endpoint={endPointUrl} />
						</div>
						<div className="post__content">
							<h3>{post.title}</h3>
							<p className="post-text">{post.content}</p>
							{post.image && (
								<div className="post__image-container">
									<img
										className="post__image"
										src={post.image}
									/>
								</div>
							)}
						</div>
						<div className="post__info">
							<div className="post__tags">
								{postOrComment ? (
									<Categories post={post} />
								) : (
									""
								)}
							</div>
							<div className="post__likes">
								{endPointUrl === "comments" ||
								endPointUrl === "post" ? (
									""
								) : (
									<CommentsIcon post={post} />
								)}

								<LikesAndDislikes
									post={post}
									endPointUrl={endpoint}
								/>
							</div>
						</div>
					</div>
				))}
			</div>
		</div>
	)
}
