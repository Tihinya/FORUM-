import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../core/framework"
// import ErrorWindow from "../error-window/err or-window.jsx"
import { registrationRequest } from "../../additional-funcitons/authorization"

export default function Registation() {
	const navigate = useNavigate()
	// const [errorArr, setErrorArr] = useState([])

	const [formData, setFormData] = useState({
		email: "",
		username: "",
		password: "",
		password_confirmation: "",
	})

	const handleInputChange = (e) => {
		const { name, value } = e.target
		setFormData((prevData) => ({
			...prevData,
			[name]: value,
		}))
	}
	const handleSubmit = async (e) => {
		e.preventDefault()
		console.log(formData)

		try {
			const resultInJson = await registrationRequest(formData)

			if (resultInJson.status === "success") {
				console.log("Registration successful:", resultInJson.message)
			} else if (resultInJson.status === "error") {
				console.error("Registration error:", resultInJson.message)
			}
		} catch (error) {
			console.error("Error during registration:", error)
		}
	}

	return (
		<div className="main__block">
			<div className="sign-up__block">
				<div className="small_part-sign-up">
					<div className="small_part_content">
						<p>Already have an account?</p>
						<h6>Sign in for better experience!</h6>
						<a
							className="sign__button"
							onClick={() => navigate("/login")}
						>
							Sign In
						</a>
						<h6>Or return to</h6>
						<a
							className="sign__button"
							onClick={() => navigate("/")}
						>
							Main Page
						</a>
					</div>
				</div>
				<div className="big_part">
					<div className="big_part_content">
						<p>Sign Up</p>
						<div className="auth">
							<img src="/src/img/git.svg" />
							<img src="/src/img/goggle.svg" />
						</div>
						<h3>Sign in for better experience!</h3>
						<form className="form" onSubmit={handleSubmit}>
							<div className="input-fields">
								<input
									name="email"
									className="input-design"
									type="text"
									placeholder="Email"
									value={formData.email}
									onChange={handleInputChange}
								/>
								<input
									name="username"
									className="input-design"
									type="text"
									placeholder="Nickname"
									value={formData.username}
									onChange={handleInputChange}
								/>
								<input
									name="password"
									className="input-design"
									placeholder="Password"
									type="text"
									value={formData.password}
									onChange={handleInputChange}
								/>
								<input
									name="passwordConfirmation"
									className="input-design"
									placeholder="Repeat Password"
									type="text"
									value={formData.password_confirmation}
									onChange={handleInputChange}
								/>
							</div>
							{/* <ErrorWindow errorArr={errorArr} /> */}
							<button
								className="sign__button"
								type="submit"
								onClick={() => navigate("/login")}
							>
								Sign Up
							</button>
						</form>
					</div>
				</div>
			</div>
		</div>
	)
}
