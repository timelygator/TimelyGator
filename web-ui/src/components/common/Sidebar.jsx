import { BarChart2, NotebookTabs, Menu, Settings, TrendingUp } from "lucide-react";
import { useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { Link } from "react-router-dom";

const SIDEBAR_ITEMS = [
    {
        name: "Overview",
        icon: BarChart2,
        color: "#6366f1",
        href: "/",
    },
    { 
        name: "Tabs", 
        icon: NotebookTabs, 
        color: "#EC4899", 
        href: "/tabs" 
    },
    { 
        name: "Analytics", 
        icon: TrendingUp, 
        color: "#3B82F6", 
        href: "/analytics" 
    },
    { 
        name: "Settings", 
        icon: Settings, 
        color: "#6EE7B7", 
        href: "/settings" 
    },
];

const Sidebar = () => {
    const [isSidebarOpen, setIsSidebarOpen] = useState(true);

return (
    <motion.div
        className={`relative z-10 transition-all duration-300 ease-in-out flex-shrink-0 ${isSidebarOpen ? "w-64" : "w-20"
            }`}
        animate={{ width: isSidebarOpen ? 256 : 80 }}
        data-cy="sidebar"
    >
        <div className='h-full bg-backgroundSecondary bg-opacity-100 backdrop-blur-md p-4 flex flex-col border-r border-border'>
            <motion.button
                whileHover={{ scale: 1.1 }}
                whileTap={{ scale: 0.9 }}
                onClick={() => setIsSidebarOpen(!isSidebarOpen)}
                className='p-2 rounded-full hover:bg-actionHover transition-colors max-w-fit'
                data-cy="collapse-button"
            >
                <Menu size={24} />
            </motion.button>

            <nav className='mt-8 flex-grow'>
                {SIDEBAR_ITEMS.map((item) => (
                    <Link key={item.href} to={item.href}>
                        <motion.div
                            className='flex items-center p-4 text-sm font-medium rounded-lg hover:bg-actionHover transition-colors mb-2'
                            data-cy={`sidebar-item-${item.name.toLowerCase()}`}
                        >
                            <item.icon size={20} style={{ color: item.color, minWidth: "20px" }} />
                            <AnimatePresence>
                                {isSidebarOpen && (
                                    <motion.span
                                        className='ml-4 whitespace-nowrap'
                                        initial={{ opacity: 0, width: 0 }}
                                        animate={{ opacity: 1, width: "auto" }}
                                        exit={{ opacity: 0, width: 0 }}
                                        transition={{ duration: 0.2, delay: 0.3 }}
                                    >
                                        {item.name}
                                    </motion.span>
                                )}
                            </AnimatePresence>
                        </motion.div>
                    </Link>
                ))}
            </nav>
        </div>
    </motion.div>
);
};

export default Sidebar;
