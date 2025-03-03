import { Sun, MoonStar, GalleryVerticalEnd, Eye } from "lucide-react";
import { motion } from "framer-motion";

import Header from "../components/common/Header";
import StatCard from "../components/common/StatCard";
import TabTable from "../components/tabs/TabTable";
// import UserGrowthChart from "../components/tabs/TabGrowthChart";
// import UserActivityHeatmap from "../components/tabs/TabActivityHeatmap";
import WebsiteDistChart from "../components/tabs/TabDemographicsChart";



//                             <    API CALLING FOR TABS DATA TEST FORMAT   >





// const [tabStats, setTabStats] = useState({
// 	totaltabs: 0,
// 	activetabs: 0,
// 	hybernatetabs: 0,
// 	awakerate: "0%",
// });

// useEffect(() => {
// 	Axios.get("http://localhost:5000/api/tabs")
// 		.then((res) => {
// 			setTabStats({
// 				totaltabs: res.data.totaltabs,
// 				activetabs: res.data.activetabs,
// 				hybernatetabs: res.data.hybernatetabs,
// 				awakerate: res.data.awakerate,
// 			});
// 		})
// 		.catch((error) => {
// 			console.error("Error fetching tab stats:", error);
// 		});
// }, []);

const TabStats = {
    totaltabs: 56,
    activetabs: 29,
    hybernatetabs: 27,
    awakerate: "51.7%",
};

const TabsPage = () => {
    return (
        <div className='flex-1 overflow-auto relative z-10'>
            <Header title='Tabs' />
            <main className='max-w-7xl mx-auto py-6 px-4 lg:px-8'>
                {/* STATS */}
                <motion.div
                    className='grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-8'
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 1 }}
                >
                <StatCard
                    name='Total Tabs'
                    icon={GalleryVerticalEnd}
                    value={TabStats.totaltabs.toLocaleString()}
                    color='#EC4899'
                />
                <StatCard 
                    name='Active Tabs' 
                    icon={Sun} 
                    value={TabStats.activetabs} 
                    color='#10B981' 
                />
                <StatCard
                    name='Hybernating Tabs'
                    icon={MoonStar}
                    value={TabStats.hybernatetabs.toLocaleString()}
                    color='#F59E0B'
                />
                <StatCard 
                    name='Awake Rate' 
                    icon={Eye} 
                    value={TabStats.awakerate} 
                    color='#EF4444' 
                />
                </motion.div>
                <TabTable />

                {/* USER CHARTS */}
                <div className='grid grid-cols-1 lg:grid-cols-2 gap-6 mt-8'>
                    {/* <UserGrowthChart />
                    <UserActivityHeatmap /> */}
                    <WebsiteDistChart />
                </div>
            </main>
        </div>
    );
};

export default TabsPage;
