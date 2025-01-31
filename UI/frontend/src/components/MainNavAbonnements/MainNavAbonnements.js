import 'react-multi-carousel/lib/styles.css';
import MainNav from "../Main/MainNav/mainNav";
import MainAbonnements from "../Main/MainAbonements/mainAbonnements";

export default function MainNavAbonnements() {

    return (
        <div style={{display: 'flex', height: '100vh', overflow: 'hidden', color: 'white'}}>
            <MainNav/>
            <MainAbonnements/>
        </div>

    );
}