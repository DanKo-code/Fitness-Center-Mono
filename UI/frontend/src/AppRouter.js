import {Routes, Route, Navigate} from 'react-router-dom';

import SignIn from './components/signIn/signIn';
import MainNavHome from './components/MainNavHome/MainNavHome'
import SignUp from "./components/signUp/signUp";
import MainNavAbonnements from './components/MainNavAbonnements/MainNavAbonnements';
import MainNavCoaches from './components/MainNavCoaches/MainNavCoaches';
import MainNavProfile from "./components/MainNavProfile/MainNavProfile";
import {useContext} from "react";
import {AuthContext} from "./context/AuthContext";
import CoachDetailsCard from "./components/Main/MainCoaches/CoachesDetailedCard/coachDetailCard";
import MainNavCoachDetailCard from "./components/MainNavAbonementDetail/MainNavAbonementDetail";
import MainNavAdminPanel from "./components/AdminPanel/AdminPanel";

export function AppRouter() {
    const {isUserLogged} = useContext(AuthContext);

    return (
        <Routes>
            {!isUserLogged ? (
                <>
                    <Route
                        path='/signin'
                        element={<SignIn/>}
                    />
                    <Route
                        path='/signup'
                        element={<SignUp/>}
                    />
                </>
            ) : (
                <>
                    <Route
                        path='/main/home'
                        element={<MainNavHome/>}
                    />
                    <Route
                        path='/main/abonnements'
                        element={<MainNavAbonnements/>}
                    />

                    <Route
                        path='/main/profile'
                        element={<MainNavProfile/>}
                    />

                    <Route
                        path='/main/coaches'
                        element={<MainNavCoaches/>}
                    />

                    <Route
                        path='/main/coaches/details'
                        element={<MainNavCoachDetailCard/>}
                    />

                    <Route
                        path='/adminPanel'
                        element={<MainNavAdminPanel/>}
                    />
                </>
            )}

            <Route path="*" element={<Navigate to={isUserLogged ? "/main/home" : "/signin"}/>}></Route>

        </Routes>
    );
}