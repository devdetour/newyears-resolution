import {
    type RouteObject,
    Outlet,
} from "react-router-dom";
import HomeView from "./views/HomeView";
import LoginView from "./views/LoginView";
import View404 from "./views/View404";
import RegisterView from "./views/RegisterView";
import DataSourceView from "./views/DataSourceView";
import { AppBar, MenuItem, Typography } from "@mui/material";
import ExternalAuthReceiver from "./views/ExternalAuthReceiver";
import NavBar from "./components/NavBar";
import InspectTokensView from "./views/InspectTokensView";
import ContractsView from "./views/ContractsView";
import DataView from "./views/DataView";
import ContractsHistory from "./views/ContractsHistory";

const routes = [
    {
        path: "/",
        name: "Home",
        element: <HomeView />
    },
    {
        path: "/login",
        name: "Login",
        element: <LoginView />
    },
    {
        path: "/register",
        name: "Register",
        element: <RegisterView />
    },
    {
        path: "/data",
        name: "Data View",
        element: <DataView />
    },
    {
        path: "/contracts_history",
        name: "Contracts History",
        element: <ContractsHistory />
    },
    {
        path: "/link_datasource",
        name: "Link Datasource",
        element: <DataSourceView />
    },
    {
        path: "/receive_token",
        element: <ExternalAuthReceiver />
    },
    {
        path: "/inspect_tokens",
        name: "Inspect Tokens",
        element: <InspectTokensView />
    },
    {
        path: "/contracts",
        name: "Contracts",
        element: <ContractsView />
    },
    {
        path: "*",
        element: <View404 />
    }
]

const routeConfig = (): RouteObject[] => [
    {
        element: (
            <>
            <NavBar routes={routes}/>
            <Outlet/>
            </>
        ), children: routes
    }
]

export default routeConfig