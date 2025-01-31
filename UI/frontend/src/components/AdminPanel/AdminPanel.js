import 'react-multi-carousel/lib/styles.css';
import MainNav from "../Main/MainNav/mainNav";
import MainAbonnements from "../Main/MainAbonements/mainAbonnements";
import MainAdminPanel from "../Main/MainAdminPanel/mainAdminPanel";

export default function MainNavAdminPanel() {

    return (
        <div style={{display: 'flex', height: '100vh', overflow: 'hidden', color: 'white'}}>
            <MainNav/>
            <MainAdminPanel/>
        </div>
    );
}