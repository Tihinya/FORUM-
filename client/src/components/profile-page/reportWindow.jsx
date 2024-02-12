import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"

import { fetchData } from "../../additional-funcitons/api.js"
export default function ReportWindow({ endPointUrl }) {
	const { setErrorMessage } = useContext("currentErrorMessage")

	const { activeSubj } = useContext("currentCategory")
	const [reports, setReports] = useState([])

	function getReportPost() {
		fetchData(null, endPointUrl, "GET").then((resultInJson) => {
			setReports(resultInJson)
		})
	}

	useEffect(() => {
		getReportPost()
	}, [activeSubj])

	if (reports == null) {
		return <h1 style={"text-align: center"}>Reports not found</h1>
	}

	function reportAnswerButton(formData) {
		fetchData(formData, "postreport/update", "PATCH").then(
			(resultInJson) => {
				if (resultInJson.status === "success") {
					setReports([])
					getReportPost()
				} else if (resultInJson.status === "error") {
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
								Person{" "}
								<span
									className="highlight"
									style={"color: var(--button-color-orange)"}
								>
									{report.user_name}
								</span>{" "}
								reported post{" "}
								<span
									className="highlight"
									style={"color: var(--button-color-orange)"}
								>
									{report.title}
								</span>
							</h3>
						</div>
						<div className="promotion__buttons">
							<div
								className="promotion-button accept"
								onClick={() => {
									const formData = {
										report_id: report.report_id,
										status: "approved",
									}
									reportAnswerButton(formData)
								}}
							>
								Approve report
							</div>
							<div
								className="promotion-button decline"
								onClick={() => {
									const formData = {
										report_id: report.report_id,
										status: "rejected",
									}
									reportAnswerButton(formData)
								}}
							>
								Decline report
							</div>
						</div>
					</div>
				))}
			</div>
		</div>
	)
}
