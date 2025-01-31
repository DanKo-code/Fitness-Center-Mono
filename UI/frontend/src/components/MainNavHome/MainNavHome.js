import React, {useState} from 'react';
import 'react-multi-carousel/lib/styles.css';
import MainNav from "../Main/MainNav/mainNav";
import MainHome from "../Main/MainHome/mainHome";
export default function MainNavHome() {

    return (
        <div style={{display: 'flex', height: '100vh', overflow: 'hidden', color: 'white'}}>
            <MainNav/>
            <MainHome/>
        </div>

    );
}