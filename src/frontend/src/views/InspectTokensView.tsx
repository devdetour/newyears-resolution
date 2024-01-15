import { ReactElement, useState } from "react";
import { CookieStore } from "../util/CookieStore";
import { Constants } from "../constants";

const ALL_TOKENS_PATH = "/auth/token/all";

function InspectTokensView (): ReactElement {
    const [token, setToken] = useState([])
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get("code") ?? "";
    const scope = urlParams.get("scope") ?? "";

    const jwt = CookieStore.getCookie(Constants.JWTCookieName)

    async function fetchTokens() {
        try {
            const response = await fetch(ALL_TOKENS_PATH, {
              method: "GET",
              headers: {
                "Content-Type": "application/json",
                  "Authorization": `Bearer ${jwt}`
              },
            });
      
            if (!response.ok) {
              throw new Error(`Error: ${response.statusText}`);
            }
      
            const responseData = await response.json();
            console.log("Response data:", responseData);
            // Set JWT in context TODO probably remove this once cookie works
            // Set JWT as a same-site cookie
            setToken(responseData.data)
          } catch (error: any) {
            console.error("Error posting data:", error.message);
            // setErr(error)
          }
    }

    return (
    <div>
        <h1>All Tokens for User</h1>
        <p>Code: {code}</p>
        <p>Scope: {scope}</p>
        <p>Tokens: {JSON.stringify(token)} </p>
        <button onClick={() => {fetchTokens()}}>Session Inspect</button>
    </div>
    )
}

export default InspectTokensView