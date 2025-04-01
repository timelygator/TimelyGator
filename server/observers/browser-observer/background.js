console.log("TimelyGator Browser Observer - Background script loaded.");

let relayUrl = null;
let relayToken = null;

async function loadConfig() {
  try {
    // Fetch the relay URL and token from sync storage
    const items = await chrome.storage.sync.get({ relayUrl: null, relayToken: null });
    relayUrl = items.relayUrl;
    relayToken = items.relayToken;
    if (relayUrl) {
        console.log("Configuration loaded. Relay URL:", relayUrl, "Token provided:", !!relayToken);
    } else {
        console.log("Configuration loaded. Relay URL not set.");
    }
  } catch (error) {
    console.error("Error loading configuration:", error);
    relayUrl = null;
    relayToken = null;
  }
}

chrome.storage.onChanged.addListener((changes, namespace) => {
  if (namespace === 'sync') {
    let configChanged = false;
    if (changes.relayUrl) {
      relayUrl = changes.relayUrl.newValue || null;
      console.log("Relay URL configuration changed to:", relayUrl);
      configChanged = true;
    }
    if (changes.relayToken) {
      relayToken = changes.relayToken.newValue || null;
      console.log("Relay Token configuration changed.");
      configChanged = true;
    }
  }
});

// --- Event Sending Function ---
async function sendEventToServer(eventType, data) {
  if (relayUrl === null) {
      console.log("Relay URL is null, attempting to load configuration.");
      await loadConfig();
  }

  if (!relayUrl) {
      console.warn("Cannot send event: Relay URL is not configured.");
      return;
  }

  const payload = {
    observer: "browser",
    event_type: eventType,
    timestamp: new Date().toISOString(),
    data: data,
  };

  console.log(`Sending event to ${relayUrl}:`, payload);

  try {
    const headers = {
      "Content-Type": "application/json",
    };
    if (relayToken) {
        headers["Authorization"] = `Bearer ${relayToken}`;
    }

    const response = await fetch(relayUrl, {
      method: "POST",
      headers: headers,
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      console.error(`Error sending event: ${response.status} ${response.statusText}`);
    } else {
      console.log("Event sent successfully.");
    }
  } catch (error) {
    console.error("Network error or server unreachable:", relayUrl, error);
  }
}

// --- Tab Event Listeners ---

// Get initial state when the extension loads/wakes up
async function getInitialState() {
    await loadConfig();
    try {
        const tabs = await chrome.tabs.query({});
        const activeTabs = await chrome.tabs.query({ active: true, lastFocusedWindow: true });
        const activeTab = activeTabs.length > 0 ? activeTabs[0] : null;

        sendEventToServer("initial_state", {
            totalTabs: tabs.length,
            activeTab: activeTab ? { id: activeTab.id, title: activeTab.title, url: activeTab.url } : null,
        });
    } catch (error) {
        console.error("Error getting initial state:", error);
    }
}


// Fired when the active tab in a window changes
chrome.tabs.onActivated.addListener(async (activeInfo) => {
  console.log("Tab activated:", activeInfo);
  try {
    const tab = await chrome.tabs.get(activeInfo.tabId);
    // The tab might not have loaded its title/url yet, especially if just created.
    // onUpdated will catch the title/url when they are ready.
    if (tab) {
        const allTabs = await chrome.tabs.query({});
        sendEventToServer("tab_activated", {
            activeTab: { id: tab.id, title: tab.title, url: tab.url },
            totalTabs: allTabs.length
        });
    }
  } catch (error) {
    console.error("Error getting activated tab details:", error);
     // Might happen if the tab is closed quickly after activation
     const allTabs = await chrome.tabs.query({});
     sendEventToServer("tab_activated_error", { error: error.message, totalTabs: allTabs.length });
  }
});

// Fired when a tab is updated (e.g., URL change, title change, loading state)
chrome.tabs.onUpdated.addListener(async (tabId, changeInfo, tab) => {
  // We only care about updates to the title or URL of the *active* tab
  if (tab.active && (changeInfo.title || changeInfo.url)) {
    console.log("Active tab updated:", tabId, changeInfo);
    const allTabs = await chrome.tabs.query({});
sendEventToServer("tab_updated", {
    activeTab: { id: tab.id, title: tab.title, url: tab.url },
    totalTabs: allTabs.length,
    changes: changeInfo // Include what changed (e.g., {title: "New Title"})
});
}
// We could also track total tabs here if a tab finishes loading,
// but onCreated/onRemoved are more direct for tab count changes.
});

// Fired when a new tab is created
chrome.tabs.onCreated.addListener(async (tab) => {
    console.log("Tab created:", tab.id);
    // Query for total tabs *after* the new tab is potentially registered
    setTimeout(async () => { // Use a small delay as querying immediately might not include the new tab yet
        const allTabs = await chrome.tabs.query({});
        sendEventToServer("tab_created", { newTabId: tab.id, totalTabs: allTabs.length });
    }, 100); // 100ms delay, adjust if needed
});

// Fired when a tab is closed
chrome.tabs.onRemoved.addListener(async (tabId, removeInfo) => {
    console.log("Tab removed:", tabId);
    // If the window is closing, we might not be able to query accurately
    if (!removeInfo.isWindowClosing) {
        // Query for total tabs *after* the tab is potentially removed
       setTimeout(async () => { // Use a small delay
           const allTabs = await chrome.tabs.query({});
           sendEventToServer("tab_removed", { removedTabId: tabId, totalTabs: allTabs.length });
       }, 100);
   } else {
       console.log("Window is closing, not sending tab removed count event.");
        sendEventToServer("window_closing", { removedTabId: tabId });
   }
});

// --- Initialization ---
chrome.runtime.onInstalled.addListener(async (details) => {
  console.log("Extension installed/updated.", details.reason);
  await loadConfig();
  getInitialState();
});

chrome.runtime.onStartup.addListener(async () => {
    console.log("Browser startup detected.");
    await loadConfig();
    getInitialState();
});

// Listen for messages from the popup (e.g., config update notification)
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    console.log("Message received:", message);
    if (message.type === 'CONFIG_UPDATED') {
        console.log("Received config update message, reloading configuration...");
        loadConfig().then(() => {
            // You could add logic here, e.g., retry sending queued events
            sendResponse({ status: "Config reloaded" });
        }).catch(error => {
            console.error("Failed to reload config after message:", error);
            sendResponse({ status: "Error reloading config", error: error.message });
        });
        return true; // Indicates response will be sent asynchronously
    }
    return false;
});

getInitialState(); 
