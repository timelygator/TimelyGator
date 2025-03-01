import { motion } from "framer-motion";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend, Cell } from "recharts";
import { useState, useEffect } from "react";
import Axios from "axios";

const COLORS = ["#6366F1", "#8B5CF6", "#EC4899", "#10B981", "#F59E0B"];

// Default data for website distribution
// const defaultWebsiteDistribution = [
//     { name: "Educational", value: 0 },
//     { name: "eCommerce", value: 0 },
//     { name: "Entertainment", value: 0 },
//     { name: "Social Media", value: 0 }
// ];

// const [websiteDistribution, setWebsiteDistribution] = useState(defaultWebsiteDistribution);

// useEffect(() => {
//     Axios.get("http://localhost:5000/api/website-distribution")
//         .then((res) => {
//             // Extract the 'Category' attribute from the API response
//             const data = res.data.map(item => ({
//                 name: item.Category,
//                 value: item.value // Assuming 'value' is the attribute you want to display
//             }));
//             setWebsiteDistribution(data);
//         })
//         .catch((error) => {
//             console.error("Error fetching website distribution data:", error);
//         });
// }, []);

const WEBSITE_DISTRIBUTION = [
    { name: "Educational", value: 4.9 },
    { name: "eCommerce", value: 0.5 },
    { name: "Entertainment", value: 3.7 },
    { name: "Social Media", value: 1.7 }
];

const WebsiteDistributionChart = () => {
    return (
        <motion.div
            className='bg-gray-800 bg-opacity-50 backdrop-blur-md shadow-lg rounded-xl p-6 lg:col-span-2 border border-gray-700'
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.4 }}
        >
            <h2 className='text-lg font-medium mb-4 text-gray-100'>Website Distribution</h2>

            <div className='h-80'>
                <ResponsiveContainer>
                    <BarChart data={WEBSITE_DISTRIBUTION}>
                        <CartesianGrid strokeDasharray='3 3' stroke='#4B5563' />
                        <XAxis dataKey='name' stroke='#9CA3AF' />
                        <YAxis stroke='#9CA3AF' />
                        <Tooltip
                            contentStyle={{
                                backgroundColor: "rgba(31, 41, 55, 0.8)",
                                borderColor: "#4B5563",
                            }}
                            itemStyle={{ color: "#E5E7EB" }}
                        />
                        <Bar dataKey="value" fill='#8884d8'>
                            {WEBSITE_DISTRIBUTION.map((entry, index) => (
                                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                            ))}
                        </Bar>
                        <Legend />
                    </BarChart>
                </ResponsiveContainer>
            </div>
        </motion.div>
    );
};

export default WebsiteDistributionChart;