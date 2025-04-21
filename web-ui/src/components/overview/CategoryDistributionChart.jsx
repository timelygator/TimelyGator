import { motion } from "framer-motion";
import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer, Legend } from "recharts";
import { FaYoutube, FaCode, FaGithub, FaInstagram, FaShoppingCart } from "react-icons/fa";
import { useState, useEffect } from "react";
import Axios from "axios";

const COLORS = ["#6366F1", "#8B5CF6", "#EC4899", "#10B981", "#F59E0B"];

// Updated getIcon function that uses the category name
const getIcon = (category) => {
    switch (category) {
        case "YouTube":
            return <FaYoutube />;
        case "LeetCode":
            return <FaCode />;
        case "Github":
            return <FaGithub />;
        case "Instagram":
            return <FaInstagram />;
        case "Amazon":
            return <FaShoppingCart />;
        default:
            return null;
    }
};

const CategoryDistributionChart = () => {
    const [categoryData, setCategoryData] = useState([]);

    useEffect(() => {
        Axios.get("http://192.168.0.155:8080/api/v1/v1/buckets/tg-observer-window_yigirus/events")
            .then((res) => {
                // First, aggregate total duration and durations per category
                let totalDuration = 0;
                const categoryTotals = {};

                res.data.forEach((item) => {
                    totalDuration += item.duration;
                    // Determine the category.
                    // If the title includes "YouTube", then force the category to be YouTube.
                    let category = item.app;
                    if (item.title.includes("YouTube")) {
                        category = "YouTube";
                    }
                    categoryTotals[category] = (categoryTotals[category] || 0) + item.duration;
                });

                // Calculate percentage for each category based on its duration
                const data = Object.keys(categoryTotals).map((cat) => ({
                    name: cat,
                    // Percentage value
                    value: (categoryTotals[cat] / totalDuration) * 100,
                    icon: getIcon(cat),
                }));

                setCategoryData(data);
            })
            .catch((error) => {
                console.error("Error fetching category distribution data:", error);
            });
    }, []);

    return (
        <motion.div
            className='bg-backgroundSecondary bg-opacity-100 backdrop-blur-md shadow-lg rounded-xl p-6 border border-border'
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
        >
            <h2 className='text-lg font-medium mb-4 text-text'>Time Distribution</h2>
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
