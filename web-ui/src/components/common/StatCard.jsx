import React from 'react';
import { motion } from 'framer-motion';

const StatCard = ({ name, icon: Icon, value, color }) => {
    return (
        <motion.div
            className='bg-gray-800 bg-opacity-50 backdrop-blur-md shadow-lg rounded-xl p-6 border border-gray-700'
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 1 }}
            data-cy={`stat-card-${name.toLowerCase().replace(/\s+/g, '-')}`}
        >
            <div className='flex items-center justify-between'>
                <div>
                    <h3 className='text-sm font-medium text-gray-400'>{name}</h3>
                    <p className='mt-1 text-3xl font-semibold text-gray-100'>{value}</p>
                </div>
                <div className='p-3 rounded-full' style={{ backgroundColor: color }}>
                    <Icon size={23} className='text-white' /> {/* Set the size prop to 32 */}
                </div>
            </div>
        </motion.div>
    );
};

export default StatCard;
