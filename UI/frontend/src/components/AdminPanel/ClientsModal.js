import IconButton from "@mui/material/IconButton";
import CloseIcon from "@mui/icons-material/Close";
import React, { useEffect, useState } from "react";
import { Resource } from "../../context/AuthContext";
import showErrorMessage from "../../utils/showErrorMessage";
import Button from "@mui/material/Button";
import ShowSuccessMessage from "../../utils/showSuccessMessage";
import ShowErrorMessage from "../../utils/showErrorMessage";
import noAva from "../../images/no_ava.png";

export default function ClientsModal({ onClose }) {
  const [clients, setClients] = useState([]);
  const [currentClient, setCurrentClient] = useState("");

  const handleCloseClientsModal = () => {
    onClose();
  };

  const handleClientSelection = async (selectedClient) => {
    if (currentClient?.id === selectedClient.id) {
      setCurrentClient("");
    } else {
      setCurrentClient(selectedClient);
    }
  };

  const handleDelete = async (event) => {
    event.preventDefault();

    const clientId = currentClient.id;

    const url = `/users/${clientId}`;

    try {
      if (currentClient) {
        const response = await Resource.delete(url);

        if (response.status === 200) {
          setClients(
            clients.filter((item) => item.id !== response.data.client.id),
          );

          ShowSuccessMessage("Клиент успешно удален");
        }
      }
    } catch (error) {
      ShowErrorMessage(error);
      console.error("response.status: " + JSON.stringify(error, null, 2));
    }
  };

  useEffect(() => {
    const fetchClients = async () => {
      try {
        const response = await Resource.get("/users");

        if (response.status === 200) {
          if (response?.data?.clients?.length > 0) {
            setClients(response.data.clients);
          }
        }
      } catch (error) {
        if (error?.response?.data?.error) {
          ShowErrorMessage(error.response.data.error);
        } else {
          ShowErrorMessage("Can't get clients");
        }

        console.error("Can't get clients: " + JSON.stringify(error, null, 2));
      }
    };

    fetchClients();
  }, []);

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
        height: "600px",
      }}
    >
      <div>
        <IconButton
          onClick={handleCloseClientsModal}
          size="large"
          sx={{ position: "absolute", top: 10, right: 10 }}
        >
          <CloseIcon />
        </IconButton>
      </div>

      <div
        style={{
          marginTop: "10px",
          width: "100%",
          height: "100%",
        }}
      >
        <div
          style={{
            marginTop: "40px",
            display: "flex",
            justifyContent: "center",
            fontWeight: "bold",
            fontSize: "24px",
          }}
        >
          Все Клиенты
        </div>
        <div
          style={{ marginTop: "20px", height: "400px", position: "relative" }}
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
            <div style={{ width: "200px", textAlign: "center" }}>Фото</div>
            <div style={{ flex: 1 }}>Email</div>
            <div style={{ flex: 1 }}>Имя</div>
            <div style={{ flex: 1 }}>Время Создания</div>
            <div style={{ flex: 1 }}>Время Обновления</div>
          </div>

          {/* Список клиентов */}
          {clients.length > 0 ? (
            <div
              style={{
                height: "400px",
                overflowY: "scroll",
              }}
            >
              {clients.map((client) => (
                <div
                  key={client.id}
                  style={{
                    display: "flex",
                    marginTop: "10px",
                    gap: "20px",
                    alignItems: "center",
                    padding: "5px 0",
                    border:
                      currentClient === client ? "3px solid yellow" : "none",
                  }}
                  onClick={() => handleClientSelection(client)}
                >
                  {/* Photo */}
                  <img
                    style={{
                      width: "200px",
                      height: "110px",
                    }}
                    src={client.photo || noAva}
                    alt="client"
                  />

                  {/* Email */}
                  <div style={{ flex: 1 }}>{client.email}</div>

                  {/* Name */}
                  <div style={{ flex: 1 }}>{client.name}</div>

                  {/* Created Time */}
                  <div style={{ flex: 1 }}>{client.created_time}</div>

                  {/* Updated Time */}
                  <div style={{ flex: 1 }}>{client.updated_time}</div>
                </div>
              ))}
            </div>
          ) : (
            <div>Нету никаких клиентов</div>
          )}
        </div>

        <div
          style={{
            marginTop: "60px",
            display: currentClient ? "flex" : "none",
            justifyContent: "center",
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
            УДАЛИТЬ
          </Button>
        </div>
      </div>
    </div>
  );
}
