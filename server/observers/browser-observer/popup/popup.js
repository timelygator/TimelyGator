function relay_status(url, token) {
    path = url + "/status";
    res = "Not Connected";
    fetch(url, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
    })
        .then((response) => {
            if (response.status == 200) {
                res = "Connected";
            } else if (response.status == 401) {
                res = "Unauthorized";
            } else {
                res = "Not Connected";
            }
        })
        .catch((error) => {
            res = "Not Connected";
        });
    return res;
}

window.onload = (event) => {
    let url = document.getElementById("relay-url");
    let token = document.getElementById("relay-token");
    let label = document.getElementById("status");
    let submit = document.getElementById("submit");

    submit.addEventListener("click", () => {
        label.innerHTML = relay_status(url.value, token.value);
    });

    let val = relay_status(url.value, token.value);
    label.innerHTML = val;
};
