import React, { useState, type ReactElement, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Constants } from '../constants'
import { CookieStore } from '../util/CookieStore'

const TOKEN_PATH = "/auth/token/create"
const SESSION_PATH = "/session"

// POST tokens to server
async function storeToken(code: string, scope: string, jwt: string): Promise<boolean> {
// Call login route of API then set JWT
    console.log("trying to store token!")

    const body = {
        text: code,
        scope
    }
    try {
      const response = await fetch(TOKEN_PATH, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${jwt}`
        },
        body: JSON.stringify(body),
      });

      if (!response.ok) {
        throw new Error(`Error: ${response.statusText}`);
      }

      const responseData = await response.json();
      console.log("Response data:", responseData);
      // Handle the response data as needed
      return Promise.resolve(true)
    } catch (error: any) {
      console.error("Error posting data:", error.message);
      // Handle the error (e.g., show an error message)
      return Promise.resolve(false)
    }
}

function ExternalAuthReceiver (): ReactElement {
    const [msg, setMsg] = useState("")
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get("code") ?? "";
    const scope = urlParams.get("scope") ?? "";
    const navigate = useNavigate();
    const cookieVal = CookieStore.getCookie(Constants.JWTCookieName)

    useEffect(() => {
      const fetchData = async () => {
        let result = await storeToken(code, scope, cookieVal)
        if (result) {
          setMsg("Stored token! Redirecting to dashboard...")
          navigate("/data")
        } else {
          setMsg("Failed to store token! That sucks huh")
        }
      }
      fetchData();
    }, [])

    return (
    <div>
        <h1>Receiving token...</h1>
        { msg.length > 0 ? <p>{msg}</p> : null}
    </div>
    )
}

export default ExternalAuthReceiver