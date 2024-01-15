import React, { useState, type ReactElement, useEffect } from 'react'
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@mui/material'
import { CookieStore } from '../util/CookieStore'
import { Constants } from '../constants'

const jwt = CookieStore.getCookie(Constants.JWTCookieName)

const CONTRACTS_HISTORY_PATH = "/api/contracts/history"

async function GetContractHistory() {
  try {
    const response = await fetch(CONTRACTS_HISTORY_PATH, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${jwt}`
      },
    })
    const responseData = await response.json();
    console.log("Response data:", responseData);
    return responseData
    
    // Set JWT in context TODO probably remove this once cookie works
    // Set JWT as a same-site cookie
  } catch (error: any) {
    console.error("Error getting data:", error.message);
    return []
  }
}

function ContractsHistory (): ReactElement {
    const [contracts, setContracts] = React.useState([])

    useEffect(() => {
      let fetchContracts = async () => {
        const data = await GetContractHistory()
        console.log(data)
        setContracts(data.data)
      }
      fetchContracts()
    }, [])

  return (
    <div>
      <h1>Contract Evaluation History</h1>

      <TableContainer component={Paper}>
          <Table>
          <TableHead>
            <TableRow>
              <TableCell>Contract ID</TableCell>
              <TableCell>Evaluation Time</TableCell>
              <TableCell>Threshold Met?</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
              {contracts != undefined && contracts.length > 0 ? contracts.map((item: any) => (
              <TableRow key={item.ID}>
                  <TableCell>{item.ID}</TableCell>
                  <TableCell>{item.EvaluationTime}</TableCell>
                  <TableCell>{"" + item.ThresholdMet}</TableCell>
              </TableRow>
              )) : "LOADING"}
          </TableBody>
          </Table>
      </TableContainer>
    </div>
  )
}

export default ContractsHistory