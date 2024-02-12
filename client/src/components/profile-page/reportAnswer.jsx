import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"

import { fetchData } from "../../additional-funcitons/api.js"
export default function ReportAnswer({ endPointUrl }) {
	const { setErrorMessage } = useContext("currentErrorMessage")

	const { activeSubj } = useContext("currentCategory")
	const [reports, setReports] = useState([])

	function getReportPostAnswer() {
		fetchData(null, endPointUrl, "GET").then((resultInJson) => {
			setReports(resultInJson)
		})
	}

	useEffect(() => {
		getReportPostAnswer()
	}, [activeSubj])

	if (reports == null) {
		return <h1 style={"text-align: center"}>Not found Reports</h1>
	}

	function reportAnswerButton(formData) {
		fetchData(formData, "postreport/answer/update", "PATCH").then(
			(resultInJson) => {
				if (resultInJson.status === "success") {
					setReports([])
					getReportPostAnswer()
				}
				if (resultInJson.status === "error") {
					setErrorMessage(resultInJson.message)
				}
			}
		)
	}

	return (
		<div>
			<div className="post__container">
				{reports.map((report) => (
					<div className="post__box">
						<div
							className="post__report_info"
							style={"font-family: 'Roboto', sans-serif"}
						>
							<h3>
								Admin said this post{" "}
								<span
									className="highlight"
									style={"color: var(--button-color-orange)"}
								>
									{report.title}
								</span>{" "}
								is{" "}
								<span
									className="highlight"
									style={"color: var(--button-color-orange)"}
								>
									{report.message}
								</span>
								, now live with that!
							</h3>
						</div>
						<div className="promotion__buttons">
							<div
								className="promotion-button accept"
								onClick={() => {
									const formData = {
										report_id: report.report_id,
										seen: true,
									}
									reportAnswerButton(formData)
								}}
								style={
									"font-family: 'Roboto', sans-serif; font-weight: bold"
								}
							>
								OK!
							</div>
						</div>
					</div>
				))}
			</div>
		</div>
	)
}
