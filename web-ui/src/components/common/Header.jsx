import { useState, useEffect } from 'react';
import { Sun, Moon } from 'lucide-react';

const Header = ({ title }) => {
    const [isDarkMode, setIsDarkMode] = useState(() => {
        // Retrieve the dark mode preference from localStorage
        const savedTheme = localStorage.getItem('isDarkMode');
        return savedTheme === 'true'; // Default to false if not set
    });

    useEffect(() => {
        // Apply the theme class to the document
        if (isDarkMode) {
            document.documentElement.classList.remove('light');
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
            document.documentElement.classList.add('light');
        }

        // Save the dark mode preference to localStorage
        localStorage.setItem('isDarkMode', isDarkMode);
    }, [isDarkMode]);

    const toggleTheme = () => {
        setIsDarkMode(!isDarkMode);
    };

    return (
        <header className='bg-backgroundSecondary bg-opacity-100 backdrop-blur-md shadow-lg border-b border-border'>
            <div className='max-w-7xl mx-auto py-4 px-4 sm:px-6 lg:px-8 flex justify-between items-center'>
                <h1 className='text-2xl font-semibold text-text'>{title}</h1>
                <button
                    onClick={toggleTheme}
                    className='bg-actionDefault hover:bg-actionHover text-text font-bold py-2 px-4 rounded transition duration-200'
                >
                    {isDarkMode ? <Sun className='w-5 h-5' /> : <Moon className='w-5 h-5' />}
                </button>
            </div>
        </header>
    );
};

export default Header;