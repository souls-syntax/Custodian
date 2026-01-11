// Configuration
const API_URL = "https://api.custodianbackend.shop/v1/verify";

// DOM Elements
const inputField = document.getElementById('user-input');
const submitBtn = document.getElementById('submit-btn');
const responseArea = document.getElementById('response-area');
const loader = document.getElementById('loader');
const container = document.getElementById('main-container');

// Event Listeners
inputField.addEventListener("keypress", function(event) {
    if (event.key === "Enter") {
        event.preventDefault();
        sendQuery();
    }
});

submitBtn.addEventListener("click", sendQuery);

// Main Logic
async function sendQuery() {
    const queryText = inputField.value.trim();
    if (!queryText) return;

    // UI State: Loading
    container.classList.add('has-results');
    loader.style.display = 'block';
    responseArea.style.display = 'none';
    responseArea.innerHTML = ''; 

    // Payload Construction
    const payload = {
        query: queryText,
        userid: null // Nullable per struct definition
    };

    try {
        const res = await fetch(API_URL, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!res.ok) throw new Error(`Status Code: ${res.status}`);

        const data = await res.json();
        renderResponse(data);

    } catch (error) {
        console.error(error);
        responseArea.innerHTML = `
            <div style="color: #ef4444; padding: 10px;">
                <strong>Connection Error:</strong> ${error.message}<br>
                <small>Is the backend running on ${API_URL}?</small>
            </div>`;
        responseArea.style.display = 'block';
    } finally {
        loader.style.display = 'none';
    }
}

function renderResponse(data) {
    // 1. Calculate Confidence Percentage
    const confidencePct = (data.decision_confidence * 100).toFixed(1);
    
    // 2. Determine Styling based on Verdict
    let verdictDisplay = data.decision_verdict;
    let verdictColor = "var(--text-color)";

    switch (data.decision_verdict) {
        case "likely_true":
            verdictDisplay = "likely True"
            verdictColor = "var(--accent-green)";
            break;
        case "likely_false":
            verdictDisplay = "likely False"
            verdictColor = "var(--accent-red)";
            break;
        default:
            verdictDisplay = "confusing and is unsure";
            verdictColor = "var(--accent-amber)";
            break;
    }

    // 3. Construct HTML
    const htmlContent = `
        <div class="bot-message">
            <div class="verdict-text">
                Our Intelligence thinks that this query is 
                <strong style="color: ${verdictColor}">${verdictDisplay}</strong> 
                with this much confidence (<strong>${confidencePct}%</strong>).
            </div>
            <br>
            <div style="font-size: 0.9rem; color: #888;">
                <em>Tip: If this was your first time asking it you would get better result if asked second time.</em>
            </div>
        </div>

        <div class="debug-toggle" onclick="toggleJson(this)">[+] View Raw Intel Packet</div>
        <div class="json-dump">
            <pre>${JSON.stringify(data, null, 2)}</pre>
        </div>
    `;

    responseArea.innerHTML = htmlContent;
    responseArea.style.display = 'block';
}

// Helper to toggle the JSON view
window.toggleJson = function(element) {
    const nextSibling = element.nextElementSibling;
    if (nextSibling.style.display === "block") {
        nextSibling.style.display = "none";
        element.innerText = "[+] View Raw Intel Packet";
    } else {
        nextSibling.style.display = "block";
        element.innerText = "[-] Hide Raw Intel Packet";
    }
};
