import { useState } from "react";
import { motion } from "framer-motion";
import * as Icons from "lucide-react";

// Add API data below in tabData
const tabData = [
    { id: 1, website: "YouTube", name: "React tutorial", Type: "Entertainment", status: "Active" },
    { id: 2, website: "YouTube", name: "Git tutorial", Type: "Entertainment", status: "Active" },
    { id: 3, website: "Github", name: "Github/PulkitGarg777", Type: "Educational", status: "Inactive" },
    { id: 4, website: "Instagram", name: "Instagram.com", Type: "Social", status: "Active" },
    { id: 5, website: "Google", name: "Homemade Chocolate Cake", Type: "Entertainment", status: "Active" },
];

const getIconForWebsite = (website) => {
    const normalizedWebsite = website.charAt(0).toUpperCase() + website.slice(1).toLowerCase().replace(/\s+/g, '');
    return Icons[normalizedWebsite] || Icons.Globe; // Default to Globe icon if no mapping found
};

const TabTable = () => {
    const [searchTerm, setSearchTerm] = useState("");
    const [filteredUsers, setFilteredUsers] = useState(tabData);

    const handleSearch = (e) => {
        const term = e.target.value.toLowerCase();
        setSearchTerm(term);
        const filtered = tabData.filter(
            (tab) => tab.website.toLowerCase().includes(term) || tab.name.toLowerCase().includes(term)
        );
        setFilteredUsers(filtered);
    };

    const handleDelete = (id) => {
        const updatedTabs = filteredUsers.filter((tab) => tab.id !== id);
        setFilteredUsers(updatedTabs);
    };

    return (
        <motion.div
            className='bg-backgroundSecondary bg-opacity-100 backdrop-blur-md shadow-lg rounded-xl p-6 border border-border'
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
        >
            <div className='flex justify-between items-center mb-6'>
                <h2 className='text-xl font-semibold text-text'>Tabs</h2>
                <div className='relative'>
                    <input
                        type='text'
                        placeholder='Search tabs...'
                        className='bg-actionDefault text-text placeholder-gray-400 rounded-lg pl-10 pr-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500'
                        value={searchTerm}
                        onChange={handleSearch}
                    />
                    <Icons.Search className='absolute left-3 top-2.5 text-gray-400' size={18} />
                </div>
            </div>

            <div className='overflow-x-auto'>
                <table className='min-w-full divide-y divide-gray-700'>
                    <thead>
                        <tr>
                            <th className='px-6 py-3 text-left text-xs font-medium text-cardTextHeading uppercase tracking-wider'>
                                Website
                            </th>
                            <th className='px-6 py-3 text-left text-xs font-medium text-cardTextHeading uppercase tracking-wider'>
                                Name
                            </th>
                            <th className='px-6 py-3 text-left text-xs font-medium text-cardTextHeading uppercase tracking-wider'>
                                Type
                            </th>
                            <th className='px-6 py-3 text-left text-xs font-medium text-cardTextHeading uppercase tracking-wider'>
                                Status
                            </th>
                            <th className='px-6 py-3 text-left text-xs font-medium text-cardTextHeading uppercase tracking-wider'>
                                Actions
                            </th>
                        </tr>
                    </thead>

                    <tbody className='divide-y divide-gray-700'>
                        {filteredUsers.map((tab) => {
                            const Icon = getIconForWebsite(tab.website);
                            return (
                                <motion.tr
                                    key={tab.id}
                                    initial={{ opacity: 0 }}
                                    animate={{ opacity: 1 }}
                                    transition={{ duration: 0.3 }}
                                >
                                    <td className='px-6 py-4 whitespace-nowrap'>
                                        <div className='flex items-center'>
                                            <div className='flex-shrink-0 h-10 w-10'>
                                                <Icon className='h-10 w-10 text-text' />
                                            </div>
                                            <div className='ml-4'>
                                                <div className='text-sm font-medium text-text'>{tab.website}</div>
                                            </div>
                                        </div>
                                    </td>

                                    <td className='px-6 py-4 whitespace-nowrap'>
                                        <div className='text-sm text-cardText'>{tab.name}</div>
                                    </td>
                                    <td className='px-6 py-4 whitespace-nowrap'>
                                        <span
                                            className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                                tab.Type === "Educational"
                                                    ? "bg-green-800 text-green-100"
                                                    : tab.Type === "Entertainment"
                                                    ? "bg-red-800 text-red-100"
                                                    : tab.Type === "Social"
                                                    ? "bg-pink-800 text-pink-100"
                                                    : tab.Type === "eCommerce"
                                                    ? "bg-yellow-800 text-yellow-100"
                                                    : "bg-gray-800 text-gray-100"
                                            }`}
                                        >
                                            {tab.Type}
                                        </span>
                                    </td>

                                    <td className='px-6 py-4 whitespace-nowrap'>
                                        <span
                                            className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                                tab.status === "Active"
                                                    ? "bg-green-800 text-green-100"
                                                    : "bg-red-800 text-red-100"
                                            }`}
                                        >
                                            {tab.status}
                                        </span>
                                    </td>

                                    <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-300'>
                                        {/* <button className='text-indigo-400 hover:text-indigo-300 mr-2'>Edit</button> */}
                                        <button
                                            className='text-red-400 hover:text-red-300'
                                            onClick={() => handleDelete(tab.id)}
                                        >
                                            Delete
                                        </button>
                                    </td>
                                </motion.tr>
                            );
                        })}
                    </tbody>
                </table>
            </div>
        </motion.div>
    );
};

export default TabTable;
