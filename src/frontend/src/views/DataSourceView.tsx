import React, { useState, type ReactElement } from 'react'
import LinkDataSourceForm from '../components/LinkDataSourceForm'

const CLIENT_ID = "REPLACE_"

function DataSourceView (): ReactElement {
  return (
    <div>
      <h1>Data Sources</h1>
      <LinkDataSourceForm name="Strava" link={`http://www.strava.com/oauth/authorize?client_id=${CLIENT_ID}&response_type=code&redirect_uri=http://localhost:3000/receive_token&approval_prompt=force&scope=read,activity:read&state=strava`}></LinkDataSourceForm>
    </div>
  )
}

export default DataSourceView