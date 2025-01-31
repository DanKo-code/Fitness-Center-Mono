import React, {useEffect, useState} from 'react';
import {FormControl, InputLabel, MenuItem, Select, Slider} from "@mui/material";
import TextField from "@mui/material/TextField";
import axios from 'axios';
import HouseIcon from "@mui/icons-material/House";
import Button from "@mui/material/Button";

import AccessTimeIcon from '@mui/icons-material/AccessTime';
import CalendarMonthIcon from '@mui/icons-material/CalendarMonth';
import LocalAtmIcon from '@mui/icons-material/LocalAtm';

import sad_doing_abonnements_card from '../../../images/sad_doing_abonnements_card.jpg'
import AbonnementCard from "./AbonementsCard/abonementCard";
import showErrorMessage from "../../../utils/showErrorMessage";
import config from "../../../config";
import inMemoryJWT from "../../../services/inMemoryJWT";
import {Resource} from "../../../context/AuthContext";


export default function MainAbonnements() {

    const [age, setAge] = useState('');
    const [titleSearch, setTitleSearch] = useState('');
    const [abonnements, setAbonnements] = useState([]);
    const [searchedAbonnements, setSearchedAbonnements] = useState([]);
    const [showAbonnementsList, setShowAbonnementsList] = useState(true);
    const [maxAbonementPrice, setMaxAbonementPrice] = useState(0);
    const [minAbonementPrice, setMinAbonementPrice] = useState(0);


    useEffect(() => {

        Resource.get('/abonements')
            .then(response => {

                if (response.data.abonements?.abonementsWithServices?.length > 0) {

                    let abonements = response.data.abonements.abonementsWithServices

                    setAbonnements(abonements);
                    console.log('abonements: ' + JSON.stringify(abonements, null, 2))



                    const sortPrices = abonements.map(abonement => abonement.abonement.price).sort((a, b) => a - b);

                    console.log('sortPrices: ' + JSON.stringify(sortPrices, null, 2));

                    setMaxAbonementPrice(sortPrices[sortPrices.length - 1] + 100)
                    setMinAbonementPrice(sortPrices[0])

                    setPriceRange([0, sortPrices[sortPrices.length - 1] + 100])


                    setSearchedAbonnements(abonements);
                }
            })
            .catch(error => {
                showErrorMessage(error);
                console.error('Failed to fetch abonnements:', error);
            });
    }, []);

    const handleTitleSearchChange = (event) => {

        setTitleSearch(event.target.value);

        if (!event.target.value) {
            setSearchedAbonnements(abonnements.filter(abonnement => abonnement.abonement.price >= priceRange[0] && abonnement.abonement.price <= priceRange[1]));
        } else {
            setSearchedAbonnements(abonnements.filter(abonnement => abonnement.abonement.title.toLowerCase().includes(event.target.value.toLowerCase())));
        }
    };

    useEffect(() => {
        if (searchedAbonnements.length === 0) {
            setShowAbonnementsList(false);
        } else {
            setShowAbonnementsList(true);
        }
    }, [searchedAbonnements]);

    const [priceRange, setPriceRange] = useState([0, 0]);

    const handlePriceRangeChange = (event, newValue) => {
        setPriceRange(newValue);

        setSearchedAbonnements(abonnements.filter(abonnement => abonnement.abonement.price >= priceRange[0] && abonnement.abonement.price <= priceRange[1]));
    };


    return (
        <div style={{
            width: '70%',
            height: '100vh',
            background: 'rgba(117,100,163,255)',
            display: "flex",
            justifyContent: "center",
            alignItems: "center"

        }}>
            <div style={{
                marginLeft: '5%',
                marginRight: '5%'
            }}>
                <div style={{display: 'flex', justifyContent: 'space-between'}}>

                    <div style={{width: '40%'}}>
                        <div style={{textAlign: 'center'}}>Price:</div>
                        <Slider
                            aria-labelledby="range-slider"
                            value={priceRange}
                            onChange={handlePriceRangeChange}
                            valueLabelDisplay="auto"
                            min={0}
                            max={maxAbonementPrice}
                        />
                    </div>

                    <TextField
                        style={{width: '40%'}}
                        fullWidth
                        id="search"
                        label="Name Search"
                        name="search"
                        autoComplete="search"
                        value={titleSearch}
                        onChange={handleTitleSearchChange}
                    />
                </div>
                <div style={{display: 'flex', justifyContent: 'center'}}>
                    {showAbonnementsList ? <div style={{marginTop: '40px', height: '400px', overflowY: 'scroll'}}>
                        {searchedAbonnements.map(abonnement => (
                            <AbonnementCard abonnement={abonnement} width={'600px'} height={'400px'}
                                            buyButton={{buttonState: true}}/>
                        ))}
                    </div> : <div style={{
                        marginTop: '40px',
                        width: "600px",
                        height: "400px",
                        display: "flex",
                        justifyContent: "center",
                        alignItems: "center",
                        fontSize: "20px"
                    }}>There are no such abonnements</div>
                    }
                </div>
            </div>
        </div>
    );
}