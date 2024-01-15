import React, { useState, type ReactElement, useEffect } from 'react'
import { Box, Button, FormControl, InputLabel, MenuItem, Paper, Select, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, TextField } from '@mui/material'
import { CookieStore } from '../util/CookieStore'
import { Constants } from '../constants'
import { Unstable_NumberInput as NumberInput } from '@mui/base';

const jwt = CookieStore.getCookie(Constants.JWTCookieName)

const CONTRACTS_CREATE_PATH = "/api/contracts/create"
const CONTRACTS_GET_PATH = "/api/contracts/get"
const CONTRACTS_DELETE_PATH = "/api/contracts/delete"

const contractTypes = {
    RECURRING: "recurring"
}

// TODO call these subtypes? a bit confusing
let goalTypes = [
  "Distance",
  "Time"
]

const goalCategories = ["strava"]

const units = ['hours']; //, 'days'];

// POST contract to server
async function PostContract(contractType: string, schedule: string, goalCategory: string, goalType: number, goal: number, lookback: number, lookbackUnit: string) {
  const body = {
    type: contractType,
    schedule,
    goalCategory,
    goalType,
    goal: typeof goal === "string" ? parseFloat(goal) : goal, // convert from string to number
    lookback: typeof lookback === "string" ? parseFloat(lookback) : lookback, // convert from string to number
    lookbackUnit 
  }

  try {
      const response = await fetch(CONTRACTS_CREATE_PATH, {
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
      
      // Set JWT in context TODO probably remove this once cookie works
      // Set JWT as a same-site cookie
    } catch (error: any) {
      console.error("Error posting data:", error.message);
      // setErr(error)
  }
}

async function GetContracts() {
  try {
    const response = await fetch(CONTRACTS_GET_PATH, {
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

async function DeleteContract(contractId: number) {
  const body = {
    contractId
  }

  try {
      const response = await fetch(CONTRACTS_DELETE_PATH, {
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
      
      // Set JWT in context TODO probably remove this once cookie works
      // Set JWT as a same-site cookie
    } catch (error: any) {
      console.error("Error deleting contract:", error.message);
      // setErr(error)
  }
}

function formatNanosecondsToHours(nanoseconds: number): number {
    const nanosecondsPerSecond = 1e9;
    const nanosecondsPerMinute = nanosecondsPerSecond * 60;
    const nanosecondsPerHour = nanosecondsPerMinute * 60;
  
    const hours = Math.floor(nanoseconds / nanosecondsPerHour);
    const minutes = Math.floor((nanoseconds % nanosecondsPerHour) / nanosecondsPerMinute);
    const seconds = Math.floor((nanoseconds % nanosecondsPerMinute) / nanosecondsPerSecond);
    const remainingNanoseconds = nanoseconds % nanosecondsPerSecond;
  
    return hours;
}




function ContractsView (): ReactElement {
    const [scheduleValue, setSchedulerValue] = React.useState("* * * * *")
    const [goalCategory, setGoalCategory] = React.useState(goalCategories[0])
    const [goalType, setGoalType] = React.useState(0) // INDEX, not value of the array.
    const [goal, setGoal] = React.useState(0)
    const [lookback, setLookback] = React.useState(0)
    const [unit, setUnit] = React.useState(units[0])

    const [contracts, setContracts] = React.useState([])

    useEffect(() => {
      let fetchContracts = async () => {
        const data = await GetContracts()
        console.log(data)
        setContracts(data.data)
      }
      fetchContracts()
    }, [])

    // TODO refactor so handleChange and handleSelectChange are same fn
    const handleChange = (event: any, setter: Function) => {
      setter(event.target.value)
    }

    const contractTypeList = [
        contractTypes.RECURRING
    ]

  return (
    <div>
      <h1>Create Recurring Contract</h1>

      {contractTypeList.map(elt => {
          return <>
          <TableContainer component={Paper}>
              <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Contract ID</TableCell>
                  <TableCell>Contract Type</TableCell>
                  <TableCell>Schedule</TableCell>
                  <TableCell>Lookback Hours</TableCell>
                  <TableCell>Goal Type</TableCell>
                  <TableCell>Goal Threshold</TableCell>
                  <TableCell>Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                  {contracts != undefined && contracts.length > 0 ? contracts.map((item: any) => (
                  <TableRow key={item.ID}>
                      <TableCell>{item.ID}</TableCell> 
                      <TableCell>{item.EvaluationFunctionName}</TableCell> 
                      <TableCell>{item.EvaluationSchedule}</TableCell>
                      <TableCell>{formatNanosecondsToHours(item.EvaluationLookback)}</TableCell>
                      <TableCell>{item.Goal.GoalType}</TableCell>
                      <TableCell>{item.Goal.GoalThreshold}</TableCell>
                      <TableCell><button onClick={() => DeleteContract(item.ID)}>Delete</button></TableCell>
                  </TableRow>
                  )) : "LOADING"}
              </TableBody>
              </Table>
          </TableContainer>


          <form>
            <p>Choose a category of goal.</p>
            <TextField
              value={goalCategory}
              label="Goal Category"
              onChange={e => handleChange(e, setGoalCategory)}
              select
              hidden={true}
            >
              {goalCategories.map(v => {
                return <MenuItem value={v}>{v}</MenuItem>
              })}
            </TextField>


            <p>Goal Type: choose your goal type (distance or time).</p>
            <TextField
              value={goalType}
              label="Goal Type"
              onChange={e => handleChange(e, setGoalType)}
              select
            >
              {goalTypes.map((val, idx) => {
                return <MenuItem value={idx}>{val}</MenuItem>
              })}
            </TextField>

            <br />

            <p>Goal Magnitude ({ goalType === 0 ? "Distance in meters" : `Time in minutes`})</p>
            <NumberInput
              aria-label="Goal Magnitude (distance or time)"
              placeholder="Goal"
              value={goal}
              min={0}
              onChange={e => handleChange(e, setGoal)}
            />

            <br />

            <p>Schedule Evaluation Interval (CRON syntax)</p>
            <p>How often the contract will be evaluated.</p>
            <TextField 
              label="Contract schedule"
              value={scheduleValue}
              onChange={e => handleChange(e, setSchedulerValue)}
            />

            <br />
            <p>Evaluation Lookback Interval ({unit})</p>
            <p>How many {unit} to look back when evaluating your goal. Use to set goal for e.g. exercising once a WEEK/once a DAY, etc.</p>
            <NumberInput
              aria-label="Lookback Time"
              placeholder="Lookback"
              value={lookback}
              min={0}
              onChange={e => handleChange(e, setLookback)}
            />

            {/* only show units for time goals */}
            <FormControl>
              <InputLabel id="unit-label">Unit</InputLabel>
              <Select
                labelId="unit-label"
                id="unit-select"
                value={unit}
                onChange={e => handleChange(e, setUnit)}
              >
              {units.map((unitOption) => (
                <MenuItem key={unitOption} value={unitOption}>
                  {unitOption.charAt(0).toUpperCase() + unitOption.slice(1)}
                </MenuItem>
              ))}
              </Select>
            </FormControl>

            <br />

            <Button variant="contained" onClick={() => { PostContract(elt, scheduleValue, goalCategory, goalType, goal, lookback, unit) } }>Create Contract</Button>
            </form>
          </>
      })}
    </div>
  )
}

export default ContractsView