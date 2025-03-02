import MainNav from "../Main/MainNav/mainNav";
import MainTrainingRoom from "../Main/MainTrainingRoom/mainTrainingRoom";

export default function MainNavTrainingRoom() {
  return (
    <div
      style={{
        display: "flex",
        height: "100vh",
        overflow: "hidden",
        color: "white",
      }}
    >
      <MainNav />
      <MainTrainingRoom />
    </div>
  );
}
