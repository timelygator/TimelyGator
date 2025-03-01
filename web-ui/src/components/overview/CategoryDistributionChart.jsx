import { motion } from "framer-motion";
import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer, Legend } from "recharts";
// import { useState, useEffect } from "react";
// import Axios from "axios";

const categoryData = [
    { name: "YouTube", value: 3.7 },
    { name: "LeetCode", value: 2.6 },
    { name: "Github", value: 2.3 },
    { name: "Instagram", value: 1.7 },
    { name: "Amazon", value: 1.6 },
];

const COLORS = ["#6366F1", "#8B5CF6", "#EC4899", "#10B981", "#F59E0B"];

//            <   API CALLING FOR CATEGORY DISTRIBUTION TEST FORMAT   >

// Default data for category distribution
// const defaultCategoryData = [
// 	{ name: "YouTube", value: 0 },
// 	{ name: "LeetCode", value: 0 },
// 	{ name: "Github", value: 0 },
// 	{ name: "Instagram", value: 0 },
// 	{ name: "Amazon", value: 0 },
// ];

// const [categoryData, setCategoryData] = useState(defaultCategoryData);

// useEffect(() => {
// 	Axios.get("http://localhost:5000/api/category-distribution")
// 		.then((res) => {
// 			// Calculate the time spent on each site
// 			const data = res.data.map(item => {
// 				const startTime = new Date(item.Starttime);
// 				const endTime = new Date(item.EndTime);
// 				const timeSpent = (endTime - startTime) / (1000 * 60 * 60); // Convert milliseconds to hours
// 				return {
// 					name: item.Category,
// 					value: timeSpent
// 				};
// 			});
// 			setCategoryData(data);
// 		})
// 		.catch((error) => {
// 			console.error("Error fetching category distribution data:", error);
// 		});
// }, []);

//                     <     API format that may be used:    >
// [
//     {
//         "Category": "YouTube",
//         "Starttime": "2025-02-28T08:00:00Z",
//         "EndTime": "2025-02-28T09:00:00Z",
//         "activityId": "12345"
//     },
//     ...
// ]

const CategoryDistributionChart = () => {
    return (
        <motion.div
            className='bg-gray-800 bg-opacity-50 backdrop-blur-md shadow-lg rounded-xl p-6 border border-gray-700'
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
        >
            <h2 className='text-lg font-medium mb-4 text-gray-100'>Time Distribution</h2>
            <div className='h-80'>
                <ResponsiveContainer width={"100%"} height={"100%"}>
                    <PieChart>
                        <Pie
                            data={categoryData}
                            cx={"50%"}
                            cy={"50%"}
                            labelLine={false}
                            outerRadius={80}
                            fill='#8884d8'
                            dataKey='value'
                            label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                        >
                            {categoryData.map((entry, index) => (
                                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                            ))}
                        </Pie>
                        <Tooltip
                            contentStyle={{
                                backgroundColor: "rgba(31, 41, 55, 0.8)",
                                borderColor: "#4B5563",
                            }}
                            itemStyle={{ color: "#E5E7EB" }}
                        />
                        <Legend />
                    </PieChart>
                </ResponsiveContainer>
            </div>
        </motion.div>
    );
};

export default CategoryDistributionChart;
