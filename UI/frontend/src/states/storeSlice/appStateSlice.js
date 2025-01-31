import {createSlice} from "@reduxjs/toolkit"

const initialState = {
    appState: "signUp",
    user: ""
}

const appSateSlice = createSlice({
    name: "appSate",
    initialState,
    reducers: {
        setAppState(state, action){
            state.appState = action.payload
        },

        setUser(state, action){
            state.user = action.payload
        }
    }
})

export const {setAppState, setUser} = appSateSlice.actions
export default appSateSlice.reducer