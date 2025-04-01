document.addEventListener('DOMContentLoaded', () => {
    const urlInput = document.getElementById('relay-url');
    const tokenInput = document.getElementById('relay-token');
    const connectButton = document.getElementById('connect');
    const statusSpan = document.getElementById('status');

    // Load saved configuration on popup open
    chrome.storage.sync.get(['relayUrl', 'relayToken'], (result) => {
        urlInput.value = result.relayUrl || ''; // Use empty string if not set
        tokenInput.value = result.relayToken || ''; // Use empty string if not set
        // Initial status could be improved later by checking connection
        statusSpan.textContent = (result.relayUrl && result.relayToken) ? 'Configured' : 'Not Configured';
        console.log('Popup loaded config:', result);
    });

    // Save configuration on button click
    connectButton.addEventListener('click', () => {
        const relayUrl = urlInput.value.trim();
        const relayToken = tokenInput.value.trim();

        if (!relayUrl) {
             statusSpan.textContent = 'Error: Relay URL is required.';
             statusSpan.style.color = 'red';
             return;
        }
        // Basic URL validation (optional but recommended)
        try {
            new URL(relayUrl); // Check if it's a valid URL structure
        } catch (_) {
            statusSpan.textContent = 'Error: Invalid Relay URL format.';
            statusSpan.style.color = 'red';
            return;
        }


        chrome.storage.sync.set({ relayUrl, relayToken }, () => {
            console.log('Configuration saved:', { relayUrl, relayToken });
            statusSpan.textContent = 'Saved!';
            statusSpan.style.color = 'green';

            // Optionally, send a message to background script to update immediately
            chrome.runtime.sendMessage({ type: 'CONFIG_UPDATED' }, (response) => {
                if (chrome.runtime.lastError) {
                    console.warn("Could not send CONFIG_UPDATED message:", chrome.runtime.lastError.message);
                } else {
                    console.log("Background script notified of config update.", response);
                }
            });


            // Maybe close the popup after a short delay
            // setTimeout(() => window.close(), 1000);
        });
    });
});
