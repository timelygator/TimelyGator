import { Sun, MoonStar, GalleryVerticalEnd, Eye } from "lucide-react";
import { motion } from "framer-motion";

import Header from "../components/common/Header";
import StatCard from "../components/common/StatCard";
import TabTable from "../components/tabs/TabTable";
// import UserGrowthChart from "../components/tabs/TabGrowthChart";
// import UserActivityHeatmap from "../components/tabs/TabActivityHeatmap";
import WebsiteDistChart from "../components/tabs/TabDemographicsChart";

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
					color='#6366F1'
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
