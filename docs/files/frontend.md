
# Core Concepts & Features

### Usage Tracking
* **What is Tracked:** The app tracks the time spent on websites (primarily via browser activity) and the usage of standalone browsers (like Mozilla, Ch, etc.).
* **Active Screen Time (`Time elapsed`):** Measures the total time you are actively using your device (with keyboard or mouse input) while the app's tracking is running. This metric resets automatically each day.
* **Idle Time (`Idle Time`):** Measures the time when no keyboard or mouse activity is detected while the app is tracking (essentially AFK/Away From Keyboard time).
* **Data Usage:** Tracks the amount of network data consumed during the monitored period.

### Activity Categorization
* **Automatic Categorization:** When you visit a website, the app automatically assigns it a category (e.g., Educational, Entertainment, Social Media, eCommerce, Blog) based on backend flags and AI analysis.
* **Manual Editing:** You can manually change the category assigned to any website via the **Tabs** section.
* **Scope:** Categorization primarily applies to websites visited within browsers. Standalone applications are tracked but generally not assigned these types of categories.
* **Limitations:** Currently, there is no option for bulk-categorizing websites or creating custom categorization rules.

### Tab Management (Tabs Section)
* **Tracking:** Monitors your open browser tabs.
* **States:**
    * `Active Tabs`: Tabs currently loaded and potentially using system resources.
    * `Inactive / Hybernating Tabs`: Tabs that haven't been viewed or used for a while. The app actively determines this status. The `Hybernating Tabs` count reflects the total number of tabs currently in this 'Inactive' state.
    * `Awake Rate`: The percentage of your total open tabs that are currently 'Active'.
* **Actions:**
    
    * `Delete`: Closes the actual tab in your web browser.

### Data Integration (Connected Accounts)
* You can connect external accounts (e.g., Google, Twitter, GitHub) via the **Settings** section.
* **Purpose:** This is used to **integrate data** from those services (e.g., calendar events, repository activity) to provide richer context and more accurate insights into your productivity patterns. It is not solely for login purposes.

## 4. Navigating the App

### Overview Tab
Your main dashboard for a quick summary of your daily activity.

* **Key Metrics:**
    * `Time elapsed`: Total active screen time tracked today (resets daily).
    * `Ideal Time`: Total idle (AFK) time tracked today. *(See definition above)*.
    * `Tabs open`: Current number of open browser tabs.
    * `Data Used`: Network data consumed today.
* **Time Distribution (Pie Chart):** Shows the percentage breakdown of your active time spent on specific top websites *for the current day*.
* **Website Distribution (Bar Chart):** Shows time spent across different website *categories* (Educational, Entertainment, etc.) for the current day. 


### Tabs Tab
Manage and analyze your browser tabs.

* **Key Metrics:** `Total Tabs`, `Active Tabs`, `Hybernating Tabs` (count of inactive tabs), `Awake Rate`.
* **Tabs List:** Detailed view of each open tab:
    * Website Icon & Domain
    * Page Title (`Name`)
    * Category (`Type`) - Editable
    * Status (`Active` / `Inactive`)
    * Actions ( `Delete` to close tab)
* **Search:** Find specific open tabs quickly.
* **Website Distribution (Pie Chart):** Shows the percentage breakdown of your *currently open tabs* by category (e.g., 30% Social, 20% Educational). This differs from the Overview chart as it's based on tab count, not time spent.

### Analytics Tab
Dive deeper into trends and progress.

* **Top Items:** Displays your most used Application, Window Title, Browser Domain, and Website Category for the selected period. 
* **Actual vs Target (Graph):** Visualizes your actual tracked hours against a dynamic target over time.
    * `Actual`: Your recorded active hours.
    * `Target`: A dynamically adjusted optimal productivity time, calculated by the app based on general recommendations and your past performance/behaviour. It is *not* manually set by the user.
    * `Time Period`: You can select the time frame for analysis (e.g., This Week, This Month, This Year, Quarter) using the dropdown menu. 
* **AI-Powered Insights:**
    * ***(Experimental Feature Note):*** *This section is designed to provide future AI-driven analysis, explanations for trends, and actionable suggestions based on your data. 

### Settings Tab
Configure your profile, preferences, and account settings.

* **Profile:** View and edit your name and email address.
* **Notifications:** Enable or disable notifications via different channels:
    * `Push Notifications`
    * `Email Notifications`
    * `SMS Notifications`
    * **Triggers:** Notifications may include general reminders, daily/weekly reports, app updates, productivity tips, and promotional messages.
* **Security:**
    * `Two-Factor Authentication`: Option to enhance account security (toggle available).
    * `Change Password`.
* **Connected Accounts:** Manage linked accounts (Google, Twitter, GitHub, etc.) used for data integration.
* **Danger Zone:**
    * `Delete Account`: Permanently removes your account and all associated data from **both the servers and your local machine**. This action is irreversible.

## 5. Understanding Your Productivity Metrics

* **Time elapsed:** Daily active screen time while tracked. Helps understand total engagement.
* **Ideal Time (Idle Time):** Daily inactive/AFK time. Useful for seeing break times or periods away.
* **Tabs open:** Current tab count. High numbers might indicate multitasking or information overload.
* **Data Used:** Network usage. Relevant for users on limited data plans.
* **Active Tabs / Hybernating Tabs / Awake Rate:** Metrics related to tab management and potential system resource usage by the browser. Lower awake rates might suggest many unused tabs are open.
* **Time Distribution (Overview):** Identifies top time-consuming websites daily.
* **Website Distribution (Overview):** Shows daily time allocation across activity types (work, learning, leisure).
* **Website Distribution (Tabs):** Shows the categorical makeup of your currently open tabs (focus vs. distraction).
* **Actual vs Target (Analytics):** Tracks progress towards a dynamically adjusted productivity goal over selected periods.
* **Top Items (Analytics):** Highlights the primary applications, domains, and categories used over the selected period.



## 6. Frequently Asked Questions (FAQ)

* **Q: What does "Ideal Time" mean?**
    * A: In the current version, "Ideal Time" represents the time you were idle or away from your keyboard/mouse while the app was tracking.
* **Q: How are websites categorized? Can I change a category?**
    * A: Websites are categorized automatically by the app using backend flags and AI. You can manually change the category for any site using the 'Edit' button in the 'Tabs' section. Bulk editing or custom rules are not currently supported.
* **Q: Does the app track applications too?**
    * A: Yes, the app tracks time spent in standalone applications as well as websites, although categorization primarily applies to websites.
* **Q: What happens when I click 'Delete' next to a tab?**
    * A: It closes the actual tab in your web browser.
* **Q: How is the 'Target' in the Analytics graph determined?**
    * A: The Target is dynamically calculated by the app based on generally accepted optimal productivity times and adjusted based on your own performance patterns. You do not set it manually.
* **Q: Is my data secure?**
    * A: _[Refer to Privacy section summary and link to full policy]_
* **Q: How do I delete my account and data?**
    * A: Go to Settings -> Danger Zone and click 'Delete Account'. This removes data from servers and your local machine. It cannot be undone.
