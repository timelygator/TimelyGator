import { motion } from "framer-motion";
import { TrendingUp, Zap, TrendingDown, ListTodo } from "lucide-react";

const INSIGHTS = [
	{
		icon: TrendingUp,
		color: "text-green-500",
		insight: "Uptime is up 15% compared to last month, driven primarily because of a new project.",
	},
	{
		icon: Zap,
		color: "text-blue-500",
		insight: "Peak productive hours are between 10pm and 2am, with a 20% increase in productivity.",
	},
	{
		icon: TrendingDown,
		color: "text-purple-500",
		insight: 'Social Media domains taking up more than 10% of total usage, can be reduced.',
	},
	{
		icon: ListTodo,
		color: "text-yellow-500",
		insight: "Target completion rate up by 20% for this week, can be improved.",
	},
];

const AIPoweredInsights = () => {
	return (
		<motion.div
			className='bg-gray-800 bg-opacity-50 backdrop-filter backdrop-blur-lg shadow-lg rounded-xl p-6 border border-gray-700'
			initial={{ opacity: 0, y: 20 }}
			animate={{ opacity: 1, y: 0 }}
			transition={{ delay: 1.0 }}
		>
			<h2 className='text-xl font-semibold text-gray-100 mb-4'>AI-Powered Insights</h2>
			<div className='space-y-4'>
				{INSIGHTS.map((item, index) => (
					<div key={index} className='flex items-center space-x-3'>
						<div className={`p-2 rounded-full ${item.color} bg-opacity-20`}>
							<item.icon className={`size-6 ${item.color}`} />
						</div>
						<p className='text-gray-300'>{item.insight}</p>
					</div>
				))}
			</div>
		</motion.div>
	);
};
export default AIPoweredInsights;
