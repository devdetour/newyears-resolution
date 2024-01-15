import React, { useState, type ReactElement, useEffect } from 'react'
import { Box, Button, FormControl, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, TextField } from '@mui/material'
import { CookieStore } from '../util/CookieStore'
import { Constants } from '../constants'

const TOKEN_PATH = "/api/data/strava"

const jwt = CookieStore.getCookie(Constants.JWTCookieName)

async function getStravaData() {
    try {
        const response = await fetch(TOKEN_PATH, {
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
        return responseData.data
        
        // Set JWT in context TODO probably remove this once cookie works
        // Set JWT as a same-site cookie
      } catch (error: any) {
        console.error("Error posting data:", error.message);
        return []
        // setErr(error)
    }
}

function DataView (): ReactElement {
    const [data, setData] = useState([] as any)
    useEffect(() => {
        getStravaData().then(result => setData(result))
    }, [])

  return (
    <>
    <h1>Recent Activities from Strava</h1>
    <TableContainer component={Paper}>
        <Table>
        <TableHead>
          <TableRow>
            <TableCell>Name</TableCell>
            <TableCell>Distance</TableCell>
            <TableCell>Moving Time</TableCell>
            <TableCell>Elapsed Time</TableCell>
            <TableCell>Type</TableCell>
            <TableCell>Start Date</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
            {data != undefined && data.length > 0 ? data.map((item: any) => (
            <TableRow key={item.id}>
                <TableCell>{item.name}</TableCell>
                <TableCell>{item.distance}</TableCell>
                <TableCell>{item.moving_time}</TableCell>
                <TableCell>{item.elapsed_time}</TableCell>
                <TableCell>{item.type}</TableCell>
                <TableCell>{item.start_date}</TableCell>
            </TableRow>
            )) : "LOADING"}
        </TableBody>
        </Table>
    </TableContainer>
    </>
  )
}

export default DataView