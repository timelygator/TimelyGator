import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { Chrome, FileCode2, Github, Layers, ArrowDownRight, ArrowUpRight } from "lucide-react";
import Axios from "axios";

const OverviewCards = () => {
    const [overviewData, setOverviewData] = useState([
        { name: "Top Application", value: "Loading...", change: 0, icon: Chrome },
        { name: "Top Window Titles", value: "Loading...", change: 0, icon: FileCode2 },
        { name: "Top Browser Domains", value: "Loading...", change: 0, icon: Github },
        { name: "Top Category", value: "Loading...", change: 0, icon: Layers },
    ]);

    useEffect(() => {
        Axios.get("http://192.168.0.166:8080/api/v1/v1/buckets/tg-observer-window_sidd-Predator-PHN16-72/events")
            .then((res) => {
                const appFrequency = {};
                const titleFrequency = {};
                const domainFrequency = {};

                res.data.forEach((event) => {
                    const app = event.app || "Unknown";
                    const title = event.title || "Unknown";

                    // Skip items with "Unknown" title
                    if (title === "Unknown") return;

                    // Increment application frequency
                    appFrequency[app] = (appFrequency[app] || 0) + 1;

                    // Extract tab name, site name, and application name from the title
                    const titleParts = title.split(" - ");
                    const tabName = titleParts[1] || "Unknown";
                    const siteName = titleParts[titleParts.length - 1]?.split("—")[0]?.trim() || "Unknown";

                    // Skip "Unknown" tab names and site names
                    if (tabName !== "Unknown") {
                        titleFrequency[tabName] = (titleFrequency[tabName] || 0) + 1;
                    }
                    if (siteName !== "Unknown") {
                        domainFrequency[siteName] = (domainFrequency[siteName] || 0) + 1;
                    }
                });

                // Find the top application, title, and domain by frequency
                const topApp = Object.entries(appFrequency).reduce(
                    (max, entry) => (entry[1] > max[1] ? entry : max),
                    ["Unknown", 0]
                );
                const topTitle = Object.entries(titleFrequency).reduce(
                    (max, entry) => (entry[1] > max[1] ? entry : max),
                    ["Unknown", 0]
                );
                const topDomain = Object.entries(domainFrequency).reduce(
                    (max, entry) => (entry[1] > max[1] ? entry : max),
                    ["Unknown", 0]
                );

                // Update the overview data dynamically
                setOverviewData((prevData) =>
                    prevData.map((item) => {
                        if (item.name === "Top Application") {
                            return { ...item, value: topApp[0], change: 0 };
                        } else if (item.name === "Top Window Titles") {
                            return { ...item, value: topTitle[0], change: 0 };
                        } else if (item.name === "Top Browser Domains") {
                            return { ...item, value: topDomain[0], change: 0 };
                        }
                        return item;
                    })
                );
            })
            .catch((error) => {
                console.error("Error fetching events data:", error);
            });
    }, []);

    return (
        <div className='grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-8'>
            {overviewData.map((item, index) => (
                <motion.div
                    key={item.name}
                    className='bg-backgroundSecondary bg-opacity-100 backdrop-filter backdrop-blur-lg shadow-lg rounded-xl p-6 border border-border'
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.1 }}
                >
                    <div className='flex items-center justify-between'>
                        <div>
                            <h3 className='text-sm font-medium text-gray-400'>{item.name}</h3>
                            <p className='mt-1 text-xl font-semibold text-text'>{item.value}</p>
                        </div>

                        <div
                            className={`p-3 rounded-full bg-opacity-20 ${item.change >= 0 ? "bg-green-500" : "bg-red-500"}`}
                        >
                            <item.icon className={`size-6  ${item.change >= 0 ? "text-green-500" : "text-red-500"}`} />
                        </div>
                    </div>
                </motion.div>
            ))}
        </div>
    );
};

export default OverviewCards;
