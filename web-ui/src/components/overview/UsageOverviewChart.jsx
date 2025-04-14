import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";
import { motion } from "framer-motion";
import { useState, useEffect } from "react";
import Axios from "axios";

const UsageData = [
    { name: "Mon", Hours: 9 },
    { name: "Tue", Hours: 8 },
    { name: "Wed", Hours: 10 },
    { name: "Thur", Hours: 6.5 },
    { name: "Fri", Hours: 9 },
    { name: "Sat", Hours: 16 },
    { name: "Sun", Hours: 14 }
];

// Default data for usage overview
// const defaultUsageData = [
// 	{ name: "Mon", Hours: 0 },
// 	{ name: "Tue", Hours: 0 },
// 	{ name: "Wed", Hours: 0 },
// 	{ name: "Thur", Hours: 0 },
// 	{ name: "Fri", Hours: 0 },
// 	{ name: "Sat", Hours: 0 },
// 	{ name: "Sun", Hours: 0 }
// ];

// const [UsageData, setUsageData] = useState(defaultUsageData);

// useEffect(() => {
// 	Axios.get("http://localhost:5000/api/usage")
// 		.then((res) => {
// 			setUsageData(res.data);
// 		})
// 		.catch((error) => {
// 			console.error("Error fetching usage data:", error);
// 		});
// }, []);

// API format that may be used:
// [
//     {
//         "name": "Mon",
//         "Hours": 9,
//         "activityId": "12345"
//     },
//     ...
// ]

const UsageOverviewChart = () => {
    return (
        <motion.div
            className='bg-backgroundSecondary bg-opacity-100 backdrop-blur-md shadow-lg rounded-xl p-6 border border-border'
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
        >
            <h2 className='text-lg font-medium mb-4 text-text'>Usage Overview</h2>

            <div className='h-80'>
                <ResponsiveContainer width={"100%"} height={"100%"}>
                    <LineChart data={UsageData}>
                        <CartesianGrid strokeDasharray='3 3' stroke='#4B5563' />
                        <XAxis dataKey={"name"} stroke='#9ca3af' />
                        <YAxis stroke='#9ca3af' />
                        <Tooltip
                            contentStyle={{
                                backgroundColor: "rgba(31, 41, 55, 0.8)",
                                borderColor: "#4B5563",
                            }}
                            itemStyle={{ color: "#E5E7EB" }}
                        />
                        <Line
                            type='monotone'
                            dataKey='Hours'
                            stroke='#6366F1'
                            strokeWidth={3}
                            dot={{ fill: "#6366F1", strokeWidth: 2, r: 6 }}
                            activeDot={{ r: 8, strokeWidth: 2 }}
                        />
                    </LineChart>
                </ResponsiveContainer>
            </div>
        </motion.div>
    );
}
};

export default UsageOverviewChart;