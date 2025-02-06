import sad_doing_coachs_card from "../../../../images/sad_doing_abonnements_card.jpg";
import Button from "@mui/material/Button";
import React, { useEffect, useState } from "react";
import axios from "axios";
import { useDispatch, useSelector } from "react-redux";
import sad_doing_abonnements_card from "../../../../images/sad_doing_abonnements_card.jpg";
import { useLocation } from "react-router-dom";
import danilaAvatar from "../../../../images/danila_avatar.jpg";
import Carousel from "react-multi-carousel";
import mainCss from "../../../MainNavHome/MainNavHome.module.css";
import AbonnementCard from "../../MainAbonements/AbonementsCard/abonementCard";
import { Modal, TextareaAutosize } from "@mui/material";
import { Resource, TrainingResource } from "../../../../context/AuthContext";
import showErrorMessage from "../../../../utils/showErrorMessage";
import ShowErrorMessage from "../../../../utils/showErrorMessage";
import ShowSuccessMessage from "../../../../utils/showSuccessMessage";
import noAva from "../../../../images/no_ava.png";

const initTrainingsState = [
  { time_from: "07:00", time_until: "08:00", status: "свободно" },
  { time_from: "08:00", time_until: "09:00", status: "свободно" },
  { time_from: "09:00", time_until: "10:00", status: "свободно" },
  { time_from: "10:00", time_until: "11:00", status: "свободно" },
  { time_from: "11:00", time_until: "12:00", status: "свободно" },
  { time_from: "12:00", time_until: "13:00", status: "свободно" },
  { time_from: "13:00", time_until: "14:00", status: "свободно" },
  { time_from: "14:00", time_until: "15:00", status: "свободно" },
  { time_from: "15:00", time_until: "16:00", status: "свободно" },
  { time_from: "16:00", time_until: "17:00", status: "свободно" },
  { time_from: "17:00", time_until: "18:00", status: "свободно" },
  { time_from: "18:00", time_until: "19:00", status: "свободно" },
  { time_from: "19:00", time_until: "20:00", status: "свободно" },
  { time_from: "20:00", time_until: "21:00", status: "свободно" },
  { time_from: "21:00", time_until: "22:00", status: "свободно" },
  { time_from: "22:00", time_until: "23:00", status: "свободно" },
];

const today = new Date().toISOString().split("T")[0];

export default function CoachDetailsCard(props) {
  const location = useLocation();
  const { coach } = location.state || {}; // Получаем переданный пропс

  const [openModal, setOpenModal] = useState(false);
  const [reviewText, setReviewText] = useState("");
  const [coachComments, setCoachComments] = useState([]);

  const [date, setDate] = useState("");

  const [dayTrainings, setDayTrainings] = useState("");

  let currentUser = useSelector((state) => state.userSliceMode.user);

  useEffect(() => {
    if (coach.reviewWithUser?.length > 0) {
      setCoachComments(
        coach.reviewWithUser.sort(
          (a, b) =>
            new Date(b.reviewObject.created_time) -
            new Date(a.reviewObject.created_time),
        ),
      );
    }

    //handleDateChange(today).then((r) => console.log("ok"));
  }, []);

  const handleOpenModal = () => {
    setOpenModal(true);
  };

  const handleCloseModal = () => {
    setOpenModal(false);
  };

  const handleReviewSubmit = async () => {
    try {
      const data = {
        UserId: currentUser.id,
        Body: reviewText,
        CoachId: coach.coach.id,
      };

      const response = await Resource.post("/reviews", data);

      if (response.status === 200) {
        console.log("postComments: " + JSON.stringify(response, null, 2));

        setCoachComments((coachComments) =>
          [response.data.reviewWithUser, ...coachComments].sort(
            (a, b) =>
              new Date(b.reviewObject.created_time) -
              new Date(a.reviewObject.created_time),
          ),
        );
        ShowSuccessMessage("Комментарий успешно добавлен");
        setReviewText("");
        handleCloseModal();
      }
    } catch (error) {
      if (error?.response?.data?.error) {
        ShowErrorMessage(error?.response?.data?.error);
      } else {
        ShowErrorMessage("Can't send review");
      }

      console.error("Can't send review: " + JSON.stringify(error, null, 2));
    }
  };

  const handleRevieChange = async (e) => {
    const inputValue = e.target.value;

    if (inputValue.length > 255) {
      const truncatedText = inputValue.slice(0, 255);
      setReviewText(truncatedText);
    } else {
      setReviewText(inputValue);
    }
  };

  const handleTrainingSelect = async (timeFrom, timeUntil) => {
    try {
      const data = {
        coach_id: coach.coach.id,
        time_from: timeFrom,
        time_until: timeUntil,
        date: date,
      };

      const response = await TrainingResource.post("/training/book", data);

      if (response.status === 200) {
        console.log("postTraining: " + JSON.stringify(response, null, 2));

        setDayTrainings((prevTrainings) =>
          prevTrainings.map((training) =>
            training.time_from === timeFrom && training.time_until === timeUntil
              ? { ...training, status: "забронировано" }
              : training,
          ),
        );

        ShowSuccessMessage("Тренировка успешно забронирована");
      }
    } catch (error) {
      if (error?.response?.data?.error) {
        ShowErrorMessage(error?.response?.data?.error);
      } else {
        ShowErrorMessage("Can't book training");
      }

      console.error("Can't book training: " + JSON.stringify(error, null, 2));
    }
  };

  /*comments schedule buttons*/
  const [active, setActive] = useState("schedule");

  const handleDateChange = async (dateFromInput) => {
    setDate(dateFromInput);

    try {
      const response = await TrainingResource.get(
        "/training/day/" + dateFromInput + "/coach/" + coach.coach.id,
      );

      if (response.status === 200) {
        console.log("get Trainings: " + JSON.stringify(response, null, 2));

        let trainingsFromResp = response.data.trainings;

        let dayTrainingsReact = [...initTrainingsState];

        for (const dayTrainingsReactElement of dayTrainingsReact) {
          dayTrainingsReactElement.status = "свободно";
        }

        setDayTrainings([...dayTrainingsReact]);

        if (trainingsFromResp.length > 0) {
          for (let i = 0; i < dayTrainingsReact.length; i++) {
            //
            for (let j = 0; j < trainingsFromResp.length; j++) {
              if (
                dayTrainingsReact[i].time_from ===
                  trainingsFromResp[j].TimeFrom.split("T")[1].slice(0, 5) &&
                dayTrainingsReact[i].time_until ===
                  trainingsFromResp[j].TimeUntil.split("T")[1].slice(0, 5)
              ) {
                if (
                  trainingsFromResp[j].Status === "booked" &&
                  trainingsFromResp[j].ClientId === currentUser.id
                ) {
                  dayTrainingsReact[i].status = "забронировано";
                }

                if (
                  trainingsFromResp[j].Status === "booked" &&
                  trainingsFromResp[j].ClientId !== currentUser.id
                ) {
                  dayTrainingsReact[i].status = "недоступно";
                }

                if (
                  trainingsFromResp[j].Status === "active" &&
                  trainingsFromResp[j].ClientId === currentUser.id
                ) {
                  dayTrainingsReact[i].status = "активно";
                }

                if (
                  trainingsFromResp[j].Status === "active" &&
                  trainingsFromResp[j].ClientId !== currentUser.id
                ) {
                  dayTrainingsReact[i].status = "недоступно";
                }

                if (trainingsFromResp[j].Status === "passed") {
                  dayTrainingsReact[i].status = "недоступно";
                }
              }
            }
          }

          setDayTrainings([...dayTrainingsReact]);
        }

        /*setCoachComments((coachComments) =>
                [response.data.reviewWithUser, ...coachComments].sort(
                  (a, b) =>
                    new Date(b.reviewObject.created_time) -
                    new Date(a.reviewObject.created_time),
                ),
              );*/
      }
    } catch (error) {
      if (error?.response?.data?.error) {
        ShowErrorMessage(error?.response?.data?.error);
      } else {
        ShowErrorMessage("Ошибка взятия тренировок");
      }

      console.error(
        "Ошибка взятия тренировок: " + JSON.stringify(error, null, 2),
      );
    }
  };

  return (
    <div
      style={{
        width: "70%",
        height: "100vh",
        background: "rgba(117,100,163,255)",
        overflowY: "scroll",
      }}
    >
      {/* Модальное окно */}
      <Modal
        open={openModal}
        onClose={handleCloseModal}
        aria-labelledby="modal-title"
        aria-describedby="modal-description"
      >
        <div
          style={{
            color: "white",
            position: "absolute",
            top: "50%",
            left: "50%",
            transform: "translate(-50%, -50%)",
            background: "rgba(160, 147, 197, 1)",
            borderRadius: "8px",
            padding: "20px",
            width: "300px",
          }}
        >
          <h2 id="modal-title">Оставить отзыв</h2>

          <TextareaAutosize
            style={{
              width: "100%",
              minHeight: "100px", // Минимальная высота для текстового поля
              resize: "none", // Отключаем изменение размера по умолчанию
              color: "white",
              backgroundColor: "transparent", // Прозрачный фон
              border: "1px solid white", // Белая рамка
              borderRadius: "4px",
              padding: "8px",
              boxSizing: "border-box",
            }}
            value={reviewText}
            onChange={handleRevieChange}
            maxRows={10} // Максимальное количество строк для отображения
          />

          <div style={{ display: "flex", justifyContent: "space-between" }}>
            <Button
              style={{
                marginTop: "10px",
                color: "white",
                background: "rgba(117, 100, 163, 255)",
              }}
              onClick={handleReviewSubmit}
            >
              Отправить
            </Button>

            <Button
              style={{
                marginTop: "10px",
                color: "white",
                background: "rgba(117, 100, 163, 255)",
              }}
              onClick={handleCloseModal}
            >
              Закрыть
            </Button>
          </div>
        </div>
      </Modal>

      <div
        style={{
          marginLeft: "10%",
          marginRight: "10%",
          background: "rgba(117,100,163,255)",
          marginBottom: "10px",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "column",
          marginTop: "50px",
        }}
      >
        <div style={{ width: "400px", paddingRight: "20px" }}>
          {/*{abonnement.Photo}*/}
          <img
            style={{ width: "100%", height: "auto" }}
            src={coach.coach.photo}
          />
        </div>

        <div style={{ marginTop: "5px", fontSize: "24px" }}>
          {coach.coach.name}
        </div>

        <div
          style={{ marginTop: "5px", fontSize: "18px", marginBottom: "20px" }}
        >
          {coach.coach.description}
        </div>

        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            width: "100%",
            marginBottom: "20px",
          }}
        >
          <div>
            <div style={{ paddingBottom: "15px", fontSize: "18px" }}>
              Услуги:
            </div>
            <div style={{ display: "flex" }}>
              {coach.services.map((Service) => (
                <div style={{ marginRight: "10px" }}>
                  <div style={{ width: "80px", height: "60px" }}>
                    <img
                      style={{ width: "100%", height: "auto" }}
                      src={Service.photo}
                    />
                  </div>
                  <div>
                    {Service.title === "swimming-pool"
                      ? "Бассейн"
                      : Service.title === "sauna"
                        ? "Сауна"
                        : Service.title === "gym"
                          ? "Тренажерный зал"
                          : Service.title}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {currentUser.role === "client" && (
            <Button
              style={{
                marginTop: "20px",
                color: "white",
                background: "rgba(160, 147, 197, 1)",
                width: "270px",
                height: "50px",
                marginBottom: "50px",
              }}
              onClick={handleOpenModal}
            >
              Оставить отзыв
            </Button>
          )}
        </div>

        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            gap: "200px",
          }}
        >
          <Button
            onClick={() => setActive("schedule")}
            disabled={active === "schedule"}
            style={{
              color: "white",
              background: "rgba(160, 147, 197, 1)",
              height: "62px",
              marginTop: "5%",
            }}
          >
            Расписание
          </Button>

          <Button
            onClick={() => setActive("comments")}
            disabled={active === "comments"}
            style={{
              color: "white",
              background: "rgba(160, 147, 197, 1)",
              height: "62px",
              marginTop: "5%",
            }}
          >
            Комментарии
          </Button>
        </div>

        {active === "comments" && (
          <>
            <div style={{ marginTop: "5px", fontSize: "18px" }}>
              Комментарии:
            </div>

            <div
              style={{
                display: "flex",
                justifyContent: "center",
                marginTop: "5px",
                width: "100%",
              }}
            >
              {coachComments.length > 0 ? (
                <div style={{ marginTop: "40px" }}>
                  {coachComments
                    .sort(
                      (a, b) =>
                        new Date(b.reviewObject.updated_time) -
                        new Date(a.reviewObject.updated_time),
                    )
                    .map((comment) => (
                      <div
                        style={{
                          display: "flex",
                          marginBottom: "50px",
                          background: "rgba(160, 147, 197, 1)",
                          padding: "10px",
                          borderRadius: "10px",
                          width: "500px",
                          alignItems: "flex-start",
                        }}
                      >
                        <div style={{ marginRight: "10px" }}>
                          <div style={{ width: "100px" }}>
                            <img
                              style={{ width: "100%", height: "auto" }}
                              src={comment.userObject.photo || noAva}
                            />
                          </div>
                          <div style={{ display: "flex", gap: "4px" }}>
                            <div>{comment.userObject.name}</div>
                          </div>
                        </div>

                        <div
                          style={{
                            flex: 1,
                            whiteSpace: "normal",
                            wordBreak: "break-word",
                          }}
                        >
                          {comment.reviewObject.body}
                        </div>
                      </div>
                    ))}
                </div>
              ) : (
                <div>There are no comments</div>
              )}
            </div>
          </>
        )}

        {active === "schedule" && (
          <div style={{ marginTop: "20px" }}>
            <div style={{ display: "flex", justifyContent: "center" }}>
              <input
                type="date"
                value={date}
                onChange={(e) => handleDateChange(e.target.value)}
                min={today}
                style={{
                  padding: "10px",
                  fontSize: "16px",
                  borderRadius: "8px",
                  border: "2px solid #A093C5",
                  backgroundColor: "rgba(160, 147, 197, 1)",
                  color: "white",
                  outline: "none",
                  cursor: "pointer",
                  transition: "0.3s",
                }}
              />
            </div>

            <div
              style={{
                marginTop: "20px",
                height: "400px",
                position: "relative",
              }}
            >
              {/* Заголовки колонок */}
              <div
                style={{
                  display: "flex",
                  position: "sticky",
                  top: 0,
                  zIndex: 1, // Поверх контента
                  fontWeight: "bold",
                  padding: "10px 0",
                  gap: "20px",
                }}
              >
                <div style={{ flex: 1 }}>Время</div>
                <div style={{ flex: 1 }}>Статус</div>
              </div>

              {/* Список клиентов */}
              {dayTrainings.length > 0 ? (
                <div>
                  {dayTrainings.map((training) => (
                    <div
                      key={training.time_from}
                      style={{
                        display: "flex",
                        marginTop: "10px",
                        gap: "20px",
                        alignItems: "center",
                        padding: "5px 0",
                      }}
                    >
                      {/* Время */}
                      <div style={{ flex: 1 }}>
                        {training.time_from + " - " + training.time_until}
                      </div>

                      {/* Статус */}
                      <Button
                        style={{
                          flex: 1,
                          color: "white",
                          background:
                            training.status === "забронировано"
                              ? "gray"
                              : training.status === "недоступно"
                                ? "red"
                                : training.status === "активно"
                                  ? "blue"
                                  : "rgba(160, 147, 197, 1)", // Цвет по умолчанию
                        }}
                        disabled={
                          training.status === "забронировано" ||
                          training.status === "недоступно" ||
                          training.status === "активно" ||
                          (today === date && training.status === "свободно")
                        }
                        onClick={() =>
                          handleTrainingSelect(
                            training.time_from,
                            training.time_until,
                          )
                        }
                      >
                        {today === date && training.status === "свободно"
                          ? "недоступно"
                          : training.status}
                      </Button>
                    </div>
                  ))}
                </div>
              ) : (
                <div>Нету никаких тренировок</div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
