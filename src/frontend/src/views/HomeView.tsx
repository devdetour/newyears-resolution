import React, { type ReactElement } from 'react'
import { Link } from 'react-router-dom'

function HomeView (): ReactElement {

  return (
        <div className="container">
            <h1>Home View</h1>
            <Link to="/login">Login</Link>
            <br></br>
            <Link to="/link_datasource">Link Datasource</Link>
        </div>
  )
}

export default HomeView