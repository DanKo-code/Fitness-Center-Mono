import {AuthClient} from "../context/AuthContext";
import config from "../config"

const inMemoryJWTService = () => {
    let inMemoryJWT = null;
    let refreshTimeoutId = null;

    const refreshToken = (expiration) => {
        const timeoutTrigger = expiration - 10000;

        refreshTimeoutId = setTimeout(() => {
            AuthClient.post("/refresh")
                .then((res) => {
                    const {accessToken, accessTokenExpiration} = res.data;
                    setToken(accessToken, accessTokenExpiration);
                })
                .catch(console.error);
        }, timeoutTrigger);
    }

    const abortRefreshToken = () => {
      if(refreshTimeoutId){
          clearInterval(refreshTimeoutId);
      }
    };

    const getToken = () => inMemoryJWT;

    const setToken = (token, tokenExpiration) => {
        inMemoryJWT = token;
        refreshToken(tokenExpiration);
    }

    const deleteToken = () => {
        inMemoryJWT = null;
        abortRefreshToken();
        localStorage.setItem(config.LOGOUT_STORAGE_KEY, Date.now());
    }

    return {getToken, setToken, deleteToken}
}



export default inMemoryJWTService();