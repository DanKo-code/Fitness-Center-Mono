import {createSlice} from "@reduxjs/toolkit"


const initialState = {
    mainComponentsState: "main"
}

const mainComponentsSlice = createSlice({
    name: "mainComponents",
    initialState,
    reducers: {
        setMainComponentsState(state, action){
            state.mainComponentsState = action.payload
        }
    }
})

export const {setMainComponentsState} = mainComponentsSlice.actions
export default mainComponentsSlice.reducer