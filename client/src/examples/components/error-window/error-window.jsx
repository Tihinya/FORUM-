import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../core/framework"

export default function ErrorWindow({ errorArr }) {
	return (
		<div className="error-window">
			{errorArr.length > 0 && (
				<ul className="error-window__list">
					{errorArr.map((v, i) => (
						<li className="error-window__list-item" key={i}>
							{v.message}
						</li>
					))}
				</ul>
			)}
		</div>
	)
}
