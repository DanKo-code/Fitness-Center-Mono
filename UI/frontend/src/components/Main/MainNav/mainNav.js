import {FitnessCenter} from "@mui/icons-material";
import Button from "@mui/material/Button";
import HouseIcon from "@mui/icons-material/House";
import {setAppState} from "../../../states/storeSlice/appStateSlice";
import CardMembershipIcon from "@mui/icons-material/CardMembership";
import GroupsIcon from "@mui/icons-material/Groups";
import ManageAccountsIcon from "@mui/icons-material/ManageAccounts";
import ExitToAppIcon from "@mui/icons-material/ExitToApp";
import React, {useContext} from "react";
import {useNavigate} from 'react-router-dom'
import axios from 'axios';
import config from "../../../config";
import inMemoryJWT from "../../../services/inMemoryJWT";
import showErrorMessage from "../../../utils/showErrorMessage";
import {AuthContext} from "../../../context/AuthContext";
import {useSelector} from "react-redux";

export const AuthClient = axios.create({
    baseURL: `${config.API_URL}/auth`,
    withCredentials: true,
})

export default function MainNav() {

    let user = useSelector((state) => state.userSliceMode.user);

    const {LogOut} = useContext(AuthContext);
    const navigate = useNavigate();

    const handleLogOut = async () =>{
        await LogOut();
    }

    return (
        <div style={{width: '30%', height: '100vh', background: 'rgba(160, 147, 197, 1)'}}>
            <div style={{marginLeft: '10%', marginRight: '10%'}}>
                <div style={{
                    display: 'flex',
                    justifyContent: "center",
                    alignItems: 'center',
                    marginTop: '50px',
                    gap: '10%'
                }}>
                    <div style={{fontSize: '30px'}}>FitLab</div>
                    <div style={{
                        width: '60px',
                        height: '60px',
                        background: 'rgba(117,100,163,255)',
                        borderRadius: '50%',
                        display: 'flex',
                        justifyContent: "center",
                        alignItems: 'center'
                    }}>
                        <FitnessCenter sx={{width: 50, height: 50}}/>
                    </div>

                </div>

                <div>
                    <div style={{display: 'flex', flexDirection: 'column'}}>
                        <Button style={{color: 'white', background: 'rgba(117,100,163,255)', marginTop: '5%'}}
                                startIcon={< HouseIcon sx={{width: 50, height: 50}}/>}
                                sx={{
                                    justifyContent: 'flex-start', // Выравнивание контента по левому краю
                                    paddingLeft: '25%', // Добавляем отступ слева для текста
                                }}

                                onClick={()=>{navigate('/main/home')}}
                                >
                            ГЛАВНАЯ
                        </Button>

                        <Button style={{color: 'white', background: 'rgba(117,100,163,255)', marginTop: '5%'}}
                                startIcon={< CardMembershipIcon sx={{width: 50, height: 50}}/>}
                                sx={{
                                    justifyContent: 'flex-start', // Выравнивание контента по левому краю
                                    paddingLeft: '25%', // Добавляем отступ слева для текста
                                }}

                                onClick={()=>{navigate('/main/abonnements')}}
                                >
                            АБОНЕМЕНТЫ
                        </Button>

                        <Button style={{color: 'white', background: 'rgba(117,100,163,255)', marginTop: '5%'}}
                                startIcon={< GroupsIcon sx={{width: 50, height: 50}}/>}
                                sx={{
                                    justifyContent: 'flex-start', // Выравнивание контента по левому краю
                                    paddingLeft: '25%', // Добавляем отступ слева для текста
                                }}

                                onClick={()=>{navigate('/main/coaches')}}
                                >
                            ТРЕНЕРА
                        </Button>

                        { user.role === "client" && (<Button style={{color: 'white', background: 'rgba(117,100,163,255)', marginTop: '5%'}}
                                startIcon={< ManageAccountsIcon sx={{width: 50, height: 50}}/>}
                                sx={{
                                    justifyContent: 'flex-start', // Выравнивание контента по левому краю
                                    paddingLeft: '25%', // Добавляем отступ слева для текста
                                }}

                                onClick={()=>{navigate('/main/profile')}}>

                            Профиль
                        </Button>)
                        }

                        { user.role === "admin" && (
                            <Button style={{color: 'white', background: 'rgba(117,100,163,255)', height: '62px', marginTop: '5%'}}
                                    onClick={()=>{navigate('/adminPanel')}}>

                                ПАНЕЛЬ АДМИНИСТРАТОРА
                            </Button>
                            )
                        }

                        <Button style={{color: 'white', background: 'rgba(117,100,163,255)', marginTop: '30%'}}
                                startIcon={< ExitToAppIcon sx={{width: 50, height: 50, transform: 'scaleX(-1)'}}/>}
                                sx={{
                                    justifyContent: 'flex-start', // Выравнивание контента по левому краю
                                    paddingLeft: '25%', // Добавляем отступ слева для текста
                                }}
                                onClick={handleLogOut}
                                >
                            ВЫХОД
                        </Button>
                    </div>

                </div>
            </div>

        </div>
    );
}