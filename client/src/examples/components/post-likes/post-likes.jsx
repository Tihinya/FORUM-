import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

import { fetchData } from "../../additional-funcitons/api.js"

export default function LikesAndDislikes({ post, endPointUrl }) {
	const postLikeUrl = "user/liked"
	const postDisLikeUrl = "user/disliked"
	const commentLikeUrl = "user/likedcomments"
	const commentDislikeUrl = "user/dislikedcomments"
	
	const [contextType, setContextType] = useState(undefined)
	const { setPosts } = useContext("currentPosts")
	const { setErrorMessage } = useContext("currentErrorMessage")
	const { setComments } = useContext("currentComments")
	const [likedPosts, setLikedPosts] = useState([])
	const [dislikedPosts, setDislikedPosts] = useState([])
	const [likedComments, setLikedComments] = useState([])
	const [dislikedComments, setDislikedComments] = useState([])

	const isLoggin = useContext("isAuthenticated").isAuthenticated

	const fetchPostLikes = () => {
		fetchData(null, postLikeUrl, "GET").then((resultInJson) => {
			setLikedPosts(resultInJson)
		})
	}

	const fetchPostDislikes = () => {
		fetchData(null, postDisLikeUrl, "GET").then((resultInJson) => {
			setDislikedPosts(resultInJson)
		})
	}

	const fetchCommentLikes = () => {
		fetchData(null, commentLikeUrl, "GET").then((resultInJson) => {
			setLikedComments(resultInJson)
		})
	}

	const fetchCommentDislikes = () => {
		fetchData(null, commentDislikeUrl, "GET").then((resultInJson) => {
			setDislikedComments(resultInJson)
		})
	}

	if (post.hasOwnProperty("comment_count")) {
        setContextType("post")
    } else if (post.hasOwnProperty("post_id")) {
        setContextType("comment")
    }
	

	useEffect(() => {
		if (isLoggin) {
			fetchPostLikes()
			fetchPostDislikes()
			fetchCommentLikes()
			fetchCommentDislikes()
		}
	}, [])

	const handleLike = (type, postId) => {
		let endpoint

		const isLikingPosts = !likedPosts.some((post) => (post.id == postId)) && //true
		!dislikedPosts.includes(postId)

		const isLikingComments = !likedComments.some((id) => (id == postId)) && //true
		!dislikedComments.includes(postId)

		if (contextType == "post") {
			endpoint = isLikingPosts ? `${postId}/${type}` : `${postId}/un${type}`
		} else {
			endpoint = isLikingComments ? `${postId}/${type}` : `${postId}/un${type}`
		}

		fetchData(null, `${contextType}/${endpoint}`, "POST").then(
			(resultInJson) => {
				if (resultInJson.status === "success") {
					fetchPostLikes()
					fetchPostDislikes()
					fetchCommentLikes()
					fetchCommentDislikes()

					if (endPointUrl == "posts") {
						fetchData(null, `posts`, "GET").then(
							(resultInJson) => {
								setPosts(resultInJson)
							}
 						)
					} else if (contextType == "comment") {
						fetchData(null, `comments/${post.post_id}`, "GET").then(
							(resultInJson) => {
								setComments(resultInJson)
							}
						)
					} else {
							fetchData(null, `post/${post.id}`, "GET").then(
								(resultInJson) => {
									setPosts(resultInJson)
								}
							)
						}

				} else if (resultInJson.status === "error") {
					setErrorMessage(resultInJson.message)
				}
			}
		)
	}

	return (
		<div>
			{isLoggin ? (
				<div className="post__likes">
					<img
						onClick={() => handleLike("like", post.id)}
						src="../img/thumbs-up.svg"
					/>
					<p onClick={() => handleLike("like", post.id)}>
						{post.likes}
					</p>
					<img
						onClick={() => handleLike("dislike", post.id)}
						src="../img/thumbs-down.svg"
					/>
					<p onClick={() => handleLike("dislike", post.id)}>
						{post.dislikes}
					</p>
				</div>
			) : (
				<div className="post__likes">
					<img src="../img/thumbs-up.svg" />
					<p>{post.likes}</p>
					<img src="../img/thumbs-down.svg" />
					<p>{post.dislikes}</p>
				</div>
			)}
		</div>
	)
}
