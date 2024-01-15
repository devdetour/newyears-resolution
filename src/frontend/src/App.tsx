import React, { useState } from 'react';
import logo from './logo.svg';
import './App.css';
import { Link, RouterProvider, createBrowserRouter, useNavigate } from 'react-router-dom';
import routeConfig from './Routes';
import { JwtContext } from './JwtContext';
import { AppBar, MenuItem, Typography } from '@mui/material';

const router = createBrowserRouter(routeConfig());

function App() {

  return (
    <div className="App">
        <RouterProvider router={router}></RouterProvider>
    </div>
  );
}

export default App;
