import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import state from './states/store'
import {Provider} from "react-redux";
import {BrowserRouter} from 'react-router-dom';
import {AppRouter} from './AppRouter'
import {SnackbarProvider} from "notistack";
import AuthProvider from "./context/AuthContext";
import { Elements } from "@stripe/react-stripe-js";
import { loadStripe } from "@stripe/stripe-js";

const stripePromise = loadStripe('pk_test_51PxOL0A75DCPwyUvRKsHO8bzbcpXZj7bkXSHOEVeCUrbFGj7WrvG7OPmaKMoFHPeXDX9OLNdLOHgUkLgGBK4fhpA00kmxf34Ke');

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <Elements stripe={stripePromise}>
    <Provider store={state}>
    <AuthProvider>
        <React.StrictMode>
                <SnackbarProvider/>
                <BrowserRouter>
                    <AppRouter/>
                </BrowserRouter>
        </React.StrictMode>
    </AuthProvider>
    </Provider>
    </Elements>

);
