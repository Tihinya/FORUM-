export default function HandleInputChange(e) {
	const { name, value } = e.target
	setFormData((prevData) => ({
		...prevData,
		[name]: value,
	}))
}

export default function HandleSubmit(e) {
    e.preventDefault()
    console.log(formData)

    try {
        const resultInJson =  RegistrationRequest(formData)

        if (resultInJson.status === "success") {
            console.log("Registration successful:", resultInJson.message)
        } else if (resultInJson.status === "error") {
            console.error("Registration error:", resultInJson.message)
        }
    } catch (error) {
        console.error("Error during registration:", error)
    }
}

export async function authReturnHandler(
    r,
    { setErrorArr, navigate },
    isRegistration
) {
    if (r.status === 200) {
        try {
            const responseJson = await r.json();

            if (responseJson.errors && responseJson.errors.length !== 0) {
                const errArr = fetchErrorChecker(responseJson.errors, navigate);
                if (errArr) setErrorArr(errArr);
                return;
            }

            if (responseJson.data && responseJson.data.id) {
                localStorage.setItem("id", String(responseJson.data.id));
            }

            navigate(
                isRegistration
                    ? `/additional-registration`
                    : `/user/${localStorage.getItem("id")}`
            );
        } catch (error) {
            navigate(`/internal-error`);
        }
        return;
    }

    throw new Error();
}
