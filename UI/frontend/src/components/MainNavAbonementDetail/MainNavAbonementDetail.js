import 'react-multi-carousel/lib/styles.css';
import MainNav from "../Main/MainNav/mainNav";
import MainAbonnements from "../Main/MainAbonements/mainAbonnements";
import CoachDetailsCard from "../Main/MainCoaches/CoachesDetailedCard/coachDetailCard";

export default function MainNavCoachDetailCard() {

    return (
        <div style={{display: 'flex', height: '100vh', overflow: 'hidden', color: 'white'}}>
            <MainNav/>
            <CoachDetailsCard/>
        </div>

    );
}