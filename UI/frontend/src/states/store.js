import {configureStore} from "@reduxjs/toolkit"
import appSateSliceReducer from "./storeSlice/appStateSlice"
import mainComponentsSliceReducer from "./storeSlice/mainComponentsSlice"
import userSliceReducer from "./storeSlice/appStateSlice"

const store = configureStore({

    reducer: {
        appSateSliceMode : appSateSliceReducer,
        mainComponentsSliceMode : mainComponentsSliceReducer,
        userSliceMode: userSliceReducer
    },
})

export default store;