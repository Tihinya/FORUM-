import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"
export default function Categories({ post }) {
	if (post.categories === undefined) {
		return
	}

	return (
		<div className="post__tags">
			{post.categories.map((categories) => (
				<p className="tag">{categories}</p>
			))}
		</div>
	)
}
