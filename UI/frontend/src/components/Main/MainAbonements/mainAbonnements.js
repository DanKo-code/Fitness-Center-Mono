import React, {useEffect, useState} from 'react';
import {FormControl, InputLabel, ListItemText, MenuItem, OutlinedInput, Select, Slider} from "@mui/material";
import TextField from "@mui/material/TextField";
import AbonnementCard from "./AbonementsCard/abonementCard";
import showErrorMessage from "../../../utils/showErrorMessage";
import {Resource} from "../../../context/AuthContext";
import {useSelector} from "react-redux";
import {Checkbox} from "@mui/joy";

const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
    PaperProps: {
        style: {
            maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
            width: 250,
        },
    },
};

export default function MainAbonnements() {

    const [age, setAge] = useState('');
    const [titleSearch, setTitleSearch] = useState('');
    const [validityPeriodSearch, setValidityPeriodSearch] = useState('Any');
    const [abonnements, setAbonnements] = useState([]);
    const [searchedAbonnements, setSearchedAbonnements] = useState([]);
    const [showAbonnementsList, setShowAbonnementsList] = useState(true);
    const [maxAbonementPrice, setMaxAbonementPrice] = useState(0);
    const [minAbonementPrice, setMinAbonementPrice] = useState(0);
    const [priceRange, setPriceRange] = useState([0, 0]);
    const [sortOrder, setSortOrder] = useState("asc");
    const [sortFilter, setSortFilter] = useState("title");
    const [orders, setOrders] = useState([]);
    const [visitingTime, setVisitingTime] = useState('Any');
    const [allServices, setAllServices] = useState([]);
    const [currentServices, setCurrentServices] = useState([]);

    let user = useSelector((state) => state.userSliceMode.user);

    useEffect(() => {
        Resource.get('/abonements')
            .then(response => {
                if (response.data.abonements?.abonementsWithServices?.length > 0) {
                    let abonements = response.data.abonements.abonementsWithServices;

                    setAbonnements(abonements);

                    const sortPrices = abonements.map(abonement => abonement.abonement.price).sort((a, b) => a - b);

                    console.log('sortPrices: ' + JSON.stringify(sortPrices, null, 2));

                    setMaxAbonementPrice(sortPrices[sortPrices.length - 1])
                    setMinAbonementPrice(sortPrices[0])

                    setPriceRange([sortPrices[0], sortPrices[sortPrices.length - 1]])


                    filterData(abonnements, '', 'Any', 'Any', sortPrices[0], [], sortPrices[sortPrices.length - 1])
                    setSearchedAbonnements(sortAbonnements(abonements, sortFilter, sortOrder));
                }

                return Resource.get('/orders/'+user.id)
            })
            .then(response => {
                setOrders(response.data.orders)

                return Resource.get('/services')
            })
            .then(response => {
                if(response.data.services.serviceObject.length > 0){
                    setAllServices(response.data.services.serviceObject);
                }
            })
            .catch(error => {
                showErrorMessage(error);
                console.error('Failed to fetch data:', error);
            });
    }, []);

    const handleTitleSearchChange = (event) => {
        setTitleSearch(event.target.value);

        const filtered = filterData(abonnements, event.target.value, validityPeriodSearch, visitingTime, currentServices, priceRange[0], priceRange[1])

        setSearchedAbonnements(sortAbonnements(filtered, sortFilter, sortOrder));
    };

    const handlePriceRangeChange = (event, newValue) => {
        setPriceRange(newValue);

        const filtered = filterData(abonnements, titleSearch, validityPeriodSearch, visitingTime, currentServices, newValue[0], newValue[1])

        setSearchedAbonnements(sortAbonnements(filtered, sortFilter, sortOrder));
    };

    const handleSortOrderChange = (event) => {
        setSortOrder(event.target.value);
        const filtered = filterData(abonnements, titleSearch, validityPeriodSearch, visitingTime, currentServices, priceRange[0], priceRange[1])

        const sorted = sortAbonnements(filtered, sortFilter, event.target.value);

        console.log("sorted: "+JSON.stringify(sorted, null, 2))

        setSearchedAbonnements(sorted);
    };

    const sortAbonnements = (abonnements, sortingFilter, order) => {
        return [...abonnements].sort((a, b) => {
            if (sortingFilter === "price") {
                return order === "asc"
                    ? a.abonement.price - b.abonement.price
                    : b.abonement.price - a.abonement.price;
            }

            if (sortingFilter === "title") {
                const titleA = a.abonement.title.toLowerCase();
                const titleB = b.abonement.title.toLowerCase();

                if (order === "asc") {
                    return titleA.localeCompare(titleB);
                } else {
                    return titleB.localeCompare(titleA);
                }
            }

            return 0;
        });
    };

    function filterData(data, searchName, validityPeriod, visitingTime, currentServices, minPrice, maxPrice) {
        return data.filter(item => {
            const matchesName = searchName
                ? item.abonement.title.toLowerCase().includes(searchName.toLowerCase())
                : true;

            let matchesValidityPeriod = validityPeriod
                ? item.abonement.validity === validityPeriod
                : true;
            if (validityPeriod === 'Any') {
                matchesValidityPeriod = true
            }

            let matchesVisitingTime = visitingTime
                ? item.abonement.visiting_time === visitingTime
                : true
            if (visitingTime === 'Any') {
                matchesVisitingTime = true
            }

            const matchesPrice = (minPrice === undefined || item.abonement.price >= minPrice) &&
                (maxPrice === undefined || item.abonement.price <= maxPrice);

            const containsAllValues = currentServices.map(column => column.title).every(value => item.services.map(column => column.title).includes(value));

            return matchesName && matchesPrice && matchesValidityPeriod && matchesVisitingTime && containsAllValues;
        });
    }

    useEffect(() => {
        if (searchedAbonnements.length === 0) {
            setShowAbonnementsList(false);
        } else {
            setShowAbonnementsList(true);
        }
    }, [searchedAbonnements]);

    function handleValidityPeriodChange(event) {
        const value = event.target.value;

        setValidityPeriodSearch(value);

        const filtered = filterData(abonnements, titleSearch, value, visitingTime, currentServices, priceRange[0], priceRange[1])
        setSearchedAbonnements(sortAbonnements(filtered, sortFilter, sortOrder));
    }

    function handleSortFilterChange(event) {
        setSortFilter(event.target.value);
        const filtered = filterData(abonnements, titleSearch, validityPeriodSearch, visitingTime, currentServices, priceRange[0], priceRange[1])
        setSearchedAbonnements(sortAbonnements(filtered, event.target.value, sortOrder));
    }

    const handleVisitingTimeChange = async (event) => {
        const value = event.target.value;

        setVisitingTime(value);

        const filtered = filterData(abonnements, titleSearch, validityPeriodSearch, value, currentServices, priceRange[0], priceRange[1])
        setSearchedAbonnements(sortAbonnements(filtered, sortFilter, sortOrder));
    }

    const handleServicesChange = async (service) => {
        if (currentServices.map(serviceObg => serviceObg.id).includes(service.id)) {
            const updatedServices = currentServices.filter(serviceObj => serviceObj.id !== service.id);
            setCurrentServices(updatedServices);

            const filtered = filterData(abonnements, titleSearch, validityPeriodSearch, visitingTime, updatedServices, priceRange[0], priceRange[1])
            setSearchedAbonnements(sortAbonnements(filtered, sortFilter, sortOrder));
        } else {
            setCurrentServices([service, ...currentServices]);

            const filtered = filterData(abonnements, titleSearch, validityPeriodSearch, visitingTime, [service, ...currentServices], priceRange[0], priceRange[1])
            setSearchedAbonnements(sortAbonnements(filtered, sortFilter, sortOrder));
        }
    }

    const translations = {
        "swimming-pool": "Бассейн",
        "sauna": "Сауна",
        "gym": "Тренажерный зал"
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
                <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center'}}>
                    <TextField
                        style={{width: '30%'}}
                        fullWidth
                        id="search"
                        label="Поиск по имени"
                        name="search"
                        autoComplete="search"
                        value={titleSearch}
                        onChange={handleTitleSearchChange}>
                    </TextField>

                    <div style={{width: '30%'}}>
                        <div style={{textAlign: 'center'}}>Цена:</div>
                        <Slider
                            aria-labelledby="range-slider"
                            value={priceRange}
                            onChange={handlePriceRangeChange}
                            valueLabelDisplay="auto"
                            min={minAbonementPrice}
                            max={maxAbonementPrice}
                        />
                    </div>

                    <FormControl style={{width: '30%'}}>
                        <InputLabel id="sort-order-label">Фильтр сортировки</InputLabel>
                        <Select
                            labelId="sort-order-label"
                            id="sort-order"
                            value={sortFilter}
                            onChange={handleSortFilterChange}
                        >
                            <MenuItem value="title">Название</MenuItem>
                            <MenuItem value="price">Цена</MenuItem>
                        </Select>
                    </FormControl>
                </div>

                <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: '20px'}}>
                    <FormControl fullWidth style={{width: '33%'}}>
                        <InputLabel id="validity-period-label">Период валидности в месяцах</InputLabel>
                        <Select
                            id="validity-period-label"
                            value={validityPeriodSearch}
                            label="Validity Period"
                            onChange={handleValidityPeriodChange}
                        >
                            <MenuItem value={'1'}>1</MenuItem>
                            <MenuItem value={'3'}>3</MenuItem>
                            <MenuItem value={'6'}>6</MenuItem>
                            <MenuItem value={'12'}>12</MenuItem>
                            <MenuItem value={'Any'}>Любой</MenuItem>
                        </Select>
                    </FormControl>

                    <FormControl fullWidth style={{width: '30%'}}>
                        <InputLabel>Время посещения</InputLabel>
                        <Select
                            value={visitingTime}
                            label="Visiting Time"
                            onChange={handleVisitingTimeChange}
                        >
                            <MenuItem value={'7.00 - 14.00'}>7.00 - 14.00</MenuItem>
                            <MenuItem value={'14.00 - 24.00'}>14.00 - 24.00</MenuItem>
                            <MenuItem value={'Any Time'}>Любое время</MenuItem>
                            <MenuItem value={'Any'}>Любое</MenuItem>
                        </Select>
                    </FormControl>

                    <FormControl style={{width: '30%'}}>
                        <InputLabel id="sort-order-label">Порядок сортировки</InputLabel>
                        <Select
                            labelId="sort-order-label"
                            id="sort-order"
                            value={sortOrder}
                            onChange={handleSortOrderChange}
                        >
                            <MenuItem value="asc">Восходящий</MenuItem>
                            <MenuItem value="desc">Нисходящий</MenuItem>
                        </Select>
                    </FormControl>
                </div>

                <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: '20px'}}>
                    <FormControl fullWidth>
                        <InputLabel fullWidth>Услуги</InputLabel>
                        <Select
                            fullWidth
                            multiple
                            value={currentServices}
                            input={<OutlinedInput label="Tag"/>}
                            renderValue={(selected) =>
                                selected.map(sel => (translations[sel.title] || sel.title) + ' ')
                            }
                            MenuProps={MenuProps}
                        >
                            {allServices.map((service) => (
                                <MenuItem key={service.id} value={service.title}>
                                    <Checkbox onChange={() => handleServicesChange(service)}
                                              checked={currentServices.map(ser => ser.id).includes(service.id)}/>
                                    <ListItemText
                                        primary={
                                            service.title === "swimming-pool" ? "Бассейн" :
                                                service.title === "sauna" ? "Сауна" :
                                                    service.title === "gym" ? "Тренажерный зал" :
                                                        service.title
                                        }
                                    />
                                </MenuItem>
                            ))}
                        </Select>
                    </FormControl>
                </div>

                <div style={{display: 'flex', justifyContent: 'center'}}>
                    {showAbonnementsList ? (
                        <div style={{marginTop: '40px', height: '400px', overflowY: 'scroll'}}>
                            {searchedAbonnements
                                .map(abonnement => (
                                <AbonnementCard key={abonnement?.abonement?.id} abonnement={abonnement} width={'600px'} height={'400px'}
                                                buyButton={!orders?.some(order => abonnement?.abonement?.id === order?.orderObject?.abonement_id && order?.orderObject?.status === "Valid")}/>
                            ))}
                        </div>
                    ) : (
                        <div style={{
                            marginTop: '40px',
                            width: "600px",
                            height: "400px",
                            display: "flex",
                            justifyContent: "center",
                            alignItems: "center",
                            fontSize: "20px"
                        }}>
                            Таких абонентов не существует
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
