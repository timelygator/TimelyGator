import { Route, Routes } from "react-router-dom";

import Sidebar from "./components/common/Sidebar";

import OverviewPage from "./pages/OverviewPage";
import TabsPage from "./pages/TabsPage";
import AnalyticsPage from "./pages/AnalyticsPage";
import SettingsPage from "./pages/SettingsPage";

function App() {
	return (
		<div className='flex h-screen bg-gray-300 text-gray-100 overflow-hidden dark:bg-gray-900 '>
			{/* BG */}
			<div className='fixed inset-0 z-0'>
				{/* <div className='absolute inset-0 bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 opacity-80' /> */}
				<div className='absolute inset-0 backdrop-blur-sm' />
			</div>

			<Sidebar />
			<Routes>
				<Route path='/' element={<OverviewPage />} />
				<Route path='/tabs' element={<TabsPage />} />
				<Route path='/analytics' element={<AnalyticsPage />} />
				<Route path='/settings' element={<SettingsPage />} />
			</Routes>
		</div>
	);
}

export default App;
