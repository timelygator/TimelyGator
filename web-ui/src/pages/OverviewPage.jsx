import { useState, useEffect } from "react";
import { GalleryVerticalEnd, Radio, TimerOff, Timer } from "lucide-react";
import { motion } from "framer-motion";
import Axios from "axios";

import Header from "../components/common/Header";
import StatCard from "../components/common/StatCard";
import UsageOverviewChart from "../components/overview/UsageOverviewChart";
import CategoryDistributionChart from "../components/overview/CategoryDistributionChart";
import WebsiteDistributionChart from "../components/overview/WebsiteDis";

const OverviewPage = () => {
    const [timeElapsed, setTimeElapsed] = useState("Loading...");
    const [idleTime, setIdleTime] = useState("Loading...");

    useEffect(() => {
        Axios.get("http://192.168.0.166:8080/api/v1/v1/buckets/tg-observer-afk_sidd-Predator-PHN16-72/events")
            .then((res) => {
                let totalTime = 0;
                let afkTime = 0;

                res.data.forEach((event) => {
                    totalTime += event.duration;
                    if (event.status === "afk") {
                        afkTime += event.duration;
                    }
                });

                // Convert time to hours and format it
                const formatTime = (timeInSeconds) => {
                    const hours = Math.floor(timeInSeconds / 3600);
                    const minutes = Math.floor((timeInSeconds % 3600) / 60);
                    return `${hours} h ${minutes} m`;
                };

                setTimeElapsed(formatTime(totalTime));
                setIdleTime(formatTime(afkTime));
            })
            .catch((error) => {
                console.error("Error fetching time data:", error);
            });
    }, []);

    return (
        <div className='flex-1 overflow-auto relative z-10'>
            <Header title='Overview' />

            <main className='max-w-7xl mx-auto py-6 px-4 lg:px-8'>
                {/* STATS */}
                <motion.div
                    className='grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-8'
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 1 }}
                >
                    <StatCard name='Time elapsed' icon={Timer} value={timeElapsed} color='#6366F1' />
                    <StatCard name='Idle Time' icon={TimerOff} value={idleTime} color='#8B5CF6' />
                    <StatCard name='Tabs open' icon={GalleryVerticalEnd} value='56' color='#EC4899' />
                    <StatCard name='Data Used' icon={Radio} value='1254 MB' color='#10B981' />
                </motion.div>

                {/* CHARTS */}
                <div className='grid grid-cols-1 lg:grid-cols-2 gap-8'>
                    <UsageOverviewChart />
                    <CategoryDistributionChart />
                    <WebsiteDistributionChart />
                </div>
            </main>
        </div>
    );
};

export default OverviewPage;
