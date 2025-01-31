import { enqueueSnackbar } from "notistack";

export default (error) =>{


  if(error){
    enqueueSnackbar(error, {variant: "error"})
    return;
  }

  let yupError = JSON.parse(error.response.data.error).errors;

  if(yupError){
    let fullError = '';

    yupError.forEach(error => {

      console.log('error: '+ error)

      fullError += error + '\n';
    });

    enqueueSnackbar(fullError, { variant: "error" })
  }
  else{
    enqueueSnackbar(error.response.data.error, { variant: "error" })
  }



}