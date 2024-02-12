import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"

import messageSquare from "../../img/message-square.svg"

export default function CommentsIcon({ post }) {
	const navigate = useNavigate()
	return (
		<div className="post__info">
			<a
				onClick={() => {
					navigate(`/comments-authorized/${post.id}`)
				}}
			>
				<img src={messageSquare} />
			</a>
			<p
				onClick={() => {
					navigate(`/comments-authorized/${post.id}`)
				}}
			>
				{post.comment_count}
			</p>
		</div>
	)
}
