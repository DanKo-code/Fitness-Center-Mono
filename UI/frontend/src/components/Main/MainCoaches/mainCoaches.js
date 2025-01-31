import React, { useEffect, useState } from "react";
import { Resource } from "../../../context/AuthContext";
import showErrorMessage from "../../../utils/showErrorMessage";
import CoachCard from "./CoachesCard/coachCard";

export default function MainCoaches() {
    const [coaches, setCoaches] = useState([]);

    useEffect(() => {
        Resource.get('/coaches')
            .then(response => {
                if(response?.data?.coaches?.coachWithServicesWithReviewsWithUsers?.length > 0){
                    setCoaches(response.data.coaches.coachWithServicesWithReviewsWithUsers);
                }
            })
            .catch(error => {
                showErrorMessage(error);
                console.error('Failed to fetch coaches:', error);
            });
    }, []);

    return (
        <div
            style={{
                width: '70%',
                height: '100vh',
                background: 'rgba(117, 100, 163, 255)',
                display: "flex",
                justifyContent: "center",
                alignItems: "center"
            }}
        >
            <div style={{ display: 'flex', justifyContent: 'center' }}>
                {coaches.length > 0 ? (
                    <div style={{ marginTop: '40px', height: '400px', overflowY: 'scroll' }}>
                        {coaches
                            .sort((a, b) => new Date(b.coach.updated_time) - new Date(a.coach.updated_time))
                            .map(coach => <CoachCard key={coach.coach.Id} coach={coach} width={'600px'} height={'400px'}></CoachCard>)}
                    </div>
                ) : (
                    <div>Нет никаких тренеров.</div>
                )}
            </div>
        </div>
    );
}
