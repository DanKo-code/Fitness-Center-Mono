import { enqueueSnackbar } from "notistack";

export default (message) =>{
    enqueueSnackbar(message, { variant: "success" })}

