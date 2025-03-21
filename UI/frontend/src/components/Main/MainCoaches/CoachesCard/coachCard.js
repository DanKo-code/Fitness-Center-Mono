import sad_doing_coachs_card from "../../../../images/sad_doing_abonnements_card.jpg";
import Button from "@mui/material/Button";
import React, { useEffect } from "react";
import axios from "axios";
import { useDispatch, useSelector } from "react-redux";
import sad_doing_abonnements_card from "../../../../images/sad_doing_abonnements_card.jpg";
import { useNavigate } from "react-router-dom";
import noAva from "../../../../images/no_ava.png";

export default function CoachCard(props) {
  const { coach, width, height, imageSize, button, onClick } = props;
  const navigate = useNavigate();

  const handleDetails = async () => {
    navigate("/main/coaches/details", { state: { coach } });
  };

  return (
    <div
      onClick={onClick}
      style={{
        background: "rgb(160, 147, 197)",
        width: width,
        height: height,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        borderRadius: "20px",
      }}
      key={coach.Id}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          flexDirection: "column",
        }}
      >
        {/*<div style={{width: '320px', paddingRight: '20px'}}>
                    {abonnement.Photo}
                    <img style={{width: '100%', height: 'auto'}} src={sad_doing_abonnements_card}/>
                </div>*/}

        <div
          style={{
            display: "flex",
            justifyContent: "center",
            paddingBottom: "20px",
            borderRadius: "20px",
          }}
        >
          <div
            style={{
              minWidth: "180px",
              paddingRight: "20px",
              borderRadius: "20px",
            }}
          >
            {/*{abonnement.Photo}*/}
            <img
              style={{ width: "100%", height: "120px", borderRadius: "20px" }}
              src={coach.coach.photo || noAva}
            />
          </div>

          <div
            style={{
              display: "flex",
              flexDirection: "column",
              justifyContent: "space-between",
              gap: "10px",
            }}
          >
            <div
              title={coach.coach.name}
              style={{
                marginTop: "5px",
                fontSize: "24px",
                overflow: "hidden",
                width: "120px",
                display: "-webkit-box",
                WebkitBoxOrient: "vertical",
                WebkitLineClamp: 2, // Ограничивает до двух строк
                lineHeight: "1.2em", // Высота строки
                maxHeight: "2.4em", // Две строки по 1.2em
                minHeight: "2.4em", // Минимальная высота для фиксации двух строк
              }}
            >
              {coach.coach.name}
            </div>

            {button ?? (
              <Button
                style={{
                  color: "white",
                  background: "rgba(117,100,163,255)",
                  width: "180px",
                  height: "50px",
                  marginBottom: "50px",
                }}
                onClick={handleDetails}
              >
                Подробности
              </Button>
            )}
          </div>
        </div>

        <div style={{ paddingBottom: "15px", fontSize: "18px" }}>Услуги:</div>
        <div style={{ display: "flex" }}>
          {coach.services.map((Service) => (
            <div style={{ marginRight: "10px" }}>
              <div style={{ width: "80px", height: "60px" }}>
                <img
                  style={{ width: "100%", height: "auto" }}
                  src={Service.photo || noAva}
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
    </div>
  );
}
