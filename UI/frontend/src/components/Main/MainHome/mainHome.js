import Avatar from "@mui/material/Avatar";
import danilaAvatar from "../../../images/danila_avatar.jpg";
import Carousel from "react-multi-carousel";
import mainCss from "../../MainNavHome/MainNavHome.module.css";
import React from "react";
import { useSelector } from "react-redux";
import swimming_pool from "../../../images/swimming_pool.jpg";
import gym from "../../../images/gym.jpg";
import tennis from "../../../images/tennis.jpg";

export default function MainHome() {
  let props = "tablet";

  const pic1 = swimming_pool;
  const pic2 = gym;
  const pic3 = tennis;

  let user = useSelector((state) => state.userSliceMode.user);

  const responsive = {
    superLargeDesktop: {
      // the naming can be any, depends on you.
      breakpoint: { max: 4000, min: 3000 },
      items: 5,
    },
    desktop: {
      breakpoint: { max: 3000, min: 1024 },
      items: 3,
    },
    tablet: {
      breakpoint: { max: 1024, min: 464 },
      items: 2,
    },
    mobile: {
      breakpoint: { max: 464, min: 0 },
      items: 1,
    },
  };

  return (
    <div
      style={{
        width: "70%",
        height: "100vh",
        background: "rgba(117,100,163,255)",
      }}
    >
      <div
        style={{
          marginTop: "50px",
          display: "flex",
          justifyContent: "end",
          alignItems: "center",
          gap: "5%",
          fontSize: "22px",
          marginRight: "5%",
        }}
      >
        <Avatar src={user.photo} sx={{ width: 100, height: 100 }} />
        <div>{user.name}</div>
      </div>
      <div style={{ height: "100vh", marginRight: "5%", marginLeft: "5%" }}>
        <div style={{ marginTop: "50px" }}>
          Добро пожаловать в приложение FitLab для фитнес-центра, где вы можете
          приобрести абонементы, просматривать тренеров и оставлять отзывы о
          них, редактировать свой профиль и просматривать историю покупок.
        </div>
        <div>
          <div style={{ marginTop: "100px" }}>
            <Carousel
              responsive={responsive}
              infinite={true}
              autoPlay={true}
              autoPlaySpeed={3000}
              transitionDuration={500}
              itemClass="carousel-item-padding-20-px"
            >
              <div
                className={mainCss.sliderItemDiv}
                style={{ height: "200px" }}
              >
                <img src={pic1} className={mainCss.sliderItemImg} />
              </div>

              <div
                className={mainCss.sliderItemDiv}
                style={{ height: "200px" }}
              >
                <img src={pic2} className={mainCss.sliderItemImg} />
              </div>

              <div
                className={mainCss.sliderItemDiv}
                style={{ height: "200px" }}
              >
                <img src={pic3} className={mainCss.sliderItemImg} />
              </div>
            </Carousel>
          </div>
        </div>
      </div>
    </div>
  );
}
