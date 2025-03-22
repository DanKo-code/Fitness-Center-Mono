import React, { useEffect, useRef, useState } from "react";
import { Resource } from "../../context/AuthContext";
import showErrorMessage from "../../utils/showErrorMessage";
import IconButton from "@mui/material/IconButton";
import CloseIcon from "@mui/icons-material/Close";
import AbonnementCard from "../Main/MainAbonements/AbonementsCard/abonementCard";
import sad_doing_abonnements_card from "../../images/sad_doing_abonnements_card.jpg";
import TextField from "@mui/material/TextField";
import {
  FormControl,
  InputAdornment,
  InputLabel,
  ListItemText,
  MenuItem,
  OutlinedInput,
  Select,
  TextareaAutosize,
} from "@mui/material";
import { Checkbox } from "@mui/joy";
import Button from "@mui/material/Button";
import CoachCard from "../Main/MainCoaches/CoachesCard/coachCard";
import ShowErrorMessage from "../../utils/showErrorMessage";
import ShowSuccessMessage from "../../utils/showSuccessMessage";
import noAva from "../../images/no_ava.png";

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

export default function CoachesModal({ onClose }) {
  const [coaches, setCoaches] = useState([]);
  const [currentCoach, setCurrentCoach] = useState("");

  const [photoUrl, setPhotoUrl] = useState("");
  const [selectedFile, setSelectedFile] = useState(null);
  const fileInputRef = useRef(null);

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");

  const [currentServices, setCurrentServices] = useState([]);
  const [currentShift, setCurrentShift] = useState("");
  const [allServices, setAllServices] = useState([]);

  const handleFileChange = (e) => {
    const file = e.target.files[0]; // Получаем выбранный файл
    setSelectedFile(e.target.files[0]);

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

  const handleServicesChange = async (service) => {
    if (
      currentServices.map((serviceObg) => serviceObg.id).includes(service.id)
    ) {
      const updatedServices = currentServices.filter(
        (serviceObj) => serviceObj.id !== service.id,
      );
      setCurrentServices(updatedServices);
    } else {
      setCurrentServices([service, ...currentServices]);
    }
  };

  useEffect(() => {
    Resource.get("/coaches")
      .then((response) => {
        if (
          response.data.coaches?.coachWithServicesWithReviewsWithUsers?.length >
          0
        ) {
          setCoaches(
            response.data.coaches.coachWithServicesWithReviewsWithUsers,
          );
        }
      })
      .catch((error) => {
        showErrorMessage(error);
        console.error("Failed to fetch coaches:", error);
      });

    Resource.get("/services")
      .then((response) => {
        if (response.data.services?.serviceObject?.length > 0) {
          setAllServices(response.data.services.serviceObject);
        }
      })
      .catch((error) => {
        showErrorMessage(error);
        console.error("Failed to fetch services:", error);
      });
  }, []);

  const handleCloseCoachesModal = () => {
    setName("");
    setDescription("");
    setCurrentServices([]);
    setPhotoUrl("");

    onClose();
  };

  const handleCoachSelection = async (selectedCoach) => {
    setSelectedFile(null);

    if (currentCoach?.coach?.id === selectedCoach.coach.id) {
      setCurrentCoach("");
    } else {
      setCurrentCoach(selectedCoach);
    }
  };

  const handleNameChange = async (event) => {
    setName(event.target.value);
  };
  const handleDescriptionChange = async (e) => {
    const inputValue = e.target.value;

    if (inputValue.length > 500) {
      const truncatedText = inputValue.slice(0, 500);
      setDescription(truncatedText);
    } else {
      setDescription(inputValue);
    }
  };

  //Current coach changing
  useEffect(() => {
    if (currentCoach) {
      setName(currentCoach.coach.name);
      setDescription(currentCoach.coach.description);
      setPhotoUrl(currentCoach.coach.photo);
      setCurrentServices(currentCoach.services.map((as) => as));
      setCurrentShift(currentCoach.coach.shift);
    } else {
      setName("");
      setDescription("");
      setCurrentServices([]);
      setPhotoUrl("");
      setCurrentShift("");
    }
  }, [currentCoach]);

  const handleDelete = async (event) => {
    event.preventDefault();

    const coachId = currentCoach.coach.id;

    const url = `/coaches/${coachId}`;

    try {
      if (currentCoach) {
        const response = await Resource.delete(url);

        if (response.status === 200) {
          setCoaches(
            coaches.filter((item) => item.coach.id !== response.data.coach.id),
          );

          setName("");
          setDescription("");
          setCurrentServices([]);
          setCurrentCoach("");
          setPhotoUrl("");

          ShowSuccessMessage("Тренер успешно удален");
        }
      }
    } catch (error) {
      ShowErrorMessage(error);
      console.error(
        "response.status: " +
          JSON.stringify(error.response.data.message, null, 2),
      );
    }
  };
  const handleSubmit = async (event) => {
    event.preventDefault();

    if (currentCoach) {
      //update coach
      try {
        const formData = new FormData();
        formData.append("id", currentCoach.coach.id);
        formData.append("name", name);
        formData.append("description", description);
        formData.append("photo", selectedFile);
        formData.append(
          "services",
          currentServices.map((service) => service.id),
        );
        formData.append("shift", currentShift);

        const response = await Resource.put("/coaches", formData);

        if (response.status === 200) {
          let coachWithServices = response.data.coach;

          setName(coachWithServices.coach.name);
          setDescription(coachWithServices.coach.description);
          setCurrentServices(coachWithServices.services);
          setCurrentCoach(coachWithServices);
          setCurrentShift(coachWithServices.coach.shift);

          const newArray = coaches.map((coach) => {
            if (coach.coach.id === coachWithServices.coach.id) {
              return coachWithServices;
            }
            return coach;
          });

          setCoaches(newArray);
          ShowSuccessMessage("Coach updated successfully");
        }
      } catch (error) {
        if (error?.response?.data?.error) {
          ShowErrorMessage(error?.response?.data?.error);
        } else {
          ShowErrorMessage("Can't update coach");
        }

        console.error("Can't update coach: " + JSON.stringify(error, null, 2));
      }
    } else {
      //create coach
      try {
        const formData = new FormData();
        formData.append("name", name);
        formData.append("description", description);
        formData.append("photo", selectedFile);
        formData.append(
          "services",
          currentServices.map((service) => service.id),
        );
        formData.append("shift", currentShift);

        const response = await Resource.post("/coaches", formData);

        if (response.status === 200) {
          let coachWithServices = response.data.coach;

          console.log("coachWithServices:", coachWithServices);

          setName(coachWithServices.coach.name);
          setDescription(coachWithServices.coach.description);
          setPhotoUrl(coachWithServices.coach.photo);
          setCoaches([coachWithServices, ...coaches]);
          setCurrentServices(coachWithServices.services);
          setCurrentCoach(coachWithServices);
          setCurrentShift(coachWithServices.coach.shift);

          ShowSuccessMessage("Coach created successfully");
        }
      } catch (error) {
        if (error?.response?.data?.error) {
          ShowErrorMessage(error?.response?.data?.error);
        } else {
          ShowErrorMessage("Can't create coach");
        }

        console.error("Can't create coach: " + JSON.stringify(error, null, 2));
      }
    }
  };

  const translations = {
    "swimming-pool": "Бассейн",
    tennis: "Сауна",
    gym: "Тренажерный зал",
  };

  return (
    <div
      style={{
        color: "white",
        position: "absolute",
        top: "50%",
        left: "50%",
        transform: "translate(-50%, -50%)",
        background: "rgba(117,100,163,255)",
        borderRadius: "8px",
        padding: "20px",
        width: "1000px",
        height: "620px",
      }}
    >
      <div>
        <IconButton
          onClick={handleCloseCoachesModal}
          size="large"
          sx={{ position: "absolute", top: 10, right: 10 }}
        >
          <CloseIcon />
        </IconButton>
      </div>

      <div style={{ display: "flex", marginTop: "10px", gap: "20px" }}>
        <div style={{ width: "40%" }}>
          <div
            style={{
              display: "flex",
              justifyContent: "center",
              marginBottom: "10px",
              fontSize: "18px",
            }}
          >
            Все тренера
          </div>
          <div style={{ height: "550px", overflowY: "scroll" }}>
            {coaches.length > 0 ? (
              <div>
                {coaches
                  .sort(
                    (a, b) =>
                      new Date(b.coach.updated_time) -
                      new Date(a.coach.updated_time),
                  )
                  .map((coach) => (
                    <div key={coach.coach.id}>
                      <div
                        style={{
                          border:
                            currentCoach === coach
                              ? "4px solid yellow"
                              : "none",
                          borderRadius: "20px",
                          marginBottom: "10px",
                        }}
                      >
                        {/*<AbonnementCard
                                            abonnement={abonnement}
                                            onClick={() => handleAbonementSelection(abonnement)}
                                        />*/}

                        <CoachCard
                          key={coach.coach.id}
                          coach={coach}
                          width={"370px"}
                          height={"310px"}
                          imageSize={"100px"}
                          button={false}
                          onClick={() => handleCoachSelection(coach)}
                        />
                      </div>
                    </div>
                  ))}
              </div>
            ) : (
              <div>Нет никаких тренеров</div>
            )}
          </div>
        </div>
        <div style={{ width: "60%", height: "100%" }}>
          <div
            style={{
              display: "flex",
              justifyContent: "center",
              marginBottom: "10px",
              fontSize: "18px",
            }}
          >
            {currentCoach ? "Обновить Тренера" : "Создать Тренера"}
          </div>

          <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
            <input
              type="file"
              accept="image/*" // Разрешаем только фото
              ref={fileInputRef} // Привязываем ссылку
              style={{ display: "none" }} // Скрываем элемент
              onChange={handleFileChange}
            />

            <div style={{ marginBottom: "20px" }}>
              {/*{abonnement.Photo}*/}

              <img
                style={{
                  width: "200px",
                  height: "150px",
                  objectFit: "cover",
                  border: "1px solid #ccc",
                  borderRadius: "8px",
                  cursor: "pointer", // Курсор как у кнопки
                  boxShadow: "0px 4px 6px rgba(0, 0, 0, 0.1)",
                }}
                onClick={handleImageClick}
                src={photoUrl || noAva}
              />

              {/*
                            <img style={{width: '200px', height: 'auto'}} src={sad_doing_abonnements_card}/>
*/}
            </div>

            <TextField
              style={{ marginBottom: "10px" }}
              fullWidth
              label="Имя"
              value={name}
              onChange={handleNameChange}
              sx={{
                input: { color: "white" }, // Цвет вводимого текста
                "& .MuiInputLabel-root": { color: "white" }, // Обычное состояние лейбла
                "& .MuiInputLabel-root.Mui-focused": { color: "white" }, // Фокусное состояние лейбла
              }}
            />
          </div>

          <TextareaAutosize
            style={{
              width: "100%",
              minHeight: "100px",
              resize: "none",
              backgroundColor: "transparent",
              borderRadius: "4px",
              padding: "8px",
              boxSizing: "border-box",
              marginBottom: "10px",
              border: "1px solid rgba(0, 0, 0, 0.23)", // Темная граница по умолчанию
              outline: "none",
              color: "white", // Белый текст при вводе
              caretColor: "white", // Белый курсор
            }}
            placeholder="Описание"
            value={description}
            onChange={handleDescriptionChange}
            maxRows={10}
          />

          <style>
            {`
    textarea::placeholder {
      color: white;
      font-size: 16px;
    }
  `}
          </style>

          <div style={{ marginBottom: "10px" }}>
            <FormControl fullWidth>
              <InputLabel
                fullWidth
                sx={{
                  color: "white !important", // Принудительно белый цвет
                  "&.Mui-focused": { color: "white !important" }, // При фокусе
                }}
              >
                Услуги
              </InputLabel>
              <Select
                fullWidth
                multiple
                value={currentServices}
                sx={{ color: "white" }}
                input={<OutlinedInput label="Tag" />}
                renderValue={(selected) =>
                  selected.map(
                    (sel) => (translations[sel.title] || sel.title) + " ",
                  )
                }
                MenuProps={MenuProps}
              >
                {allServices.map((service) => (
                  <MenuItem key={service.id} value={service.title}>
                    <Checkbox
                      onChange={() => handleServicesChange(service)}
                      checked={currentServices
                        .map((ser) => ser.id)
                        .includes(service.id)}
                    />
                    <ListItemText
                      primary={
                        service.title === "swimming-pool"
                          ? "Бассейн"
                          : service.title === "tennis"
                            ? "Теннис"
                            : service.title === "gym"
                              ? "Тренажерный зал"
                              : service.title
                      }
                    />
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </div>

          <FormControl fullWidth>
            <InputLabel
              fullWidth
              sx={{
                color: "white !important", // Принудительно белый цвет
                "&.Mui-focused": { color: "white !important" }, // При фокусе
              }}
            >
              Смена
            </InputLabel>
            <Select
              fullWidth
              sx={{ color: "white" }}
              value={currentShift}
              onChange={(e) => setCurrentShift(e.target.value)}
              input={<OutlinedInput label="Tag" />}
            >
              <MenuItem key={1} value={"1"}>
                <ListItemText primary={"1"} />
              </MenuItem>
              <MenuItem key={2} value={"2"}>
                <ListItemText primary={"2"} />
              </MenuItem>
            </Select>
          </FormControl>

          {currentCoach && (
            <div
              style={{
                display: "flex",
                justifyContent: "center",
                marginTop: "25px",
              }}
            >
              <Button
                style={{
                  color: "white",
                  background: "rgb(160, 147, 197)",
                  height: "50px",
                  width: "80%",
                }}
                onClick={handleDelete}
              >
                Удалить
              </Button>
            </div>
          )}

          <div
            style={{
              display: "flex",
              justifyContent: "center",
              marginTop: currentCoach ? "20px" : "100px",
            }}
          >
            <Button
              style={{
                color: "white",
                background: "rgb(160, 147, 197)",
                height: "50px",
                width: "80%",
              }}
              onClick={handleSubmit}
            >
              Подтвердить
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
