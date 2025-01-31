import 'react-multi-carousel/lib/styles.css';
import MainNav from "../Main/MainNav/mainNav";
import MainCoaches from "../Main/MainCoaches/mainCoaches";

export default function MainNavCoaches() {

    return (
        <div style={{display: 'flex', height: '100vh', overflow: 'hidden', color: 'white'}}>
            <MainNav/>
            <MainCoaches/>
        </div>
    );
}