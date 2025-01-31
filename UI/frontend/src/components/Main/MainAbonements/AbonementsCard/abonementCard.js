import AccessTimeIcon from "@mui/icons-material/AccessTime";
import CalendarMonthIcon from "@mui/icons-material/CalendarMonth";
import LocalAtmIcon from "@mui/icons-material/LocalAtm";
import sad_doing_abonnements_card from "../../../../images/sad_doing_abonnements_card.jpg";
import Button from "@mui/material/Button";
import React, {useEffect} from "react";
import axios from 'axios';
import {useDispatch, useSelector} from "react-redux";
import config from "../../../../config";
import inMemoryJWT from "../../../../services/inMemoryJWT";
import {Resource} from "../../../../context/AuthContext";
import io from 'socket.io-client';
import {setUser} from "../../../../states/storeSlice/appStateSlice";
import showSuccessMessage from "../../../../utils/showSuccessMessage";
import ShowErrorMessage from "../../../../utils/showErrorMessage";
import { useStripe } from "@stripe/react-stripe-js";
import noAva from "../../../../images/no_ava.png"

export default function AbonnementCard(props) {

    const stripe = useStripe();


    const dispatch = useDispatch();
    let user = useSelector((state) => state.userSliceMode.user);

    const {abonnement, width, height, buyButton, onClick, status} = props;

    const handleBuy = async (event) => {
        try {

            let postOrdersData = {
                "ClientId": user.id,
                "AbonementId": abonnement.abonement.id,
                "StripePriceId": abonnement.abonement.stripe_price_id,
            }

            const response = await Resource.post('/create-checkout-session', postOrdersData);

            if (response.status === 200) {

                const result = await stripe.redirectToCheckout({
                    sessionId: response.data.sessionUrl,
                });

                if (result.error) {
                    console.error(result.error.message);
                }

                console.log("CreatedOrder: "+ JSON.stringify(response.data.order, null, 2))
            }
        } catch (e) {
            ShowErrorMessage(e);
            console.error('err: ' + JSON.stringify(e, null, 2))
        }
    }

    const styles = {
        overlay: {
            position: 'absolute',
            top: 0,
            left: 0,
            width: '100%',
            height: '100%',
            backgroundColor: 'rgba(0, 0, 0, 0.5)',
            color: 'white',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 1,
            fontSize: '24px',
        },


    }

    return (
        <div onClick={onClick} style={{
            background: 'rgb(160, 147, 197)',
            width: width,
            height: height,
            marginBottom: '10px',
            position: "relative"
        }}>
            {status === "Expired" && <div style={styles.overlay}>Expired</div>}

            <div style={{display: 'flex', justifyContent: 'center'}}>
                <div style={{marginTop: '5px', fontSize: '24px'}}>
                    {abonnement.abonement.title}
                </div>
            </div>

            <div style={{display: 'flex', marginTop: '20px', justifyContent: 'space-between'}}>
                <div style={{paddingLeft: '20px'}}>
                    <div style={{
                        marginBottom: '20px',
                        display: 'flex',
                        justifyContent: 'start',
                        alignItems: 'center'
                    }}>
                        <AccessTimeIcon fontSize="large" style={{marginRight: '10px'}}/>
                        <div style={{fontSize: '14px'}}>Периол валидности(в месяцах):</div>
                        <div style={{marginLeft: '5px', marginRight: '5px', fontSize: '14px'}}>
                            {abonnement.abonement.validity}
                        </div>
{/*
                        <div style={{fontSize: '14px'}}>months</div>
*/}
                    </div>

                    <div style={{
                        marginBottom: '20px',
                        display: 'flex',
                        justifyContent: 'start',
                        alignItems: 'center'
                    }}>
                        <CalendarMonthIcon fontSize="large" style={{marginRight: '10px'}}/>
                        <div style={{fontSize: '14px'}}>Время посещения:</div>
                        <div style={{fontSize: '14px', marginLeft: '5px', marginRight: '5px'}}>
                            {abonnement.abonement.visiting_time === "Any Time" ? "Любое Время" : abonnement.abonement.visiting_time}
                        </div>

                    </div>

                    <div style={{
                        marginBottom: '20px',
                        display: 'flex',
                        justifyContent: 'start',
                        alignItems: 'center'
                    }}>
                        <LocalAtmIcon fontSize="large" style={{marginRight: '10px'}}/>
                        <div style={{fontSize: '14px'}}>Цена(BYN):</div>
                        <div style={{
                            fontSize: '14px',
                            marginLeft: '5px',
                            marginRight: '5px'
                        }}>{abonnement.abonement.price}</div>
                    </div>
                </div>

                <div style={{width: '200px', height: '210px', paddingRight: '20px'}}>
                    {/*{abonnement.Photo}*/}
                    <img style={{width: '100%', height: 'auto'}} src={abonnement.abonement.photo}/>
                </div>
            </div>

            <div style={{paddingLeft: '20px'}}>
                Услуги:
            </div>
            <div style={{
                paddingLeft: '20px',
                paddingRight: '40px',
                marginTop: '20px',
                display: 'flex',
                justifyContent: 'space-between'
            }}>
                <div style={{display: 'flex',}}>
                    {abonnement?.services?.length > 0 &&
                        abonnement.services.map(Service => (
                        <div style={{marginRight: '10px'}}>
                            <div style={{width: '80px', height: '60px'}}>
                                <img style={{width: '100%', height: 'auto'}}
                                     src={Service.photo || noAva}/>
                            </div>
                            <div>{
                                Service.title === "swimming-pool" ? "Бассейн" :
                                Service.title === "sauna" ? "Сауна" :
                                Service.title === "gym" ? "Тренажерный зал" :
                                Service.title
                            }</div>
                        </div>
                    ))}
                </div>

                {buyButton && user.role !== 'admin' ? <Button
                    style={{
                        color: 'white',
                        background: 'rgba(117,100,163,255)',
                        width: '170px',
                        height: '50px'
                    }}

                    onClick={handleBuy}
                >
                    Купить
                </Button> : <div/>}

            </div>
        </div>
    )


}