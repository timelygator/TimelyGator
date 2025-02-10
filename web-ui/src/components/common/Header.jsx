import { useState, useEffect } from 'react';
import { Sun, Moon } from 'lucide-react';

const Header = ({ title }) => {
    const [isDarkMode, setIsDarkMode] = useState(false);

    useEffect(() => {
        if (isDarkMode) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    }, [isDarkMode]);

    const toggleTheme = () => {
        setIsDarkMode(!isDarkMode);
    };

    return (
        <header className='bg-gray-800 bg-opacity-50 backdrop-blur-md shadow-lg border-b border-gray-700'>
            <div className='max-w-7xl mx-auto py-4 px-4 sm:px-6 lg:px-8 flex justify-between items-center'>
                <h1 className='text-2xl font-semibold text-gray-100'>{title}</h1>
                <button
                    onClick={toggleTheme}
                    className='bg-gray-700 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded transition duration-200'
                >
                    {isDarkMode ? <Sun className='w-5 h-5' /> : <Moon className='w-5 h-5' />}
                </button>
            </div>
        </header>
    );
};

export default Header;