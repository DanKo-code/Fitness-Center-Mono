import MainNav from "../Main/MainNav/mainNav";
    import MainProfile from "../Main/MainProfile/mainProfile";

export default function MainNavProfile() {
    return (
        <div style={{display: 'flex', height: '100vh', overflow: 'hidden', color: 'white'}}>
            <MainNav/>
            <MainProfile/>
        </div>
    )
}