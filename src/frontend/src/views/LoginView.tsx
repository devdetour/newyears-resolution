import React, { useState, type ReactElement, useContext } from 'react'
import { Box, Button, FormControl, TextField } from '@mui/material'
import LoginForm from '../components/LoginForm'
import { Link, useNavigate } from 'react-router-dom'
import { JwtContext } from '../JwtContext'
import { Constants } from '../constants'
import { CookieStore } from '../util/CookieStore'

const LOGIN_PATH = "/auth/login"

function LoginView (): ReactElement {
  const [fields, _] = useState({
    identity: '',
    password: ''
  })

  const navigate = useNavigate();

  const [err, setErr] = useState(null as Error | null)
  const [msg, setMsg] = useState("");

  // Call login route of API then set JWT
  async function submit(formState: any) {
    console.log("trying to login!")
    console.log(formState)
    setErr(null)
    setMsg("")

    // TODO validate
    const body = formState

    try {
      const response = await fetch(LOGIN_PATH, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          // Add any additional headers if needed
        },
        body: JSON.stringify(body),
      });

      if (!response.ok) {
        throw new Error(`Error: ${response.statusText}`);
      }

      const responseData = await response.json();
      console.log("Response data:", responseData);
      // Set JWT in context TODO probably remove this once cookie works
      // setJwt(responseData.data)
      // Set JWT as a same-site cookie
      CookieStore.setCookie(Constants.JWTCookieName, responseData.data, 10)
      setMsg("Logged in! Redirecting to datasource view...")
      setTimeout(() => { navigate("/link_datasource") }, 3000)

    } catch (error: any) {
      console.error("Error posting data:", error.message);
      setErr(error)
      // Handle the error (e.g., show an error message)
    }
  }

  return (
    <div>
      <h1>Login</h1>
      {/* TODO actual good react here? */}
      { err != null ? `Failed to login! ${err}` :
        msg.length > 0 ? `Login succeeded! ${msg}` :
        null }
      {LoginForm(fields, submit)}
      <Link to="/register">Register</Link>
    </div>
  )
}

export default LoginView