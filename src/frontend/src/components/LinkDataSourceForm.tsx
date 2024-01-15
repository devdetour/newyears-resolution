import React, { useState, type ReactElement, useContext } from 'react'
import { Box, Button, FormControl, TextField } from '@mui/material'
import { JwtContext } from '../JwtContext';

interface LinkDataSourceFormProps {
    name: string
    link: string
}

function LinkDataSourceForm (props: LinkDataSourceFormProps): ReactElement {
    const jwtCtx = useContext(JwtContext);
    return (
    <div className="container">
        {/* {JSON.stringify(jwtCtx)} */}
        <Box component="section" sx={{p: 2, border: '1px dashed grey'}}>
            <h1>Link to {props.name}</h1>
            <p>Click the below link to authenticate with {props.name}.</p>
            <a href={props.link}>Authenticate with {props.name}</a>
        </Box>
    </div>
    )
}

export default LinkDataSourceForm