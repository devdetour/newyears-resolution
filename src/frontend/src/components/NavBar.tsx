import { AppBar, MenuItem, Typography } from "@mui/material"
import { ReactElement } from "react"
import { Link } from "react-router-dom"

type NavBarProps = {
    routes: any[] // TODO real type for this.. it is RouteObject + name
}

function NavBar(props: NavBarProps): ReactElement {
    return (
        <AppBar position="static">
            {props.routes.map((route) => (
                route.name != undefined
                ? <Link to={route.path}>
                    <MenuItem key={route.name} onClick={() => {} }>
                        <Typography textAlign="center">{route.name}</Typography>
                    </MenuItem>
                  </Link>
                : null
            ))}
        </AppBar>
    )
}


export default NavBar