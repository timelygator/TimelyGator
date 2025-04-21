import { useState } from "react";
import { User } from "lucide-react";
import SettingSection from "./SettingSection";

const Profile = () => {
    const [editing, setEditing] = useState(false);
    const [name, setName] = useState("ABC XYZ");
    const [email, setEmail] = useState("abc.xyz@example.com");
    const [profilePic, setProfilePic] = useState("/profile pic for various.jpg");
    const [previewPic, setPreviewPic] = useState(profilePic);

    const handleEdit = () => setEditing(true);
    const handleCancel = () => {
        // Reset preview image and editing state when canceled
        setPreviewPic(profilePic);
        setEditing(false);
    };
    const handleSave = () => {
        // might update the backend here.
        setProfilePic(previewPic);
        setEditing(false);
    };

    const handleImageChange = (e) => {
        if (e.target.files && e.target.files[0]) {
            const file = e.target.files[0];
            const localImageUrl = URL.createObjectURL(file);
            setPreviewPic(localImageUrl);
        }
    };

    return (
        <SettingSection icon={User} title={"Profile"}>
            <div className="flex flex-col sm:flex-row items-center mb-6">
                <img
                    src={previewPic}
                    alt="Profile"
                    className="rounded-full w-20 h-20 object-cover mr-4 border border-border shadow-sm"
                />

                <div>
                    {editing ? (
                        <>
                            <input
                                type="text"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                placeholder="Enter your name"
                                className="text-lg font-semibold text-text mb-2 p-2 border border-border rounded shadow-sm focus:outline-none focus:ring-2 focus:ring-buttonDefaultAccent1"
                            />
                            <input
                                type="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                placeholder="Enter your email"
                                className="text-cardTextHeading p-2 border border-border rounded shadow-sm focus:outline-none focus:ring-2 focus:ring-buttonDefaultAccent1"
                            />
                        </>
                    ) : (
                        <>
                            <h3 className="text-lg font-semibold text-text">{name}</h3>
                            <p className="text-cardTextHeading">{email}</p>
                        </>
                    )}
                </div>
            </div>

            {editing && (
                <div className="mb-6">
                    <label className="block mb-2 text-text font-medium">
                        Change Profile Picture
                    </label>
                    <input
                        type="file"
                        accept="image/*"
                        onChange={handleImageChange}
                        className="p-2 border border-border rounded shadow-sm focus:outline-none focus:ring-2 focus:ring-buttonDefaultAccent1"
                    />
                </div>
            )}

            {editing ? (
                <div className="flex gap-4">
                    <button
                        onClick={handleSave}
                        className="bg-buttonDefaultAccent1 text-white font-bold py-2 px-4 rounded transition duration-200 w-full sm:w-auto"
                    >
                        Save
                    </button>
                    <button
                        onClick={handleCancel}
                        className="bg-buttonHoverAccent1 text-white font-bold py-2 px-4 rounded transition duration-200 w-full sm:w-auto"
                    >
                        Cancel
                    </button>
                </div>
            ) : (
                <button
                    onClick={handleEdit}
                    data-cy="edit-profile-button"
                    className="bg-buttonDefaultAccent1 text-white font-bold py-2 px-4 rounded transition duration-200 w-full sm:w-auto"
                >
                    Edit Profile
                </button>
            )}
        </SettingSection>
    );
};

export default Profile;
