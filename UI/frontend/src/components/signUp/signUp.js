import React, {useState, useContext} from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import { useNavigate, Link } from 'react-router-dom';
import showErrorMessage from "../../utils/showErrorMessage";
import {AuthContext} from "../../context/AuthContext";
import ShowErrorMessage from "../../utils/showErrorMessage";
import {setUser} from "../../states/storeSlice/appStateSlice";
import {useDispatch} from "react-redux";

export default function SignUp() {

    const {handleSignUp} = useContext(AuthContext);

    const navigate = useNavigate();

    const dispatch = useDispatch();

    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const handleSubmit = async (event) => {


        event.preventDefault();

        const data = {
            name: name,
            email: email,
            password: password,
        }

        try {
             const response = await handleSignUp(data);

            if (response.status === 200) {
                dispatch(setUser(response.data.user));
                navigate('/signin');
            }
        } catch (error) {
            if(error?.response?.data?.errors) {

                const errorMessages = error.response.data.errors;

                const formattedErrors = Object.entries(errorMessages)
                    .map(([field, message]) => `${message}`)
                    .join('\n');

                ShowErrorMessage(formattedErrors)
            } else {
                ShowErrorMessage("Неверные данные для регистрации")
            }

            console.error('Incorrect SignUp Data: ' + JSON.stringify(error, null, 2))
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
                            Регистрация
                        </h2>
                    </div>
                </div>

                <div style={{
                    display: "flex",
                    gap: '2%',
                    marginBottom: '2%'
                }}>
                    <TextField
                        fullWidth
                        id="name"
                        label="Имя"
                        name="name"
                        autoComplete="family-name"
                        value={name}
                        onChange={(event) => {
                            setName(event.target.value)
                        }}
                    />
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
                    Регистрация
                </Button>
                <div style={{
                    display: 'flex',
                    justifyContent: "end"
                }}>
                    <Link to={'/signin'}>
                        У вас уже есть учетная запись? Вход
                    </Link>
                </div>
            </div>
        </div>

    );
}