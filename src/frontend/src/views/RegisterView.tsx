import React, { useState, type ReactElement } from 'react'
import { Box, Button, FormControl, TextField } from '@mui/material'
import LoginForm from '../components/LoginForm'
import { Link, useNavigate } from 'react-router-dom'


const LOGIN_PATH = "/auth/register"

function RegisterView (): ReactElement {
  const navigate = useNavigate();
  
  const [fields, setFields] = useState({
    username: '',
    email: '',
    password: ''
  })

  const [msg, setMsg] = useState("")
  const [err, setErr] = useState("")

  // Call login route of API
  async function submit(formState: any) {
    console.log("Submitting")
    console.log(formState)

    setMsg("")
    setErr("")

    // TODO validate this
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

      const responseData = await response.json();
      console.log("Response data:", responseData);

      if (!response.ok) {
        setErr(JSON.stringify(responseData))
      } else {
        setMsg("Registered! Redirecting to login view...")
        setTimeout(() => { navigate("/login") }, 3000)
      }
    } catch (error: any) {
      console.error("Error posting data:", error.message);
      setErr(error)
    }
  }

  return (
    <div>
      <h1>Register</h1>
      {err.length > 0 ? `An error occurred! ${err}` : msg.length > 0 ? msg : null}
      {LoginForm(fields, submit)}
      <Link to="/login">Login</Link>
    </div>
  )
}

export default RegisterView