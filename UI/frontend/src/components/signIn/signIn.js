import React, {useContext, useState} from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import axios from 'axios';
import {useDispatch, useSelector} from "react-redux";
import {setAppState} from '../../states/storeSlice/appStateSlice'
import {setUser} from '../../states/storeSlice/appStateSlice'
import { useNavigate, Link } from 'react-router-dom';
import inMemoryJWT from "../../services/inMemoryJWT";
import showErrorMessage from "../../utils/showErrorMessage";
import {AuthContext} from "../../context/AuthContext";
import ShowErrorMessage from "../../utils/showErrorMessage";


export default function SignIn() {
    const {handleSignIn} = useContext(AuthContext);
    const navigate = useNavigate();

    const dispatch = useDispatch();

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const handleSubmit = async (event) => {


        event.preventDefault();

        const data = {
            email: email,
            password: password,
        }

        try {
            const response = await handleSignIn(data);

            if (response.status === 200) {
                dispatch(setUser(response.data.user));
                navigate('/main');
            }
        } catch (error) {
            if(error?.response?.data?.errors) {

                const errorMessages = error.response.data.errors;

                const formattedErrors = Object.entries(errorMessages)
                    .map(([field, message]) => `${message}`)
                    .join('\n');

                ShowErrorMessage(formattedErrors)
            } else {
                ShowErrorMessage("Неверные данные для входа")
            }

            console.error('Incorrect SignIn Data: ' + JSON.stringify(error, null, 2))
        }
    };

    return (
        <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh',
        }}>
            <div style={{
                width: '40%',
                padding: '20px',
            }}>
                <div style={{display: 'flex', justifyContent: "center"}}>
                    <div style={{
                        display: 'flex',
                        justifyContent: 'center',
                        flexDirection: "column"
                    }}>
                        <div style={{display: 'flex', justifyContent: 'center'}}>
                            <Avatar sx={{m: 1, bgcolor: 'secondary.main'}}>
                                <LockOutlinedIcon/>
                            </Avatar>
                        </div>

                        <h2>
                            Вход
                        </h2>
                    </div>
                </div>

                <div style={{
                    display: "flex",
                    gap: '2%',
                    marginBottom: '2%'
                }}>

                </div>

                <TextField style={{
                    marginBottom: '2%'
                }}
                           fullWidth
                           id="email"
                           label="Email Адрес"
                           name="email"
                           autoComplete="email"
                           value={email}
                           onChange={(event) => {
                               setEmail(event.target.value)
                           }}
                />
                <TextField style={{
                    marginBottom: '2%'
                }}
                           fullWidth
                           name="password"
                           label="Пароль"
                           type="password"
                           id="password"
                           autoComplete="new-password"
                           value={password}
                           onChange={(event) => {
                               setPassword(event.target.value)
                           }}
                />
                <Button style={{
                    marginBottom: '2%'
                }}
                        onClick={handleSubmit}
                        fullWidth
                        variant="contained"
                >
                    Вход
                </Button>
                <div style={{
                    display: 'flex',
                    justifyContent: "end"
                }}>
                    <Link

                          to={'/signup'}

                          >
                        У вас нет учетной записи? Зарегистрироваться
                    </Link>
                </div>
            </div>
        </div>

    );
}