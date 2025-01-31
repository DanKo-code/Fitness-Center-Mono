import sad_doing_abonnements_card from "../../../images/sad_doing_abonnements_card.jpg";
import React, {useEffect, useState, useRef} from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import axios from "axios";
import AbonnementCard from "../MainAbonements/AbonementsCard/abonementCard";
import {setUser} from "../../../states/storeSlice/appStateSlice";
import {useDispatch, useSelector} from "react-redux";
import config from "../../../config";
import inMemoryJWT from "../../../services/inMemoryJWT";
import {Resource} from "../../../context/AuthContext";
import ShowErrorMessage from "../../../utils/showErrorMessage";
import ShowSuccessMessage from "../../../utils/showSuccessMessage";
import noAva from "../../../images/no_ava.png"



export default function MainProfile() {

    const dispatch = useDispatch();

    const [photoUrl, setPhotoUrl] = useState('');
    const [selectedFile, setSelectedFile] = useState(null)
    const fileInputRef = useRef(null);
    const [name, setName] = useState('');
    const [showAbonnementsList, setShowAbonnementsList] = useState(true);
    const [orders, setOrders] = useState([]);

    let user = useSelector((state) => state.userSliceMode.user);

    const handleFileChange = (e) => {
        const file = e.target.files[0]; // Получаем выбранный файл
        setSelectedFile(e.target.files[0])

        if (file && file.type.startsWith("image/")) {
            const reader = new FileReader();
            reader.onload = (event) => {
                setPhotoUrl(event.target.result); // Устанавливаем URL изображения
            };
            reader.readAsDataURL(file); // Читаем файл как DataURL
        } else {
            alert("Please select a valid image file!"); // Сообщение, если файл не фото
        }
    };

    const handleImageClick = () => {
        fileInputRef.current.click(); // Активируем скрытый input
    };

    // useEffect для установки данных пользователя
    useEffect(() => {

        if (user) {
            setName(user.name);
            setPhotoUrl(user.photo);

            Resource.get('/orders/'+user.id)
                .then(response => {
                    console.log('get.ordersByUser response.data: '+JSON.stringify(response.data.orders, null, 2))
                    setOrders(response.data.orders)
                })
                .catch(error => {
                    console.error('Failed to fetch orders:', error);
                });
        }
    }, [user]); // Выполнится только при изменении `user`

    const handleNameChange = (event) => {
        setName(event.target.value);
    }

    const handleSubmit = async (event) => {
        event.preventDefault();

        const formData =new FormData();
        formData.append('photo', selectedFile)
        formData.append('name', name)

        try {
            const response = await Resource.put('/users/'+user.id, formData)

            if (response.status === 200) {

                const newUser = response.data.user

                dispatch(setUser(newUser));
                ShowSuccessMessage('Клиент успешно обновлен')
            }
        } catch (error) {
            if(error?.response?.data?.error) {
                ShowErrorMessage(error?.response?.data?.error)
            } else {
                ShowErrorMessage("Can't update user")
            }

            console.error('Can\'t update user: ' + JSON.stringify(error, null, 2))
        }
    }

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
                marginRight: '5%',

                height: '600px',

                display: 'flex',
                justifyContent: 'space-between'
            }}>

                {/*Change User Data Form*/}
                <div style={{width: '40%', display: 'flex', flexDirection: 'column'}}>

                    <input
                        type="file"
                        accept="image/*" // Разрешаем только фото
                        ref={fileInputRef} // Привязываем ссылку
                        style={{display: "none"}} // Скрываем элемент
                        onChange={handleFileChange}
                    />

                    <div style={{marginBottom: '20px'}}>
                        <img
                            style={{
                                width: '100%',
                                height: "300px",
                                objectFit: "cover",
                                border: "1px solid #ccc",
                                borderRadius: "8px",
                                cursor: "pointer", // Курсор как у кнопки
                                boxShadow: "0px 4px 6px rgba(0, 0, 0, 0.1)",
                            }}
                            onClick={handleImageClick}
                            src={photoUrl || noAva}/>
                    </div>
                    <TextField style={{marginBottom: '20px'}}
                               fullWidth
                               label="Имя"
                               value={name}
                               onChange={handleNameChange}
                    />
                    <Button
                        style={{
                            color: 'white',
                            background: 'rgb(160, 147, 197)',
                            height: '50px'
                        }}

                        onClick={handleSubmit}
                    >
                        Изменить
                    </Button>
                </div>

                {/*Orders List*/}
                <div style={{width: '40%'}}>
                    <div style={{
                        display: 'flex',
                        justifyContent: 'center',
                        marginBottom: '20px',
                        fontSize: '24px'
                    }}>Приобретенные абоненты
                    </div>
                    <div style={{height: '550px', overflowY: 'scroll'}}>
                        {orders ? <div>
                            {orders
                                .sort((a, b) => new Date(b.orderObject.created_time) - new Date(a.orderObject.created_time))
                                .map(order => (
                                <AbonnementCard abonnement={
                                    {
                                        abonement: order.abonementObject,
                                        services: order.serviceObjects,
                                    }
                                } status={order.orderObject.status}/>
                            ))}
                        </div> : <div>Нет никакаих заказов</div>}
                    </div>
                </div>

            </div>
        </div>
    )
}