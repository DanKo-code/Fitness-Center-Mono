import { useDispatch, useSelector } from "react-redux";
import React, { useEffect, useRef, useState } from "react";
import { Resource } from "../../../context/AuthContext";
import { setUser } from "../../../states/storeSlice/appStateSlice";
import ShowSuccessMessage from "../../../utils/showSuccessMessage";
import ShowErrorMessage from "../../../utils/showErrorMessage";
import noAva from "../../../images/no_ava.png";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import config from "../../../config";
import AbonnementCard from "../MainAbonements/AbonementsCard/abonementCard";
import { useNavigate, useParams } from "react-router-dom";
import inMemoryJWT from "../../../services/inMemoryJWT";

const Room = (props) => {
  const { trainingId, coachId, coachName, clientName } = useParams();

  let currentUser = useSelector((state) => state.userSliceMode.user);

  const navigate = useNavigate();

  //our camera and micro
  const userVideo = useRef();
  const userStream = useRef();

  const partnerVideo = useRef();
  const peerRef = useRef();

  //ws obj
  const webSocketRef = useRef();

  const [cameras, setCameras] = useState([]);
  const [microphones, setMicrophones] = useState([]);
  const [selectedCamera, setSelectedCamera] = useState("");
  const [selectedMicrophone, setSelectedMicrophone] = useState("");
  const [isOpenMediaSelected, setIsOpenMediaSelected] = useState("");

  const [isCameraOn, setIsCameraOn] = useState(true);
  const [isMicrophoneOn, setIsMicrophoneOn] = useState(true);

  const toggleCamera = () => {
    if (userStream.current) {
      userStream.current.getVideoTracks().forEach((track) => {
        track.enabled = !isCameraOn;
      });
      setIsCameraOn(!isCameraOn);
    }
  };

  const toggleMicrophone = () => {
    if (userStream.current) {
      userStream.current.getAudioTracks().forEach((track) => {
        track.enabled = !isMicrophoneOn;
      });
      setIsMicrophoneOn(!isMicrophoneOn);
    }
  };

  const initDevices = async () => {
    try {
      // Запрашиваем разрешение на доступ к камере и микрофону
      await navigator.mediaDevices.getUserMedia({ video: true, audio: true });

      // После получения доступа перечисляем устройства
      const allDevices = await navigator.mediaDevices.enumerateDevices();

      setCameras(allDevices.filter((d) => d.kind === "videoinput"));
      setMicrophones(allDevices.filter((d) => d.kind === "audioinput"));
    } catch (error) {
      console.error("Ошибка доступа к устройствам:", error);
      ShowErrorMessage("Необходимо дать доступ к устройствам ввода и вывода");
    }
  };

  const openMedia = async (videoDeviceId, audioDeviceId) => {
    const constraints = {
      video: videoDeviceId ? { deviceId: { exact: videoDeviceId } } : true,
      audio: {
        deviceId: audioDeviceId ? { exact: audioDeviceId } : undefined,
        echoCancellation: true, // Включение подавления эха
        noiseSuppression: true, // Подавление шума
      },
    };

    navigator.mediaDevices
      .getUserMedia(constraints)
      .then((stream) => {
        userVideo.current.srcObject = stream;
        userStream.current = stream;

        setIsOpenMediaSelected(true);

        const accessToken = inMemoryJWT.getToken();

        webSocketRef.current = new WebSocket(
          config.TRAINING_JOIN_API_URL +
            `/training/join/${trainingId}?coachId=${coachId}&token=${accessToken}`,
        );

        webSocketRef.current.addEventListener("open", () => {
          webSocketRef.current.send(JSON.stringify({ join: true }));
        });

        webSocketRef.current.addEventListener("message", async (e) => {
          const message = JSON.parse(e.data);

          console.log("message: " + JSON.stringify(e));

          if (message.join) {
            callUser();
          }

          if (message.offer) {
            handleOffer(message.offer);
          }

          if (message.answer) {
            console.log("Receiving Answer");
            peerRef.current.setRemoteDescription(
              new RTCSessionDescription(message.answer),
            );
          }

          if (message.iceCandidate) {
            console.log("Receiving and Adding ICE Candidate");
            try {
              await peerRef.current.addIceCandidate(message.iceCandidate);
            } catch (err) {
              console.log("Error Receiving ICE Candidate", err);
            }
          }

          if (message.end) {
            if (userStream.current) {
              userStream.current.getTracks().forEach((track) => track.stop());
              userStream.current = null;
            }

            if (userVideo.current) userVideo.current.srcObject = null;
            if (partnerVideo.current) partnerVideo.current.srcObject = null;

            if (peerRef.current) {
              peerRef.current.close();
              peerRef.current = null;
            }

            if (webSocketRef.current) {
              webSocketRef.current.close();
              webSocketRef.current = null;
            }

            navigate("/main");
          }
        });
      })
      .catch((e) => {
        ShowErrorMessage("Невалидные устройства ввода или вывода");
      });
  };

  useEffect(() => {
    initDevices().then((r) => console.log("devices received"));
  }, []);

  const callUser = () => {
    console.log("Calling Other User");
    peerRef.current = createPeer();

    userStream.current.getTracks().forEach((track) => {
      peerRef.current.addTrack(track, userStream.current);
    });
  };

  const handleOffer = async (offer) => {
    console.log("Received Offer, Creating Answer");
    peerRef.current = createPeer();

    await peerRef.current.setRemoteDescription(
      new RTCSessionDescription(offer),
    );

    userStream.current.getTracks().forEach((track) => {
      peerRef.current.addTrack(track, userStream.current);
    });

    const answer = await peerRef.current.createAnswer();
    await peerRef.current.setLocalDescription(answer);

    webSocketRef.current.send(
      JSON.stringify({ answer: peerRef.current.localDescription }),
    );
  };

  const createPeer = () => {
    console.log("Creating Peer Connection");
    const peer = new RTCPeerConnection({
      iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
    });

    peer.onnegotiationneeded = handleNegotiationNeeded;
    peer.onicecandidate = handleIceCandidateEvent;
    peer.ontrack = handleTrackEvent;

    return peer;
  };

  const handleNegotiationNeeded = async () => {
    console.log("Creating Offer");

    try {
      const myOffer = await peerRef.current.createOffer();
      await peerRef.current.setLocalDescription(myOffer);

      webSocketRef.current.send(
        JSON.stringify({ offer: peerRef.current.localDescription }),
      );
    } catch (err) {}
  };

  const handleIceCandidateEvent = (e) => {
    console.log("Found Ice Candidate");
    if (e.candidate) {
      console.log(e.candidate);
      webSocketRef.current.send(JSON.stringify({ iceCandidate: e.candidate }));
    }
  };

  const handleTrackEvent = (e) => {
    console.log("Received Tracks");
    partnerVideo.current.srcObject = e.streams[0];
  };

  const videoStyle = {
    width: "500px",
    height: "300px",
    borderWidth: "4px",
    borderColor: "rgba(160, 147, 197, 1)",
    borderStyle: "solid",
    borderRadius: "10px",
    objectFit: "cover",
  };

  return currentUser.role === "client" || currentUser.role === "coach" ? (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        background: "rgba(117,100,163,255)",
        flexDirection: "column",
        width: "100%",
      }}
    >
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          gap: "40px",
          background: "rgba(117,100,163,255)",
          width: "100%",
        }}
      >
        <div style={{ display: "flex", flexDirection: "column" }}>
          <video
            style={{ ...videoStyle, transform: "scaleX(-1)" }}
            autoPlay
            ref={userVideo}
          ></video>
          <div
            style={{
              paddingTop: "10px",
              display: "flex",
              justifyContent: "center",
              fontSize: "24px",
            }}
          >
            {currentUser.role === "client" ? clientName : coachName}
          </div>
        </div>

        <div style={{ display: "flex", flexDirection: "column" }}>
          <video
            style={{ ...videoStyle, transform: "scaleX(-1)" }}
            autoPlay
            ref={partnerVideo}
          ></video>
          <div
            style={{
              paddingTop: "10px",
              display: "flex",
              justifyContent: "center",
              fontSize: "24px",
            }}
          >
            {currentUser.role === "client" ? coachName : clientName}
          </div>
        </div>
      </div>
      {!isOpenMediaSelected && (
        <div
          style={{
            marginTop: "10px",
            display: "flex",
            width: "89%",
            gap: "15px",
          }}
        >
          <select
            style={{
              width: "150px",
              height: "40px",
              color: "white",
              background: "rgba(160, 147, 197, 1)",
              border: "none",
              borderRadius: "5px",
              padding: "5px 10px",
              fontSize: "16px",
              cursor: "pointer",
              outline: "none",
            }}
            onChange={(e) => setSelectedCamera(e.target.value)}
          >
            {cameras.map((camera) => (
              <option key={camera.deviceId} value={camera.deviceId}>
                {camera.label || `Камера ${camera.deviceId}`}
              </option>
            ))}
          </select>
          <select
            style={{
              width: "150px",
              height: "40px",
              color: "white",
              background: "rgba(160, 147, 197, 1)",
              border: "none",
              borderRadius: "5px",
              padding: "5px 10px",
              fontSize: "16px",
              cursor: "pointer",
              outline: "none",
            }}
            onChange={(e) => setSelectedMicrophone(e.target.value)}
          >
            {microphones.map((mic) => (
              <option key={mic.deviceId} value={mic.deviceId}>
                {mic.label || `Микрофон ${mic.deviceId}`}
              </option>
            ))}
          </select>
          <Button
            style={{
              width: "160px",
              height: "40px",
              color: "white",
              background: "rgba(160, 147, 197, 1)",
              border: "none",
              borderRadius: "5px",
              padding: "5px 10px",
              fontSize: "16px",
              cursor: "pointer",
              outline: "none",
            }}
            onClick={() => openMedia(selectedCamera, selectedMicrophone)}
          >
            Начать звонок
          </Button>
        </div>
      )}
      {isOpenMediaSelected && (
        <div
          style={{
            marginTop: "10px",
            display: "flex",
            width: "89%",
            gap: "15px",
          }}
        >
          <Button
            onClick={toggleCamera}
            style={{ background: "rgba(160, 147, 197, 1)", color: "white" }}
          >
            {isCameraOn ? "Выключить камеру" : "Включить камеру"}
          </Button>
          <Button
            onClick={toggleMicrophone}
            style={{
              background: "rgba(160, 147, 197, 1)",
              color: "white",
              fontSize: "16px",
            }}
          >
            {isMicrophoneOn ? "Выключить микрофон" : "Включить микрофон"}
          </Button>
        </div>
      )}
    </div>
  ) : null;
};

export default Room;
