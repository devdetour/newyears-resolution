import React, { useState, type ReactElement } from 'react'
import { Box, Button, FormControl, TextField } from '@mui/material'

// TODO not any
function LoginForm (fields: any, submitCb: Function): ReactElement {
  const [formState, setformState] = useState(fields)

    const handleChange = (e: { target: { id: any; value: any } }) => {
        // console.log(e.target)
        console.log(e.target.id)
        setformState({
            ...formState,
            [e.target.id]: e.target.value,
        })
    }

    const entryFields = (fields: any) => {
        return (<div>
            {Object.keys(fields).map((key) => {
                return <TextField
                required
                id={key}
                label={key}
                value={fields.key}
                onChange={handleChange}
                />
            })}
        </div>)
    }

  return (
    <div className="container">
        <Box component="section" sx={{p: 2, border: '1px dashed grey'}}>
            <FormControl variant="standard">
                {entryFields(fields)}
            <Button variant="contained" onClick={() => submitCb(formState)}>SUBMIT</Button>
            </FormControl>
        </Box>
    </div>
  )
}

export default LoginForm