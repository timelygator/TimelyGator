import React, { useState } from 'react';

const ToggleSwitch = ({ label, isOn, onToggle }) => {
    const [isToggled, setIsToggled] = useState(isOn);

    const handleToggle = () => {
        setIsToggled(!isToggled);
        onToggle(!isToggled);
    };

    return (
        <div className='flex items-center justify-between py-3'>
            <span className='text-gray-300'>{label}</span>
            <button
                data-cy="toggle-switch"
                className={`
        relative inline-flex items-center h-6 rounded-full w-11 transition-colors focus:outline-none
        ${isToggled ? "bg-indigo-600" : "bg-gray-600"}
        `}
                onClick={handleToggle}
            >
                <span
                    className={`inline-block w-4 h-4 transform transition-transform bg-white rounded-full 
            ${isToggled ? "translate-x-6" : "translate-x-1"}
            `}
                />
            </button>
        </div>
    );
};

export default ToggleSwitch;
