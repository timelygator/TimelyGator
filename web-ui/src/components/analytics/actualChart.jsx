import { useState } from "react";
import { motion } from "framer-motion";
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from "recharts";

// Daily data for the entire year (365 days)
const actualData = Array.from({ length: 365 }, (_, index) => {
    const date = new Date(new Date().getFullYear(), 0, index + 1); // Generate date for each day of the year
    return {
        day: date.toISOString().split("T")[0], // Format date as YYYY-MM-DD
        weekday: date.toLocaleString("en-US", { weekday: "long" }), // Day of the week (e.g., Sunday, Monday)
        week: Math.ceil((index + 1) / 7), // Week number
        month: date.toLocaleString("en-US", { month: "short" }), // Month name (e.g., Jan, Feb)
        actual: Math.floor(Math.random() * 8) + 1, // Random hours between 1 and 8
        target: Math.floor(Math.random() * 8) + 1, // Random hours between 1 and 8
    };
});

// Aggregate data for weeks, months, and quarters
const aggregateData = (data, key) => {
    const aggregated = {};
    data.forEach((item) => {
        if (!aggregated[item[key]]) {
            aggregated[item[key]] = { [key]: item[key], actual: 0, target: 0 };
        }
        aggregated[item[key]].actual += item.actual;
        aggregated[item[key]].target += item.target;
    });

    if (key === "week") {
        return Object.values(aggregated).map((item, index) => ({
            ...item,
            weekLabel: `Week ${index + 1}`,
        }));
    }

    if (key === "month") {
        const weeksInMonth = ["Week 1", "Week 2", "Week 3", "Week 4"];
        return Object.values(aggregated).map((item, index) => ({
            ...item,
            weekLabel: weeksInMonth[index % 4] || `Week ${index + 1}`, // Ensure 4 weeks per month
        }));
    }

    return Object.values(aggregated);
};

const filterDataByTimeRange = (timeRange) => {
    const currentDate = new Date();
    const currentDayOfYear = Math.ceil(
        (currentDate - new Date(currentDate.getFullYear(), 0, 0)) / (24 * 60 * 60 * 1000)
    ); // Calculate the current day of the year
    const currentMonth = currentDate.getMonth(); // 0-based index for months
    const currentQuarter = Math.floor(currentMonth / 3); // 0-based index for quarters

    switch (timeRange) {
        case "This Week":
            const startOfWeek = currentDayOfYear - (currentDate.getDay() || 7) + 1; // Start of the current week
            return actualData.filter(
                (_, index) => index + 1 >= startOfWeek && index + 1 < startOfWeek + 7
            ); // Data for the current week
        case "This Month":
            const startOfMonth = Math.ceil(
                (new Date(currentDate.getFullYear(), currentMonth, 1) - new Date(currentDate.getFullYear(), 0, 1)) /
                    (24 * 60 * 60 * 1000)
            );
            const endOfMonth = Math.ceil(
                (new Date(currentDate.getFullYear(), currentMonth + 1, 0) - new Date(currentDate.getFullYear(), 0, 1)) /
                    (24 * 60 * 60 * 1000)
            );
            return aggregateData(
                actualData.filter((_, index) => index + 1 >= startOfMonth && index + 1 <= endOfMonth),
                "week"
            ); // Aggregate data by week for the current month
        case "This Quarter":
            const startOfQuarter = Math.ceil(
                (new Date(currentDate.getFullYear(), currentQuarter * 3, 1) - new Date(currentDate.getFullYear(), 0, 1)) /
                    (24 * 60 * 60 * 1000)
            );
            const endOfQuarter = Math.ceil(
                (new Date(currentDate.getFullYear(), currentQuarter * 3 + 3, 0) - new Date(currentDate.getFullYear(), 0, 1)) /
                    (24 * 60 * 60 * 1000)
            );
            return aggregateData(
                actualData.filter((_, index) => index + 1 >= startOfQuarter && index + 1 <= endOfQuarter),
                "month"
            ); // Aggregate data by month for the current quarter
        case "This Year":
            return aggregateData(actualData, "month"); // Aggregate data by month for the year
        default:
            return actualData; // Default to all data
    }
};

const ActualChart = () => {
    const [selectedTimeRange, setSelectedTimeRange] = useState("This Week"); // Default to "This Week"

    // Filter data based on the selected time range
    const filteredData = filterDataByTimeRange(selectedTimeRange);

    // Determine the X-axis key based on the selected time range
    const xAxisKey =
        selectedTimeRange === "This Week"
            ? "weekday" // Days of the week
            : selectedTimeRange === "This Month"
            ? "weekLabel" // Week labels (e.g., Week 1, Week 2)
            : "month"; // Months for quarters and years

    return (
        <motion.div
            className='bg-backgroundSecondary bg-opacity-100 backdrop-filter backdrop-blur-lg shadow-lg rounded-xl p-6 border border-border mb-8'
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
        >
            <div className='flex justify-between items-center mb-6'>
                <h2 className='text-xl font-semibold text-text'>Actual vs Target</h2>
                <select
                    className='bg-actionDefault text-text rounded-md px-3 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500'
                    value={selectedTimeRange}
                    onChange={(e) => setSelectedTimeRange(e.target.value)} // Update state on selection
                >
                    <option>This Week</option>
                    <option>This Month</option>
                    <option>This Quarter</option>
                    <option>This Year</option>
                </select>
            </div>

            <div style={{ width: "100%", height: 400 }}>
                <ResponsiveContainer>
                    <AreaChart data={filteredData}> {/* Render filtered data */}
                        <CartesianGrid strokeDasharray='3 3' stroke='#374151' />
                        <XAxis dataKey={xAxisKey} stroke='#9CA3AF' /> {/* Dynamic X-axis key */}
                        <YAxis stroke='#9CA3AF' />
                        <Tooltip
                            contentStyle={{ backgroundColor: "rgba(31, 41, 55, 0.8)", borderColor: "#4B5563" }}
                            itemStyle={{ color: "#E5E7EB" }}
                        />
                        <Legend />
                        <Area type='monotone' dataKey='actual' stroke='#8B5CF6' fill='#8B5CF6' fillOpacity={0.3} />
                        <Area type='monotone' dataKey='target' stroke='#10B981' fill='#10B981' fillOpacity={0.3} />
                    </AreaChart>
                </ResponsiveContainer>
            </div>
        </motion.div>
    );
};

export default ActualChart;
